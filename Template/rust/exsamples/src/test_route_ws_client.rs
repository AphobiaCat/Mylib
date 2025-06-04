use route_websocket_client::WsClient;
use public::DBG_LOG;

async fn hello(code: i16, payload: String){
	DBG_LOG!("[hello] code[", code, "] payload[", payload, "]");
}

async fn world(code: i16, payload: String){
	DBG_LOG!("[hello] code[", code, "] payload[", payload, "]");
}

async fn big_payload(code: i16, payload: String, big_payload: String){
	DBG_LOG!("[hello] code[", code, "] payload[", payload, "] big_payload[", big_payload, "]");
}

pub async fn test_route_ws_client() {
    let ws = WsClient::new(
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiZHVudHkiLCJleHAiOjE3NDkwNDIyMjN9.Yjqoay7piBjxW16vYiJqL0IGtt_ownmFOTBk2OKNAkc".to_string(), 
		"ws://127.0.0.1:1234/".to_string()
	);

    ws.route_ws("hello", hello);

    ws.route_ws("world", world);
	ws.route_ws_big_payload("big-payload", big_payload);

	ws.start_ws();

    ws.send("hello".to_string(), "hello world".to_string(), "".to_string()).await;
	ws.send("big-payload".to_string(), "hello world".to_string(), "hello big payload".to_string()).await;
}
