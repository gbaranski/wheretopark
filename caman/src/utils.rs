use std::sync::Arc;

use image::buffer::ConvertBuffer;
use image::Pixel;
use image::RgbImage;
use image::Rgba;
use image::RgbaImage;
use imageproc::contours::Contour;
use imageproc::drawing::draw_hollow_rect_mut;
use imageproc::drawing::draw_polygon_mut;
use imageproc::drawing::Blend;
use imageproc::drawing::Canvas;
use imageproc::rect::Rect;
use num::traits::AsPrimitive;

#[derive(Debug, Clone)]
pub struct Point {
    pub x: f32,
    pub y: f32,
}

#[derive(Debug, Clone)]
pub struct BoundingBox {
    pub min: Point,
    pub max: Point,
}

impl BoundingBox {
    pub fn width(&self) -> f32 {
        self.max.x - self.min.x
    }

    pub fn height(&self) -> f32 {
        self.max.y - self.min.y
    }

    // Calculate the area of the bounding box
    pub fn area(&self) -> f32 {
        (self.max.x - self.min.x) * (self.max.y - self.min.y)
    }

    // Calculate the intersection area between two bounding boxes
    pub fn intersection(&self, other: &BoundingBox) -> f32 {
        let x_overlap = f32::max(
            0.0,
            f32::min(self.max.x, other.max.x) - f32::max(self.min.x, other.min.x),
        );
        let y_overlap = f32::max(
            0.0,
            f32::min(self.max.y, other.max.y) - f32::max(self.min.y, other.min.y),
        );
        x_overlap * y_overlap
    }

    // Calculate the IoU (Intersection over Union) between two bounding boxes
    pub fn iou(&self, other: &BoundingBox) -> f32 {
        let intersection = self.intersection(other);
        let union = self.area() + other.area() - intersection;
        intersection / union
    }
}

pub fn compute_iou(
    position1: &SpotPosition,
    positions2: &[SpotPosition],
    area1: f32,
    areas2: &Vec<f32>,
) -> Vec<f32> {
    let mut overlaps: Vec<f32> = vec![0.0; positions2.len()];

    for (i, box2) in positions2.iter().enumerate() {
        let x1 = f32::max(position1.bbox.min.x, box2.bbox.min.x);
        let y1 = f32::max(position1.bbox.min.y, box2.bbox.min.y);
        let x2 = f32::min(position1.bbox.max.x, box2.bbox.max.x);
        let y2 = f32::min(position1.bbox.max.y, box2.bbox.max.y);

        let intersection = f32::max(0.0, x2 - x1) * f32::max(0.0, y2 - y1);
        let union = area1 + areas2[i] - intersection;

        overlaps[i] = intersection / union;
    }

    overlaps
}

pub fn compute_overlaps(pos1: &[SpotPosition], pos2: &[SpotPosition]) -> Vec<Vec<f32>> {
    let mut overlaps: Vec<Vec<f32>> = vec![vec![0.0; pos2.len()]; pos1.len()];

    let areas1: Vec<f32> = pos1.iter().map(|box1| box1.bbox.area()).collect();
    let areas2: Vec<f32> = pos2.iter().map(|box2| box2.bbox.area()).collect();

    for (i, box1) in pos1.iter().enumerate() {
        let overlaps_row = compute_iou(box1, pos2, areas1[i], &areas2);
        overlaps[i] = overlaps_row;
    }

    overlaps
}

#[derive(Debug)]
pub struct Object {
    pub bbox: BoundingBox,
    pub label: &'static str,
    pub score: f32,
    pub contours: Vec<Contour<u32>>,
}

#[derive(Debug, Clone)]
pub struct SpotPosition {
    pub bbox: BoundingBox,
    pub contours: Arc<Vec<Contour<u32>>>,
}

#[derive(Debug, PartialEq, Eq)]
pub enum SpotState {
    Vacant,
    Occupied,
}

#[derive(Debug)]
pub struct Spot {
    pub position: SpotPosition,
    pub state: SpotState,
}

fn draw_contours<T: AsPrimitive<i32> + std::ops::Add<Output = T>, P: Pixel>(
    canvas: &mut impl Canvas<Pixel = P>,
    x: T,
    y: T,
    contours: &[Contour<T>],
    color: P,
) {
    contours.iter().for_each(|contour| {
        let mut points = contour
            .points
            .iter()
            .map(|point| imageproc::point::Point {
                x: (x + point.x).as_(),
                y: (y + point.y).as_(),
            })
            .collect::<Vec<_>>();
        if points.first() == points.last() {
            points.pop();
        }
        draw_polygon_mut(canvas, &points, color);
    });
}

const GREEN: Rgba<u8> = Rgba([0, 255, 0, 255]);
const RED: Rgba<u8> = Rgba([255, 0, 0, 255]);
const TRANSLUCENT_GREEN: Rgba<u8> = Rgba([0, 128, 0, 128]);

pub fn visualize_objects(image: &RgbImage, objects: &[Object]) -> RgbImage {
    let image: RgbaImage = image.convert();
    let mut image = Blend(image);
    objects.iter().for_each(|object| {
        let rect = Rect::at(object.bbox.min.x as i32, object.bbox.min.y as i32)
            .of_size(object.bbox.width() as u32, object.bbox.height() as u32);
        draw_hollow_rect_mut(&mut image, rect, GREEN);
        draw_contours(
            &mut image,
            object.bbox.min.x as u32,
            object.bbox.min.y as u32,
            &object.contours,
            TRANSLUCENT_GREEN,
        );
    });
    image.0.convert()
}

pub fn visualize_spots(image: &RgbImage, spots: &[Spot]) -> RgbImage {
    let image: RgbaImage = image.convert();
    let mut image = Blend(image);
    spots.iter().for_each(|spot| {
        let rect = Rect::at(
            spot.position.bbox.min.x as i32,
            spot.position.bbox.min.y as i32,
        )
        .of_size(
            spot.position.bbox.width() as u32,
            spot.position.bbox.height() as u32,
        );

        let color = if spot.state == SpotState::Vacant { GREEN } else { RED };
        draw_hollow_rect_mut(&mut image, rect, color);
    });
    image.0.convert()
}


pub fn crop_by_polygon(image: &RgbImage, points: Vec<u32>) -> RgbImage {
    

}