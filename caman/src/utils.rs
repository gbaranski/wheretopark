use image::buffer::ConvertBuffer;
use image::RgbImage;
use image::Rgba;
use image::RgbaImage;
use imageproc::drawing::draw_hollow_rect_mut;
use imageproc::drawing::draw_polygon_mut;
use imageproc::drawing::Blend;
use imageproc::point::Point;
use imageproc::rect::Rect;

use crate::model::Vehicle;

const GREEN: Rgba<u8> = Rgba([0, 255, 0, 255]);
const TRANSLUCENT_GREEN: Rgba<u8> = Rgba([0, 128, 0, 128]);

pub fn visualise(image: RgbImage, vehicles: &[Vehicle]) -> RgbImage {
    let image: RgbaImage = image.convert();
    let mut image = Blend(image);
    vehicles.iter().for_each(|object| {
        let rect = Rect::at(object.bbox.min.x as i32, object.bbox.min.y as i32)
            .of_size(object.bbox.width() as u32, object.bbox.height() as u32);
        draw_hollow_rect_mut(&mut image, rect, GREEN);
        object.contours.iter().for_each(|contour| {
            let points = contour
                .points
                .iter()
                .map(|point| Point {
                    x: object.bbox.min.x as i32 + point.x as i32,
                    y: object.bbox.min.y as i32 + point.y as i32,
                })
                .collect::<Vec<_>>();
            draw_polygon_mut(&mut image, &points, TRANSLUCENT_GREEN);
        });
    });
    image.0.convert()
}
