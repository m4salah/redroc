use std::net::SocketAddr;

#[derive(clap::Parser, Clone)]
pub struct Config {
    #[clap(long, env)]
    pub download_backend_addr: SocketAddr,

    #[clap(long, env)]
    pub upload_backend_addr: SocketAddr,

    #[clap(long, env)]
    pub search_backend_addr: SocketAddr,

    #[clap(long, env)]
    pub port: u16,
}
