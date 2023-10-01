use crate::CameraMetadata;
use crate::Worker;
use axum::extract::Path;
use axum::extract::State;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::response::Json;
use axum::Form;
use axum::routing::get;
use serde::Serialize;
use std::sync::Arc;

#[derive(Debug, Clone)]
pub struct ServerState {
    worker: Arc<Worker>,
}

impl ServerState {
    pub fn new(worker: Arc<Worker>) -> Self {
        Self { worker }
    }
}

async fn post_camera(
    State(app_state): State<ServerState>,
    Path(id): Path<String>,
    Form(metadata): Form<CameraMetadata>,
) -> impl IntoResponse {
    tracing::info!(id=%id, metadata=?metadata, "insert camera");
    app_state.worker.add(id, metadata);
    StatusCode::CREATED
}

async fn get_camera(
    State(app_state): State<ServerState>,
    Path(id): Path<String>,
) -> impl IntoResponse {
    tracing::info!(id=%id, "get camera");
    let state = app_state.worker.state_of(&id).map(|state| state.clone());
    let status = if state.is_none() {
        StatusCode::NOT_FOUND
    } else {
        StatusCode::OK
    };
    (status, Json(state))
}

async fn visualize_occupancy() {}
async fn visualize_spots() {}

#[derive(Debug, Serialize)]
struct AppStatus {
    cameras: usize,
}

async fn status(State(app_state): State<ServerState>) -> impl IntoResponse {
    return Json(AppStatus {
        cameras: app_state.worker.cameras(),
    });
}

pub async fn run(app_state: ServerState) -> anyhow::Result<()> {
    let app = axum::Router::new()
        .route("/cameras/:id", get(get_camera).post(post_camera))
        .route("/cameras/:id/visualize/occupancy", get(visualize_occupancy))
        .route("/cameras/:id/visualize/spots", get(visualize_spots))
        .route("/status", get(status))
        .with_state(app_state);
    let url = "0.0.0.0:3000".parse().unwrap();
    tracing::info!("starting server at {}", url);
    axum::Server::bind(&url)
        .serve(app.into_make_service())
        .await
        .map_err(|err| err.into())
}
