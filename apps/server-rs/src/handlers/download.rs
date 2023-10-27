use crate::AppState;

use super::grpc::download_photo_client::DownloadPhotoClient;
use super::grpc::DownloadPhotoRequest;
use axum::extract::{Path, State};
use axum::response::IntoResponse;

pub async fn get_img(
    State(app_state): State<AppState>,
    Path(img_name): Path<String>,
) -> impl IntoResponse {
    // TODO: We need a way to handle the error better
    // TODO: We need to make this url into env variable
    let mut client = DownloadPhotoClient::connect(app_state.config.download_backend_addr)
        .await
        .unwrap();
    let request = tonic::Request::new(DownloadPhotoRequest { img_name });
    let response = client.download(request).await.unwrap();
    (
        ([(axum::http::header::CONTENT_TYPE, "image/png")]),
        response.into_inner().img_blob,
    )
}
