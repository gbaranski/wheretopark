mod classification;
mod detection;

use classification::ClassificationModel;
use detection::DetectionModel;
use directories::ProjectDirs;
use image::imageops;
use image::Pixel;
use image::RgbImage;
use itertools::Itertools;
use ndarray::ArrayBase;
use ndarray::CowArray;
use ndarray::CowRepr;
use ndarray::Dim;
use ndarray::IxDynImpl;
use ort::Environment;
use ort::ExecutionProvider;
use ort::LoggingLevel;

use crate::Object;

#[derive(Debug)]
pub struct Model {
    classification: ClassificationModel,
    detection: DetectionModel,
}

impl Model {
    pub fn new(project_directories: ProjectDirs) -> anyhow::Result<Self> {
        let environment = Environment::builder()
            .with_name("caman")
            .with_log_level(LoggingLevel::Verbose)
            .with_execution_providers([ExecutionProvider::CPU(Default::default())])
            .build()?
            .into_arc();

        let classification = ClassificationModel::new(
            &environment,
            project_directories
                .data_dir()
                .join("Inception-ResNet-v2.onnx"),
        )?;
        let detection = DetectionModel::new(
            &environment,
            project_directories.data_dir().join("yolov8x.onnx"),
        )?;
        Ok(Self {
            classification,
            detection,
        })
    }

    pub fn detect(&self, images: &[RgbImage]) -> anyhow::Result<Vec<Vec<Object>>> {
        self.detection.infere(images)
    }

    pub fn classify(&self, images: &[RgbImage]) -> anyhow::Result<Vec<f32>> {
        self.classification.infere(images)
    }
}

fn generate_input(
    channels: usize,
    width: usize,
    height: usize,
    images: &[RgbImage],
) -> anyhow::Result<ArrayBase<CowRepr<'_, f32>, Dim<IxDynImpl>>> {
    let images = images
        .into_iter()
        .map(|img| {
            imageops::resize(
                img,
                width as u32,
                height as u32,
                imageops::FilterType::Lanczos3,
            )
        })
        .collect_vec();
    let image = ndarray::Array::from_shape_fn(
        [images.len(), channels, height, width],
        |(idx, channel, y, x)| {
            let image = &images[idx];
            let pixel = image.get_pixel(x as u32, y as u32);
            let channels = pixel.channels();
            channels[channel] as f32 / 255.0
        },
    );
    let input = CowArray::from(image.into_dyn());
    Ok(input)
}
