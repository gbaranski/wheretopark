mod model;
mod utils;

use image::imageops::{self, FilterType};
use model::Model;

use crate::utils::visualise;

fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let model = Model::new()?;
    let image = image::open("demo.png")?.into_rgb8();
    let objects = model.infere(&image)?;
    dbg!(&objects);
    let image = imageops::resize(&image, 1280, 32 * 22, FilterType::Triangle);
    let image = visualise(image, &objects);
    image.save("output.jpeg")?;

    Ok(())
}
