use crate::stream;
use crate::utils;
use crate::utils::SpotPosition;
use crate::utils::SpotState;
use crate::CameraID;
use crate::CameraMetadata;
use crate::CameraState;
use crate::Model;
use crate::Spot;
use anyhow::Ok;
use dashmap::DashMap;
use image::RgbImage;
use itertools::Itertools;
use rand::Rng;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::watch;
use tokio::sync::Mutex;

#[derive(Debug)]
struct CameraWorker {
    images: watch::Receiver<Option<RgbImage>>,
    // metadata: CameraMetadata,
    positions: Vec<SpotPosition>,
    // Hash map of random indexes and its incoming positions
    // Used to later count how many SpotPosition exists in this whole HashMap. Compared by IOU.
    incoming_positions: HashMap<usize, SpotPosition>,
}

// Minimum amount of times that SpotPosition has to appear in the pre-staging list in order to be staged.
const MINIMUM_SCORE: usize = 5;

// Minimum IOU for two SpotPositions to be considered overlapping.
const OVERLAPPING_MIN_IOU: f32 = 0.7;

// Maximum IOU for two SpotPositions to be considered non-overlapping and thus unjoint.
const UNJOINT_MAX_IOU: f32 = 0.15;

impl CameraWorker {
    fn create(metadata: CameraMetadata) -> anyhow::Result<Self> {
        let images = stream::images(&metadata.url)?;
        Ok(Self {
            images,
            // metadata,
            positions: Vec::new(),
            incoming_positions: HashMap::new(),
        })
    }

    async fn image(&mut self) -> anyhow::Result<RgbImage> {
        loop {
            self.images.changed().await?;
            if let Some(image) = self.images.borrow().clone() {
                return Ok(image);
            }
        }
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

#[derive(Debug)]
pub struct Worker {
    model: Model,
    cameras: DashMap<CameraID, Arc<Mutex<CameraWorker>>>,
    state: DashMap<CameraID, CameraState>,
    visualizations: DashMap<CameraID, RgbImage>,
}

impl Worker {
    pub fn create(
        model: Model,
        cameras: impl Iterator<Item = (CameraID, CameraMetadata)>,
    ) -> anyhow::Result<Self> {
        let cameras = cameras
            .map(|(id, metadata)| {
                let worker = CameraWorker::create(metadata)?;
                Ok((id, Arc::new(Mutex::new(worker))))
            })
            .collect::<anyhow::Result<_>>()?;
        Ok(Self {
            model,
            cameras,
            state: DashMap::new(),
            visualizations: DashMap::new(),
        })
    }

    pub fn add(&self, id: CameraID, metadata: CameraMetadata) {
        let worker = CameraWorker::create(metadata).unwrap();
        self.cameras.insert(id, Arc::new(Mutex::new(worker)));
    }

    pub fn cameras(&self) -> usize {
        self.cameras.len()
    }

    pub fn state_of(&self, id: &CameraID) -> Option<CameraState> {
        self.state.get(id).map(|state| state.value().clone())
    }

    pub fn visualization_of(&self, id: &CameraID) -> Option<RgbImage> {
        self.visualizations
            .get(id)
            .map(|image| image.value().clone())
    }

    pub async fn run(&self) -> anyhow::Result<()> {
        loop {
            self.work().await?;
        }
    }

    pub async fn work(&self) -> anyhow::Result<()> {
        let ids = self
            .cameras
            .iter()
            .map(|camera| camera.key().clone())
            .collect_vec();
        let images = ids.iter().map(|id| async move {
            let worker = self.cameras.get(id).unwrap();
            let mut worker = worker.lock().await;
            let image = worker.image().await?;
            anyhow::Ok(image)
        });
        let images = futures::future::try_join_all(images).await?;
        tracing::info!("collected {} images", images.len());
        let start = std::time::Instant::now();
        let predictions = self.model.infere(&images)?;
        tracing::debug!(
            "inference for {} images finished after {}ms",
            images.len(),
            start.elapsed().as_millis()
        );
        let state = predictions
            .into_iter()
            .enumerate()
            .map(|(i, objects)| (&ids[i], objects))
            .map(|(id, objects)| async move {
                let positions = objects
                    .into_iter()
                    .filter(|object| ["car", "bus", "truck"].contains(&object.label))
                    .map(|vehicle| SpotPosition {
                        bbox: vehicle.bbox,
                        contours: Arc::new(vehicle.contours),
                    })
                    .collect::<Vec<_>>();

                let spots = {
                    let worker = self.cameras.get(id).unwrap();
                    let mut worker = worker.lock().await;
                    worker.update(positions)?
                };
                tracing::debug!(%id, "spots update");
                Ok((id.clone(), spots))
            });
        let spots = futures::future::try_join_all(state).await?;

        for (id, spots) in spots {
            let available_spots = spots
                .iter()
                .filter(|spot| spot.state == SpotState::Vacant)
                .count();
            let state = CameraState {
                total_spots: spots.len() as u32,
                available_spots: available_spots as u32,
            };
            tracing::debug!(%id, ?state, "state update");
            self.state.insert(id.clone(), state);
            let index = ids.iter().position(|i| *i == id).unwrap();
            let visualization = utils::visualize_spots(&images[index], &spots);
            self.visualizations.insert(id, visualization);
        }

        Ok(())
    }
}
