use axum::{response::Html, routing::get};

pub fn router() -> axum::Router {
    axum::Router::new().route("/", get(index))
}
async fn index() -> Html<&'static str> {
    Html("<h1>Welcome to Redroc</h1>")
}
