package websocketc_manager

import (
	"net/http"
	"github.com/gorilla/websocket"

	"mylib/src/public"
)

var ws_clients_map map[string] WebSocket_Client_Manager = make(map[string] WebSocket_Client_Manager)

type WebSocket_Client_Manager struct {
	url					string
	conn				*websocket.Conn
	isConnected			bool
	reconnectInterval	int
	underReconnect		bool

	recv_msg			chan string
	send_msg			chan string
}

func (client *WebSocket_Client_Manager) connect() error {
	header := http.Header{}
	// can add auth msg
	// header.Set("Authorization", "Bearer your_token_here")

	public.DBG_LOG("try to connect ws", client.url)

	var err error
	client.conn, _, err = websocket.DefaultDialer.Dial(client.url, header)
	if err != nil {
		return err
	}

	client.isConnected		= true
	client.underReconnect	= false
	public.DBG_LOG("Connected to WebSocket server")

	return nil
}

func (client *WebSocket_Client_Manager) init_chan() (chan string, chan string){
	client.recv_msg = make(chan string)
	client.send_msg = make(chan string)

	return client.recv_msg, client.send_msg
}

func (client *WebSocket_Client_Manager) reconnect() {
	client.isConnected		= false
	client.underReconnect	= true
	for {
	
		public.DBG_LOG("Attempting to reconnect...")
		err := client.connect()
		if err == nil {
			return
		}
		public.DBG_ERR("Reconnect failed:", err)
		public.Sleep(client.reconnectInterval)
		
	}
}


func (client *WebSocket_Client_Manager) ReadMessages() {
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			public.DBG_ERR("Read error:", err)
			client.conn.Close()

			if !client.underReconnect{
				client.reconnect()
			}			
			
			public.Sleep(1000)
			continue
		}
		client.recv_msg <- string(message)
	}
}


func (client *WebSocket_Client_Manager) SendMessages() {
	for {
		message := <- client.send_msg

		for !client.isConnected {
			public.DBG_ERR("not connected to WebSocket server")
			public.Sleep(1000)
		}
		
		err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			public.DBG_ERR("Write error:", err)
			client.conn.Close()

			if !client.underReconnect{
				client.reconnect()
			}
		}
	}
}

func WebsocketC_Init(serverURL string, reconnectInterval_ms int) (chan string, chan string){

	client := &WebSocket_Client_Manager{
		url					: serverURL,
		reconnectInterval	: reconnectInterval_ms,
		underReconnect		: true,
	}

	err := client.connect()
	if err != nil {
		public.DBG_ERR("Initial connection failed: %v", err)
	}

	recv_chan, send_chan := client.init_chan()

	go client.ReadMessages()
	go client.SendMessages()

	return recv_chan, send_chan
}

