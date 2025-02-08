package socket_manager


import (
	"net"
	"mylib/src/public"
	"os"
	"context"
	"github.com/quic-go/quic-go"
	"crypto/tls"	
	"crypto/x509"
)

const(
	socket_tcp = iota
	socket_udp
	socket_quic
)

type Socket_Client_Manager struct {
	url					string
	conn				net.Conn
	isConnected			bool
	reconnectInterval	int
	underReconnect		bool
	socket_type			int

	//quic use
	session				quic.Connection
	stream				quic.Stream
	ca_path				string

	recv_msg			chan string
	send_msg			chan string
}

func (client *Socket_Client_Manager) connect() error {

	public.DBG_LOG("try to connect :", client.url)

	var err error

	switch client.socket_type{
		case socket_tcp:
			client.conn, err = net.Dial("tcp", client.url)
			if err != nil {
				public.DBG_ERR("tcp connect failed: ", err)
				return err
			}
			
		case socket_udp:
			remoteAddr, err := net.ResolveUDPAddr("udp", client.url)
			if err != nil {
				public.DBG_ERR("address parser failed: ", err)
				return err
			}

			// create udp socket
			client.conn, err = net.DialUDP("udp", nil, remoteAddr)
			if err != nil {
				public.DBG_ERR("udp connect failed: ", err)
				return err
			}

		case socket_quic:

			tlsConf, err := generateClientTLSConfig(client.ca_path)

			if err != nil{
				public.DBG_ERR("connect quic failed")
				return err
			}

			client.session, err = quic.DialAddr(context.Background(), client.url, tlsConf, nil)
			if err != nil {
				public.DBG_ERR("quic connect failed: ", err)
				return err
			}

			client.stream, err = client.session.OpenStreamSync(context.Background())
			if err != nil {
				public.DBG_ERR("open quic stream failed: ", err)
				return err
			}
	}

	client.isConnected		= true
	client.underReconnect	= false
	public.DBG_LOG("Connected to socket server")

	return nil	
}

func (client *Socket_Client_Manager) close() {
	switch client.socket_type{
		case socket_tcp, socket_udp:
			client.conn.Close()
		case socket_quic:
			defer client.session.CloseWithError(0, "normal close")
	}
}


func (client *Socket_Client_Manager) init_chan() (chan string, chan string){
	client.recv_msg = make(chan string, 10)
	client.send_msg = make(chan string, 10)

	return client.recv_msg, client.send_msg
}

func (client *Socket_Client_Manager) reconnect() {

	if client.underReconnect == false{
		client.isConnected		= false
		client.underReconnect	= true
		for {
		
			public.DBG_LOG("Attempting to reconnect...")
			err := client.connect()
			if err == nil {
				return
			}
			public.DBG_ERR("Reconnect failed:", err)
			public.Sleep(client.reconnectInterval)
		}
	}else{
		for client.underReconnect{
			public.Sleep(1000)

			public.DBG_LOG("wait reconnect")
		}
	}
}


func (client *Socket_Client_Manager) read_messages() {

	defer func(){
		if r := recover(); r != nil {
			public.DBG_ERR("ReadMessages err:", r)

			client.reconnect()

			go client.read_messages()
		}
	}()

	buffer := make([]byte, 1024)

	switch client.socket_type{
		case socket_tcp, socket_udp:
			for {
				n, err := client.conn.Read(buffer)
				if err != nil {
					public.DBG_ERR("Read error:", err)
					client.close()

					client.reconnect()
					
					public.Sleep(1000)
					continue
				}
				client.recv_msg <- string(buffer[:n])
			}

		case socket_quic:
			for {

				n, err := client.stream.Read(buffer)
				if err != nil {
					public.DBG_ERR("Read error:", err)
					client.close()

					client.reconnect()
					
					public.Sleep(1000)
					continue
				}
				
				client.recv_msg <- string(buffer[:n])
			}
	}
}


func (client *Socket_Client_Manager) send_messages() {

	var message string

	defer func(){
		if r := recover(); r != nil {
			public.DBG_ERR("SendMessages err:", r)

			client.reconnect()

			client.send_msg <- message
			
			go client.send_messages()
		}
	}()


	switch client.socket_type{
		case socket_tcp, socket_udp:
			for {

				message = <- client.send_msg
				
				_, err := client.conn.Write([]byte(message))
				if err != nil {
					public.DBG_ERR("Write error:", err)
					client.close()

					client.reconnect()

				}
			}

		case socket_quic:

			for {

				message = <- client.send_msg
				
				_, err := client.stream.Write([]byte(message))
				if err != nil {
					public.DBG_ERR("Write error:", err)
					client.close()

					client.reconnect()

				}
			}
	}	
}

func generateClientTLSConfig(ca_path string)(*tls.Config, error){
	certPool := x509.NewCertPool()
	certData, err := os.ReadFile(ca_path)
	if err != nil {
		public.DBG_ERR("unable read CA redit: ", err)
		return nil, err
	}
	certPool.AppendCertsFromPEM(certData)

	return &tls.Config{RootCAs: certPool, NextProtos: []string{"quic-example"}}, nil
}

func socket_init_client(serverURL string, reconnectInterval_ms int, connect_type int, ca_path string)(chan string, chan string){
	client := &Socket_Client_Manager{
		url					: serverURL,
		reconnectInterval	: reconnectInterval_ms,
		underReconnect		: false,
		socket_type			: connect_type,
		ca_path				: ca_path,
	}

	err := client.connect()
	if err != nil {
		public.DBG_ERR("Initial connection failed: ", err)
	}

	recv_chan, send_chan := client.init_chan()

	go client.read_messages()
	go client.send_messages()

	return recv_chan, send_chan
}

func Socket_Init_TCP_Client(serverURL string, reconnectInterval_ms int) (chan string, chan string){
	return socket_init_client(serverURL, reconnectInterval_ms, socket_tcp, "")
}

func Socket_Init_UDP_Client(serverURL string, reconnectInterval_ms int) (chan string, chan string){
	return socket_init_client(serverURL, reconnectInterval_ms, socket_udp, "")
}

func Socket_Init_QUIC_Client(serverURL string, reconnectInterval_ms int, ca_path string) (chan string, chan string){
	return socket_init_client(serverURL, reconnectInterval_ms, socket_quic, ca_path)
}

