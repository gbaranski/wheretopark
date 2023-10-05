use crate::CameraMetadata;
use crate::Worker;
use axum::extract::Path;
use axum::extract::State;
use axum::http::header;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::response::Json;
use axum::routing::get;
use axum::Form;
use axum::routing::put;
use image::codecs::png::PngEncoder;
use serde::Serialize;
use std::io::BufWriter;
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

async fn put_metadata(
    State(app_state): State<ServerState>,
    Path(id): Path<String>,
    Form(metadata): Form<CameraMetadata>,
) -> impl IntoResponse {
    tracing::info!(id=%id, metadata=?metadata, "insert camera");
    app_state.worker.add(id, metadata);
    StatusCode::CREATED
}

async fn get_state(
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

async fn get_all_state(State(app_state): State<ServerState>) -> impl IntoResponse {
    tracing::info!("get all cameras");
    let state = app_state.worker.state();
    Json(state)
}

async fn visualize(
    State(app_state): State<ServerState>,
    Path(id): Path<String>,
) -> impl IntoResponse {
    tracing::info!(%id, "visualize occupancy");
    let visualization = app_state.worker.visualization_of(&id);
    match visualization {
        Some(image) => {
            let mut buf = vec![];
            let writer = BufWriter::new(&mut buf);
            let encoder = PngEncoder::new(writer);
            image.write_with_encoder(encoder).unwrap();
            (StatusCode::OK, [(header::CONTENT_TYPE, "image/png")], buf)
        }
        None => (
            StatusCode::NOT_FOUND,
            [(header::CONTENT_TYPE, "text/plain")],
            "visualization not found".as_bytes().to_vec(),
        ),
    }
}
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
        .route("/cameras/state", get(get_all_state))
        .route("/cameras/:id/state", get(get_state))
        .route("/cameras/:id/metadata", put(put_metadata))
        .route("/cameras/:id/visualize", get(visualize))
        .route("/status", get(status))
        .with_state(app_state);
    let url = "0.0.0.0:3000".parse().unwrap();
    tracing::info!("starting server at {}", url);
    axum::Server::bind(&url)
        .serve(app.into_make_service())
        .await
        .map_err(|err| err.into())
}
