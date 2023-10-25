use super::grpc::download_photo_client::DownloadPhotoClient;
use super::grpc::DownloadPhotoRequest;
use axum::extract::Path;
use axum::response::IntoResponse;

pub async fn get_img(Path(img_name): Path<String>) -> impl IntoResponse {
    // TODO: We need a way to handle the error better
    // TODO: We need to make this url into env variable
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
