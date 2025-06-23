package example

import(
	"mylib/src/public"
	"mylib/src/module/socket_manager"
)


func Example_tcp_socket_client(){

	recv_msg, send_msg := socket_manager.Socket_Init_TCP_Client("127.0.0.1:8080", 1000)

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

func Example_udp_socket_client(){

	recv_msg, send_msg := socket_manager.Socket_Init_UDP_Client("127.0.0.1:8082", 1000)

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

func Example_quic_socket_client(){
	recv_msg, send_msg := socket_manager.Socket_Init_QUIC_Client("127.0.0.1:8086", 1000, "./src/example/test_socket_credit/ca.pem")

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

func Example_socket_client(){
	//Example_tcp_socket_client()
	//Example_udp_socket_client()
	Example_quic_socket_client()
}

