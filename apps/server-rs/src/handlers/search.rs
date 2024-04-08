use crate::handlers::grpc::get_thumbnail_client::GetThumbnailClient;
use crate::handlers::grpc::GetThumbnailImagesRequest;
use crate::AppState;

use axum::extract::{Query, State};
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum::routing::get;
use axum::Json;
use serde::Deserialize;

pub fn router(app_state: AppState) -> axum::Router {
    axum::Router::new()
        .route("/search", get(search))
        .with_state(app_state)
}

#[derive(Debug, Deserialize)]
struct Params {
    #[serde(alias = "q")]
    query: String,
}

async fn search(
    State(app_state): State<AppState>,
    Query(params): Query<Params>,
) -> Result<Response, StatusCode> {
    tracing::info!("params passed for search: {:?}", params);
    // TODO: We need a way to handle the error better
    // TODO: We need to make this url into env variable
    let mut client = GetThumbnailClient::connect(app_state.config.search_backend_addr.to_string())
        .await
        .map_err(|e| {
            tracing::error!("error while connecting to search grpc service: {e:?}");
            StatusCode::INTERNAL_SERVER_ERROR
        })?;
    let request = tonic::Request::new(GetThumbnailImagesRequest {
        search_keyword: params.query,
    });
    let response = client
        .get_thumbnail(request)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    Ok(Json(response.into_inner().storage_url).into_response())
}
