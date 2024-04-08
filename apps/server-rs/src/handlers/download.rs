use crate::handlers::grpc::download_photo_client::DownloadPhotoClient;
use crate::handlers::grpc::DownloadPhotoRequest;
use crate::AppState;

use axum::extract::{Path, State};
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum::routing::get;

pub fn router(app_state: AppState) -> axum::Router {
    axum::Router::new()
        .route("/download/:img_name", get(get_img))
        .with_state(app_state)
}

async fn get_img(
    State(app_state): State<AppState>,
    Path(img_name): Path<String>,
) -> Result<Response, StatusCode> {
    // TODO: We need a way to handle the error better
    // TODO: We need to make this url into env variable
    tracing::info!("Downloading image with name {}", img_name.as_str());
    let mut client =
        DownloadPhotoClient::connect(app_state.config.download_backend_addr.to_string())
            .await
            .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    let request = tonic::Request::new(DownloadPhotoRequest {
        img_name: img_name.clone(),
    });
    let response = client
        .download(request)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    tracing::info!("Image Downloaded {}", img_name.as_str());
    Ok((
        ([(axum::http::header::CONTENT_TYPE, "image/png")]),
        response.into_inner().img_blob,
    )
        .into_response())
}
