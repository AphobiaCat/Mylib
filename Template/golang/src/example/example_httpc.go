package example

import(
	"mylib/src/public"
	hm "mylib/src/module/http_manager"
)

type HTTP_Test_Data struct {
	Uid			string	`json:"uid"`
	Name		string	`json:"name"`
	Age			int 	`json:"age"`
}

func Example_HTTP(){

	default_headers := make(map[string]string)

	default_headers["auth"] = "hello"
	default_headers["auth2"] = "world"

	hm.Set_Default_Headers(default_headers)

	get_params := make(map[string]string)

	get_params["one"] = "hello world"

	ret := hm.Get("http://127.0.0.1:7001/get_test", get_params)
	public.DBG_LOG(ret)

	get_params["two"] = "hello world"
	ret = hm.Get("http://127.0.0.1:7001/get_test2", get_params)
	public.DBG_LOG(ret)

	ret = hm.Post("http://127.0.0.1:7001/post_test3", HTTP_Test_Data{Uid:"123", Name:"Dunty", Age:25})
	public.DBG_LOG(ret)


	ret = hm.Post("http://127.0.0.1:7001/post_test4", HTTP_Test_Data{Uid:"123", Name:"Dunty", Age:25})
	public.DBG_LOG(ret)

	
}
