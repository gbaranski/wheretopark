use crate::model::HEIGHT;
use crate::model::WIDTH;
use anyhow::Context;
use image::RgbImage;
use std::process::Stdio;
use tokio::io::AsyncBufRead;
use tokio::io::AsyncBufReadExt;
use tokio::io::AsyncRead;
use tokio::io::AsyncReadExt;
use tokio::io::BufReader;
use tokio::process::Child;
use tokio::process::ChildStderr;
use tokio::process::ChildStdout;
use tokio::process::Command;
use tokio::sync::watch;
use tracing::span;
use tracing::Instrument;
use tracing::Level;

fn command(url: impl AsRef<str>) -> Command {
    let mut command = Command::new("ffmpeg");
    command.arg("-hwaccel");
    command.arg("auto");

    // set input URL
    command.arg("-i");
    command.arg(url.as_ref());
    // set video filters
    command.arg("-vf");
    command.arg("fps=1/20,format=rgb24");
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
    command.kill_on_drop(true);
    command
}

async fn read_image(
    buf: &mut [u8],
    mut reader: impl AsyncRead + AsyncBufRead + Unpin + Send,
) -> anyhow::Result<RgbImage> {
    let n = reader.read_exact(buf).await?;
    tracing::debug!("read {n} bytes");
    if n == 0 {
        return Err(anyhow::anyhow!("no bytes read"));
    }
    let image =
        RgbImage::from_raw(WIDTH as u32, HEIGHT as u32, buf.to_vec()).context("create image")?;
    Ok(image)
}

async fn wait_for_exit(mut child: Child) {
    let status = child.wait().await.expect("ffmpeg exit");
    if !status.success() {
        tracing::error!("ffmpeg exited with status {}", status);
    }
}

async fn process_stderr(stderr: BufReader<ChildStderr>) {
    let mut lines = stderr.lines();
    while let Some(line) = lines.next_line().await.unwrap() {
        tracing::trace!("ffmpeg: {}", line)
    }
}

async fn process_stdout(
    mut stdout: BufReader<ChildStdout>,
    sender: watch::Sender<Option<RgbImage>>,
) {
    let mut buf = vec![0; WIDTH * HEIGHT * 3];
    loop {
        tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;
        let result = read_image(&mut buf, &mut stdout).await;
        match result {
            Ok(image) => {
                sender.send(Some(image)).unwrap();
            }
            Err(err) => {
                tracing::error!("read image failed: {:#}", err);
            }
        }
    }
}

pub fn images(url: impl AsRef<str>) -> anyhow::Result<watch::Receiver<Option<RgbImage>>> {
    let url = url.as_ref();
    tracing::debug!("connect to {url}");
    let span = span!(Level::INFO, "images()", url = %url);
    let mut command = command(&url);

    let mut child = command.spawn()?;
    let stderr = BufReader::new(child.stderr.take().unwrap());
    let stdout = BufReader::new(child.stdout.take().unwrap());
    let (tx, rx) = watch::channel(None);
    tokio::spawn(wait_for_exit(child).instrument(span.clone()));
    tokio::spawn(process_stderr(stderr).instrument(span.clone()));
    tokio::spawn(process_stdout(stdout, tx).instrument(span.clone()));

    Ok(rx)
}
