mod model;
mod server;
mod stream;
mod utils;
mod worker;

pub use model::Model;
pub use utils::BoundingBox;
pub use utils::Object;
pub use utils::Point;
pub use utils::Spot;
pub use worker::Worker;

use anyhow::Context;
use serde::Deserialize;
use serde::Serialize;
use server::ServerState;
use std::collections::HashMap;
use std::path::PathBuf;
use std::sync::Arc;
use url::Url;

pub type CameraID = String;

#[derive(Debug, Deserialize)]
pub struct CameraMetadata {
    pub url: Url,
}

#[derive(Debug, Default, Clone, Serialize)]
pub struct CameraState {
    pub total_spots: u32,
    pub available_spots: u32,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let project_directories =
        directories::ProjectDirs::from("app", "wheretopark", "caman").unwrap();

    let model_path: PathBuf;
    if let Some(path) = std::option_env!("MODEL_PATH") {
        model_path = path.parse().unwrap();
    } else {
        model_path = project_directories.data_dir().join("yolov8x.onnx");
    }
    let model = Model::new(model_path)?;
    let mut cameras = HashMap::new();
    // cameras.insert(
    //     "u35s2sey91_1".to_string(),
    //     CameraMetadata {
    //         url: Url::parse("https://cam5out.klemit.net/hls/cammt852.m3u8").unwrap(),
    //     },
    // );
    // cameras.insert(
    //     "u35s3nvprd_0".to_string(),
    //     CameraMetadata {
    //         url: Url::parse("https://cam5out.klemit.net/hls/cammm841.m3u8").unwrap(),
    //     },
    // );
    // cameras.insert(
    //     "u35krvemdk_0".to_string(),
    //     CameraMetadata {
    //         url: Url::parse("https://cam4out.klemit.net/hls/camn583.m3u8").unwrap(),
    //     },
    // );
    cameras.insert(
        "u2gyfvc23d_0".to_string(),
        CameraMetadata {
            url: Url::parse(
                "http://91.238.55.4:5080/LiveApp/streams/435465478973256862461988.m3u8?token=null",
            )
            .unwrap(),
        },
    );

    let worker = Worker::create(model, cameras.into_iter())?;
    let worker = Arc::new(worker);
    tokio::select! {
        result = worker.run() => {
            result.context("worker failure")?;
        }
        result = server::run(ServerState::new(worker.clone())) => {
            result.context("server failure")?;
        }
    }

    Ok(())
}
