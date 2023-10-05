use std::collections::HashMap;

use chrono::Utc;
use chrono::DateTime;
use serde::Serialize;
use serde::Deserialize;
use url::Url;

pub type CameraID = String;

pub type SpotType = String;

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Spot {
    pub points: Vec<[u32; 4]>,
    pub r#type: SpotType,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CameraMetadata {
    pub url: Url,
    pub spots: Vec<Spot>,
}

#[derive(Debug, Clone, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct CameraState {
    pub last_updated: DateTime<Utc>,
    pub total_spots: HashMap<SpotType, u32>,
    pub available_spots: HashMap<SpotType, u32>,
}
