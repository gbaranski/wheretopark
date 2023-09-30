use image::imageops::resize;
use image::imageops::FilterType;
use image::GrayImage;
use image::Pixel;
use image::RgbImage;
use imageproc::contours::Contour;
use imageproc::contours::find_contours_with_threshold;
use itertools::izip;
use itertools::Itertools;
use ndarray::s;
use ndarray::ArrayBase;
use ndarray::Axis;
use ndarray::CowArray;
use ndarray::CowRepr;
use ndarray::Dim;
use ndarray::IxDynImpl;
use ndarray::OwnedRepr;
use ort::tensor::TensorDataToType;
use ort::Environment;
use ort::ExecutionProvider;
use ort::GraphOptimizationLevel;
use ort::LoggingLevel;
use ort::Session;
use ort::SessionBuilder;
use ort::Value;
use std::collections::HashMap;
use std::path::Path;

use crate::BoundingBox;
use crate::Point;
use crate::Object;

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

fn generate_input(image: &RgbImage) -> anyhow::Result<ArrayBase<CowRepr<'_, f32>, Dim<IxDynImpl>>> {
    image.save("/tmp/temp.png").unwrap();
    assert_eq!(image.width() as usize, WIDTH);
    assert_eq!(image.height() as usize, HEIGHT);
    let image = ndarray::Array::from_shape_fn([CHANNELS, HEIGHT, WIDTH], |(channel, y, x)| {
        let pixel = image.get_pixel(x as u32, y as u32);
        let channels = pixel.channels();
        channels[channel] as f32 / 255.0
    });
    let image = image.insert_axis(Axis(0));
    let input = CowArray::from(image.into_dyn());
    Ok(input)
}

fn try_extract<'a, T: TensorDataToType>(
    value: &'a Value,
) -> anyhow::Result<ArrayBase<OwnedRepr<T>, Dim<IxDynImpl>>> {
    let tensor = value.try_extract::<T>()?;
    let view = tensor.view();
    let transposed = view.t();
    Ok(transposed.into_owned())
}

const VEHICLE_LABELS: [u8; 3] = [3, 6, 8];

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

    pub fn infere(&self, image: &RgbImage) -> anyhow::Result<Vec<Object>> {
        let input = generate_input(image)?;
        let outputs = self
            .session
            .run(vec![Value::from_array(self.session.allocator(), &input)?])?;

        let output = try_extract::<f32>(&outputs[0])?;
        let output = output.slice(s![.., .., 0]);
        dbg!(output.shape());
        let objects = output
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
            .filter(|(_, (_, score))| *score > 0.1)
            .map(|(dim, (label, score))| {
                let center = Point {
                    x: dim[0],
                    y: dim[1],
                };
                dbg!(&center);
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
            });

            Ok(objects.collect_vec())
    }
}
