use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct SettingsResponse {
    pub loglevel: String,
    pub traceids: Vec<String>,
}

#[derive(Serialize, Deserialize)]
pub struct SpecialRandomResponse {
    pub start: String,
    pub goal: String,
}

#[derive(Serialize, Deserialize)]
pub struct WikiPageRequest {
    pub subject: String,
}

#[derive(Serialize, Deserialize)]
pub struct WikiPageResponse {
    pub page: String,
    pub error: String,
}
