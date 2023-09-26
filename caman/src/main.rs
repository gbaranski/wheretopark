mod model;
mod server;
mod stream;
mod utils;
mod worker;

pub use model::Model;
use url::Url;
pub use utils::BoundingBox;
pub use utils::Point;
pub use utils::Spot;
pub use utils::Vehicle;

use dashmap::DashMap;
use serde::Deserialize;
use serde::Serialize;
use std::sync::Arc;
use std::time::Duration;
use worker::Worker;

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
    let model = Model::new()?;
    let cameras = Arc::new(DashMap::new());
    cameras.insert(
        "u35s2dpn4t_0".to_string(),
        Camera {
            metadata: CameraMetadata {
                url: Url::parse("https://cam4out.klemit.net/hls/camn826.m3u8").unwrap(),
            },
            state: CameraState::default(),
        },
    );

    let worker = Worker::new(model, cameras.clone());
    // server::run(ServerState::new(cameras)).await;
    loop {
        if let Err(err) = worker.work().await {
            tracing::error!(%err, "work fail");
        }
        tokio::time::sleep(Duration::from_secs(5)).await;
    }
}
