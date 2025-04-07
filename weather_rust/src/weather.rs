use rmcp::{
    Error as McpError, ServerHandler,
    model::{CallToolResult, Content, ServerCapabilities, ServerInfo},
    tool,
};
use serde_json::Value;

#[derive(Debug, Clone)]
pub struct Weather {
    pub api_key: String,
}

#[derive(schemars::JsonSchema, serde::Serialize, serde::Deserialize, Debug)]
struct Location {
    lat: f64,
    lon: f64,
}

#[derive(serde::Deserialize)]
struct QueriedLocation {
    lat: f64,
    lon: f64,
}

const URL: &str = "http://api.openweathermap.org";

#[tool(tool_box)]
impl Weather {
    pub fn new(api_key: String) -> Self {
        Self { api_key }
    }

    #[tool(description = "Get latitude and longitude of a location")]
    async fn get_location(
        &self,
        #[tool(param)]
        #[schemars(description = "location name")]
        location: String,
    ) -> Result<CallToolResult, McpError> {
        let api_key = &self.api_key;
        let url = format!("{URL}/geo/1.0/direct?q={location}&limit=1&appid={api_key}");

        let resp = reqwest::get(url).await.map_err(|e| {
            McpError::invalid_request(
                format!("OpenWeatherMap Api request failed {:?}", e.to_string()),
                None,
            )
        })?;

        if let Some(err) = handle_request_resp(&resp) {
            return Err(err);
        }

        let items = resp.json::<Vec<QueriedLocation>>().await.map_err(|e| {
            McpError::internal_error(format!("Parse failed: {:?}", e.to_string()), None)
        })?;

        if items.is_empty() {
            return Err(McpError::resource_not_found("Location not found", None));
        }

        let item = items.first().unwrap();

        let location = Location {
            lat: item.lat,
            lon: item.lon,
        };

        Ok(CallToolResult::success(vec![Content::json(location)?]))
    }

    #[tool(description = "Get current location weather")]
    async fn get_current_weather(
        &self,
        #[tool(aggr)] Location { lat, lon }: Location,
    ) -> Result<CallToolResult, McpError> {
        let api_key = &self.api_key;
        let resp = reqwest::get(format!(
            "{URL}/data/2.5/weather?lat={lat}&lon={lon}&appid={api_key}"
        ))
        .await
        .map_err(|e| {
            McpError::invalid_request(
                format!("OpenWeatherMap Api request failed {:?}", e.to_string()),
                None,
            )
        })?;

        if let Some(err) = handle_request_resp(&resp) {
            return Err(err);
        }

        let data = resp.json::<Value>().await.map_err(|e| {
            McpError::internal_error(format!("Parse failed: {:?}", e.to_string()), None)
        })?;

        Ok(CallToolResult::success(vec![Content::json(data)?]))
    }

    #[tool(description = "Get 5-day 3-hour forecast")]
    async fn get_5day_3hour_forecast(
        &self,
        #[tool(aggr)] Location { lat, lon }: Location,
    ) -> Result<CallToolResult, McpError> {
        let api_key = &self.api_key;

        let resp = reqwest::get(format!(
            "{URL}/data/2.5/forecast?lat={lat}&lon={lon}&appid={api_key}"
        ))
        .await
        .map_err(|e| {
            McpError::invalid_request(
                format!("OpenWeatherMap Api request failed {:?}", e.to_string()),
                None,
            )
        })?;

        if let Some(err) = handle_request_resp(&resp) {
            return Err(err);
        }

        let data = resp.json::<Value>().await.map_err(|e| {
            McpError::internal_error(format!("Parse failed: {:?}", e.to_string()), None)
        })?;

        Ok(CallToolResult::success(vec![Content::json(data)?]))
    }
}

fn handle_request_resp(resp: &reqwest::Response) -> Option<McpError> {
    if !resp.status().is_success() {
        let status = resp.status();
        return match status.as_u16() {
            401 => Some(McpError::invalid_request("Invalid API key", None)),
            404 => Some(McpError::resource_not_found("API endpoint not found", None)),
            429 => Some(McpError::invalid_request("Rate limit exeeded", None)),
            _ => Some(McpError::internal_error(
                format!("OpenWeatherMap API error: {}", status),
                None,
            )),
        };
    }
    None
}

#[tool(tool_box)]
impl ServerHandler for Weather {
    fn get_info(&self) -> ServerInfo {
        ServerInfo {
            instructions: Some("Get weather from you input location".to_string()),
            capabilities: ServerCapabilities::builder().enable_tools().build(),
            ..Default::default()
        }
    }
}
