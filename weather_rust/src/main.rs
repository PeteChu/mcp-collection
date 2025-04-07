use std::env;

use anyhow::Result;
use rmcp::ServiceExt;
use tokio::io::{stdin, stdout};
use tracing_subscriber::EnvFilter;
use weather_rust::weather::Weather;

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env().add_directive(tracing::Level::DEBUG.into()))
        .with_writer(std::io::stderr)
        .with_ansi(false)
        .init();

    tracing::info!("Starting MCP server...");

    let api_key = env::var("OPENWEATHER_API_KEY").unwrap_or_else(|_| {
        tracing::error!("Please set the OPENWEATHER_API_KEY environment variable.");
        std::process::exit(1);
    });

    let transport = (stdin(), stdout());
    let service = Weather::new(api_key).serve(transport).await?;

    service.waiting().await?;
    Ok(())
}
