package example

import(
	"mylib/src/public"
	"mylib/src/module/socket_manager"
)


func socket_test_client(client socket_manager.Socket_Client){

	for{
		select {
			case recv := <- client.Recv_msg:
				public.DBG_LOG("recv client msg:", recv)
	
				send := "recv client msg:" + recv
				
				client.Send_msg <- send
	
			case err := <- client.Err_msg:
				public.DBG_LOG("client error:", err)
				return
		}

	}
}

func Example_tcp_socket(){

	clients :=  socket_manager.Socket_TCP_Listen("8080")

	for {	
		client := <- clients

		go socket_test_client(client)
	}
}

func Example_udp_socket(){

	clients :=  socket_manager.Socket_UDP_Listen("8082")

	for {	
		client := <- clients

		go socket_test_client(client)
	}
}

func Example_udp_socket_config_timeout(){

	clients :=  socket_manager.Socket_UDP_Listen("8084", 5)

	for {	
		client := <- clients

		go socket_test_client(client)
	}
}

func Example_quic_socket(){
	clients :=  socket_manager.Socket_QUIC_Listen("8086", "./src/example/test_socket_credit/cert.pem", "./src/example/test_socket_credit/key.pem")

	for {	
		client := <- clients

		go socket_test_client(client)
	}
}

func Example_socket_server(){
	//Example_tcp_socket()
	//Example_udp_socket()
	//go Example_udp_socket_config_timeout()
	Example_quic_socket()
}

