package route_manager

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	
	"net/http"

	"mylib/src/public"

	redis "mylib/src/module/redis_manager"
)

var allow_origins = []string{"*"}
var allow_methods = []string{"*"} //[]string{"GET", "POST", "PUT", "DELETE"} 
var allow_headers = []string{"*"}

const mid_data_key = "MidData"

type postCallback	func(body string)(interface{}, bool)
type getCallback	func(params map[string]string)(interface{}, bool)
type mitCallback	func(params map[string]string)(map[string]string, bool)

type RouteManager struct{
	http_service	*gin.Engine

	routes			[]Route
	routes_len		uint32
}

type Route struct{
	api				string
	post_route		postCallback
	get_route		getCallback
	recv_params		[]string
	alert			string
	mid_callbacks	[]mitCallback
	mid_params		[][]string
	mid_alert		[]string
	mid_index		int
	req_limit		int
	reload_limit_s	int64
	need_user_ip	bool
}

func New() *RouteManager{
	ret := &RouteManager{}

	return ret
}

func (rm *RouteManager) Route_Post(api string, call_back postCallback) *Route{
	rm.routes	= append(rm.routes, Route{})

	ret			:= &(rm.routes[rm.routes_len])
	rm.routes_len += 1

	ret.api				= api
	ret.post_route		= call_back
	ret.mid_index		= -1
	ret.req_limit		= 0
	ret.reload_limit_s	= 60
	
	return ret
}

func (rm *RouteManager) Route_Get(api string, call_back getCallback) *Route{
	rm.routes = append(rm.routes, Route{})

	ret := &(rm.routes[rm.routes_len])
	rm.routes_len += 1
	
	ret.api				= api
	ret.get_route		= call_back
	ret.mid_index		= -1
	ret.req_limit		= 0
	ret.reload_limit_s	= 60

	return ret
}

func (r *Route) RecvParams(params ...string) *Route{
	r.recv_params = append(r.recv_params, params...)
	
	return r
}

func (r *Route) Alert(alert string) *Route{

	r.alert	= alert
	
	return r
}

func (r *Route) Middle(middle mitCallback) *Route{
	r.mid_callbacks	= append(r.mid_callbacks, middle)
	r.mid_index		+= 1
	r.mid_params	= append(r.mid_params, []string{})
	r.mid_alert		= append(r.mid_alert, "")

	return r
}

func (r *Route) MiddleParams(params ...string) *Route{
	if r.mid_index >= 0{
		r.mid_params[r.mid_index] = append(r.mid_params[r.mid_index], params...)
	}

	return r
}

func (r *Route) MiddleAlert(alert string) *Route{
	if r.mid_index >= 0{
		r.mid_alert[r.mid_index] = alert
	}

	return r
}

func (r *Route) ReqLimit(count int, reload_time ...int64) *Route{
	r.req_limit = count
	if len(reload_time) != 0{
		r.reload_limit_s = reload_time[0]
	}

	return r
}

func (r *Route) NeedUserIp() *Route{
	r.need_user_ip = true

	return r
}


func stream_control(api string, ip string, call_limit int, reload_time int64)bool{

	if call_limit == 0{
		return true
	}

	// public.DBG_LOG(api, " request:", ip)

	redis_key	:= "stream_control_" + api + "_" + ip
	count := redis.Timer_Count(redis_key, int64(call_limit), reload_time)

	if count >= 0{
		return true
	}else{
		return false
	}
}

func process_route_middleware_module(process mitCallback, need_header []string, err_info string) gin.HandlerFunc{
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
			c.JSON(http.StatusUnauthorized, gin.H{"code": -1, "error": err_info})
			c.Abort()
		}     
	}
}

func (rm *RouteManager) Init_Route(bind_addr string){
	gin.SetMode(gin.ReleaseMode)
	rm.http_service	= gin.New()
	
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = allow_origins  
	corsConfig.AllowMethods = allow_methods
	corsConfig.AllowHeaders = allow_headers 

	rm.http_service.Use(cors.New(corsConfig))

	rm.http_service.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.1"})	//only trust local proxy


	for _, route := range rm.routes{
		if route.get_route != nil{
			
			get_route_process := func(context *gin.Context){
		
				defer func(){
					if err := recover(); err != nil{
						public.DBG_ERR("err:", err)
					}
				}()
		
				clientIP := context.ClientIP()
		
				if !stream_control(route.api, clientIP, route.req_limit, route.reload_limit_s){
					context.JSON(http.StatusOK, gin.H{
						"code": -429,
						"error": "too many requests",
					})
					return
				}
		
				params := make(map[string]string)
		
				for _, key_val := range route.recv_params{
					if val, exists := context.GetQuery(key_val); exists {
						params[key_val] = val
					} else {
						public.DBG_ERR("key[", key_val, "] no exist")
					}
				}
		
				if route.need_user_ip {
					params["ip"] = clientIP
				}

				mid_params_i, _ := context.Get(mid_data_key)

				if mid_params, ok := mid_params_i.(map[string]string); ok{
					for key, val := range mid_params{
						params[key] = val
					}
				}

				ret, succ := route.get_route(params)
		
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
		
					public.DBG_ERR(route.api, " err:", route.alert)
				}
			}
		
		
			if len(route.mid_callbacks) > 0{
				mids_func := []gin.HandlerFunc{}
				
				for index, mid_process := range route.mid_callbacks{
					mids_func = append(mids_func, process_route_middleware_module(mid_process, route.mid_params[index], route.mid_alert[index]))
				}
				mids_func = append(mids_func, get_route_process)
		
				rm.http_service.GET(route.api, mids_func...) 
			}else{
				rm.http_service.GET(route.api, get_route_process)	
			}

			public.DBG_LOG("Get  --> ", route.api)
		}else if route.post_route != nil{
		
			post_route_process := func(context *gin.Context){

				defer func(){
					if err := recover(); err != nil{
						public.DBG_ERR("err:", err)
					}
				}()

				clientIP := context.ClientIP()

				if !stream_control(route.api, clientIP, route.req_limit, route.reload_limit_s){
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
				
				if route.need_user_ip || len(route.mid_callbacks) > 0{					
					var tmp_map map[string]interface{}
					public.Parser_Json(body_str, &tmp_map)	

					if route.need_user_ip {				
						tmp_map["ip"] = clientIP
					}

					mid_params_i, _ := context.Get(mid_data_key)

					if mid_params, ok := mid_params_i.(map[string]string); ok{
						for key, val := range mid_params{
							if key == "jwt"{
								tmp_map[key] = json.RawMessage(val)
							}else{
								tmp_map[key] = val
							}
						}
					}

					body_str = public.Build_Json(tmp_map)
				}

				ret, succ := route.post_route(body_str)
				
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

					public.DBG_ERR(route.api, " err:", route.alert)
				}
			}


			if len(route.mid_callbacks) > 0{
				mids_func := []gin.HandlerFunc{}
				
				for index, mid_process := range route.mid_callbacks{
					mids_func = append(mids_func, process_route_middleware_module(mid_process, route.mid_params[index], route.mid_alert[index]))
				}
				mids_func = append(mids_func, post_route_process)

				rm.http_service.POST(route.api, mids_func...)	
			}else{
				rm.http_service.POST(route.api, post_route_process)
			}

			public.DBG_LOG("Post --> ", route.api)
		}else{
			public.DBG_ERR("route no define.")
		}
	}

	public.DBG_LOG("bind addr :", bind_addr)
	if err := rm.http_service.Run(bind_addr); err != nil {
		panic(err)
	}
}

