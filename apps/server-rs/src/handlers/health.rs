use axum::response::IntoResponse;

pub async fn health() -> impl IntoResponse {
    String::from("healthy!")
}
