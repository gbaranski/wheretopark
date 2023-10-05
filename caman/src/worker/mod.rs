use crate::model::Model;
use crate::utils::SpotPosition;
use crate::utils::SpotState;
use Position;
use dashmap::DashMap;
use futures::Future;
use futures::StreamExt;
use futures::TryStreamExt;
use futures::stream::FuturesUnordered;
use tokio::sync::watch;
use crate::utils::Spot
use image::RgbImage;
use itertools::Either;
use itertools::Itertools;
use std::collections::HashMap;
use std::sync::Arc;
use std::sync::Mutex;
use wheretopark_caman::CameraID;
use wheretopark_caman::CameraMetadata;
use wheretopark_caman::CameraState;

mod classification;
mod detection;

use classification::ClassificationWorker;
use detection::DetectionWorker;

#[derive(Debug)]
enum CameraWorker {
    Detection(DetectionWorker),
    Classification(ClassificationWorker),
}

async fn get_image(rx: &mut watch::Receiver<Option<RgbImage>>) -> anyhow::Result<RgbImage> {
    loop {
        rx.changed().await?;
        if let Some(image) = rx.borrow().clone() {
            return Ok(image);
        }
    }
}

impl CameraWorker {
    async fn image(&mut self) -> anyhow::Result<RgbImage> {
        let images = match self {
            Self::Detection(worker) => &mut worker.images,
            Self::Classification(worker) => &mut worker.images,
        };
        get_image(images).await
    }
}

#[derive(Debug)]
pub struct Master {
    model: Model,
    workers: Mutex<HashMap<CameraID, CameraWorker>>,
    state: DashMap<CameraID, CameraState>,
    visualizations: DashMap<CameraID, RgbImage>,
}

impl Master {
    pub fn create(model: Model) -> anyhow::Result<Self> {
        Ok(Self {
            model,
            workers: Mutex::new(HashMap::new()),
            state: DashMap::new(),
            visualizations: DashMap::new(),
        })
    }

    pub fn add(&self, id: CameraID, metadata: CameraMetadata) {
        // let worker = CameraDetectionWorker::create(metadata).unwrap();
        // self.cameras.insert(id, Arc::new(Mutex::new(worker)));
    }

    pub fn state(&self) -> HashMap<CameraID, CameraState> {
        self.state
            .iter()
            .map(|entry| (entry.key().clone(), entry.clone()))
            .collect()
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

    async fn work_detection(
        &self,
        workers: Vec<(&CameraID, &mut DetectionWorker)>,
    ) -> anyhow::Result<()> {
        // let images = futures::future::try_join_all(images).await?;
        // let image_date = chrono::Utc::now();
        // tracing::info!("collected {} images", images.len());
        // let start = std::time::Instant::now();
        // let predictions = self.model.detect(&images)?;
        // tracing::debug!(
        //     "object detection inference for {} images finished after {}ms",
        //     images.len(),
        //     start.elapsed().as_millis()
        // );

        // let state = predictions
        //     .into_iter()
        //     .enumerate()
        //     .map(|(i, objects)| (&workers[i], objects))
        //     .map(|((id, worker), objects)| async move {
        //         let positions = objects
        //             .into_iter()
        //             .filter(|object| ["car", "bus", "truck"].contains(&object.label))
        //             .map(|vehicle| SpotPosition {
        //                 bbox: vehicle.bbox,
        //                 contours: Arc::new(vehicle.contours),
        //             })
        //             .collect::<Vec<_>>();

        //         let spots = worker.update(positions)?;
        //         Ok((id.clone(), spots))
        //     });
        // let spots = futures::future::try_join_all(state).await?;

        // for (id, spots) in spots {
        //     let available_spots = spots
        //         .iter()
        //         .filter(|spot| spot.state == SpotState::Vacant)
        //         .count();
        //     let state = CameraState {
        //         last_updated: image_date,
        //         total_spots: spots.len() as u32,
        //         available_spots: available_spots as u32,
        //     };
        //     tracing::debug!(%id, ?state, "state update");
        //     self.state.insert(id.clone(), state);
        //     let index = ids.iter().position(|i| *i == id).unwrap();
        //     let visualization = utils::visualize_spots(&images[index], &spots);
        //     self.visualizations.insert(id, visualization);
        // }

        // Ok(())
        todo!()
    }

    async fn work_classification(&self, mut workers: HashMap<CameraID, &mut ClassificationWorker>) {
        let mut images = workers.iter_mut().map(|(id, worker)| async move {
            let result = get_image(&mut worker.images).await;
            (id, result)
        }).collect::<FuturesUnordered<_>>();

        while let Some((id, result)) = images.next().await {
            let worker = &workers[id];

        }
    }

    pub async fn work(&self) -> anyhow::Result<()> {
        let mut workers = self.workers.lock().unwrap();
        let (classification, detection): (Vec<_>, Vec<_>) =
            workers.iter_mut().partition_map(|(id, worker)| match worker {
                CameraWorker::Classification(worker) => Either::Left((id, worker)),
                CameraWorker::Detection(worker) => Either::Right((id, worker)),
            });
        self.work_classification(classification).await;


        // let images = workers.iter_mut().map(|(id, worker)| async move {
        //     let image = worker.image().await?;
        //     anyhow::Ok(image)
        // });
        // let images = images.collect::<FuturesUnordered<_>>();
        // let images = futures::future::try_join_all(images).await?;
        // let image_date = chrono::Utc::now();
        // tracing::info!("collected {} images", images.len());

        Ok(())
    }
}
