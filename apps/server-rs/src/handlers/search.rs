use super::grpc::get_thumbnail_client::GetThumbnailClient;
use super::grpc::GetThumbnailImagesRequest;
use axum::extract::Query;
use axum::response::IntoResponse;
use axum::Json;
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Params {
    #[serde(alias = "q")]
    query: String,
}

pub async fn search(Query(params): Query<Params>) -> impl IntoResponse {
    // TODO: We need a way to handle the error better
    // TODO: We need to make this url into env variable
    let mut client = GetThumbnailClient::connect("https://redroc-search-jo7doiawta-uc.a.run.app")
        .await
        .unwrap();
    let request = tonic::Request::new(GetThumbnailImagesRequest {
        search_keyword: params.query,
    });
    let response = client.get_thumbnail(request).await.unwrap();
    Json(response.into_inner().storage_url)
}
