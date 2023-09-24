mod model;
mod utils;
mod stream;

use image::imageops;
use image::imageops::FilterType;
use model::Model;
use stream::capture;

use crate::utils::visualise;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let model = Model::new()?;
    let url = reqwest::Url::parse("https://cam4out.klemit.net/hls/camn826.m3u8")?;
    let image = capture(url).await?;
    let image = image.into_rgb8();
    let objects = model.infere(&image)?;
    let image = imageops::resize(&image, 1280, 32 * 22, FilterType::Lanczos3);
    let image = visualise(image, &objects);
    image.save("output.jpeg")?;

    Ok(())
}
