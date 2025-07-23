use serde::Serialize;

pub struct HttpClient<'a> {
	url: &'a str,
	headers: Vec<(&'a str, &'a str)>,
	params: Vec<(&'a str, &'a str)>,
	body: Option<String>,
	is_get: bool,
}

impl<'a> HttpClient<'a> {
	pub fn get(url: &'a str) -> Self {
		Self {
			url,
			headers: Vec::new(),
			params: Vec::new(),
			body: None,
			is_get: true,
		}
	}

	pub fn post(url: &'a str) -> Self {
		Self {
			url,
			headers: Vec::new(),
			params: Vec::new(),
			body: None,
			is_get: false,
		}
	}

	pub fn header(mut self, k: &'a str, v: &'a str) -> Self {
		self.headers.push((k, v));
		self
	}

	pub fn param(mut self, k: &'a str, v: &'a str) -> Self {
		self.params.push((k, v));
		self
	}

	pub fn body<T: Serialize>(mut self, data: &T) -> Self {
		match serde_json::to_string(data) {
			Ok(json) => self.body = Some(json),
			Err(_) => self.body = None,
		}
		self
	}

	pub async fn send(self) -> Result<String, String> {
		if self.is_get{
			http_get(self.url, &self.headers, &self.params).await
		}else{
			http_post(self.url, &self.headers, &self.params, self.body).await
		}
	}
}

pub async fn http_get(
	base_url: &str,
	headers: &[(&str, &str)],
	params: &[(&str, &str)],
) -> Result<String, String> {
	#[cfg(target_arch = "wasm32")]
	{
		use gloo_net::http::Request;

		let query_string = if !params.is_empty() {
			let joined: String = params
				.iter()
				.map(|(k, v)| format!("{}={}", k, v))
				.collect::<Vec<_>>()
				.join("&");
			format!("{}?{}", base_url, joined)
		} else {
			base_url.to_string()
		};

		let mut req = Request::get(&query_string);

		for (k, v) in headers {
			req = req.header(k, v);
		}

		let res = req.send().await.map_err(|e| e.to_string())?;
		let text = res.text().await.map_err(|e| e.to_string())?;
		Ok(text)
	}

	#[cfg(not(target_arch = "wasm32"))]
	{
		use reqwest::Client;

		let client = Client::new();
		let mut req = client.get(base_url);

		for (k, v) in headers {
			req = req.header(*k, *v);
		}

		req = req.query(&params);

		let res = req.send().await.map_err(|e| e.to_string())?;
		let text = res.text().await.map_err(|e| e.to_string())?;
		Ok(text)
	}
}


pub async fn http_post(
	base_url: &str,
	headers: &[(&str, &str)],
	params: &[(&str, &str)],
	body: Option<String>,
) -> Result<String, String> {
	#[cfg(target_arch = "wasm32")]
	{
		use gloo_net::http::Request;

		let query_string = if !params.is_empty() {
			let joined: String = params
				.iter()
				.map(|(k, v)| format!("{}={}", k, v))
				.collect::<Vec<_>>()
				.join("&");
			format!("{}?{}", base_url, joined)
		} else {
			base_url.to_string()
		};

		let mut req = Request::post(&query_string);

		for (k, v) in headers {
			req = req.header(k, v);
		}

		if let Some(json) = body {
			req = req.header("Content-Type", "application/json");
			let res = req.body(json).send().await.map_err(|e| e.to_string())?;
			let text = res.text().await.map_err(|e| e.to_string())?;
			Ok(text)
		} else {
			let res = req.send().await.map_err(|e| e.to_string())?;
			let text = res.text().await.map_err(|e| e.to_string())?;
			Ok(text)
		}
	}

	#[cfg(not(target_arch = "wasm32"))]
	{
		use reqwest::Client;

		let client = Client::new();
		let mut req = client.post(base_url).query(&params);

		for (k, v) in headers {
			req = req.header(*k, *v);
		}

		let req = if let Some(json) = body {
			req.header("Content-Type", "application/json").body(json)
		} else {
			req
		};

		let res = req.send().await.map_err(|e| e.to_string())?;
		let text = res.text().await.map_err(|e| e.to_string())?;
		Ok(text)
	}
}