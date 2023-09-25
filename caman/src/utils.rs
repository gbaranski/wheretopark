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
    pub fn intersection_area(&self, other: &BoundingBox) -> f32 {
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
        let intersection = self.intersection_area(other);
        let union = self.area() + other.area() - intersection;
        intersection / union
    }
}

pub fn compute_iou(
    box1: &BoundingBox,
    boxes2: &Vec<BoundingBox>,
    area1: f32,
    areas2: &Vec<f32>,
) -> Vec<f32> {
    let mut overlaps: Vec<f32> = vec![0.0; boxes2.len()];

    for (i, box2) in boxes2.iter().enumerate() {
        let x1 = f32::max(box1.min.x, box2.min.x);
        let y1 = f32::max(box1.min.y, box2.min.y);
        let x2 = f32::min(box1.max.x, box2.max.x);
        let y2 = f32::min(box1.max.y, box2.max.y);

        let intersection = f32::max(0.0, x2 - x1) * f32::max(0.0, y2 - y1);
        let union = area1 + areas2[i] - intersection;

        overlaps[i] = intersection / union;
    }

    overlaps
}

pub fn compute_overlaps(boxes1: &Vec<BoundingBox>, boxes2: &Vec<BoundingBox>) -> Vec<Vec<f32>> {
    let mut overlaps: Vec<Vec<f32>> = vec![vec![0.0; boxes2.len()]; boxes1.len()];

    let areas1: Vec<f32> = boxes1.iter().map(|box1| box1.area()).collect();
    let areas2: Vec<f32> = boxes2.iter().map(|box2| box2.area()).collect();

    for (i, box1) in boxes1.iter().enumerate() {
        let overlaps_row = compute_iou(box1, boxes2, areas1[i], &areas2);
        overlaps[i] = overlaps_row;
    }

    overlaps
}

#[derive(Debug)]
pub struct Vehicle {
    pub bbox: BoundingBox,
    pub label: i64,
    pub score: f32,
    pub contours: Vec<Contour<u32>>,
}

#[derive(Debug)]
pub struct ParkingSpot {
    pub bbox: BoundingBox,
    pub score: f32,
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

pub fn visualize_vehicles(image: &RgbImage, vehicles: &[Vehicle]) -> RgbImage {
    let image: RgbaImage = image.convert();
    let mut image = Blend(image);
    vehicles.iter().for_each(|object| {
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

pub fn visualize_spots(image: &RgbImage, spots: &[ParkingSpot]) -> RgbImage {
    let image: RgbaImage = image.convert();
    let mut image = Blend(image);
    spots.iter().for_each(|spot| {
        let rect = Rect::at(spot.bbox.min.x as i32, spot.bbox.min.y as i32)
            .of_size(spot.bbox.width() as u32, spot.bbox.height() as u32);

        let color = if spot.score < 0.15 { GREEN } else { RED };
        draw_hollow_rect_mut(&mut image, rect, color);
    });
    image.0.convert()
}
