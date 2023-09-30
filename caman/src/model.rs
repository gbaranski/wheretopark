use image::imageops::resize;
use image::imageops::FilterType;
use image::GrayImage;
use image::Pixel;
use image::RgbImage;
use imageproc::contours::find_contours_with_threshold;
use itertools::izip;
use itertools::Itertools;
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
use crate::Vehicle;

#[derive(Debug)]
pub struct Model {
    session: Session,
}

pub const HEIGHT: usize = 1024;
pub const WIDTH: usize = 1024;

fn generate_input(image: &RgbImage) -> anyhow::Result<ArrayBase<CowRepr<'_, u8>, Dim<IxDynImpl>>> {
    assert_eq!(image.width() as usize, WIDTH);
    assert_eq!(image.height() as usize, HEIGHT);
    let image =
        ndarray::Array::from_shape_fn([HEIGHT, WIDTH, 3], |(y, x, channel)| {
            let pixel = image.get_pixel(x as u32, y as u32);
            let channels = pixel.channels();
            channels[channel]
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
    let array = view.view();
    Ok(array.into_owned())
}

const VEHICLE_LABELS: [u8; 3] = [3, 6, 8];

impl Model {
    pub fn new(model_path: impl AsRef<Path>) -> anyhow::Result<Self> {
        let environment = Environment::builder()
            .with_name("MaskRCNN")
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

    pub fn infere(&self, image: &RgbImage) -> anyhow::Result<Vec<Vehicle>> {
        let input = generate_input(image)?;
        let outputs = self
            .session
            .run(vec![Value::from_array(self.session.allocator(), &input)?])?;

        let outputs: HashMap<&str, ArrayBase<OwnedRepr<f32>, Dim<IxDynImpl>>> = self
            .session
            .outputs
            .iter()
            .enumerate()
            .map(|(i, o)| {
                let value = try_extract::<f32>(&outputs[i])?;
                anyhow::Ok((o.name.as_str(), value))
            })
            .collect::<anyhow::Result<_>>()?;

        let boxes = &outputs["detection_boxes"];
        let classes = &outputs["detection_classes"];
        let masks = &outputs["detection_masks"];
        let scores = &outputs["detection_scores"];

        let objects = izip!(
            boxes.axis_iter(Axis(1)),
            classes.axis_iter(Axis(1)),
            scores.axis_iter(Axis(1)),
            masks.axis_iter(Axis(1))
        )
        .map(|(bbox, class, score, mask)| {
            let bbox = bbox.as_slice().unwrap();
            let bbox = BoundingBox {
                min: Point {
                    x: bbox[1] * image.width() as f32,
                    y: bbox[0] * image.height() as f32,
                },
                max: Point {
                    x: bbox[3] * image.width() as f32,
                    y: bbox[2] * image.height() as f32,
                },
            };
            let class = class[0] as u8;
            let score = score[0];
            (bbox, class, score, mask)
        })
        .filter(|(_, class, _, _)| VEHICLE_LABELS.contains(class))
        .filter(|(_, _, score, _)| *score > 0.5)
        .map(|(bbox, _, score, mask)| {
            assert_eq!(mask.shape(), [1, 33, 33]);
            let mut mask_image = GrayImage::new(33, 33);

            for (x, y, pixel) in mask_image.enumerate_pixels_mut() {
                let value = mask[[0, x as usize, y as usize]];
                *pixel = image::Luma([(value * 255.0) as u8]);
            }
            // mask_image.save("mask.png").unwrap();
            let mask_image = resize(
                &mask_image,
                bbox.width() as u32,
                bbox.height() as u32,
                FilterType::Lanczos3,
            );
            // mask_image.save("mask-resized.png").unwrap();
            let contours = find_contours_with_threshold(&mask_image, 128u8);

            Vehicle {
                bbox,
                score,
                contours,
            }
        });

        Ok(objects.collect())
    }
}
