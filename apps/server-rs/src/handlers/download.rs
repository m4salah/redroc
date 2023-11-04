use crate::AppState;

use super::grpc::download_photo_client::DownloadPhotoClient;
use super::grpc::DownloadPhotoRequest;
use axum::extract::{Path, State};
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};

pub async fn get_img(
    State(app_state): State<AppState>,
    Path(img_name): Path<String>,
) -> Result<Response, StatusCode> {
    // TODO: We need a way to handle the error better
    // TODO: We need to make this url into env variable
    tracing::info!("Downloading image with name {}", img_name.as_str());
    let mut client = DownloadPhotoClient::connect(app_state.config.download_backend_addr)
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
