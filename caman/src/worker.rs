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
use std::sync::Arc;
use tokio::sync::watch;
use tokio::sync::Mutex;

#[derive(Debug)]
struct CameraWorker {
    images: watch::Receiver<Option<RgbImage>>,
    // metadata: CameraMetadata,
    positions: Vec<SpotPosition>,
}

impl CameraWorker {
    fn create(metadata: CameraMetadata) -> anyhow::Result<Self> {
        let images = stream::images(&metadata.url)?;
        Ok(Self {
            images,
            // metadata,
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

    fn update(&mut self, positions: Vec<SpotPosition>) -> anyhow::Result<CameraState> {
        if self.positions.is_empty() {
            self.positions = positions;
            return Ok(CameraState::default());
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

        let available_spots = spots.iter().filter(|spot| spot.state.score < 0.15).count();
        let state = CameraState {
            total_spots: self.positions.len() as u32,
            available_spots: available_spots as u32,
        };
        Ok(state)
    }
}

#[derive(Debug)]
pub struct Worker {
    model: Model,
    cameras: DashMap<CameraID, Arc<Mutex<CameraWorker>>>,
    state: DashMap<CameraID, CameraState>,
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

                let state = {
                    let worker = self.cameras.get(id).unwrap();
                    let mut worker = worker.lock().await;
                    worker.update(positions)?
                };
                tracing::debug!(%id, ?state, "state update");
                Ok((id.clone(), state))
            });
        let state = futures::future::try_join_all(state).await?;

        for (id, state) in state {
            self.state.insert(id, state);
        }

        Ok(())
    }
}
