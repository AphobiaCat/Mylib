package example

import(
	"mylib/src/public"
	ws_m "mylib/src/module/websocketc_manager"
)

func Example_Webscoketc(){
	recv_msg, send_msg := ws_m.Init_WebSocket_Client("ws://127.0.0.1:1234/ws", 1000)

	i := int64(1)
	
	for {

		public.DBG_LOG("send")
		send_msg <- ("hello world" + public.ConvertNumToHexStr(i))
		public.DBG_LOG("recv")
		recv := <- recv_msg
	
		public.DBG_LOG("recv msg:", recv)
		i++

		public.Sleep(10)
	}
}

