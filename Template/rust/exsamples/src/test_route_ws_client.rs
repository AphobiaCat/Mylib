use route_websocket_client::WsClient;


fn hello(code: i16, payload: String){
	println!("[hello] code = {}, payload = {}", code, payload);
}

fn world(code: i16, payload: String){
	println!("[world] code = {}, payload = {}", code, payload);
}

pub async fn test_route_ws_client() {
    let ws = WsClient::new(
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiZHVudHkiLCJleHAiOjE3NDc0MDM2Mzl9.L_EIUBzrFhHakplWBOZz2bttzIL4KWpbs23q3sXth2U".to_string(), 
		"ws://127.0.0.1:1234/".to_string()
	);

    ws.route_ws("hello", hello);

    ws.route_ws("world", world);

    ws.start_ws();

    ws.send("hello".to_string(), "hello world".to_string()).await;
}
