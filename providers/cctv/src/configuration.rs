use serde::Deserialize;
use serde::Serialize;
use std::collections::HashMap;
use url::Url;
use wheretopark::LanguageCode;

use wheretopark::parking_lot::{Currency, PaymentMethod};
use wheretopark::parking_lot::Feature;
use wheretopark::parking_lot::{Coordinate, Rule};

#[derive(Debug, Clone, Serialize, Deserialize, Hash, Eq, PartialEq)]
#[serde(rename_all = "kebab-case")]
pub struct ParkingSpot {
    pub points: Vec<(u32, u32)>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
pub struct ParkingLot {
    pub name: String,
    pub address: String,
    pub location: Coordinate,
    pub resources: Vec<Url>,
    #[serde(default)]
    pub max_width: Option<i32>,
    #[serde(default)]
    pub max_height: Option<i32>,
    #[serde(default)]
    pub payment_methods: Vec<PaymentMethod>,
    pub features: Vec<Feature>,
    pub rules: Vec<Rule>,
    #[serde(default)]
    pub comment: HashMap<LanguageCode, String>,
    pub currency: Currency,

    pub camera_url: Url,
    pub spots: Vec<ParkingSpot>,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
pub struct Configuration {
    pub parking_lots: Vec<ParkingLot>,
}
