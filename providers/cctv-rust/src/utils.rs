use chrono::{DateTime, Utc};
use geo::{Coordinate, LineString, Point, Polygon};
use image::DynamicImage;
use imageproc::drawing::{self, Blend};
use mat2image::ToImage;
use miette::{IntoDiagnostic, Result};
use opencv::core::RotatedRect;
use opencv::{core::Mat, prelude::RotatedRectTraitConst};
use rusttype::Font;
use std::str::FromStr;

use crate::configuration::ParkingSpot;

impl ParkingSpot {
    pub fn points(&self) -> Mat {
        let points = self
            .points
            .clone()
            .into_iter()
            .map(|(x, y)| [x as i32, y as i32])
            .collect::<Vec<_>>();
        Mat::from_slice_2d(&points).unwrap()
    }

    pub fn min_area_rectangle(&self) -> RotatedRect {
        let points = self.points();
        opencv::imgproc::min_area_rect(&points).unwrap()
    }

    pub fn coordinates(&self) -> Vec<Coordinate<f32>> {
        self.points
            .iter()
            .map(|p| Coordinate {
                x: p.0 as f32,
                y: p.1 as f32,
            })
            .collect::<Vec<_>>()
    }

    pub fn line_string(&self) -> LineString<f32> {
        let coordinates = self.coordinates();
        LineString::new(coordinates)
    }

    pub fn polygon(&self) -> Polygon<f32> {
        let line_string = self.line_string();
        let polygon = Polygon::new(line_string, vec![]);
        polygon
    }

    pub fn center(&self) -> Point<f32> {
        let center = self.min_area_rectangle().center();
        Point::new(center.x, center.y)
    }

    pub fn crop(&self, image: Mat) -> Result<DynamicImage> {
        let min_area_rectangle = self.min_area_rectangle();
        let size = min_area_rectangle.size();
        let mut box_points = Mat::default();
        opencv::imgproc::box_points(min_area_rectangle, &mut box_points).into_diagnostic()?;
        let destination_points = Mat::from_slice_2d(&[
            [0., size.height - 1.],
            [0., 0.],
            [size.width - 1., 0.],
            [size.width - 1., size.height - 1.],
        ])
        .into_diagnostic()?;
        let transformation = opencv::imgproc::get_perspective_transform(
            &box_points,
            &destination_points,
            opencv::core::DECOMP_LU,
        )
        .into_diagnostic()?;

        let mut output = Mat::default();
        let size = opencv::core::Size::new(size.width as i32, size.height as i32);

        opencv::imgproc::warp_perspective(
            &image,
            &mut output,
            &transformation,
            size,
            opencv::imgproc::INTER_LINEAR,
            opencv::core::BORDER_CONSTANT,
            opencv::core::Scalar::default(),
        )
        .into_diagnostic()?;
        let image = output.to_image().into_diagnostic()?;
        Ok(image)
    }
}

pub fn save_spot_image(
    image: &DynamicImage,
    date: &DateTime<Utc>,
    parking_lot_name: &str,
    index: usize,
) {
    let path = format!(
        "/tmp/wheretopark/{}/{}/{index}.png",
        parking_lot_name.replace(' ', "_"),
        date.to_rfc3339(),
    );
    let path = std::path::PathBuf::from_str(&path).unwrap();
    std::fs::create_dir_all(path.parent().unwrap()).unwrap();
    image.save(&path).unwrap();
    tracing::info!("saved spot {index} at {}", path.display());
}

const FONT: &[u8] = include_bytes!("../assets/fonts/november.ttf");

pub fn visualise<'a>(
    image: DynamicImage,
    predictions: impl Iterator<Item = &'a (ParkingSpot, f32)>,
) -> Result<DynamicImage> {
    let mut canvas = Blend(image.to_rgb8());
    let font: Font<'static> = Font::try_from_bytes(FONT).unwrap();
    for (spot, prediction) in predictions {
        let occupied = *prediction < 0.5;
        let color = if occupied { [255, 0, 0] } else { [0, 255, 0] };
        let color = image::Rgb(color);

        let polygon = spot.polygon();
        for line in polygon.exterior().lines() {
            drawing::draw_line_segment_mut(
                &mut canvas,
                (line.start.x as _, line.start.y as _),
                (line.end.x as _, line.end.y as _),
                color,
            );
        }
        let (center_x, center_y) = spot.center().x_y();
        drawing::draw_text_mut(
            &mut canvas,
            color,
            center_x as i32,
            center_y as i32,
            rusttype::Scale::uniform(15.),
            &font,
            &format!("{:.2}", prediction),
        )
    }
    Ok(DynamicImage::ImageRgb8(canvas.0))
}

pub fn save_visualisation(image: &DynamicImage, date: &DateTime<Utc>, parking_lot_name: &str) {
    let path = format!(
        "/tmp/wheretopark/{}/{}.jpg",
        parking_lot_name.replace(' ', "_"),
        date.to_rfc3339()
    );
    let path = std::path::PathBuf::from_str(&path).unwrap();
    std::fs::create_dir_all(path.parent().unwrap()).unwrap();
    image.save(&path).unwrap();
    tracing::info!("saved visualisation at {}", path.display());
}
