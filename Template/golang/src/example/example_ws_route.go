package example

import(
	"mylib/src/public"
	ws_route "mylib/src/module/websocket_route_manager"	
	route "mylib/src/module/route_manager"
)


func test_ws_route_hello(uid string, payload string)(interface{}, bool){
	public.DBG_LOG("recv uid[", uid, "] payload[", payload, "]")

	return "succ", true
}

func test_ws_route_world(uid string, payload string)(interface{}, bool){
	public.DBG_LOG("recv uid[", uid, "] payload[", payload, "] 2")

	ws_route.WS_Send_Msg(uid, "hello world test", "hello/world")

	return "failed", false
}

func test_ws_exit(uid string){
	public.DBG_ERR(uid, " logout")
}

func Example_Ws_Route(){

	test_jwt_str, _ := route.Route_Generate_Jwt_By_Str("dunty", 3600)

	public.DBG_LOG(test_jwt_str)

	ws_route.Route_WS("hello", test_ws_route_hello)
	ws_route.Route_WS("world", test_ws_route_world)
	ws_route.Route_WS_Exit(test_ws_exit)
	
	ws_route.Init_Ws_Route("0.0.0.0:1234")
}

