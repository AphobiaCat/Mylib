package main

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

	Parser_Jason(body_json, &login_data)

	DBG_LOG(login_data)

	return login_data, true
}


func example_route(){

	get_process := Route_Get_Processer_Info{get_process:   get_test, get_params: []string{"one", "two", "three"}, err_msg: "get_test err"}
	post_process:= Route_Post_Processer_Info{post_process: post_test, err_msg: "post_test err"}
	mid_process := Route_Mid_Processer_Info{process: auth_mid, get_headers: []string{"auth", "auth2"}, err_msg: "mid auth err"}

	Route_Get("get_test", 	 get_process)
	Route_Get("get_test2", 	 get_process, mid_process)
	Route_Post("post_test3", post_process)
	Route_Post("post_test4", post_process, mid_process)
	
	Init_Route("0.0.0.0:7001")
}


