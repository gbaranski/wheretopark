use anyhow::Context;
use dashmap::DashMap;
use image::RgbImage;
use itertools::Itertools;
use std::collections::HashMap;
use tokio::sync::watch;

use crate::stream;
use crate::utils;
use crate::utils::SpotPosition;
use crate::CameraMap;
use crate::Model;

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
        let mut streams = self
            .cameras
            .iter()
            .map(|camera| {
                let metadata = &camera.metadata;
                let url = metadata.url.as_str();
                let images = stream::images(url.to_string())
                    .with_context(|| format!("create image stream for {url}"))?;
                Ok((camera.key().to_string(), images))
            })
            .collect::<anyhow::Result<HashMap<String, watch::Receiver<Option<RgbImage>>>>>()?;
        loop {
            let images =
                futures::future::join_all(streams.iter_mut().map(|(id, stream)| async move {
                    (
                        id.to_string(),
                        stream
                            .changed()
                            .await
                            .map(|_| stream.borrow().clone())
                            .unwrap(),
                    )
                }))
                .await;
            let images = images
                .into_iter()
                .filter_map(|(id, image)| image.map(|image| (id, image)))
                .collect::<HashMap<String, RgbImage>>();
            tracing::info!("collected {} images", images.len());
            for (id, image) in images {
                let start = std::time::Instant::now();
                let objects = self.model.infere(&image)?;
                tracing::debug!(
                    "inference for {id} finished after {}ms",
                    start.elapsed().as_millis()
                );
                dbg!(&objects);
                let visualization = utils::visualize_objects(&image, &objects);
                visualization.save(format!("{id}-vehicles.jpeg"))?;
                tracing::debug!("saved visualization for {id}");
            }
        }
    }

    // async fn work_camera(
    //     &self,
    //     id: &str,
    //     metadata: &CameraMetadata,
    // ) -> anyhow::Result<CameraState> {
    //     let vehicles = self.model.infere(&image)?;
    //     tracing::debug!(id=%id, "inference finished");
    //     let image_vehicles = utils::visualize_vehicles(&image, &vehicles);
    //     image_vehicles.save(format!("{id}-vehicles.jpeg"))?;

    //     let current_positions: Vec<SpotPosition> = vehicles
    //         .into_iter()
    //         .map(|vehicle| SpotPosition {
    //             bbox: vehicle.bbox,
    //             contours: Arc::new(vehicle.contours),
    //         })
    //         .collect();

    //     let positions = if let Some(positions) = self.positions.get(id) {
    //         positions
    //     } else {
    //         self.positions.insert(id.to_string(), current_positions);
    //         return Ok(CameraState::default());
    //     };

    //     let overlaps = utils::compute_overlaps(&positions, &current_positions);
    //     let spots = overlaps
    //         .iter()
    //         .zip(positions.iter().cloned())
    //         .map(|(overlap, position)| {
    //             let score = overlap.iter().cloned().reduce(f32::max).unwrap();
    //             Spot {
    //                 position,
    //                 state: SpotState { score },
    //             }
    //         })
    //         .collect::<Vec<_>>();

    //     let image_occupancy = utils::visualize_spots(&image, &spots);
    //     image_occupancy.save(format!("{id}-occupancy.jpeg"))?;
    //     let available_spots = spots.iter().filter(|spot| spot.state.score < 0.15).count();
    //     tracing::info!("available spots: {}", available_spots);
    //     Ok(CameraState {
    //         total_spots: current_positions.len() as u32,
    //         available_spots: available_spots as u32,
    //     })
    // }
}
