package route_manager

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	
	"net/http"

	"mylib/src/public"

	redis "mylib/src/module/redis_manager"
)


var route_manager Route_Manager

type Route_Manager struct{
	http_service	*gin.Engine
	have_init		bool
}

type mid_array struct{
	mid_process 	Route_Mid_Process
	mid_headers		[]string
	mid_err_msg 	string
}
	

type Route_Post_Process func(string)(interface{}, bool)
type Route_Get_Process func(map[string]string)(interface{}, bool)
type Route_Post_Recv_Mid_Process func(string, ...map[string]string)(interface{}, bool)
type Route_Get_Recv_Mid_Process func(map[string]string, ...map[string]string)(interface{}, bool)
type Route_Mid_Process func(map[string]string)(map[string]string, bool)

var allow_origins = []string{"*"}
var allow_methods = []string{"*"} //[]string{"GET", "POST", "PUT", "DELETE"} 
var allow_headers = []string{"*"}

const mid_data_key = "MidData"

func stream_control(api string, ip string, call_per_minute_rate int)bool{

	if call_per_minute_rate == 0{
		return true
	}

	// public.DBG_LOG(api, " request:", ip)

	redis_key	:= "stream_control_" + api + "_" + ip
	count := redis.Timer_Count(redis_key, int64(call_per_minute_rate), 60)

	if count >= 0{
		return true
	}else{
		return false
	}
}

func Process_Route_Middleware_Module(process Route_Mid_Process, need_header []string, err_info string) gin.HandlerFunc{
	return func(c *gin.Context) {

		use_header_array := make(map[string]string)

		for _, val := range need_header{
			use_header_array[val] = c.GetHeader(val)
		}

		user_data, ret := process(use_header_array)

		if len(user_data) != 0{
			user_info_interface, exist := c.Get(mid_data_key)

			var new_user_info map[string]string

			if exist{
				new_user_info = user_info_interface.(map[string]string)				
			}else{
				new_user_info = make(map[string]string)
			}

			for key, val := range user_data{
				new_user_info[key] = val
			}

			c.Set(mid_data_key, new_user_info)
		}

		if ret{
			c.Next()

		}else{
			c.JSON(http.StatusUnauthorized, gin.H{"error": err_info})
			c.Abort()
		}     
	}
}

func (rm *Route_Manager) Post(api_path string, processer_infos ...interface{}){

	rm.Init_Gin()

	var post_process			Route_Post_Process
	var post_process_recv_mid	Route_Post_Recv_Mid_Process
	
	var post_err_msg			string
	var post_count_per_min		int

	var need_ip bool

	mids_index := -1
	mids := []mid_array{}

	who_params := ""

	for _, val := range processer_infos{

		switch val.(type){
			case func(string)(interface{}, bool):
				post_process = val.(func(string)(interface{}, bool))
				who_params = "post"

			case func(string, ...map[string]string)(interface{}, bool):
				post_process_recv_mid = val.(func(string, ...map[string]string)(interface{}, bool))
				who_params = "post"

			case func(map[string]string)(map[string]string, bool):
				mids_index += 1

				mids = append(mids, mid_array{mid_process: val.(func(map[string]string)(map[string]string, bool))})

				who_params = "mid"

			case []string:
				if who_params == "post"{
					
				}else if who_params == "mid"{
					mids[mids_index].mid_headers = val.([]string)
				}

			case string:
				if who_params == "post"{
					post_err_msg = val.(string)
				}else if who_params == "mid"{
					mids[mids_index].mid_err_msg = val.(string)
				}

			case int:
				post_count_per_min = val.(int)

			case bool:
				need_ip = val.(bool)

			default:
				public.DBG_ERR("err info:", val)
		}

	}

	if post_process == nil && post_process_recv_mid == nil{
		public.DBG_ERR("post process no exist")
		return 
	}

	post_route_process := func(context *gin.Context){

		defer func(){
			if err := recover(); err != nil{
				public.DBG_ERR("err:", err)
			}
		}()

		clientIP := context.ClientIP()

		if !stream_control(api_path, clientIP, post_count_per_min){
			context.JSON(http.StatusOK, gin.H{
				"code": -429,
				"error": "too many requests",
			})
			return
		}

		body, err := context.GetRawData()

		if err != nil{
			public.DBG_ERR("input data no exist:", body)
		}

		body_str := string(body)

		if need_ip {
			var tmp_map map[string]interface{}

			public.Parser_Json(body_str, &tmp_map)

			tmp_map["ip"] = clientIP

			body_str = public.Build_Json(tmp_map)
		}

		var ret interface{}
		var succ bool

		if post_process != nil {
			ret, succ = post_process(body_str)
		} else {
			mid_params, exist := context.Get(mid_data_key)

			if exist {
				ret, succ = post_process_recv_mid(body_str, mid_params.(map[string]string))
			} else {
				ret, succ = post_process_recv_mid(body_str)
			}
		}
			

		if succ{
			context.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": ret,
			})
		}else{
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"error": ret,
			})

			public.DBG_ERR(post_err_msg)
		}
	}


	if mids_index != -1{
		mids_func := []gin.HandlerFunc{}
		
		for _, val := range mids{
			mids_func = append(mids_func, Process_Route_Middleware_Module(val.mid_process, val.mid_headers, val.mid_err_msg))
		}
		mids_func = append(mids_func, post_route_process)

		rm.http_service.POST(api_path, mids_func...)	
	}else{
		rm.http_service.POST(api_path, post_route_process)
	}

	public.DBG_LOG("Post --> ", api_path)
}

