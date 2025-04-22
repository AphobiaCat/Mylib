##how to use this lib

let client = Rc::new(RefCell::new(WasmWsClient::new("wss://echo.websocket.events")));  

let client_clone = Rc::clone(&client);  
client.borrow_mut().on_message(move |msg| {  
    web_sys::console::log_1(&format!("Rust recv msg: {}", msg).into());  
 
    // also send msg here  
    client_clone.borrow().send("Rust : recv");  
});  

// manual send msg  
client.borrow().send("Hello from Rust!");  