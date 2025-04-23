package websockets_manager

import(
	"mylib/src/public"

	"log"
	"net"
	"net/http"
    "github.com/gorilla/websocket"

)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}


type WS_Client_Process func(recv_msg_chan chan string, send_msg_chan chan string, close_client chan bool)
var global_ws_client_process WS_Client_Process


func ws_handler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

	recv_msg_chan	:= make(chan string, 3)
	send_msg_chan	:= make(chan string)
	close_client	:= make(chan bool, 2)

    defer conn.Close()

	go func() {
        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
            	public.DBG_ERR(err)
                close_client <- true
                close_client <- true
                return
            }
			recv_msg_chan <- string(msg)
        }
    }()

    go func() {
        for {

			select{
				case <- close_client:
					conn.Close()
					return

				case msg := <- send_msg_chan:
					if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				    	public.DBG_ERR(err)
				    	conn.Close()
					}
			}
        }
    }()
    
	global_ws_client_process(recv_msg_chan, send_msg_chan, close_client)
}


func Init_Websocket_Server(bind_ip_port string, client_process WS_Client_Process) {
	listener, err := net.Listen("tcp4", bind_ip_port)
	if err != nil {
	    public.DBG_ERR("Failed to listen on port ", bind_ip_port, ":", err)
	}

	global_ws_client_process = client_process

	http.HandleFunc("/ws", ws_handler)

	public.DBG_LOG("Websocket Server started on :", bind_ip_port)
    
	ret := http.Serve(listener, nil)
	
    //log.Fatal(http.ListenAndServe("0.0.0.0:" + data_service_port, nil))
	public.DBG_ERR(ret)
}

