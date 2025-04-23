package example

import(
	"mylib/src/public"
	ws_m "mylib/src/module/websockets_manager"
)

func test_ws_client(recv_msg_chan chan string, send_msg_chan chan string, close_client chan bool){

	public.DBG_LOG("client")

	for{
		select{
			case msg:= <- recv_msg_chan:
				public.DBG_LOG("recv:", msg)
				send_msg_chan <- msg
	
			case <- close_client:
				public.DBG_LOG("client close")
				return
		}
	}

	close_client <- true
}

func Example_Webscokets(){
	ws_m.Init_Websocket_Server("0.0.0.0:1234", test_ws_client)
}

