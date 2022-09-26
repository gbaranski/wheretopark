use std::{collections::HashMap, time::Duration};
use std::collections::HashSet;
use std::sync::Arc;

use arc_swap::ArcSwap;
use async_trait::async_trait;
use chrono::DateTime;
use chrono::NaiveDateTime;
use chrono::Utc;
use miette::{Context, IntoDiagnostic};
use miette::Result;
use tracing::Instrument;
use url::Url;

use crate::{AuthorizationClient, StorekeeperClient};
use crate::authorization::AccessType;
use crate::parking_lot::{ID, Metadata, State};

const METADATA_POLL_INTERVAL: Duration = Duration::from_secs(3600);
const STATE_POLL_INTERVAL: Duration = Duration::from_secs(30);

#[async_trait]
pub trait Provider {
    async fn poll_metadatas(&self) -> Result<HashMap<ID, Metadata>>;
    async fn poll_states(&self) -> Result<HashMap<ID, State>>;
}

async fn tick_metadata(
    provider: &impl Provider,
    storekeeper_client: &StorekeeperClient,
) -> Result<()> {
    let metadatas = provider.poll_metadatas().await?;
    let count = metadatas.len();
    storekeeper_client
        .update_parking_lot_metadatas(metadatas)
        .await?;
    tracing::info!("updated {count} metadatas");
    Ok(())
}

async fn tick_state(
    provider: &impl Provider,
    storekeeper_client: &StorekeeperClient,
) -> Result<()> {
    let states = provider.poll_states().await?;
    let states = states
        .into_iter()
        .map(|(id, mut state)| {
            let timestamp = state.last_updated.timestamp();
            let naive = NaiveDateTime::from_timestamp(timestamp, 0);
            state.last_updated = DateTime::<Utc>::from_utc(naive, Utc);
            (id, state)
        })
        .collect::<HashMap<_, _>>();
    let count = states.len();

    storekeeper_client.update_parking_lot_states(states).await?;
    tracing::info!("updated {count} states");
    Ok(())
}

const STOREKEEPER_URL: &'static str = "https://storekeeper.wheretopark.app";
const AUTHORIZATION_URL: &'static str = "https://authorization.wheretopark.app";

fn get_env(name: &str) -> Result<String> {
    std::env::var(name)
        .into_diagnostic()
        .with_context(|| format!("env var: {}", name))
}

fn get_env_url_or(name: &str, default: &str) -> Result<Url> {
    let value = get_env(name).unwrap_or_else(|_| default.to_string());
    Url::parse(&value).into_diagnostic()
}

pub async fn run(provider: impl Provider + Send + Sync + Clone + 'static) -> Result<()> {
    let storekeeper_url = get_env_url_or("STOREKEEPER_URL", STOREKEEPER_URL)?;
    let authorization_url = get_env_url_or("AUTHORIZATION_URL", AUTHORIZATION_URL)?;
    let client_id = get_env("CLIENT_ID")?;
    let client_secret = get_env("CLIENT_SECRET")?;
        let authorization_client =
        AuthorizationClient::new(authorization_url, client_id, client_secret);
    let token_response = authorization_client
        .token(HashSet::from([AccessType::WriteMetadata, AccessType::WriteState]))
        .await?;
    let access_token = Arc::new(ArcSwap::new(Arc::new(token_response.access_token.clone())));
    let storekeeper_client = StorekeeperClient::new(storekeeper_url, access_token);
    let token_task = tokio::spawn({
        let authorization_client = authorization_client.clone();
        async move {
            let mut token_response = token_response;
            loop {
                let duration = Duration::from_secs(token_response.expires_in);
                tokio::time::sleep(duration).await;
                tracing::info!("refreshing token");
                token_response = authorization_client
                    .token(HashSet::from([AccessType::WriteMetadata, AccessType::WriteState]))
                    .await
                    .unwrap();
            }
        }
    });
    tick_metadata(&provider, &storekeeper_client).await?; // Prepare
    // TODO: Maybe merge those two?
    let metadata_task = tokio::spawn(
        {
            let storekeeper_client = storekeeper_client.clone();
            let provider = provider.clone();
            async move {
                let mut interval = tokio::time::interval(METADATA_POLL_INTERVAL);
                interval.tick().await;
                loop {
                    interval.tick().await;
                    if let Err(err) = tick_metadata(&provider, &storekeeper_client).await {
                        tracing::error!("tick metadata error: {err}")
                    }
                }
            }
        }
            .in_current_span(),
    );
    let state_task = tokio::spawn(
        async move {
            let mut interval = tokio::time::interval(STATE_POLL_INTERVAL);
            loop {
                interval.tick().await;
                if let Err(err) = tick_state(&provider, &storekeeper_client).await {
                    tracing::error!("tick state error: {err}")
                }
            }
        }
            .in_current_span(),
    );
    tokio::select! {
        result = metadata_task => result.into_diagnostic(),
        result = state_task => result.into_diagnostic(),
        result = token_task => result.into_diagnostic(),
        _ = tokio::signal::ctrl_c() => {
            tracing::info!("shutting down");
            Ok(())
        }
    }
}
