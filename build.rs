use std::error::Error;

fn main() -> Result<(), Box<dyn Error>> {
    tonic_build::configure()
        .build_server(false)
        .build_client(true)
        .compile(
            &[
                "libs/proto/download.proto",
                "libs/proto/upload.proto",
                "libs/proto/search.proto",
            ],
            &["proto"],
        )
        .unwrap();
    Ok(())
}
