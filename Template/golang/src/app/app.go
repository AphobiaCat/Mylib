package main

import (
	"mylib/src/example"
	"mylib/src/public"
	//"mylib/src/module/cachesql_manager"
	//"mylib/src/module/gorm_manager"
	//"mylib/src/module/http_manager"
	//"mylib/src/module/redis_manager"
	//"mylib/src/module/route_manager"
	//"mylib/src/module/timer_manager"
	//"mylib/src/module/http_manager"
	//"mylib/src/module/websocketc_manager"
	//"mylib/src/module/websockets_manager"
	//"mylib/src/module/file_manager"
	//"mylib/src/module/bignum_manager"
	//"mylib/src/module/socket_manager"
)


// ----------------Global Parameter>


// ----------------Function>

func APP_Entry() {

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
	//example.Example_file()
	//example.Example_bitnum()
	//example.Example_socket_server()
	//example.Example_socket_client()
	example.Example_thread_map()

	public.DBG_LOG("app")

	for {
		
		//public.DBG_LOG("hello world")

		public.Sleep(1000)
	}
}
