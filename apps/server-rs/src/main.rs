use std::{env, net::SocketAddr};

use clap::Parser;
use config::Config;
use tower_http::trace::TraceLayer;

use handlers::router;

mod config;
mod handlers;

#[derive(Clone)]
pub struct AppState {
    config: Config,
}
#[tokio::main]
async fn main() {
    // Get the env variables from the .env file found in the app directory not the root directory
    // This returns an error if the `.env` file doesn't exist, but that's not what we want
    let dot_env_path = env::current_dir()
        .unwrap()
        .join("apps")
        .join("server-rs")
        .join(".env");
    dotenv::from_path(dot_env_path).ok();

    // Initialize the logger.
    env_logger::init();
    let subscriber = tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .with_target(false)
        .json()
        .finish();
    // use that subscriber to process traces emitted after this point
    tracing::subscriber::set_global_default(subscriber).unwrap();

    // tracing layer
    let tracing_layer = TraceLayer::new_for_http();

    // Parse our configuration from the environment.
    // This will exit with a help message if something is wrong.
    let config = config::Config::parse();
    let app_state = AppState {
        config: config.clone(),
    };

    let app = router(app_state).layer(tracing_layer);

    let addr = SocketAddr::from(([0, 0, 0, 0], config.port));
    tracing::info!("listening on {}", addr);

    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
