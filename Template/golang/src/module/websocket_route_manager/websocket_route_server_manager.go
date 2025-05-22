package websocket_route_manager


import(
	"mylib/src/public"
	route "mylib/src/module/route_manager"

	"sync"
	"net"
	"net/http"
    "github.com/gorilla/websocket"
    "time"

)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type ws_msg struct{
	Token	string	`json:"t"`
	Route	string	`json:"r"`
	Payload	string	`json:"p"`
}

type route_process struct{
	have_ret_process	func(string, string)(interface{}, bool)
	not_ret_process		func(string, string)
	have_ret			bool
}

var send_chan_map map[string]chan string
var send_chan_map_lock sync.Mutex

var ws_route_process map[string]route_process
var ws_route_process_lock sync.Mutex

var ws_route_exit func(string)

func close_client(close_chan chan bool){
	close_chan <- true
}

func close_user_send_chan(uid string){
	public.DBG_LOG("uid[", uid, "] disconnect")
	send_chan_map_lock.Lock()
	defer send_chan_map_lock.Unlock()

	delete(send_chan_map, uid)

	go ws_route_exit(uid)
}


func ws_route_handler(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        public.DBG_ERR(err)
        return
    }

    defer conn.Close()
    defer func() {
		if r := recover(); r != nil {
			public.DBG_ERR("Recovered error:", r)
		}
	}()

	send_msg_chan		:= make(chan string)
	close_client_chan	:= make(chan bool)

    go func() {
		ticker := time.NewTicker(30 * time.Second)
    	defer ticker.Stop()
    	
        for {
			select{
				case <- close_client_chan:
					return

				case msg := <- send_msg_chan:
					if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				    	public.DBG_ERR(err)
				    	return
					}
				case <-ticker.C:
					send_msg_chan <- "p"
			}
        }
    }()

    defer close_client(close_client_chan)

    have_init := false

	local_process_map := make(map[string]route_process)

    for {
		var recv_msg ws_msg

        _, msg, err := conn.ReadMessage()
        if err != nil {
        	public.DBG_ERR(err)
            return
        }
		public.Parser_Json(string(msg), &recv_msg)

		uid, succ := route.Route_Parser_Jwt(recv_msg.Token)

		if !succ{
			public.DBG_ERR("user token[", recv_msg.Token, "] error")
			continue
		}
		
		if !have_init{
			send_chan_map_lock.Lock()
			send_chan_map[uid] = send_msg_chan
			send_chan_map_lock.Unlock()

			public.DBG_LOG("uid[", uid, "] connect")
			defer close_user_send_chan(uid)

			have_init = true
		}
		
		if !have_init{
			continue
		}

		process, exist := local_process_map[recv_msg.Route]

		if !exist{
			ws_route_process_lock.Lock()
			process, exist = ws_route_process[recv_msg.Route]
			ws_route_process_lock.Unlock()

			if exist{
				local_process_map[recv_msg.Route] = process
			}
		}

		if exist{
			go func(){
				
defer func(){
					if err := recover(); err != nil{
						public.DBG_ERR("err:", err)
					}
				}()
			
				if process.have_ret{
					ret, succ := process.have_ret_process(uid, recv_msg.Payload)
				
					var ret_s struct{
						Code	int 	`json:"c"`
						Payload string	`json:"p"`
						Route	string	`json:"r"`
					}
		
					ret_s.Payload	= public.Build_Json(ret)
					ret_s.Route 	= recv_msg.Route
					
					if succ{
						ret_s.Code = 0
					}else{
						ret_s.Code = -1
					}
		
					send_msg_chan <- public.Build_Json(ret_s)
				}else{
					process.not_ret_process(uid, recv_msg.Payload)
				}
			}()
		}else{
			public.DBG_ERR("this route[", recv_msg.Route, "] no exist")
		}		
    }
}

func Route_WS(api string, call_back interface{})bool{
	var process route_process

	switch call_back.(type){
		case func(string, string)(interface{}, bool):
			process.have_ret_process	= call_back.(func(string, string)(interface{}, bool))
			process.have_ret			= true
		case func(string, string):
			process.not_ret_process		= call_back.(func(string, string))
			process.have_ret			= false
		default:
			return false
	}
	
	ws_route_process_lock.Lock()
	ws_route_process[api] = process
	ws_route_process_lock.Unlock()

	return true
}

func Route_WS_Exit(call_back func(string)){
	ws_route_exit = call_back
}

func WS_Send_Msg(uid string, data interface{}, user_route string)bool{
	send_chan_map_lock.Lock()
	send_chan, exist := send_chan_map[uid]
	send_chan_map_lock.Unlock()

	if exist{
		var ret_s struct{
			Code 	int		`json:"c"`
			Payload string	`json:"p"`
			Route	string	`json:"r"`
		}

		ret_s.Payload	= public.Build_Json(data)
		ret_s.Route		= user_route
	
		send_chan <- public.Build_Json(ret_s)
		return true
	}else
{
		return false
	}
}


func Init_Ws_Route(bind_addr string){
	listener, err := net.Listen("tcp4", bind_addr)
	if err != nil {
	    public.DBG_ERR("Failed to listen on port ", bind_addr, ":", err)
	}

	http.HandleFunc("/", ws_route_handler)

	public.DBG_LOG("Websocket Server started on :", bind_addr)
    
	ret := http.Serve(listener, nil)
	
	public.DBG_ERR(ret)
}


func init(){
	send_chan_map = make(map[string]chan string)
	ws_route_process = make(map[string]route_process)
}

