mod model;
mod server;
mod vis;
mod utils;
mod worker;

use model::Model;
pub use utils::BoundingBox;
pub use utils::Object;
pub use utils::Point;
pub use utils::Spot;
pub use worker::Master;

use anyhow::Context;
use server::ServerState;
use std::sync::Arc;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let project_directories =
        directories::ProjectDirs::from("app", "wheretopark", "caman").unwrap();

    let model = Model::new(project_directories)?;
    let worker = Master::create(model)?;
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
