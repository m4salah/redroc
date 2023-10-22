use super::*;
use axum::extract::Path;
use axum::response::IntoResponse;
use grpc::download_photo_client::DownloadPhotoClient;
use grpc::DownloadPhotoRequest;

pub async fn get_img(Path(img_name): Path<String>) -> impl IntoResponse {
    // TODO: We need a way to handle the error better
    let mut client =
        DownloadPhotoClient::connect("https://redroc-download-jo7doiawta-uc.a.run.app")
            .await
            .unwrap();
    let request = tonic::Request::new(DownloadPhotoRequest { img_name });
    let response = client.download(request).await.unwrap();
    (
        ([(axum::http::header::CONTENT_TYPE, "image/png")]),
        response.into_inner().img_blob,
    )
}
