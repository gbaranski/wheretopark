mod model;
mod stream;
mod utils;

pub use utils::BoundingBox;
pub use utils::Point;
pub use utils::Vehicle;
pub use utils::ParkingSpot;

use image::imageops;
use image::imageops::FilterType;
use model::Model;
use utils::visualize_spots;
use utils::visualize_vehicles;
use std::time::Duration;
use stream::capture;

struct App {
    model: Model,
}

impl App {
    pub async fn process(&self, url: reqwest::Url) -> anyhow::Result<()> {
        let mut parked_vehicle_boxes: Vec<BoundingBox> = Vec::new();
        loop {
            let image = capture(url.as_str()).await?;
            let image = image.into_rgb8();
            let vehicles = self.model.infere(&image)?;
            let boxes: Vec<BoundingBox> = vehicles
                .iter()
                .map(|vehicle| vehicle.bbox.clone())
                .collect();
            if parked_vehicle_boxes.is_empty() {
                parked_vehicle_boxes = boxes;
            } else {
                let overlaps = utils::compute_overlaps(&parked_vehicle_boxes, &boxes);
                let spots = overlaps
                    .iter()
                    .zip(&parked_vehicle_boxes)
                    .map(|(overlap, bbox)| {
                        let score = overlap.iter().cloned().reduce(f32::max).unwrap();
                        ParkingSpot{
                            bbox: bbox.clone(),
                            score,
                        }
                    })
                    .collect::<Vec<_>>();
                
                let image = imageops::resize(&image, 1280, 32 * 22, FilterType::Lanczos3);
                let image_vehicles = visualize_vehicles(&image, &vehicles);
                let image_occupancy = visualize_spots(&image, &spots);
                let available_spots = spots.iter().filter(|spot| spot.score < 0.15).count();
                tracing::info!("available spots: {}", available_spots);
                image_vehicles.save("vehicles.jpeg")?;
                image_occupancy.save("occupancy.jpeg")?;
            }

            tokio::time::sleep(Duration::from_secs(5)).await;
        }
    }
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let model = Model::new()?;
    let app = App { model };
    let url = reqwest::Url::parse("https://cam4out.klemit.net/hls/camn826.m3u8")?;
    app.process(url).await?;

    Ok(())
}
