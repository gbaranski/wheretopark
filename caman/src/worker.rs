use std::sync::Arc;

use crate::stream;
use crate::utils;
use crate::utils::SpotPosition;
use crate::utils::SpotState;
use crate::CameraID;
use crate::CameraMetadata;
use crate::Model;
use crate::Spot;
use anyhow::Ok;
use image::RgbImage;
use indexmap::IndexMap;
use itertools::Itertools;
use tokio::sync::watch;
use tokio::sync::Mutex;

#[derive(Debug)]
struct CameraWorker {
    images: watch::Receiver<Option<RgbImage>>,
    metadata: CameraMetadata,
    positions: Vec<SpotPosition>,
}

impl CameraWorker {
    fn create(metadata: CameraMetadata) -> anyhow::Result<Self> {
        let images = stream::images(&metadata.url)?;
        Ok(Self {
            images,
            metadata,
            positions: vec![],
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

    async fn update(&mut self, positions: Vec<SpotPosition>) -> anyhow::Result<Vec<Spot>> {
        if self.positions.is_empty() {
            self.positions = positions;
            return Ok(vec![]);
        }

        let overlaps = utils::compute_overlaps(&positions, &self.positions);
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

        Ok(spots)
    }
}

#[derive(Debug)]
pub struct Worker {
    model: Model,
    cameras: IndexMap<CameraID, Arc<Mutex<CameraWorker>>>,
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
        Ok(Self { model, cameras })
    }

    pub async fn work(&self) -> anyhow::Result<()> {
        loop {
            let images = self.cameras.iter().map(|(_, camera)| async {
                let mut camera = camera.lock().await;
                camera.image().await
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
            for (i, objects) in predictions.into_iter().enumerate() {
                let id = &self.cameras.keys().nth(i).unwrap();
                let vehicles = objects
                    .into_iter()
                    .filter(|object| ["car", "bus", "truck"].contains(&object.label))
                    .collect_vec();
                // let visualization_objects = utils::visualize_objects(&images[i], &vehicles);
                // visualization_objects.save(format!("{id}-objects.jpeg"))?;
                // tracing::debug!(%id, "saved object visualization");

                let positions = vehicles
                    .into_iter()
                    .map(|vehicle| SpotPosition {
                        bbox: vehicle.bbox,
                        contours: Arc::new(vehicle.contours),
                    })
                    .collect::<Vec<_>>();
                let mut worker = self.cameras.get(*id).unwrap().lock().await;
                let spots = worker.update(positions).await?;
                drop(worker);
                // let visualization_occupancy = utils::visualize_spots(&images[i], &spots);
                // visualization_occupancy.save(format!("{id}-occupancy.jpeg"))?;
                // tracing::debug!("saved occupancy visualization to {}");
                let available_spots = spots.iter().filter(|spot| spot.state.score < 0.15).count();
                tracing::debug!(%id, "available spots: {}", available_spots);
            }
        }
    }
}
