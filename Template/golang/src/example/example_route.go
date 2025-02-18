package example

import(
	"mylib/src/module/route_manager"
	"mylib/src/public"
)

func auth_mid(headers map[string]string)bool{
	if headers["auth"] == "hello" && headers["auth2"] == "world"{
		return true
	}

	return false
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


func Example_Route(){

	route_manager.Route_Get("get_test", 	get_test, []string{"one", "two", "three"}, "get_test err")
	route_manager.Route_Get("get_test2", 	get_test, []string{"one", "two", "three"}, "get_test2 err", auth_mid, []string{"auth", "auth2"}, "mid auth err")
	route_manager.Route_Post("post_test3", 	post_test, "post_test err")
	route_manager.Route_Post("post_test4", 	post_test, "post_test err", auth_mid, []string{"auth", "auth2"}, "mid auth err")

	route_manager.Route_Get("get_test_1s_call", 	get_test, 60, []string{"one", "two", "three"}, "get_test err")
	route_manager.Route_Get("get_test2_2s_call", 	get_test, 30, []string{"one", "two", "three"}, "get_test2 err", auth_mid, []string{"auth", "auth2"}, "mid auth err")
	route_manager.Route_Post("post_test3_6s_call", 	post_test, 10, "post_test err")
	route_manager.Route_Post("post_test4_60s_call", post_test, 1, "post_test err", auth_mid, []string{"auth", "auth2"}, "mid auth err")
	
	route_manager.Init_Route("0.0.0.0:7001")
}


