package socket_manager

import (
	"net"

	"mylib/src/public"
	"time"
)

type Socket_Client struct {
	Err_msg			chan string
	
	Recv_msg		chan string
	Send_msg		chan string

	Last_Op_Time	int64
}

func tcp_handle_recv(conn net.Conn, client Socket_Client) {
	defer conn.Close()

	addr := conn.RemoteAddr()
	
	public.DBG_LOG("tcp client connected:", addr)

	for{

		if len(client.Err_msg) != 0{
			public.DBG_ERR("socket client[", addr, "]")
			return
		}
	
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			public.DBG_ERR("socket client[", addr, "] Error reading:", err)
			err_msg := "read failed"
			client.Err_msg <- err_msg
			return
		}	
		client.Recv_msg <- string(buffer)

		client.Last_Op_Time = public.Now_Time_S()
	}	
}

func tcp_handle_send(conn net.Conn, client Socket_Client) {
	defer conn.Close()
	//public.DBG_LOG("Client connected:", conn.RemoteAddr())

	addr := conn.RemoteAddr()

	for {
		select {
			case err := <- client.Err_msg:
				public.DBG_ERR("socket client[", addr, "] close")
				client.Err_msg <- err
				return 
		
			case send_msg := <- client.Send_msg:
				_, err := conn.Write([]byte(send_msg))

				if err != nil{
					public.DBG_ERR("socket client[", addr, "] close")
					err_msg := "send failed"
					client.Err_msg <- err_msg
					return 
				}

				client.Last_Op_Time = public.Now_Time_S()
		}
	}
}


func tcp_listen(port string, client_channel chan Socket_Client){
	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		public.DBG_ERR("Error listening:", err)
		panic(err)
	}
	defer ln.Close()

	public.DBG_LOG("TCP server is listening on port ", port, "...")
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			public.DBG_ERR("Error accepting connection:", err)
			continue
		}

		tmp_client := Socket_Client{
			Recv_msg: make(chan string, 10),
			Send_msg: make(chan string, 10),
			Err_msg	: make(chan string, 2),
		}

		client_channel <- tmp_client

		go tcp_handle_recv(conn, tmp_client)
		go tcp_handle_send(conn, tmp_client)
	}
}

func Socket_TCP_Listen(port string) chan Socket_Client{
	ret_chan := make(chan Socket_Client, 10)
	
	go tcp_listen(port, ret_chan)

	return ret_chan
}

func udp_handle_send(conn *net.UDPConn, client *net.UDPAddr, udp_client Socket_Client, client_timeout int64) {

	addr := client.String()

	public.DBG_LOG("udp client connected:", addr)

	max_sleep_time := client_timeout
	
	timeoutDuration := time.Duration(max_sleep_time) * time.Second

	for {

		timeout := time.After(timeoutDuration)
	
		select {
			case err := <- udp_client.Err_msg:
				public.DBG_ERR("socket client[", addr, "] close")
				udp_client.Err_msg <- err
				return 
		
			case send_msg := <- udp_client.Send_msg:
				_, err := conn.WriteToUDP([]byte(send_msg), client)

				if err != nil{
					public.DBG_ERR("socket client[", addr, "] close")
					err_msg := "send failed"
					udp_client.Err_msg <- err_msg
					return 
				}
			case <-timeout:

				now_time		:= public.Now_Time_S()
				last_op_time	:= udp_client.Last_Op_Time

				if (now_time - last_op_time) > max_sleep_time{
					err := "addr[" + addr +  "] long time no op"
				
					public.DBG_LOG(err)
					udp_client.Err_msg <- err
				}
				
				return
		}
	}
}

func udp_listen(port string , client_channel chan Socket_Client, udp_timeout ...int64){

	var client_timeout int64 = 60

	if len(udp_timeout) >= 1{
		client_timeout = udp_timeout[0]
	}

	addr, err := net.ResolveUDPAddr("udp", ":" + port)
	if err != nil {
		public.DBG_ERR("Error resolving address:", err)
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	
	if err != nil {
		public.DBG_ERR("Error listening:", err)
		panic(err)
	}
	defer conn.Close()

	public.DBG_LOG("UDP server is listening on port ", port, "...")


	buffer := make([]byte, 1024)
	udp_client_map := make(map[string]Socket_Client)
	
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
	
		if err != nil {
			public.DBG_ERR("Error accepting connection:", err)
			continue
		}

		client_addr := addr.String()
				
		client, exist := udp_client_map[client_addr]

		if !exist{
			tmp_client := Socket_Client{
				Recv_msg: make(chan string, 10),
				Send_msg: make(chan string, 10),
				Err_msg : make(chan string, 2),
			}
			udp_client_map[client_addr] = tmp_client
			client = tmp_client
			client_channel <- tmp_client

			go udp_handle_send(conn, addr, client, client_timeout)
		}

		client.Recv_msg <- string(buffer[:n])
		client.Last_Op_Time = public.Now_Time_S()
	}
}

func Socket_UDP_Listen(port string, udp_timeout ...int64) chan Socket_Client{
	ret_chan := make(chan Socket_Client, 10)
	
	go udp_listen(port, ret_chan, udp_timeout...)

	return ret_chan
}

