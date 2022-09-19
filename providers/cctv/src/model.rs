use image::imageops::FilterType;
use image::DynamicImage;
use miette::Result;
use std::path::Path;
use tract_onnx::prelude::*;

const WIDTH: u32 = 128;
const HEIGHT: u32 = 128;
const CHANNELS: u8 = 3;

#[derive(Debug, Clone)]
pub struct Model {
    model: SimplePlan<TypedFact, Box<dyn TypedOp>, Graph<TypedFact, Box<dyn TypedOp>>>,
}

impl Model {
    pub fn new(model_path: impl AsRef<Path>) -> Self {
        let model = tract_onnx::onnx()
            .model_for_path(model_path)
            .unwrap()
            .with_input_fact(
                0,
                f32::fact(&[1, WIDTH as _, HEIGHT as _, CHANNELS as _]).into(),
            )
            .unwrap()
            .into_optimized()
            .unwrap()
            .into_runnable()
            .unwrap();

        Self { model }
    }

    pub fn predict(&self, image: &DynamicImage) -> Result<f32> {
        let image = image.to_rgb8();
        let resized = image::imageops::resize(&image, WIDTH, HEIGHT, FilterType::Triangle);
        let image: Tensor = tract_ndarray::Array4::from_shape_fn(
            (1, WIDTH as _, HEIGHT as _, CHANNELS as _),
            |(_, x, y, c)| {
                let v = resized[(x as _, y as _)][c] as f32;
                v / 255.0
            },
        )
        .into();
        let result = self.model.run(tvec!(image)).unwrap();
        let prediction = result[0]
            .to_array_view::<f32>()
            .unwrap()
            .into_iter()
            .next()
            .unwrap();
        Ok(*prediction)
    }
}
