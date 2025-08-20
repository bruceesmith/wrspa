use axum::{ 
    Router,
    extract::{State},
    http::{StatusCode, Uri},
    response::{Html, Json},
    routing::{get, post},
};
use scraper::{Html as ScraperHtml, Selector};
use std::{fs, io, sync::Arc};
use tower_http::services::ServeDir;

use crate::api::{SettingsResponse, SpecialRandomResponse, WikiPageRequest};
use crate::client::{Client, ClientError};

#[derive(Debug)]
pub enum ServerError {
    InvalidPort(u16),
    StaticPathError(io::Error),
    StaticPathNotADirectory,
}

pub struct Server {
    client: Arc<Client>,
    port: u16,
    static_path: String,
}

impl Server {
    pub fn new(port: u16, static_path: String, client: Arc<Client>) -> Result<Self, ServerError> {
        if port <= 1024 {
            return Err(ServerError::InvalidPort(port));
        }

        match fs::metadata(&static_path) {
            Ok(metadata) => {
                if !metadata.is_dir() {
                    return Err(ServerError::StaticPathNotADirectory);
                }
            }
            Err(e) => return Err(ServerError::StaticPathError(e)),
        }

        Ok(Self {
            port,
            static_path,
            client,
        })
    }

    pub async fn serve(self) {
        let app = Router::new()
            .route("/api/settings", get(settings_handler))
            .route("/api/specialrandom", get(special_random_handler))
            .route("/api/wikipage", post(wiki_page_handler))
            .route("/w/*path", get(wikipedia_file_handler))
            .route("/static/*path", get(wikipedia_file_handler))
            .nest_service("/", ServeDir::new(&self.static_path))
            .with_state(self.client.clone());

        let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", self.port))
            .await
            .unwrap();
        axum::serve(listener, app).await.unwrap();
    }
}

async fn settings_handler() -> Json<SettingsResponse> {
    let response = SettingsResponse {
        loglevel: "info".to_string(),
        traceids: vec![],
    };
    Json(response)
}

async fn special_random_handler(
    State(client): State<Arc<Client>>,
) -> Result<Json<SpecialRandomResponse>, StatusCode> {
    let start = client.get_random().await.map_err(|e| match e {
        ClientError::Reqwest(_) => StatusCode::INTERNAL_SERVER_ERROR,
        ClientError::StatusError(s) => s,
        ClientError::MissingLocationHeader => StatusCode::INTERNAL_SERVER_ERROR,
        ClientError::InvalidLocationHeader(_) => StatusCode::INTERNAL_SERVER_ERROR,
    })?;
    let goal = client.get_random().await.map_err(|e| match e {
        ClientError::Reqwest(_) => StatusCode::INTERNAL_SERVER_ERROR,
        ClientError::StatusError(s) => s,
        ClientError::MissingLocationHeader => StatusCode::INTERNAL_SERVER_ERROR,
        ClientError::InvalidLocationHeader(_) => StatusCode::INTERNAL_SERVER_ERROR,
    })?;
    let response = SpecialRandomResponse { start, goal };
    Ok(Json(response))
}

async fn wiki_page_handler(
    State(client): State<Arc<Client>>,
    Json(request): Json<WikiPageRequest>,
) -> Result<Html<String>, StatusCode> {
    let path = format!("/wiki/{}", &request.subject);
    let page = client.get(&path).await.map_err(|e| match e {
        ClientError::Reqwest(_) => StatusCode::INTERNAL_SERVER_ERROR,
        ClientError::StatusError(s) => s,
        ClientError::MissingLocationHeader => StatusCode::INTERNAL_SERVER_ERROR,
        ClientError::InvalidLocationHeader(_) => StatusCode::INTERNAL_SERVER_ERROR,
    })?;

    let body = extract_body(&page).map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    Ok(Html(body))
}

async fn wikipedia_file_handler(
    State(client): State<Arc<Client>>,
    uri: Uri,
) -> Result<bytes::Bytes, StatusCode> {
    client.get(uri.path()).await.map_err(|e| match e {
        ClientError::Reqwest(_) => StatusCode::INTERNAL_SERVER_ERROR,
        ClientError::StatusError(s) => s,
        ClientError::MissingLocationHeader => StatusCode::INTERNAL_SERVER_ERROR,
        ClientError::InvalidLocationHeader(_) => StatusCode::INTERNAL_SERVER_ERROR,
    })
}

