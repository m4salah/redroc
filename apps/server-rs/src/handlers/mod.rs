use crate::AppState;

mod download;
mod health;
mod index;
mod search;

pub mod grpc {
    tonic::include_proto!("grpc");
}

pub fn router(app_state: AppState) -> axum::Router {
    axum::Router::new()
        .nest("/", download::router(app_state.clone()))
        .nest("/", health::router())
        .nest("/", index::router())
        .nest("/", search::router(app_state))
}
