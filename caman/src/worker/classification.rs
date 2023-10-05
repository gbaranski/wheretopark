use crate::vis;
use crate::utils::SpotPosition;
use crate::utils::SpotState;
use crate::Spot;
use anyhow::Ok;
use image::RgbImage;
use itertools::Itertools;
use rand::Rng;
use wheretopark_caman::CameraMetadata;
use std::collections::HashMap;
use tokio::sync::watch;

#[derive(Debug)]
pub struct ClassificationWorker {
    pub(super) images: watch::Receiver<Option<RgbImage>>,
    positions: Vec<SpotPosition>,
}

// Minimum amount of times that SpotPosition has to appear in the pre-staging list in order to be staged.
const MINIMUM_SCORE: usize = 5;

// Minimum IOU for two SpotPositions to be considered overlapping.
const OVERLAPPING_MIN_IOU: f32 = 0.7;

// Maximum IOU for two SpotPositions to be considered non-overlapping and thus unjoint.
const UNJOINT_MAX_IOU: f32 = 0.15;

impl ClassificationWorker {
    fn create(metadata: CameraMetadata) -> anyhow::Result<Self> {
        let images = vis::images(&metadata.url)?;
        Ok(Self {
            images,
            // metadata,
            positions: Vec::new(),
            incoming_positions: HashMap::new(),
        })
    }

    fn update(&mut self, positions: Vec<SpotPosition>) -> anyhow::Result<Vec<Spot>> {
        let mut rng = rand::thread_rng();
        tracing::debug!(
            "self.positions.len() = {}. incoming = {}",
            self.positions.len(),
            self.incoming_positions.len()
        );
        // steps:
        // 1. calculate overlaps
        // 2. mark overlapped positions as occupied
        // 3. mark unjoint positions from self.positions as vacant
        // 4. update self.positions with new positions

        let mut occupied = Vec::new();
        for pos in &positions {
            let overlap = self
                .positions
                .iter()
                .map(|other| pos.bbox.iou(&other.bbox))
                .max_by(|x, y| x.partial_cmp(y).unwrap())
                .unwrap_or(0.0);
            if overlap > OVERLAPPING_MIN_IOU {
                // occupied
                occupied.push(pos.clone());
            } else {
                // new spot, not registered in self.positions
                // self.incoming_positions.push(pos.clone());
                let id: usize = rng.gen();
                let ret = self.incoming_positions.insert(id, pos.clone());
                assert!(ret.is_none());
            }
        }

        // TODO: Maybe change that to Iterator::fold()?
        for (_, pos) in self.incoming_positions.clone() {
            let overlaps = self
                .incoming_positions
                .iter()
                .map(|(j, other)| (j, pos.bbox.iou(&other.bbox)))
                .filter(|(_, iou)| *iou > OVERLAPPING_MIN_IOU)
                .map(|(j, _)| *j)
                .collect_vec();
            if overlaps.len() > MINIMUM_SCORE {
                self.positions.push(pos.clone());
                for j in overlaps {
                    self.incoming_positions.remove(&j);
                }
            }
        }

        let mut vacant = Vec::new();
        // TODO: Maybe consider Vec::retain()?
        for pos in &self.positions {
            let overlap = positions
                .iter()
                .map(|other| pos.bbox.iou(&other.bbox))
                .max_by(|x, y| x.partial_cmp(y).unwrap())
                .unwrap_or(0.0);
            if overlap < UNJOINT_MAX_IOU {
                // vacant
                vacant.push(pos.clone());
            }
        }

        let occupied = occupied.into_iter().map(|position| Spot {
            position,
            state: SpotState::Occupied,
        });
        let vacant = vacant.into_iter().map(|position| Spot {
            position,
            state: SpotState::Vacant,
        });

        let spots = occupied.chain(vacant).collect_vec();
        Ok(spots)
    }
}
