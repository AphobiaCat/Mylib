package app

import(
	"os"
	"sync"

	"mylib/src/public"
	file "mylib/src/module/file_manager"
)


var global_app_init_params map[string]string
var global_app_params map[string]interface{}
var global_map_lock sync.Mutex

func Entry(function string, entry_interface ...interface{}){

	global_map_lock.Lock()
	defer global_map_lock.Unlock()
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

func Set_Global(key string, val interface{}){
	global_map_lock.Lock()
	defer global_map_lock.Unlock()
	global_app_params[key] = val
}

func Global[T any](key string)(T, bool){

	//public.DBG_LOG("global_app_params:", global_app_params)

	global_map_lock.Lock()
	defer global_map_lock.Unlock()
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

	
	Entry("--config", func(params string)(string, bool){
		config_json := file.File_Read(params)

		if len(config_json) != 0{

			//public.DBG_LOG("config: ", config_json)
		
			return config_json, true
		}else{
			return "file[" + params + "] no exist", false
		}
	})
}

