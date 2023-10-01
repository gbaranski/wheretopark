use image::Pixel;
use image::RgbImage;
use itertools::Itertools;
use ndarray::s;
use ndarray::ArrayBase;
use ndarray::CowArray;
use ndarray::CowRepr;
use ndarray::Dim;
use ndarray::IxDynImpl;
use ort::Environment;
use ort::ExecutionProvider;
use ort::GraphOptimizationLevel;
use ort::LoggingLevel;
use ort::Session;
use ort::SessionBuilder;
use ort::Value;
use std::path::Path;

use crate::BoundingBox;
use crate::Object;
use crate::Point;

#[derive(Debug)]
pub struct Model {
    session: Session,
}

const COCO_CLASSES_TXT: &'static str = include_str!("coco.txt");

lazy_static::lazy_static! {
    static ref COCO_CLASSES: Vec<&'static str> = COCO_CLASSES_TXT.lines().collect();
}

pub const HEIGHT: usize = 640;
pub const WIDTH: usize = 640;
const CHANNELS: usize = 3;

fn generate_input(
    images: &[RgbImage],
) -> anyhow::Result<ArrayBase<CowRepr<'_, f32>, Dim<IxDynImpl>>> {
    for image in images {
        assert_eq!(image.width() as usize, WIDTH);
        assert_eq!(image.height() as usize, HEIGHT);
    }
    let image = ndarray::Array::from_shape_fn(
        [images.len(), CHANNELS, HEIGHT, WIDTH],
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

impl Model {
    pub fn new(model_path: impl AsRef<Path>) -> anyhow::Result<Self> {
        let environment = Environment::builder()
            .with_name("YOLOv8")
            .with_log_level(LoggingLevel::Verbose)
            .with_execution_providers([ExecutionProvider::CPU(Default::default())])
            .build()?
            .into_arc();
        let session = SessionBuilder::new(&environment)?
            .with_optimization_level(GraphOptimizationLevel::Level1)?
            .with_intra_threads(1)?
            .with_model_from_file(model_path)?;

        let mut inputs = session.inputs.iter().map(|i| i.name.as_str());
        let mut outputs = session.outputs.iter().map(|o| o.name.as_str());
        tracing::debug!("inputs: {}", inputs.join(", "));
        tracing::debug!("outputs: {}", outputs.join(", "));

        Ok(Self { session })
    }

    pub fn infere(&self, images: &[RgbImage]) -> anyhow::Result<Vec<Vec<Object>>> {
        let input = generate_input(images)?;
        let outputs = self
            .session
            .run(vec![Value::from_array(self.session.allocator(), &input)?])?;

        let outputs = outputs.get(0).unwrap();
        let outputs = outputs.try_extract::<f32>()?;
        let outputs = outputs.view();
        let outputs = outputs.outer_iter().map(|output| {
            let output = output.t();
            output
                .outer_iter()
                .map(|row| {
                    let dimensions = row.slice(s![..4]).iter().cloned().collect_vec();
                    let classes = row.slice(s![4..]).iter().cloned().collect_vec();
                    assert_eq!(classes.len(), COCO_CLASSES.len());
                    (dimensions, classes)
                })
                .map(|(dim, classes)| {
                    let (class_id, score) = classes
                        .iter()
                        .enumerate()
                        .map(|(index, value)| (index, *value))
                        .reduce(|accum, row| if row.1 > accum.1 { row } else { accum })
                        .unwrap();
                    let label = COCO_CLASSES[class_id];
                    (dim, (label, score))
                })
                .filter(|(_, (_, score))| *score > 0.25)
                .map(|(dim, (label, score))| {
                    let center = Point {
                        x: dim[0],
                        y: dim[1],
                    };
                    let (width, height) = (dim[2], dim[3]);
                    let bbox = BoundingBox {
                        min: Point {
                            x: center.x - (width / 2.0),
                            y: center.y - (height / 2.0),
                        },
                        max: Point {
                            x: center.x + (width / 2.0),
                            y: center.y + (height / 2.0),
                        },
                    };
                    Object {
                        bbox,
                        label,
                        score,
                        contours: vec![],
                    }
                })
                .sorted_by(|object1, object2| object1.score.partial_cmp(&object2.score).unwrap())
                .rev()
                .fold(Vec::<Object>::new(), |mut acc, object| {
                    let overlapped = acc.iter().any(|other| object.bbox.iou(&other.bbox) > 0.7);
                    if !overlapped {
                        acc.push(object)
                    }
                    acc
                })
        });
        Ok(outputs.collect())
    }
}
