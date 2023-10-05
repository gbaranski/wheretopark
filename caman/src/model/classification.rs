use image::RgbImage;
use itertools::Itertools;
use ort::Environment;
use ort::GraphOptimizationLevel;
use ort::Session;
use ort::SessionBuilder;
use ort::Value;
use std::path::Path;
use std::sync::Arc;

#[derive(Debug)]
pub struct ClassificationModel {
    session: Session,
}

const HEIGHT: usize = 128;
const WIDTH: usize = 128;
const CHANNELS: usize = 3;

impl ClassificationModel {
    pub fn new(
        environment: &Arc<Environment>,
        model_path: impl AsRef<Path>,
    ) -> anyhow::Result<Self> {
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

    pub fn infere(&self, images: &[RgbImage]) -> anyhow::Result<Vec<f32>> {
        let input = super::generate_input(CHANNELS, WIDTH, HEIGHT, images)?;
        let mut outputs = self
            .session
            .run(vec![Value::from_array(self.session.allocator(), &input)?])?;
        let outputs = outputs.pop().unwrap();
        let outputs = outputs.try_extract::<f32>()?;
        let outputs = outputs.view();
        let outputs = outputs.iter().map(|o| *o).collect_vec();
        Ok(outputs)
    }
}
