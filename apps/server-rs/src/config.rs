#[derive(clap::Parser, Clone)]
pub struct Config {
    #[clap(long, env)]
    pub download_backend_addr: String,

    #[clap(long, env)]
    pub upload_backend_addr: String,

    #[clap(long, env)]
    pub search_backend_addr: String,
}