func (rm *Route_Manager) Get(api_path string, processer_infos ...interface{}){

	rm.Init_Gin()

	var get_process 			Route_Get_Process
	var get_process_recv_mid	Route_Get_Recv_Mid_Process
	var get_params				[]string
	var get_err_msg 			string
	var get_count_per_min		int

	var need_ip bool

	mids_index := -1
	mids := []mid_array{}
	
	who_params := ""

	for _, val := range processer_infos{

		switch val.(type){
			case func(map[string]string)(interface{}, bool):
				get_process = val.(func(map[string]string)(interface{}, bool))
				who_params = "get"

			case func(map[string]string, ...map[string]string)(interface{}, bool):
				get_process_recv_mid = val.(func(map[string]string, ...map[string]string)(interface{}, bool))
				who_params = "get"

			case func(map[string]string)(map[string]string, bool):
				mids_index += 1
				
				mids = append(mids, mid_array{mid_process: val.(func(map[string]string)(map[string]string, bool))})
				
				who_params = "mid"

			case []string:
				if who_params == "get"{
					get_params = val.([]string)
				}else if who_params == "mid"{
					mids[mids_index].mid_headers = val.([]string)
				}

			case string:
				if who_params == "get"{
					get_err_msg = val.(string)
				}else if who_params == "mid"{
					mids[mids_index].mid_err_msg = val.(string)
				}

			case int:
				get_count_per_min = val.(int)

			case bool:
				need_ip = val.(bool)

			default:
				public.DBG_ERR("err info:", val)
		}

	}

	if get_process == nil && get_process_recv_mid == nil{
		public.DBG_ERR("get process no exist")
		return 
	}

	get_route_process := func(context *gin.Context){

		defer func(){
			if err := recover(); err != nil{
				public.DBG_ERR("err:", err)
			}
		}()

		clientIP := context.ClientIP()

		if !stream_control(api_path, clientIP, get_count_per_min){
			context.JSON(http.StatusOK, gin.H{
				"code": -429,
				"error": "too many requests",
			})
			return
		}

		params := make(map[string]string)

		for _, key_val := range get_params{
			if val, exists := context.GetQuery(key_val); exists {
				params[key_val] = val
			} else {
				public.DBG_ERR("key[", key_val, "] no exist")
			}
		}

		if need_ip {
			params["ip"] = clientIP
		}

		var ret interface{}
		var succ bool

		if get_process != nil{
			ret, succ = get_process(params)
		}else{
			mid_params, exist := context.Get(mid_data_key)

			if exist{
				ret, succ = get_process_recv_mid(params, mid_params.(map[string]string))
			}else{
				ret, succ = get_process_recv_mid(params)
			}
		}

		if succ{
			context.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": ret,
			})
		}else{
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"error": ret,
			})

			public.DBG_ERR(get_err_msg)
		}
	}


	if mids_index != -1{
		mids_func := []gin.HandlerFunc{}
		
		for _, val := range mids{
			mids_func = append(mids_func, Process_Route_Middleware_Module(val.mid_process, val.mid_headers, val.mid_err_msg))
		}
		mids_func = append(mids_func, get_route_process)

		rm.http_service.GET(api_path, mids_func...)	
	}else{
		rm.http_service.GET(api_path, get_route_process)	
	}


	public.DBG_LOG("Get  --> ", api_path)
}

func (rm *Route_Manager) Init_Gin(){
	if !rm.have_init{
		gin.SetMode(gin.ReleaseMode)
		rm.http_service	= gin.New()
		
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = allow_origins  
		corsConfig.AllowMethods = allow_methods
		corsConfig.AllowHeaders = allow_headers 

		rm.http_service.Use(cors.New(corsConfig))

		rm.http_service.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.1"})	//only trust local proxy
		
		rm.have_init	= true
	}
}

func (rm *Route_Manager) Init(bind_addr string){

	if !rm.have_init{
		public.DBG_ERR("haven`t init")
		return 
	}

	public.DBG_LOG("bind addr :", bind_addr)
	if err := rm.http_service.Run(bind_addr); err != nil {
		panic(err)
	}
}

func Route_Post(api_path string, processer_infos ...interface{}){
	route_manager.Post(api_path, processer_infos...)
}

func Route_Get(api_path string, processer_infos ...interface{}){
	route_manager.Get(api_path, processer_infos...)
}

func Init_Route(bind_addr string){
	route_manager.Init(bind_addr)
}

