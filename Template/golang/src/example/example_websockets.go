package example

import(
	"mylib/src/public"
	ws_m "mylib/src/module/websockets_manager"
)

func test(client *ws_m.Ws_Client){

	public.DBG_LOG("client:", client)

	ready_return := false

	for {
		msgs, ret := client.Recv_Msg()

		public.DBG_LOG(ret)

		if ret{
			ready_return = true
			for _, val := range msgs{
				public.DBG_LOG("recv msg :", val)
				client.Send_Msg(val)
			}
		}else if ready_return{
			client.Close()
			return 
		}

		if !client.IsAlive(){
			return 
		}
		
		public.Sleep(1000)
	}
}

func Example_Webscokets(){
	ws_chan := ws_m.Init_Websocket_Server("0.0.0.0:1234")

	for{
		client := <- ws_chan

		go test(client)
	}
}

