use axum::extract::Path;
use axum::response::{Html, IntoResponse, Response};
use axum::{routing::get, Router};

async fn index() -> Html<&'static str> {
    Html("<h1>Welcome to Redroc</h1>")
}

async fn health() -> Response {
    String::from("healthy!").into_response()
}

async fn download(Path(img_name): Path<String>) -> Response {
    String::from(img_name).into_response()
}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/", get(index))
        .route("/health", get(health))
        .route("/download/:img_name", get(download));

    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
