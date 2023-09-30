use dashmap::DashMap;
use std::sync::Arc;

use crate::utils;
use crate::utils::SpotPosition;
use crate::utils::SpotState;
use crate::CameraMap;
use crate::CameraMetadata;
use crate::CameraState;
use crate::Model;
use crate::Spot;

#[derive(Debug)]
pub struct Worker {
    model: Model,
    cameras: CameraMap,
    positions: DashMap<String, Vec<SpotPosition>>,
}

impl Worker {
    pub fn new(model: Model, cameras: CameraMap) -> Self {
        Self {
            model,
            cameras,
            positions: DashMap::new(),
        }
    }

    pub async fn work(&self) -> anyhow::Result<()> {
        for camera in self.cameras.iter() {
            tracing::debug!(id=%camera.key(), "working on camera");
            let start = std::time::Instant::now();
            self.work_camera(camera.key(), &camera.metadata).await?;
            let duration = format!("{}ms", start.elapsed().as_millis());
            tracing::info!(id=%camera.key(), time=%duration, "work finished");
        }
        Ok(())
    }

    async fn work_camera(
        &self,
        id: &str,
        metadata: &CameraMetadata,
    ) -> anyhow::Result<CameraState> {
        let image = capture(metadata.url.as_str()).await?;
        tracing::debug!(id=%id, "captured image");
        let image = image.into_rgb8();
        let vehicles = self.model.infere(&image)?;
        tracing::debug!(id=%id, "inference finished");
        let image_vehicles = utils::visualize_vehicles(&image, &vehicles);
        image_vehicles.save(format!("{id}-vehicles.jpeg"))?;

        let current_positions: Vec<SpotPosition> = vehicles
            .into_iter()
            .map(|vehicle| SpotPosition {
                bbox: vehicle.bbox,
                contours: Arc::new(vehicle.contours),
            })
            .collect();

        let positions = if let Some(positions) = self.positions.get(id) {
            positions
        } else {
            self.positions.insert(id.to_string(), current_positions);
            return Ok(CameraState::default());
        };

        let overlaps = utils::compute_overlaps(&positions, &current_positions);
        let spots = overlaps
            .iter()
            .zip(positions.iter().cloned())
            .map(|(overlap, position)| {
                let score = overlap.iter().cloned().reduce(f32::max).unwrap();
                Spot {
                    position,
                    state: SpotState { score },
                }
            })
            .collect::<Vec<_>>();

        let image_occupancy = utils::visualize_spots(&image, &spots);
        image_occupancy.save(format!("{id}-occupancy.jpeg"))?;
        let available_spots = spots.iter().filter(|spot| spot.state.score < 0.15).count();
        tracing::info!("available spots: {}", available_spots);
        Ok(CameraState {
            total_spots: current_positions.len() as u32,
            available_spots: available_spots as u32,
        })
    }
}
