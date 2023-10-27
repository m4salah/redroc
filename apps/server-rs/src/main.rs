use std::env;

use axum::{routing::get, Router};
use clap::Parser;
use config::Config;
use handlers::{download, health, index, search};

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

    // Parse our configuration from the environment.
    // This will exit with a help message if something is wrong.
    let config = config::Config::parse();
    let app_state = AppState { config };

    let app = Router::new()
        .route("/", get(index))
        .route("/health", get(health))
        .route("/download/:img_name", get(download::get_img))
        .route("/search", get(search::search))
        .with_state(app_state);

    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
