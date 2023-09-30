mod model;
mod server;
mod stream;
mod utils;
mod worker;

pub use model::Model;
pub use utils::BoundingBox;
pub use utils::Point;
pub use utils::Spot;
pub use utils::Object;
pub use worker::Worker;

use dashmap::DashMap;
use serde::Deserialize;
use serde::Serialize;
use std::path::PathBuf;
use std::sync::Arc;
use std::time::Duration;
use url::Url;
// use worker::Worker;

#[derive(Debug, Deserialize)]
struct CameraMetadata {
    pub url: Url,
}

#[derive(Debug, Default, Clone, Serialize)]
pub struct CameraState {
    pub total_spots: u32,
    pub available_spots: u32,
}

#[derive(Debug)]
pub struct Camera {
    metadata: CameraMetadata,
    state: CameraState,
}

pub type CameraMap = Arc<DashMap<String, Camera>>;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let project_directories =
        directories::ProjectDirs::from("app", "wheretopark", "caman").unwrap();

    let model_path: PathBuf;
    if let Some(path) = std::option_env!("MODEL_PATH") {
        model_path = path.parse().unwrap();
    } else {
        model_path = project_directories
            .data_dir()
            .join("yolov8x.onnx");
    }
    let model = Model::new(model_path)?;
    let cameras = Arc::new(DashMap::new());
    // cameras.insert(
    //     "u35s2sey91_1".to_string(),
    //     Camera {
    //         metadata: CameraMetadata {
    //             url: Url::parse("https://cam5out.klemit.net/hls/cammt852.m3u8").unwrap(),
    //         },
    //         state: CameraState::default(),
    //     },
    // );
    cameras.insert(
        "u35s3nvprd_0".to_string(),
        Camera {
            metadata: CameraMetadata {
                url: Url::parse("https://cam5out.klemit.net/hls/cammm841.m3u8").unwrap(),
            },
            state: CameraState::default(),
        },
    );

    let worker = Worker::new(model, cameras.clone());
    // // server::run(ServerState::new(cameras)).await;
    loop {
        if let Err(err) = worker.work().await {
            tracing::error!("work fail: {:#}", err);
        }
        tokio::time::sleep(Duration::from_secs(5)).await;
    }
}
