use crate::parking_lot;
use arc_swap::ArcSwap;
use miette::{IntoDiagnostic, Result};
use reqwest::Url;
use std::collections::HashMap;
use std::sync::Arc;

#[derive(Debug, Clone)]
pub struct StorekeeperClient {
    access_token: Arc<ArcSwap<String>>,
    client: reqwest::Client,
    base_url: Url,
}

impl StorekeeperClient {
    pub fn new(base_url: Url, access_token: Arc<ArcSwap<String>>) -> Self {
        StorekeeperClient {
            client: reqwest::Client::new(),
            access_token,
            base_url,
        }
    }

    pub async fn update_parking_lot_states(
        &self,
        states: HashMap<parking_lot::ID, parking_lot::State>,
    ) -> Result<()> {
        let url = self.base_url.join(&format!("/parking-lot/state")).unwrap();
        let request = self
            .client
            .post(url)
            .json(&states)
            .bearer_auth(self.access_token.load());
        let response = request.send().await.into_diagnostic()?;
        response.error_for_status().into_diagnostic()?;
        Ok(())
    }

    pub async fn update_parking_lot_metadatas(
        &self,
        states: HashMap<parking_lot::ID, parking_lot::Metadata>,
    ) -> Result<()> {
        let url = self
            .base_url
            .join(&format!("/parking-lot/metadata"))
            .unwrap();
        let request = self
            .client
            .post(url)
            .json(&states)
            .bearer_auth(self.access_token.load());
        let response = request.send().await.into_diagnostic()?;
        response.error_for_status().into_diagnostic()?;
        Ok(())
    }
}
