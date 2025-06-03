package app

import(
	"os"
	"mylib/src/public"
)


var global_app_init_params map[string]string
var global_app_params map[string]interface{}

func Entry(function string, entry_interface ...interface{}){

	params, exist := global_app_init_params[function]

	if !exist || len(entry_interface) == 0{
		return
		
	}

	for _, entry := range entry_interface{
		switch entry.(type){
			case func(string)(string, bool):
				ret, succ := entry.(func(string)(string, bool))(params)
				if succ{
					tmp_map := make(map[string]interface{})

					public.Parser_Json(ret, &tmp_map)

					for key, val := range tmp_map{
						global_app_params[key] = val
					}

				}else{
					panic(ret)
				}
			case func(string)bool:
				succ := entry.(func(string)bool)(params)
				if !succ{
					panic("entry run failed")
				}
			case func(string):
				entry.(func(string))(params)
			default:
				panic(`entry must be 
1: func(string)(string, bool)
2: func(string)bool
3: func(string)`)
		}
	}
}

func Global[T any](key string)(T, bool){

	//public.DBG_LOG("global_app_params:", global_app_params)

	if val, exist := global_app_params[key].(T); exist{
		return val, true
	}else{
		var zero T
		return zero, false
	}
}

func init(){
	global_app_init_params	= make(map[string]string)
	global_app_params		= make(map[string]interface{})

	args_list := []string{}

	for _, val := range os.Args{
		args_list = append(args_list, val)
	}

	args_list = args_list[1:]
	
	function_index := 0
	params_index := 1

	for ; params_index < len(args_list); {
		global_app_init_params[args_list[function_index]] = args_list[params_index]

		function_index += 2
		params_index += 2
	}
}

