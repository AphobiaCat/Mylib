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
	Token			string	`json:"t"`
	Route			string	`json:"r"`
	Payload			string	`json:"p"`
}

type route_process struct{
	have_ret_process				func(string, string)(interface{}, bool)
	not_ret_process					func(string, string)
	big_payload_have_ret_process	func(string, string, string)(interface{}, bool)
	big_payload_not_ret_process		func(string, string, string)
	
	have_ret			bool
	have_big_payload	bool
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

	send_msg_chan		:= make(chan string, 8)
	close_client_chan	:= make(chan bool, 2)

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
					if len(send_msg_chan) == 0 {
				        select {
				        case send_msg_chan <- "p":
				        default:
				            // avoid block.
				        }
				    }
			}
        }
    }()

    defer close_client(close_client_chan)

    have_init := false

	local_process_map := make(map[string]route_process)

    for {
		
        _, msg, err := conn.ReadMessage()
        if err != nil {
        	public.DBG_ERR(err)
            return
        }

		origin_msg	:= string(msg)

		header_len_str	:= origin_msg[0:4]

		header_len	:= uint32(public.ConvertHEXStrToNum(header_len_str))
		header		:= origin_msg[4:(4 + header_len)]

		public.DBG_LOG("len[", header_len, "] header[", header, "]")
		
		var recv_msg ws_msg
		public.Parser_Json(header, &recv_msg)

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

				if process.have_big_payload{
					if process.have_ret{
						ret, succ := process.big_payload_have_ret_process(uid, recv_msg.Payload, origin_msg[(4 + header_len):])
						
						if succ{
							send_msg_chan <- build_ret_msg(0, recv_msg.Route, public.Build_Json(ret))
						}else{
							send_msg_chan <- build_ret_msg(-1, recv_msg.Route, public.Build_Json(ret))
						}
					}else{
						process.big_payload_not_ret_process(uid, recv_msg.Payload, origin_msg[(4 + header_len):])
					}
				}else{
					if process.have_ret{
						ret, succ := process.have_ret_process(uid, recv_msg.Payload)
					
						if succ{
							send_msg_chan <- build_ret_msg(0, recv_msg.Route, public.Build_Json(ret))
						}else{
							send_msg_chan <- build_ret_msg(-1, recv_msg.Route, public.Build_Json(ret))
						}
					}else{
						process.not_ret_process(uid, recv_msg.Payload)
					}
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
			process.have_ret_process				= call_back.(func(string, string)(interface{}, bool))
			process.have_ret						= true
		case func(string, string):
			process.not_ret_process					= call_back.(func(string, string))
			process.have_ret						= false
		case func(string, string, string)(interface{}, bool):
			process.big_payload_have_ret_process	= call_back.(func(string, string, string)(interface{}, bool))
			process.have_ret						= true
			process.have_big_payload				= true
		case func(string, string, string):
			process.big_payload_not_ret_process 	= call_back.(func(string, string, string))
			process.have_ret						= false
			process.have_big_payload				= true
			
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


func build_ret_msg(code int, user_route string, data interface{}, big_payload_option ...interface{})string{

	var ret_s struct{
		Code 			int		`json:"c"`
		Route			string	`json:"r"`
		Payload 		string	`json:"p"`
	}
	
	ret_s.Code		= code
	ret_s.Route		= user_route
	ret_s.Payload	= public.Build_Json(data)

	ret_str := public.Build_Json(ret_s)

	ret_str_len := public.ConvertNumToHexStr(int64(len(ret_str)))[2:]
	
	if len(ret_str_len) != 4{
		fill_zero := "0000"
		ret_str_len = fill_zero[4 - len(ret_str_len):] + ret_str_len
	}

	ret_msg := ret_str_len + ret_str

	if len(big_payload_option) != 0{
		if big_payload, ok := big_payload_option[0].(string); ok{
			ret_msg += big_payload
		}
	}
	return ret_msg
}

func WS_Send_Msg(uid string, user_route string, data interface{}, big_payload_option ...interface{})bool{
	send_chan_map_lock.Lock()
	send_chan, exist := send_chan_map[uid]
	send_chan_map_lock.Unlock()

	if exist{
		send_chan <- build_ret_msg(0, user_route, data, big_payload_option...)
		return true
	}else{
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
	send_chan_map		= make(map[string]chan string)
	ws_route_process	= make(map[string]route_process)
}

