use std::collections::HashMap;

use std::sync::Arc;

use async_trait::async_trait;
use chrono::Utc;
use clap::Parser;
use mat2image::ToImage;
use miette::Context;
use miette::Report;
use miette::{IntoDiagnostic, Result};
use tokio::time::Instant;

use wheretopark::parking_lot::Metadata;
use wheretopark::parking_lot::SpotType;
use wheretopark::parking_lot::ID;

use wheretopark::parking_lot::State;

use wheretopark::Provider;

use crate::configuration::Configuration;
use crate::model::Model;

mod cctv;
mod configuration;
mod model;
mod utils;

#[derive(Debug, Parser)]
#[clap(author, version, about)]
struct Args {
    #[clap(long)]
    image: Option<std::path::PathBuf>,
    #[clap(long)]
    visualise: bool,
    #[clap(long)]
    save_spots: bool,
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();
    tracing::info!("Starting program");
    let args = Args::parse();
    let xdg_dirs = xdg::BaseDirectories::with_prefix("wheretopark").into_diagnostic()?;

    let configuration = include_str!("configuration.yaml");
    let configuration = serde_yaml::from_str::<Configuration>(configuration).into_diagnostic()?;
    let model_path = xdg_dirs.get_data_file("model.onnx");
    let model = Model::new(model_path);
    let provider = CCTVProvider {
        configuration: Arc::new(configuration),
        args: Arc::new(args),
        model: Arc::new(model),
    };
    wheretopark::provider::run(provider).await.unwrap();

    Ok(())
}

#[derive(Debug, Clone)]
pub struct CCTVProvider {
    configuration: Arc<Configuration>,
    args: Arc<Args>,
    model: Arc<Model>,
}

#[async_trait]
impl Provider for CCTVProvider {
    async fn poll_metadatas(&self) -> Result<HashMap<ID, Metadata>> {
        self.configuration
            .parking_lots
            .iter()
            .map(|p| {
                let p = p.clone();
                let metadata = Metadata {
                    name: p.name,
                    address: p.address,
                    location: p.location,
                    resources: p.resources,
                    total_spots: HashMap::from([(SpotType::Car, p.spots.len() as _)]),
                    max_width: p.max_width,
                    max_height: p.max_height,
                    features: p.features,
                    payment_methods: p.payment_methods,
                    comment: p.comment,
                    currency: p.currency,
                    rules: p.rules,
                };
                let id = metadata.location.id();
                Ok((id, metadata))
            })
            .collect()
    }

    async fn poll_states(&self) -> Result<HashMap<ID, State>> {
        let futures = self
            .configuration
            .parking_lots
            .iter()
            .map(|parking_lot| async move {
                let parking_lot = parking_lot.clone();
                let model = self.model.clone();
                let args = self.args.clone();
                let (date, predictions) = tokio::task::spawn_blocking(move || {
                    let image =
                        cctv::capture(&parking_lot.camera_url).wrap_err("capture cctv frame")?;
                    let date = Utc::now();
                    let predictions = parking_lot
                        .spots
                        .into_iter()
                        .enumerate()
                        .map(|(i, spot)| {
                            let start = Instant::now();
                            let image = spot.crop(image.clone())?;
                            if args.save_spots {
                                utils::save_spot_image(&image, &date, &parking_lot.name, i);
                            }
                            let prediction = model.predict(&image)?;
                            let end = Instant::now();
                            tracing::debug!(
                                "{}/{}: {}. took {}ms",
                                parking_lot.name,
                                i,
                                prediction,
                                end.duration_since(start).as_millis()
                            );
                            Ok::<_, Report>((spot, prediction))
                        })
                        .collect::<Result<Vec<_>>>()?;

                    if args.visualise {
                        let image = image.to_image().into_diagnostic()?;
                        let visualization = utils::visualise(image, predictions.iter())?;
                        utils::save_visualisation(&visualization, &date, &parking_lot.name);
                    }
                    Ok::<_, Report>((date, predictions))
                })
                .await
                .into_diagnostic()??;
                let id = parking_lot.location.id();
                let available_spots = predictions
                    .into_iter()
                    .map(|(_, prediction)| prediction > 0.5)
                    .filter(|p| *p)
                    .count();
                let state = State {
                    last_updated: date,
                    available_spots: HashMap::from([
                        (SpotType::Car, available_spots as _),
                    ]),
                };
                Ok::<_, Report>((id, state))
            });
        let states = futures::future::try_join_all(futures).await.unwrap();
        let states = HashMap::from_iter(states.into_iter());
        Ok(states)
    }
}
