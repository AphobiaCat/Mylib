package example

import(
	"mylib/src/public"
	ws_m "mylib/src/module/websocketc_manager"
)

func Example_Webscoketc(){
	recv_msg, send_msg := ws_m.WebsocketC_Init("ws://127.0.0.1:8800", 1000)
	recv_msg2, send_msg2 := ws_m.WebsocketC_Init("ws://127.0.0.1:8800", 1000)

	i := uint32(1)
	
	for {
		send_msg <- ("hello world" + public.ConvertUint32ToHexString(i))
		recv := <- recv_msg

		send_msg2 <- ("hello world2" + public.ConvertUint32ToHexString(i))
		recv2 := <- recv_msg2
	
		public.DBG_LOG("recv msg:", recv)
		public.DBG_LOG("recv2 msg:", recv2)
		i++

		public.Sleep(10)
	}
}

