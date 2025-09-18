pub mod api;
pub mod client;
pub mod server;

use clap::Parser;
use client::Client;
use server::{Server, ServerError};
use std::sync::Arc;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    #[arg(short, long, default_value_t = 8080)]
    port: u16,

    #[arg(short, long)]
    static_path: String,
}

pub async fn run() -> Result<(), ServerError> {
    let args = Args::parse();

    let client = Arc::new(Client::new(None));
    let server = Server::new(args.port, args.static_path, client)?;

    server.serve().await;
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::process::Command;

    #[test]
    fn test_arg_parsing() {
        let args = Args::parse_from(&["rustyback", "-p", "8080", "-s", "."]);
        assert_eq!(args.port, 8080);
        assert_eq!(args.static_path, ".");
    }

    #[test]
    fn test_run_invalid_port() {
        let mut cmd = Command::new("cargo");
        cmd.args(&["run", "--", "-p", "80", "-s", "."]);
        let output = cmd.output().unwrap();
        assert!(!output.status.success());
        assert!(String::from_utf8_lossy(&output.stderr).contains("Invalid port: 80"));
    }

    #[test]
    fn test_run_invalid_static_path() {
        let mut cmd = Command::new("cargo");
        cmd.args(&["run", "--", "-p", "8080", "-s", "not-a-directory"]);
        let output = cmd.output().unwrap();
        assert!(!output.status.success());
        assert!(String::from_utf8_lossy(&output.stderr).contains("Static path error"));
    }
}