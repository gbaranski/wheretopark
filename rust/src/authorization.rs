use std::collections::HashSet;

use miette::{IntoDiagnostic, Result};
use reqwest::Url;
use serde::Deserialize;
use serde::Serialize;

#[derive(Debug, Clone)]
pub struct AuthorizationClient {
    client: reqwest::Client,
    base_url: Url,
    client_id: String,
    client_secret: String,
}

#[derive(Debug, Clone, Eq, PartialEq, Hash, Serialize, Deserialize, strum::Display)]
pub enum AccessType {
    #[strum(serialize = "read:metadata")]
    #[serde(rename = "read:metadata")]
    ReadMetadata,
    #[strum(serialize = "write:metadata")]
    #[serde(rename = "write:metadata")]
    WriteMetadata,
    #[strum(serialize = "read:state")]
    #[serde(rename = "read:state")]
    ReadState,
    #[strum(serialize = "write:state")]
    #[serde(rename = "write:state")]
    WriteState,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TokenResponse {
    pub access_token: String,
    pub expires_in: u64,
    pub scope: String,
    pub token_type: String,
}

impl AuthorizationClient {
    pub fn new(base_url: Url, client_id: String, client_secret: String) -> Self {
        Self {
            client: reqwest::Client::new(),
            base_url,
            client_id,
            client_secret,
        }
    }

    pub async fn token(&self, access_scope: HashSet<AccessType>) -> Result<TokenResponse> {
        let url = self.base_url.join(&format!("/oauth/token")).unwrap();
        let request = self.client.post(url).query(&[
            ("client_id", self.client_id.to_string()),
            ("client_secret", self.client_secret.to_string()),
            ("grant_type", "client_credentials".to_string()),
            (
                "scope",
                access_scope
                    .iter()
                    .map(|s| s.to_string())
                    .collect::<Vec<String>>()
                    .join(" "),
            ),
        ]);
        let response = request.send().await.into_diagnostic()?;
        response.error_for_status_ref().into_diagnostic()?;
        let json = response.json().await.into_diagnostic()?;
        Ok(json)
    }
}
