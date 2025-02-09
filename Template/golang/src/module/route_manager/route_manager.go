package route_manager

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"mylib/src/public"

	cache "mylib/src/module/cachesql_manager"
)


var route_manager Route_Manager

type Route_Manager struct{
	http_service	*gin.Engine
	have_init		bool
}

type Route_Post_Process func(string)(interface{}, bool)
type Route_Get_Process func(map[string]string)(interface{}, bool)
type Route_Mid_Process func(map[string]string)bool

const stream_restart_count 	int64 = 20
const stream_restart_time 	int64 = 600
const stream_restart_time_2	int64 = stream_restart_time * 2

type Stream_Control_Struct struct{
	Count	int64	`json:"count"`
}

func stream_control(ip string)bool{
	public.DBG_LOG("request:", ip)

	redis_restart	:= "lock_restart_" + ip
	redis_now		:= "lock_now_" + ip

	ret := cache.Get_Cache(redis_restart, func()interface{}{
		return 1
	}, stream_restart_time_2, stream_restart_time_2, stream_restart_time)

	if ret == "1"{
		cache.Set_Cache(redis_restart, 0, stream_restart_time)
		cache.Set_Cache(redis_now, Stream_Control_Struct{Count: stream_restart_count}, stream_restart_time)
	}

	ret = cache.Get_Cache(redis_now, func()interface{}{
		return Stream_Control_Struct{Count: 0}
	}, stream_restart_time_2, stream_restart_time_2, stream_restart_time)

	var now Stream_Control_Struct

	public.Parser_Json(ret, &now)


	public.DBG_LOG("--->:", now)
	
	if now.Count > 0 {
		now.Count -= 1
		cache.Set_Cache(redis_now, now, stream_restart_time)
		return true
	}else{
		public.DBG_LOG("ip[", ip, "] request to much")
		return false
	}
}

func Process_Route_Middleware_Module(process Route_Mid_Process, need_header []string, err_info string) gin.HandlerFunc{
	return func(c *gin.Context) {

		use_header_array := make(map[string]string)

		for _, val := range need_header{
			use_header_array[val] = c.GetHeader(val)
		}

		ret := process(use_header_array)

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

	var post_process	Route_Post_Process
	var post_err_msg	string
		
	var mid_process		Route_Mid_Process
	var mid_headers		[]string
	var mid_err_msg		string
	
	who_params := ""

	for _, val := range processer_infos{

		switch val.(type){
			case func(string)(interface{}, bool):
				post_process = val.(func(string)(interface{}, bool))
				who_params = "post"

			case func(map[string]string)bool:
				mid_process = val.(func(map[string]string)bool)		
				who_params = "mid"

			case []string:
				if who_params == "post"{
					
				}else if who_params == "mid"{
					mid_headers = val.([]string)
				}

			case string:
				if who_params == "post"{
					post_err_msg = val.(string)
				}else if who_params == "mid"{
					mid_err_msg = val.(string)
				}
		}

	}

	if post_process == nil{
		public.DBG_ERR("post process no exist")
		return 
	}

	post_route_process := func(context *gin.Context){

		clientIP := context.ClientIP()

		if !stream_control(clientIP){
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

		ret, succ := post_process(string(body))

		if succ{
			context.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": ret,
			})
		}else{
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"error": post_err_msg,
			})
		}
	}


	if mid_process != nil{
		rm.http_service.POST(api_path, Process_Route_Middleware_Module(mid_process, mid_headers, mid_err_msg), post_route_process)	
	}else{
		rm.http_service.POST(api_path, post_route_process)
	}

	public.DBG_LOG("Post --> ", api_path)
}

func (rm *Route_Manager) Get(api_path string, processer_infos ...interface{}){

	rm.Init_Gin()

	var get_process Route_Get_Process
	var get_params	[]string
	var get_err_msg string
		
	var mid_process Route_Mid_Process
	var mid_headers	[]string
	var mid_err_msg string
	
	
	who_params := ""

	for _, val := range processer_infos{

		switch val.(type){
			case func(map[string]string)(interface{}, bool):
				get_process = val.(func(map[string]string)(interface{}, bool))
				who_params = "get"

			case func(map[string]string)bool:
				mid_process = val.(func(map[string]string)bool)	
				who_params = "mid"

			case []string:
				if who_params == "get"{
					get_params = val.([]string)
				}else if who_params == "mid"{
					mid_headers = val.([]string)
				}

			case string:
				if who_params == "get"{
					get_err_msg = val.(string)
				}else if who_params == "mid"{
					mid_err_msg = val.(string)
				}

			default:
				public.DBG_ERR("err info:", val)
		}

	}

	if get_process == nil{
		public.DBG_ERR("get process no exist")
		return 
	}

	get_route_process := func(context *gin.Context){

		clientIP := context.ClientIP()

		if !stream_control(clientIP){
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

		ret, succ := get_process(params)

		if succ{
			context.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": ret,
			})
		}else{
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"error": get_err_msg,
			})
		}
	}


	if mid_process != nil{
		rm.http_service.GET(api_path, Process_Route_Middleware_Module(mid_process, mid_headers, mid_err_msg), get_route_process)	
	}else{
		rm.http_service.GET(api_path, get_route_process)	
	}


	public.DBG_LOG("Get  --> ", api_path)
}

func (rm *Route_Manager) Init_Gin(){
	if !rm.have_init{
		gin.SetMode(gin.ReleaseMode)
		rm.http_service	= gin.New()
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
