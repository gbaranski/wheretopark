use crate::model::HEIGHT;
use crate::model::WIDTH;
use anyhow::Context;
use image::DynamicImage;
use image::RgbImage;
use tokio::sync::watch;
use std::process::Stdio;
use tokio::io::AsyncBufRead;
use tokio::io::AsyncBufReadExt;
use tokio::io::AsyncRead;
use tokio::io::AsyncReadExt;
use tokio::io::BufReader;
use tokio::process::Command;
use tokio::sync::mpsc;

// #[derive(Debug)]
// pub struct ImageStream {
//     child: Child,
//     buf: Vec<u8>,
//     stdout: BufReader<ChildStdout>,
// }

// impl ImageStream {
//     fn new(child: Child) -> Self {
//         Self {
//             child,
//             buf: vec![],
//             stdout: BufReader::new(child.stdout.unwrap()),
//         }
//     }
// }

fn command(url: impl AsRef<str>) -> Command {
    let mut command = Command::new("ffmpeg");
    command.arg("-hwaccel");
    command.arg("auto");

    // set input URL
    command.arg("-i");
    command.arg(url.as_ref());
    // set video filters
    command.arg("-vf");
    command.arg("fps=1/30,format=rgb24");
    // set output image size
    command.arg("-s");
    command.arg(format!("{WIDTH}:{HEIGHT}"));
    // set output format
    command.arg("-f");
    command.arg("image2pipe");
    // set ouput codec
    command.arg("-vcodec");
    command.arg("rawvideo");
    // set output to stdout
    command.arg("-");

    command.stdout(Stdio::piped());
    command.stderr(Stdio::piped());
    command
}

async fn read(
    buf: &mut [u8],
    mut reader: impl AsyncRead + AsyncBufRead + Unpin + Send,
) -> anyhow::Result<DynamicImage> {
    let n = reader.read_exact(buf).await?;
    tracing::debug!("read {n} bytes");
    if n == 0 {
        return Err(anyhow::anyhow!("no bytes read"));
    }
    let image = RgbImage::from_raw(WIDTH as u32, HEIGHT as u32, buf.to_vec()).context("create image")?;
    Ok(image.into())
}

pub fn images(
    url: String,
) -> anyhow::Result<watch::Receiver<Option<DynamicImage>>> {
    tracing::info!("connecting to {url}");
    let mut command = command(&url);

    let child = command.spawn()?;
    tokio::spawn(async move {
        let stderr = child.stderr.unwrap();
        let reader = BufReader::new(stderr);
        let mut lines = reader.lines();
        while let Some(line) = lines.next_line().await.unwrap() {
            tracing::trace!("ffmpeg: {}", line)
        }
        tracing::info!("ffmpeg exited");
    });

    let (tx, rx) = watch::channel(None);
    let stdout = child.stdout.unwrap();
    let mut reader = BufReader::new(stdout);
    tokio::spawn(async move {
        let mut buf = vec![0; WIDTH * HEIGHT * 3];
        loop {
            tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;
            let result = read(&mut buf, &mut reader).await;
            match result {
                Ok(image) => {
                    tx.send(Some(image)).unwrap();
                }
                Err(err) => {
                    tracing::error!(url=%url, "read image failed: {:#}", err);
                }
            }
        }
    });

    Ok(rx)
}
