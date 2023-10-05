use crate::model::Model;
use crate::vis;
use crate::Master;
use axum::BoxError;
use axum::error_handling::HandleError;
use axum::error_handling::HandleErrorLayer;
use axum::extract::Path;
use axum::extract::State;
use axum::http::Uri;
use axum::http::header;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::response::Json;
use axum::routing::get;
use axum::routing::put;
use axum::Form;
use axum::http::Method;
use image::codecs::png::PngEncoder;
use serde::Deserialize;
use serde::Serialize;
use std::io::BufWriter;
use std::sync::Arc;
use url::Url;
use wheretopark_caman::CameraMetadata;
use wheretopark_caman::Spot;

#[derive(Debug, Clone)]
pub struct ServerState {
    model: Arc<Model>,
}

impl ServerState {
    pub fn new(model: Arc<Model>) -> Self {
        Self { model }
    }
}

// async fn visualize(
//     State(app_state): State<ServerState>,
//     Path(id): Path<String>,
// ) -> impl IntoResponse {
//     tracing::info!(%id, "visualize occupancy");
//     let visualization = app_state.worker.visualization_of(&id);
//     match visualization {
//         Some(image) => {
//             let mut buf = vec![];
//             let writer = BufWriter::new(&mut buf);
//             let encoder = PngEncoder::new(writer);
//             image.write_with_encoder(encoder).unwrap();
//             (StatusCode::OK, [(header::CONTENT_TYPE, "image/png")], buf)
//         }
//         None => (
//             StatusCode::NOT_FOUND,
//             [(header::CONTENT_TYPE, "text/plain")],
//             "visualization not found".as_bytes().to_vec(),
//         ),
//     }
// }

#[derive(Debug, Deserialize)]
struct ClassifyRequest {
    pub url: Url,
    pub spots: Vec<Spot>,
}

async fn classify(
    State(app_state): State<ServerState>,
    Json(request): Json<ClassifyRequest>,
) -> impl IntoResponse {
    let image = vis::capture(request.url).await.unwrap();
    // let spot_images = 

    (StatusCode::OK, String::new())
}
pub async fn run(app_state: ServerState) -> anyhow::Result<()> {
    let app = axum::Router::new()
        .route("/classify", get(classify))
        // .route("/cameras/state", get(get_all_state))
        // .route("/cameras/:id/state", get(get_state))
        // .route("/cameras/:id/metadata", put(put_metadata))
        // .route("/cameras/:id/visualize", get(visualize))
        .with_state(app_state);
    let url = "0.0.0.0:3000".parse().unwrap();
    tracing::info!("starting server at {}", url);
    axum::Server::bind(&url)
        .serve(app.into_make_service())
        .await
        .map_err(|err| err.into())
}
