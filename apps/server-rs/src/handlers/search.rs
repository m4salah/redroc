use crate::AppState;

use super::grpc::get_thumbnail_client::GetThumbnailClient;
use super::grpc::GetThumbnailImagesRequest;
use axum::extract::{Query, State};
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum::Json;
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Params {
    #[serde(alias = "q")]
    query: String,
}

pub async fn search(
    State(app_state): State<AppState>,
    Query(params): Query<Params>,
) -> Result<Response, StatusCode> {
    // TODO: We need a way to handle the error better
    // TODO: We need to make this url into env variable
    let mut client = GetThumbnailClient::connect(app_state.config.search_backend_addr)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    let request = tonic::Request::new(GetThumbnailImagesRequest {
        search_keyword: params.query,
    });
    let response = client
        .get_thumbnail(request)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    Ok(Json(response.into_inner().storage_url).into_response())
}
