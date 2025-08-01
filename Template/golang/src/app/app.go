//go:build !prod
// +build !prod

package main

import (
	"mylib/src/public"
	"mylib/src/example"
	"mylib/src/app/app_a"
)


// ----------------Global Parameter>


// ----------------Function>

func APP_Entry() {	

	//example.Example_app()
	//example.Example_Redis_Manager()
	//example.Example_Json_Op()
	//example.Example_Timer_Manager()
	//example.Example_2_Gorm()
	//example.Example_Route()
	//example.Example_Cachesql()
	//example.Example_HTTP()
	//example.Example_HTTP2()
	//example.Example_Webscoketc()
	//example.Example_Webscokets()
	//example.Example_Ws_Route()
	//example.Example_file()
	//example.Example_bitnum()
	example.Example_socket_server()
	//example.Example_socket_client()
	//example.Example_thread_map()
	//example.Example_OAuth2()
	//example.Example_Exec()
	//example.Example_msg_queue()

	app_a.Test_App_A()
	
	for {
		
		public.DBG_LOG("hello world")

		public.Sleep(1000)
	}
}

