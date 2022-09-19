pub mod parking_lot;

#[cfg(feature = "authorization")]
pub mod authorization;

#[cfg(feature = "authorization")]
pub use authorization::AuthorizationClient;

#[cfg(feature = "authorization")]
#[cfg(feature = "storekeeper")]
pub mod storekeeper;

#[cfg(feature = "storekeeper")]
pub use storekeeper::StorekeeperClient;

#[cfg(feature = "provider")]
pub mod provider;

#[cfg(feature = "provider")]
pub use provider::Provider;

use serde::Deserialize;
use serde::Serialize;

pub type LanguageCode = String;

#[derive(Debug, Clone, Serialize, Deserialize, strum::Display, strum::EnumString)]
#[serde(rename_all = "SCREAMING-KEBAB-CASE")]
#[strum(serialize_all = "SCREAMING-KEBAB-CASE")]
pub enum Weekday {
    Monday,
    Tuesday,
    Wednesday,
    Thursday,
    Friday,
    Saturday,
    Sunday,
}
