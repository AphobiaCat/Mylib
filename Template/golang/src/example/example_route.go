package example

import(
	"mylib/src/module/route_manager"
	"mylib/src/public"
)

func auth_mid(headers map[string]string)(map[string]string, bool){
	
	if headers["auth"] == "hello" && headers["auth2"] == "world"{
		ret := make(map[string]string)

		ret["auth"] = "hello_finish"
		ret["auth2"] = "world_finish"
		
		return ret, true
	}

	return map[string]string{}, false
}

func get_test(params map[string]string)(interface{}, bool){

	ret := 0
	
	if params["one"] == "hello world"{
		ret += 1
	}

	if params["two"] == "hello world"{
		ret += 2
	}

	if params["three"] == "hello world"{
		ret += 4
	}

	if ret != 0{
		return ret, true
	}else{
		return ret, false
	}
}

func get_test_need_ip(params map[string]string)(interface{}, bool){

	public.DBG_LOG("ip[", params["ip"], "]")
	return params, true
}

func post_test_need_ip(body string)(interface{}, bool){

	var body_s struct{
		IP	string	`json:"ip"`
	}

	public.Parser_Json(body, &body_s)

	public.DBG_LOG("ip[", body_s.IP, "]")
	return body_s, true
}


func get_test_recv_mid(params map[string]string)(interface{}, bool){

	public.DBG_LOG(params)
	return params, true
}

func post_test(body_json string)(interface{}, bool){

	var login_data struct {
		Uid			string	`json:"uid"`
		Name		string	`json:"name"`
		Age			int 	`json:"age"`
	}

	public.Parser_Json(body_json, &login_data)

	public.DBG_LOG(login_data)

	return login_data, true
}

func post_test_recv_mid(body_json string)(interface{}, bool){
	return body_json, true
}

func test_jwt(){
	var test_data struct{
		Name	string	`json:"name"`
		Age		int		`json:"age"`
	}

	test_data.Name	= "dunty"
	test_data.Age	= 25
	
	jwt_str, succ := route_manager.Route_Generate_Jwt(test_data, 5)

	
	public.Sleep(4500)

	if succ{
		public.DBG_LOG("jwt:", jwt_str)

		data_str, succ := route_manager.Route_Parser_Jwt(jwt_str)

		data := test_data
		data.Name	= ""
		data.Age	= 0
		
		public.Parser_Json(data_str, &data)

		if succ{
			public.DBG_LOG("result data is:", data)
		}else{
			public.DBG_ERR("parser error")
		}
		
	}else{
		public.DBG_ERR("generate jwt error")
	}

	public.Sleep(500)

	_, succ = route_manager.Route_Parser_Jwt(jwt_str)

	if succ{
		public.DBG_ERR("parser need timeout")
	}else{
		public.DBG_LOG("timeout succ")
	}
}

func show_test_jwt(){
	var test_data struct{
		Name	string	`json:"name"`
		Age		int		`json:"age"`
	}

	test_data.Name	= "dunty"
	test_data.Age	= 25
	
	jwt_str, _ := route_manager.Route_Generate_Jwt(test_data, 30)

	public.DBG_LOG(jwt_str)
}

func Example_Route(){

	//test_jwt()
	show_test_jwt()

	route := route_manager.New()

	route.Route_Get("get_test", get_test).RecvParams("one", "two", "three").Alert("get_test err")
	route.Route_Get("get_test2", get_test).RecvParams("one", "two", "three").Alert("get_test2 err")
	route.Route_Get("get_test3", get_test_recv_mid).RecvParams("one", "two", "three").Alert("get_test3 err").Middle(auth_mid).MiddleParams("auth", "auth2").MiddleAlert("mid auth err")
	route.Route_Get("get_jwt_test4", get_test_recv_mid).RecvParams("one", "two", "three").Alert("get_test4 err").Middle(auth_mid).MiddleParams("auth", "auth2").MiddleAlert("mid auth err").Middle(route_manager.Route_Get_Jwt_Mid).MiddleParams("token").MiddleAlert("jwt auth err")
	
	route.Route_Post("post_test1", post_test).Alert("post_test err")
	route.Route_Post("post_test2", post_test_recv_mid).Alert("post_test2 err").Middle(auth_mid).MiddleParams("auth", "auth2").MiddleAlert("mid auth err")
	route.Route_Post("post_jwt_test3", post_test_recv_mid).Alert("post_test3 err").Middle(auth_mid).MiddleParams("auth", "auth2").MiddleAlert("mid auth err").Middle(route_manager.Route_Get_Jwt_Mid).MiddleParams("token").MiddleAlert("jwt auth err")

	route.Route_Get("get_test_1s_call", get_test).RecvParams("one", "two", "three").Alert("get_test err").ReqLimit(1, 1)
	route.Route_Get("get_test_need_ip", get_test_need_ip).Alert("get_test_need_ip err").NeedUserIp()
	route.Route_Post("post_test_need_ip", post_test_need_ip).Alert("post_test_need_ip err").NeedUserIp()

	route.Init_Route("0.0.0.0:7001")
}


