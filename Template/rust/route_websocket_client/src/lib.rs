use tokio_tungstenite::tungstenite::protocol::Message;
use tokio_tungstenite::connect_async;
use tokio::sync::mpsc::{self, Sender};
use futures_util::{SinkExt, StreamExt};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use serde::{Deserialize, Serialize};

type RouteCallback = Arc<dyn Fn(i16, String) + Send + Sync>;

#[derive(Clone)]
pub struct WsClient {
    inner: Arc<Mutex<WsClientInner>>,
}

struct WsClientInner {
    uid: String,
    url: String,
    routes: HashMap<String, RouteCallback>,
    tx: Sender<String>,
}

#[derive(Serialize)]
struct WsRequest {
    t: String,
    r: String,
    p: String,
}

#[derive(Deserialize)]
struct WsResponse {
    c: i16,
    p: String,
    r: String,
}

impl WsClient {
    pub fn new(uid: String, url: String) -> Self {
        let (tx, _rx) = mpsc::channel::<String>(100);
        let inner = WsClientInner {
            uid,
            url,
            routes: HashMap::new(),
            tx,
        };
        Self {
            inner: Arc::new(Mutex::new(inner)),
        }
    }

    pub fn route_ws<F>(&self, api: &str, callback: F)
    where
        F: Fn(i16, String) + Send + Sync + 'static,
    {
        let mut inner = self.inner.lock().unwrap();
        inner.routes.insert(api.to_string(), Arc::new(callback));
    }

    pub fn start_ws(&self) {
        let inner = self.inner.clone();
        let (msg_tx, mut msg_rx) = mpsc::channel::<String>(100);

        {
            let mut locked = inner.lock().unwrap();
            locked.tx = msg_tx.clone();
        }

        tokio::spawn(async move {
            let url = {
                let locked = inner.lock().unwrap();
                locked.url.clone()
            };

            let (mut ws_stream, _) = connect_async(&url)
                .await
                .expect("Failed to connect");

            loop {
                tokio::select! {
                    Some(Ok(msg)) = ws_stream.next() => {
                        if let Message::Text(text) = msg {
                            if let Ok(parsed) = serde_json::from_str::<WsResponse>(&text) {
                                let cb_opt = {
                                    let locked = inner.lock().unwrap();
                                    locked.routes.get(&parsed.r).cloned()
                                };
                                if let Some(cb) = cb_opt {
                                    cb(parsed.c, parsed.p);
                                }
                            }
                        }
                    }
                    Some(msg) = msg_rx.recv() => {
                        ws_stream.send(Message::Text(msg.into())).await.unwrap();
                    }
                    else => break,
                }
            }
        });
    }

    pub async fn send(&self, route: String, payload: String) {
        let (msg, tx) = {
            let locked = self.inner.lock().unwrap();
            let req = WsRequest {
                t: locked.uid.clone(),
                r: route,
                p: payload,
            };
            let msg = serde_json::to_string(&req).unwrap();
            (msg, locked.tx.clone())
        };

        let _ = tx.send(msg).await;
    }
}
