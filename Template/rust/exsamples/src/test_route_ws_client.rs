use route_websocket_client::WsClient;
use public::DBG_LOG;

async fn hello(code: i16, payload: String){
	DBG_LOG!("[hello] code[", code, "] payload[", payload, "]");
}

async fn world(code: i16, payload: String){
	DBG_LOG!("[hello] code[", code, "] payload[", payload, "]");
}


pub async fn test_route_ws_client() {
    let ws = WsClient::new(
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiZHVudHkiLCJleHAiOjE3NDc4MTI0MDJ9.NuJWhrF7SIrI-X4mjjEx2ZP_GkfoLzsK9gWv6U66GOU".to_string(), 
		"ws://127.0.0.1:1234/".to_string()
	);

    ws.route_ws("hello", hello);

    ws.route_ws("world", world);

	ws.start_ws();

    ws.send("hello".to_string(), "hello world".to_string()).await;

}
