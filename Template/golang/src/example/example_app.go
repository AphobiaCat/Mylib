package example

import(
	"mylib/src/public"	
)


func example_get_config(params string)(string, bool){
	public.DBG_LOG(params)
	return `{"hello": "world", "age": 18}`, true
}


func Example_app(){
	public.Entry("--config", example_get_config)


	global_1, exist := public.Global[string]("hello")

	if exist{
		public.DBG_LOG("app global params ", global_1)
	}else{
		panic("global val no exist")
	}
	
	global_2, exist := public.Global[float64]("age")

	if exist{
		public.DBG_LOG("app global params ", global_2)
	}else{
		panic("global val no exist")
	}

	public.DBG_LOG("app config---> ", public.Config)
}

