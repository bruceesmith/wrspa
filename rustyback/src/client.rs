use reqwest::{header::{HeaderMap, USER_AGENT}, Client as ReqwestClient, StatusCode, redirect::Policy};
use tracing::error;
use url::Url;

const DEFAULT_WIKI_URL: &str = "https://en.wikipedia.org";

#[derive(Debug)]
pub enum ClientError {
    Reqwest(reqwest::Error),
    StatusError(StatusCode),
    MissingLocationHeader,
    InvalidLocationHeader(url::ParseError),
}

impl From<reqwest::Error> for ClientError {
    fn from(err: reqwest::Error) -> Self {
        ClientError::Reqwest(err)
    }
}

impl From<url::ParseError> for ClientError {
    fn from(err: url::ParseError) -> Self {
        ClientError::InvalidLocationHeader(err)
    }
}

pub struct Client {
    wiki_url: String,
    client: ReqwestClient,
}

impl Client {
    pub fn new(wiki_url: Option<&str>) -> Self {
        let wiki_url = wiki_url.unwrap_or(DEFAULT_WIKI_URL).to_string();
        let mut headers = HeaderMap::new();
        headers.insert(USER_AGENT, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36".parse().unwrap());
        let client = ReqwestClient::builder()
            .redirect(Policy::none())
            .default_headers(headers)
            .build()
            .unwrap();
        Self {
            wiki_url,
            client,
        }
    }

    pub async fn get(&self, path: &str) -> Result<bytes::Bytes, ClientError> {
        let url = format!("{}{}", self.wiki_url, path);
        tracing::trace!(url);
        let response = self.client.get(&url).send().await?;
        if !response.status().is_success() {
            let status = response.status();
            error!("unexpected status: {}", status);
            return Err(ClientError::StatusError(status));
        }
        Ok(response.bytes().await?)
    }

    pub async fn get_random(&self) -> Result<String, ClientError> {
        let url = format!("{}/wiki/Special:Random", self.wiki_url);
        let response = self.client.get(&url).send().await?;
        if !response.status().is_success() && !response.status().is_redirection() {
            let status = response.status();
            error!("unexpected status getting random: {}", status);
            return Err(ClientError::StatusError(status));
        }
        let location = response
            .headers()
            .get("location")
            .ok_or(ClientError::MissingLocationHeader)?;
        let url = Url::parse(location.to_str().unwrap())?;
        let path = url.path().trim_start_matches("/wiki/").to_string();
        Ok(path)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use mockito::mock;

    #[test]
    fn test_new_client() {
        let client = Client::new(None);
        assert_eq!(client.wiki_url, DEFAULT_WIKI_URL);

        let client = Client::new(Some("http://localhost:8080"));
        assert_eq!(client.wiki_url, "http://localhost:8080");
    }

    #[tokio::test]
    async fn test_get() {
        let m = mock("GET", "/test").with_status(200).with_body("hello").create();

        let client = Client::new(Some(&mockito::server_url()));
        let result = client.get("/test").await;
        assert!(result.is_ok());
        assert_eq!(result.unwrap(), "hello");
        m.assert();
    }

    #[tokio::test]
    async fn test_get_error() {
        let m = mock("GET", "/test").with_status(500).create();

        let client = Client::new(Some(&mockito::server_url()));
        let result = client.get("/test").await;
        assert!(result.is_err());
        m.assert();
    }

    #[tokio::test]
    async fn test_get_random() {
        let m = mock("GET", "/wiki/Special:Random")
            .with_status(302)
            .with_header("location", &format!("{}/wiki/Test", mockito::server_url()))
            .create();

        let client = Client::new(Some(&mockito::server_url()));
        let result = client.get_random().await;
        assert!(result.is_ok());
        assert_eq!(result.unwrap(), "Test");
        m.assert();
    }

    #[tokio::test]
    async fn test_get_random_error() {
        let m = mock("GET", "/wiki/Special:Random").with_status(500).create();

        let client = Client::new(Some(&mockito::server_url()));
        let result = client.get_random().await;
        assert!(result.is_err());
        m.assert();
    }
}