package websockets_manager

import(
	"mylib/src/public"

	"log"
	"net"
	"net/http"
    "github.com/gorilla/websocket"

)

var ws_chan chan *Ws_Client

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Ws_Client struct{
	conn		*websocket.Conn
	alive		bool
	have_msg	bool
	msg			chan string
}

func (wsc *Ws_Client) Recv_Msg() ([]string, bool){

	if wsc.alive && wsc.have_msg{

		ret := []string{}

		wsc.have_msg = false

		for i := 0; i < len(wsc.msg); i++{
			ret = append(ret, <- wsc.msg)
		}
		
		return ret, true
	}
	return []string{}, false
}

func (wsc *Ws_Client) Send_Msg(msg string) bool{

	if wsc.alive{
		if err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
	    	public.DBG_ERR(err)
	    	wsc.Close()
	    	return false
		}
		return true
	}

	return false
}

func (wsc *Ws_Client) IsAlive()bool{
	return wsc.alive
}

func (wsc *Ws_Client) Close(){
	wsc.alive = false
	wsc.conn.Close()
}


func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    
	ws_client := Ws_Client{
		conn	: conn,
		alive	: true,
		have_msg: false,
		msg		: make(chan string, 10),
	}
	
	defer ws_client.Close()

    ws_chan <- &ws_client

    for {
    	//message type
        _, message, err := conn.ReadMessage()
        if err != nil {
            public.DBG_ERR(err)

            ws_client.have_msg	= false
            
            return
        }

		ws_client.have_msg	= true
		ws_client.msg <- string(message)	
    }
}

func init_websocket_server(bind_ip_port string){
	listener, err := net.Listen("tcp4", bind_ip_port)
	if err != nil {
	    public.DBG_ERR("Failed to listen on port ", bind_ip_port, ":", err)
	}

	http.HandleFunc("/ws", wsHandler)

	public.DBG_LOG("Websocket Server started on :", bind_ip_port)
    
	ret := http.Serve(listener, nil)
	
    //log.Fatal(http.ListenAndServe("0.0.0.0:" + data_service_port, nil))
	public.DBG_ERR(ret)
}

func Init_Websocket_Server(bind_ip_port string) chan *Ws_Client{

	ws_chan = make(chan *Ws_Client)
	
	go init_websocket_server(bind_ip_port)

	return ws_chan
}

