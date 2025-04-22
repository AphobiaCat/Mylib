use dynamic_code::DynamicCode;
use wasm_bindgen::prelude::*;

use wasm_websocket_client::WasmWsClient;
use wasm_thread_manager::WasmThreadManager;

#[wasm_bindgen]
pub fn run_dynamic_test_code() -> Result<f64, JsValue> {

    let script = r#"
        pub fn add(a , b ) {
            a + b
        }

        pub fn sum_matrix(matrix) {
            let sum = 0.0;
            for row in matrix {
                for col in row {
                    sum += col;
                }
            }
            sum
        }
    "#;

    let mut dynamic_code = DynamicCode::new(script).map_err(|e| JsValue::from_str(&format!("Compile error: {e}")))?;

    // let result:f64 = dynamic_code.use_func("add", (100244, 2048)).map_err(|e| JsValue::from_str(&format!("Call error: {e}")))?;
    // let result:f64 = dynamic_code.use_func_dyn("add", "[1, 2]").map_err(|e| JsValue::from_str(&format!("Call error: {e}")))?;

    let result:f64 = dynamic_code.use_func_dyn("sum_matrix", "[[[1.4, 2.1, 3.6], [4.0, 5.1, 6.2]]]").map_err(|e| JsValue::from_str(&format!("Call error: {e}")))?;


    Ok(result)  
}


#[wasm_bindgen]
pub fn run_ws_code() -> Result<i32, JsValue> {

    let client = WasmWsClient::new("ws://192.168.136.128:1234/ws");

    let client_for_cb = client.clone(); // clone Rc
    WasmWsClient::on_message(&client, move |msg| {
        web_sys::console::log_1(&format!("Rust recv msg: {}", msg).into());
        WasmWsClient::send(&client_for_cb, "Hello from Rust!");
    });

    Ok(1_i32)
}