use anyhow::Result;
use calculator_rust::calculator::Calculator;
use rmcp::ServiceExt;
use tokio::io::{stdin, stdout};
use tracing_subscriber::EnvFilter;

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env().add_directive(tracing::Level::DEBUG.into()))
        .with_writer(std::io::stderr)
        .with_ansi(false)
        .init();

    tracing::info!("Starting RMCP server...");

    let transport = (stdin(), stdout());

    // create an instance of the calculator
    let service = Calculator::new().serve(transport).await.inspect_err(|e| {
        tracing::error!("serving error: {:?}", e);
    })?;

    service.waiting().await?;
    Ok(())
}
