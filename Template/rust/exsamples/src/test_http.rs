use std::error::Error;
use http_manager::{HttpClient};
use public::{DBG_LOG, DBG_ERR};
use serde::Serialize;

pub async fn test_http_manager() -> Result<(), Box<dyn Error>> {

	match HttpClient::get("https://google.com")
		.header("token", "a")
		.param("hello", "world")
		.param("name", "dunty")
		.send()
		.await {
		Ok(_) => DBG_LOG!("✅ Get Succ"),
		Err(err) => DBG_ERR!("❌ Get Error: ", err),
	}

	#[derive(Serialize)]
	struct MyData {
		name: String,
		age: u32,
	}

	let body = MyData {
		name: "Dunty".into(),
		age: 26,
	};

	match HttpClient::post("https://google.com")
		.header("token", "a")
		.body(&body)
		.send()
		.await {
		Ok(_) => DBG_LOG!("✅ Post succ"),
		Err(err) => DBG_ERR!("❌ Post Error: ", err),
	}

	Ok(())
}