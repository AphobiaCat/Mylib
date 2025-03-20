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

func get_test_recv_mid(params map[string]string, mid ...map[string]string)(interface{}, bool){

	for _, val := range mid{
		public.DBG_LOG("mid value:", val)
	}

	if len(mid) != 0{
		ret := mid[0]["auth"] + mid[0]["auth2"]
		return ret, true
	}
	return "null", false
}


type Login_Data struct {
	Uid			string	`json:"uid"`
	Name		string	`json:"name"`
	Age			int 	`json:"age"`
}

func post_test(body_json string)(interface{}, bool){

	var login_data Login_Data

	public.Parser_Json(body_json, &login_data)

	public.DBG_LOG(login_data)

	return login_data, true
}

func post_test_recv_mid(body_json string, mid ...map[string]string)(interface{}, bool){

	if len(mid) != 0{
		ret := mid[0]["auth"] + mid[0]["auth2"]
		return ret, true
	}
	return "null", false
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

	route_manager.Route_Get("get_test", 	get_test, []string{"one", "two", "three"}, "get_test err")
	route_manager.Route_Get("get_test2", 	get_test, []string{"one", "two", "three"}, "get_test2 err", auth_mid, []string{"auth", "auth2"}, "mid auth err")
	route_manager.Route_Get("get_test3",	get_test_recv_mid, []string{"one", "two", "three"}, "get_test2 err", auth_mid, []string{"auth", "auth2"}, "mid auth err")	
	route_manager.Route_Get("get_jwt_test4",get_test_recv_mid, []string{"one", "two", "three"}, "get_test2 err", auth_mid, []string{"auth", "auth2"}, "mid auth err", route_manager.Route_Get_Jwt_Mid, []string{"token"}, "jwt auth err")

	route_manager.Route_Post("post_test1", 	post_test, "post_test err")
	route_manager.Route_Post("post_test2", 	post_test, "post_test err", auth_mid, []string{"auth", "auth2"}, "mid auth err")	
	route_manager.Route_Post("post_test3",	post_test_recv_mid	, "post_test err", auth_mid, []string{"auth", "auth2"}, "mid auth err")
	route_manager.Route_Post("post_jwt_test4",post_test_recv_mid, "post_test err", auth_mid, []string{"auth", "auth2"}, "mid auth err", route_manager.Route_Get_Jwt_Mid, []string{"token"}, "jwt auth err")

	route_manager.Route_Get("get_test_1s_call", 	get_test, 60, []string{"one", "two", "three"}, "get_test err")
	route_manager.Route_Get("get_test2_2s_call", 	get_test, 30, []string{"one", "two", "three"}, "get_test2 err", auth_mid, []string{"auth", "auth2"}, "mid auth err")
	route_manager.Route_Post("post_test3_6s_call", 	post_test, 10, "post_test err")
	route_manager.Route_Post("post_test4_60s_call", post_test, 1, "post_test err", auth_mid, []string{"auth", "auth2"}, "mid auth err")
	
	route_manager.Init_Route("0.0.0.0:7001")
}


