use image::RgbImage;
use imageproc::{
    drawing::draw_hollow_rect_mut,
    rect::Rect,
};

use crate::model::Object;

pub fn visualise(mut image: RgbImage, objects: &[Object]) -> RgbImage {
    objects
        .iter()
        .filter(|o| o.label == 3)
        .filter(|o| o.score > 0.7)
        .for_each(|object| {
            let rect = Rect::at(object.bbox.min.x as i32, object.bbox.min.y as i32).of_size(
                (object.bbox.max.x - object.bbox.min.x) as u32,
                (object.bbox.max.y - object.bbox.min.y) as u32,
            );
            let color = image::Rgb([0, 255, 0]);
            draw_hollow_rect_mut(&mut image, rect, color);
        });
    image
}
