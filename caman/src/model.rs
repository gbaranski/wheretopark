use image::imageops::resize;
use image::imageops::FilterType;
use image::GrayImage;
use image::Pixel;
use image::RgbImage;
use imageproc::contours::find_contours_with_threshold;
use itertools::izip;
use ndarray::ArrayBase;
use ndarray::CowArray;
use ndarray::Dim;
use ndarray::IxDynImpl;
use ndarray::OwnedRepr;
use ort::download::vision::ObjectDetectionImageSegmentation;
use ort::tensor::TensorDataToType;
use ort::Environment;
use ort::ExecutionProvider;
use ort::GraphOptimizationLevel;
use ort::LoggingLevel;
use ort::Session;
use ort::SessionBuilder;
use ort::Value;

use crate::BoundingBox;
use crate::Point;
use crate::Vehicle;

#[derive(Debug)]
pub struct Model {
    session: Session,
}


fn preprocess(image: &RgbImage) -> ArrayBase<OwnedRepr<f32>, Dim<[usize; 3]>> {
    let mean = [102.9801, 115.9465, 122.7717];
    let image =
        ndarray::Array::from_shape_fn([3, image.height() as usize, image.width() as usize], |(channel, y, x)| {
            let pixel = image.get_pixel(x as u32, y as u32);
            let channels = pixel.channels();
            channels[channel] as f32 - mean[channel]
        });

    image
}

fn try_extract<'a, T: TensorDataToType>(
    value: &'a Value,
) -> anyhow::Result<ArrayBase<OwnedRepr<T>, Dim<IxDynImpl>>> {
    let tensor = value.try_extract::<T>()?;
    let view = tensor.view();
    let array = view.view();
    Ok(array.into_owned())
}

const VEHICLE_LABELS: [i64; 3] = [3, 6, 8];

impl Model {
    pub fn new() -> anyhow::Result<Self> {
        let environment = Environment::builder()
            .with_name("MaskRCNN")
            .with_log_level(LoggingLevel::Verbose)
            .with_execution_providers([ExecutionProvider::CPU(Default::default())])
            .build()?
            .into_arc();
        let session = SessionBuilder::new(&environment)?
            .with_optimization_level(GraphOptimizationLevel::Level1)?
            .with_intra_threads(1)?
            .with_model_downloaded(ObjectDetectionImageSegmentation::MaskRcnn)?;
        Ok(Self { session })
    }

    pub fn infere(&self, image: &RgbImage) -> anyhow::Result<Vec<Vehicle>> {
        let input = preprocess(image);
        let outputs = self.session.run(vec![Value::from_array(
            self.session.allocator(),
            &CowArray::from(input.into_dyn()),
        )?])?;
        let boxes = try_extract::<f32>(&outputs[0])?;
        let labels = try_extract::<i64>(&outputs[1])?;
        let scores = try_extract::<f32>(&outputs[2])?;
        let masks = try_extract::<f32>(&outputs[3])?;

        let objects = izip!(
            boxes.outer_iter(),
            labels.outer_iter(),
            scores.outer_iter(),
            masks.outer_iter()
        )
        .map(|(bbox, label, score, mask)| {
            let bbox = BoundingBox {
                min: Point {
                    x: bbox[0],
                    y: bbox[1],
                },
                max: Point {
                    x: bbox[2],
                    y: bbox[3],
                },
            };
            let label = label[[]];
            let score = score[[]];
            (bbox, label, score, mask)
        })
        .filter(|(_, label, _, _)| VEHICLE_LABELS.contains(label))
        .filter(|(_, _, score, _)| *score > 0.7)
        .map(|(bbox, label, score, mask)| {
            assert_eq!(mask.shape(), [1, 28, 28]);
            let mut mask_image = GrayImage::new(28, 28);

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
                label,
                score,
                contours,
            }
        });

        Ok(objects.collect())
    }
}
