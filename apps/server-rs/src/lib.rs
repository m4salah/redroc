use config::Config;

pub mod config;
pub mod handlers;
#[derive(Clone)]
pub struct AppState {
    config: Config,
}
