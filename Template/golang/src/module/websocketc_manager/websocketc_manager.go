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

	public.DBG_LOG("try to connect :", client.url)

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
	client.recv_msg = make(chan string, 10)
	client.send_msg = make(chan string, 10)

	return client.recv_msg, client.send_msg
}

func (client *WebSocket_Client_Manager) reconnect() {

	if client.underReconnect == false{
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
	}else{
		for client.underReconnect{
			public.Sleep(1000)

			public.DBG_LOG("wait reconnect")
		}
	}
}


func (client *WebSocket_Client_Manager) ReadMessages() {

	defer func(){
		if r := recover(); r != nil {
			public.DBG_ERR("ReadMessages err:", r)

			client.reconnect()

			go client.ReadMessages()
		}
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			public.DBG_ERR("Read error:", err)
			client.conn.Close()

			client.reconnect()
			
			public.Sleep(1000)
			continue
		}
		client.recv_msg <- string(message)
	}
}


func (client *WebSocket_Client_Manager) SendMessages() {

	var message string

	defer func(){
		if r := recover(); r != nil {
			public.DBG_ERR("SendMessages err:", r)

			client.reconnect()

			client.send_msg <- message
			
			go client.SendMessages()
		}
	}()

	for {

		message = <- client.send_msg
		
		err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			public.DBG_ERR("Write error:", err)
			client.conn.Close()

			client.reconnect()

		}
	}
}

func Init_WebSocket_Client(serverURL string, reconnectInterval_ms int) (chan string, chan string){

	client := &WebSocket_Client_Manager{
		url					: serverURL,
		reconnectInterval	: reconnectInterval_ms,
		underReconnect		: false,
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

