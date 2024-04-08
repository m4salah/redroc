use axum::{response::IntoResponse, routing::get};

pub fn router() -> axum::Router {
    axum::Router::new().route("/health", get(health))
}
async fn health() -> impl IntoResponse {
    String::from("healthy!")
}
