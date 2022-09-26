pub use iso_currency::Currency;
use std::collections::{HashMap, HashSet};

use crate::LanguageCode;
use chrono::DateTime;
use rust_decimal::Decimal;
use serde::Deserialize;
use serde::Serialize;
use url::Url;

pub type ID = String;

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
pub struct State {
    pub last_updated: DateTime<chrono::Utc>,
    pub available_spots: HashMap<SpotType, u32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
pub struct Coordinate {
    pub latitude: f64,
    pub longitude: f64,
}

impl Coordinate {
    pub fn id(&self) -> ID {
        let coordinate = geohash::Coordinate {
            x: self.latitude,
            y: self.longitude,
        };
        geohash::encode(coordinate, 12).unwrap()
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "SCREAMING-KEBAB-CASE")]
pub enum Feature {
    Uncovered,
    Covered,
    Underground,
}

pub type Weekdays = String;
pub type Hours = String;

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
pub struct PricingRule {
    pub duration: String,
    #[serde(with = "rust_decimal::serde::float")]
    pub price: Decimal,
    #[serde(default)]
    pub repeating: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
pub struct Rule {
    // https://schema.org/openingHours
    // https://wiki.openstreetmap.org/wiki/Key:opening_hours
    pub hours: String,
    #[serde(default)]
    /// If not empty, then include only those from this list
    pub includes: HashSet<SpotType>,
    pub pricing: Vec<PricingRule>,
}

#[derive(
    Debug, Clone, Serialize, Deserialize, Eq, PartialEq, Hash, strum::Display, strum::EnumString,
)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
#[strum(serialize_all = "SCREAMING_SNAKE_CASE")]
pub enum SpotType {
    Car,
    CarDisabled,
    CarElectric,
    Motorcycle,
}

pub type PaymentMethod = String;

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
pub struct Metadata {
    pub name: String,
    pub address: String,
    pub location: Coordinate,
    pub resources: Vec<Url>,
    pub total_spots: HashMap<SpotType, u32>,
    #[serde(default)]
    pub max_width: Option<i32>,
    #[serde(default)]
    pub max_height: Option<i32>,
    pub features: Vec<Feature>,
    pub payment_methods: Vec<PaymentMethod>,
    #[serde(skip_serializing_if = "HashMap::is_empty", default)]
    pub comment: HashMap<LanguageCode, String>,
    pub currency: Currency,
    pub rules: Vec<Rule>,
}
