use rustyback::run;
use rustyback::server::ServerError;

#[tokio::main]
async fn main() {
    if let Err(e) = run().await {
        match e {
            ServerError::InvalidPort(p) => eprintln!("Invalid port: {}", p),
            ServerError::StaticPathError(e) => eprintln!("Static path error: {}", e),
            ServerError::StaticPathNotADirectory => eprintln!("Static path is not a directory"),
        }
        std::process::exit(1);
    }
}