use image::io::Reader as ImageReader;
use image::DynamicImage;
use image::ImageFormat;
use std::io::Cursor;
use tokio::process::Command;

pub async fn capture(url: &str) -> anyhow::Result<DynamicImage> {
    let mut command = Command::new("ffmpeg");
    // set input URL
    command.arg("-i");
    command.arg(url.to_string());
    // set to accept only the first frame
    command.arg("-vframes");
    command.arg("1");
    // interpret input as image
    command.arg("-f");
    command.arg("image2");
    // encode to png
    command.arg("-c:v");
    command.arg("png");
    // set output to stdout
    command.arg("-");

    let output = command.output().await?;
    if !output.status.success() {
        return Err(anyhow::anyhow!(
            "ffmpeg exited with status {}",
            output.status
        ));
    }
    let cursor = Cursor::new(output.stdout);
    let mut reader = ImageReader::new(cursor);
    reader.set_format(ImageFormat::Png);
    let image = reader.decode()?;

    Ok(image)
}
