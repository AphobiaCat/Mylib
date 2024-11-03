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

	public.Parser_Jason(body_json, &login_data)

	public.DBG_LOG(login_data)

	return login_data, true
}


func Example_Route(){

	get_process := route_manager.Route_Get_Processer_Info{Get_process:   get_test	, Get_params: []string{"one", "two", "three"}	, Err_msg: "get_test err"}
	post_process:= route_manager.Route_Post_Processer_Info{Post_process: post_test	, Err_msg: "post_test err"}
	mid_process := route_manager.Route_Mid_Processer_Info{Process: auth_mid			, Get_headers: []string{"auth", "auth2"}		, Err_msg: "mid auth err"}

	route_manager.Route_Get("get_test", 	 get_process)
	route_manager.Route_Get("get_test2", 	 get_process, mid_process)
	route_manager.Route_Post("post_test3", post_process)
	route_manager.Route_Post("post_test4", post_process, mid_process)
	
	route_manager.Init_Route("0.0.0.0:7001")
}


