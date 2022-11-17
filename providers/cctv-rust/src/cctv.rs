use miette::{IntoDiagnostic, Result};
use opencv::prelude::VideoCaptureTraitConst;
use opencv::videoio;
use opencv::{core::Mat, videoio::VideoCaptureTrait};
use url::Url;

pub fn capture(url: &Url) -> Result<Mat> {
    let mut video_capture =
        videoio::VideoCapture::from_file(url.as_str(), videoio::CAP_ANY).into_diagnostic()?;
    if !video_capture.is_opened().into_diagnostic()? {
        panic!("Unable to open video capture!");
    }

    let mut image = Mat::default();
    video_capture.read(&mut image).into_diagnostic()?;
    Ok(image)
}
