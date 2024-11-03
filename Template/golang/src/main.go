package main

import (
	"time"
	"mylib/src/example"
	//"mylib/src/module/cachesql_manager"
	//"mylib/src/module/gorm_manager"
	//"mylib/src/module/http_manager"
	//"mylib/src/module/redis_manager"
	//"mylib/src/module/route_manager"
	//"mylib/src/module/timer_manager"	
)



// ----------------Global Parameter>


// ----------------Function>


func main() {

	example.Example_Redis_Manager()
	example.Example_Json_Op()
	example.Example_Timer_Manager()
	example.Example_2_Gorm()
	example.Example_Route()
	example.Example_Cachesql()	

	for {
		//DBG_LOG("hello wolrd")
		time.Sleep(1 * time.Second)
	}
	
}
