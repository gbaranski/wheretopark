use image::{
    imageops::{self, FilterType},
    Pixel, RgbImage,
};
use ndarray::{ArrayBase, CowArray, Dim, IxDynImpl, OwnedRepr};
use ort::{
    download::vision::ObjectDetectionImageSegmentation, tensor::TensorDataToType, Environment,
    ExecutionProvider, GraphOptimizationLevel, LoggingLevel, Session, SessionBuilder, Value,
};

#[derive(Debug)]
pub struct Point {
    pub x: f32,
    pub y: f32,
}

#[derive(Debug)]
pub struct BoundingBox {
    pub min: Point,
    pub max: Point,
}

#[derive(Debug)]
pub struct Object {
    pub bbox: BoundingBox,
    pub label: i64,
    pub score: f32,
}

pub struct Model {
    session: Session,
}

fn preprocess(image: &RgbImage) -> ArrayBase<OwnedRepr<f32>, Dim<[usize; 3]>> {
    // TODO: Resize images as written in https://github.com/onnx/models/tree/main/vision/object_detection_segmentation/mask-rcnn#preprocessing-steps
    let width = 1280;
    let height = 32 * 22;
    let image = imageops::resize(image, width, height, FilterType::Triangle);

    let mean = [102.9801, 115.9465, 122.7717];
    let image =
        ndarray::Array::from_shape_fn([3, height as usize, width as usize], |(channel, y, x)| {
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

    pub fn infere(&self, image: &RgbImage) -> anyhow::Result<Vec<Object>> {
        let input = preprocess(image);
        let outputs = self.session.run(vec![Value::from_array(
            self.session.allocator(),
            &CowArray::from(input.into_dyn()),
        )?])?;
        let boxes = try_extract::<f32>(&outputs[0])?;
        let labels = try_extract::<i64>(&outputs[1])?;
        let scores = try_extract::<f32>(&outputs[2])?;
        dbg!(&boxes);

        let objects = boxes
            .outer_iter()
            .enumerate()
            .map(|(i, bbox)| {
                let bbox = BoundingBox{
                    min: Point {
                        x: bbox[0],
                        y: bbox[1],
                    },
                    max: Point {
                        x: bbox[2],
                        y: bbox[3],
                    },
                };
                let label = labels[i];
                let score = scores[i];
                Object { bbox, label, score }
            })
            .collect::<Vec<_>>();

        Ok(objects)
    }
}
