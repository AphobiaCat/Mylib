use wasm_bindgen::prelude::*;
use wasm_bindgen::JsCast;
use web_sys::{WebSocket, MessageEvent, ErrorEvent, CloseEvent, console};
use std::rc::Rc;
use std::cell::RefCell;

type MessageCallback = Box<dyn Fn(String)>;

pub struct WasmWsClient {
    url: String,
    ws: Option<WebSocket>,
    on_message: Option<MessageCallback>,
}

impl WasmWsClient {
    pub fn new(url: &str) -> Rc<RefCell<Self>> {
        let client = Rc::new(RefCell::new(Self {
            url: url.to_string(),
            ws: None,
            on_message: None,
        }));

        WasmWsClient::connect(&client);
        client
    }

    pub fn send(this: &Rc<RefCell<Self>>, msg: &str) {
        let client = this.borrow();
        if let Some(ws) = &client.ws {
            if ws.ready_state() == WebSocket::OPEN {
                ws.send_with_str(msg).unwrap_or_else(|e| {
                    console::error_1(&format!("Send failed: {:?}", e).into());
                });
            }
        }
    }

    pub fn on_message<F: Fn(String) + 'static>(this: &Rc<RefCell<Self>>, callback: F) {
        this.borrow_mut().on_message = Some(Box::new(callback));
    }

    fn connect(this: &Rc<RefCell<Self>>) {
        let url = this.borrow().url.clone();
        let ws = WebSocket::new(&url).expect("Failed to create WebSocket");
        ws.set_binary_type(web_sys::BinaryType::Arraybuffer);
        this.borrow_mut().ws = Some(ws.clone());

        // onmessage
        let on_msg_cb = Rc::clone(this);
        let onmessage = Closure::wrap(Box::new(move |e: MessageEvent| {
            if let Ok(text) = e.data().dyn_into::<js_sys::JsString>() {
                if let Some(cb) = &on_msg_cb.borrow().on_message {
                    cb(text.as_string().unwrap_or_default());
                }
            }
        }) as Box<dyn FnMut(_)>);
        ws.set_onmessage(Some(onmessage.as_ref().unchecked_ref()));
        onmessage.forget();

        // onclose
        let onclose_cb = Rc::clone(this);
        let onclose = Closure::wrap(Box::new(move |_e: CloseEvent| {
            console::log_1(&"WebSocket closed. Reconnecting in 3s...".into());
            let retry_cb = Rc::clone(&onclose_cb);
            let f = Closure::wrap(Box::new(move || {
                WasmWsClient::connect(&retry_cb);
            }) as Box<dyn Fn()>);

            web_sys::window().unwrap()
                .set_timeout_with_callback_and_timeout_and_arguments_0(f.as_ref().unchecked_ref(), 3000)
                .unwrap();
            f.forget();
        }) as Box<dyn FnMut(_)>);
        ws.set_onclose(Some(onclose.as_ref().unchecked_ref()));
        onclose.forget();

        // onerror
        let onerror = Closure::wrap(Box::new(move |_e: ErrorEvent| {
            console::error_1(&"WebSocket error".into());
        }) as Box<dyn FnMut(_)>);
        ws.set_onerror(Some(onerror.as_ref().unchecked_ref()));
        onerror.forget();
    }
}