fn extract_body(page: &[u8]) -> Result<String, ()> {
    let html = String::from_utf8_lossy(page);
    let document = ScraperHtml::parse_document(&html);
    let selector = Selector::parse("body").unwrap();
    let body = document.select(&selector).next().ok_or(())?;
    Ok(body.inner_html())
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::body::{to_bytes, Body};
    use axum::http::{Request, StatusCode};
    use mockito::mock;
    use std::sync::Arc;
    use tempfile::tempdir;
    use tower::ServiceExt;

    #[test]
    fn test_new_server() {
        let client = Arc::new(Client::new(None));
        let server = Server::new(8080, ".".to_string(), client.clone());
        assert!(server.is_ok());

        let server = Server::new(80, ".".to_string(), client.clone());
        assert!(server.is_err());

        let server = Server::new(8080, "not-a-directory".to_string(), client.clone());
        assert!(server.is_err());
    }

    #[tokio::test]
    async fn test_settings_handler() {
        let client = Arc::new(Client::new(None));
        let app = Router::new().route("/api/settings", get(settings_handler)).with_state(client.clone());

        let response = app
            .oneshot(Request::builder().uri("/api/settings").body(Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = to_bytes(response.into_body(), 1_000_000).await.unwrap();
        let settings: SettingsResponse = serde_json::from_slice(&body).unwrap();
        assert_eq!(settings.loglevel, "info");
    }

    #[tokio::test]
    async fn test_special_random_handler() {
        let m1 = mock("GET", "/wiki/Special:Random")
            .with_status(302)
            .with_header("location", &format!("{}/wiki/Start", mockito::server_url()))
            .create();
        let m2 = mock("GET", "/wiki/Special:Random")
            .with_status(302)
            .with_header("location", &format!("{}/wiki/Goal", mockito::server_url()))
            .create();

        let client = Arc::new(Client::new(Some(&mockito::server_url())));
        let app = Router::new()
            .route("/api/specialrandom", get(special_random_handler))
            .with_state(client.clone());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/api/specialrandom")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = to_bytes(response.into_body(), 1_000_000).await.unwrap();
        let special_random: SpecialRandomResponse = serde_json::from_slice(&body).unwrap();
        assert_eq!(special_random.start, "Start");
        assert_eq!(special_random.goal, "Goal");
        m1.assert();
        m2.assert();
    }

    #[tokio::test]
    async fn test_wiki_page_handler() {
        let m = mock("GET", "/wiki/Test")
            .with_status(200)
            .with_body("<html><body><p>Test</p></body></html>")
            .create();

        let client = Arc::new(Client::new(Some(&mockito::server_url())));
        let app = Router::new()
            .route("/api/wikipage", post(wiki_page_handler))
            .with_state(client.clone());

        let request = WikiPageRequest { subject: "Test".to_string() };
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/wikipage")
                    .header("content-type", "application/json")
                    .body(Body::from(serde_json::to_vec(&request).unwrap()))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = to_bytes(response.into_body(), 1_000_000).await.unwrap();
        assert_eq!(body, "<p>Test</p>");
        m.assert();
    }

    #[tokio::test]
    async fn test_serve_static_files() {
        let dir = tempdir().unwrap();
        let file_path = dir.path().join("test.html");
        fs::write(&file_path, "<html><body>Test</body></html>").unwrap();

        let client = Arc::new(Client::new(None));
        let server = Server::new(8081, dir.path().to_str().unwrap().to_string(), client.clone()).unwrap();
        let app = Router::new()
            .nest_service("/", ServeDir::new(server.static_path));

        let response = app
            .oneshot(Request::builder().uri("/test.html").body(Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
        let body = to_bytes(response.into_body(), 1_000_000).await.unwrap();
        assert_eq!(body, "<html><body>Test</body></html>");
    }

    #[tokio::test]
    async fn test_wikipedia_file_handler() {
        let m = mock("GET", "/w/Test")
            .with_status(200)
            .with_body("hello")
            .create();

        let client = Arc::new(Client::new(Some(&mockito::server_url())));
        let app = Router::new()
            .route("/w/*path", get(wikipedia_file_handler))
            .with_state(client.clone());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/w/Test")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = to_bytes(response.into_body(), 1_000_000).await.unwrap();
        assert_eq!(body, "hello");
        m.assert();
    }
}