use http_types::Url;

#[derive(clap::Parser, Clone, Debug)]
pub struct Config {
    #[clap(long, env)]
    pub download_backend_addr: Url,

    #[clap(long, env)]
    pub upload_backend_addr: Url,

    #[clap(long, env)]
    pub search_backend_addr: Url,

    #[clap(long, env)]
    pub port: u16,
}
