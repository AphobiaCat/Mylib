package route_manager

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"mylib/src/public"
)


var route_manager Route_Manager

type Route_Manager struct{
	http_service	*gin.Engine
	have_init		bool
}

type Route_Post_Processer_Info struct{
	Post_process	Route_Post_Process

	Err_msg			string
}

type Route_Get_Processer_Info struct{
	Get_process		Route_Get_Process
	Get_params		[]string

	Err_msg			string
}

type Route_Mid_Processer_Info struct{
	Process			Route_Mid_Process

	Get_headers		[]string

	Err_msg			string
}

type Route_Post_Process func(string)(interface{}, bool)
type Route_Get_Process func(map[string]string)(interface{}, bool)
type Route_Mid_Process func(map[string]string)bool


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
			c.Abort() // 阻止请求继续传递
		}     
	}
}

func (rm *Route_Manager) Post(api_path string, processer_infos ...interface{}){

	rm.Init_Gin()

	if len(processer_infos) != 1 && len(processer_infos) != 2{
		public.DBG_ERR("process num error")
		return 
	}

	processer_info := processer_infos[0].(Route_Post_Processer_Info)

	post_process := func(context *gin.Context){

		body, err := context.GetRawData()

		if err != nil{
			public.DBG_ERR("input data no exist:", body)
		}

		ret, succ := processer_info.Post_process(string(body))

		if succ{
			context.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": ret,
			})
		}else{
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"error": processer_info.Err_msg,
			})
		}
	}


	if len(processer_infos) == 1{
		rm.http_service.POST(api_path, post_process)	
	}else if len(processer_infos) == 2{	//mid
		mid_processer_info := processer_infos[1].(Route_Mid_Processer_Info)
		rm.http_service.POST(api_path, Process_Route_Middleware_Module(mid_processer_info.Process, mid_processer_info.Get_headers, mid_processer_info.Err_msg), post_process)	
	}else{
		public.DBG_ERR("null process")
	}

	public.DBG_LOG("Post --> ", api_path)
}

func (rm *Route_Manager) Get(api_path string, processer_infos ...interface{}){

	rm.Init_Gin()

	if len(processer_infos) != 1 && len(processer_infos) != 2{
		public.DBG_ERR("process num error")
		return 
	}

	processer_info := processer_infos[0].(Route_Get_Processer_Info)

	get_process := func(context *gin.Context){

		params := make(map[string]string)

		for _, key_val := range processer_info.Get_params{
			if val, exists := context.GetQuery(key_val); exists {
				params[key_val] = val
			} else {
				public.DBG_ERR("key[", key_val, "] no exist")
			}
		}

		ret, succ := processer_info.Get_process(params)

		if succ{
			context.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": ret,
			})
		}else{
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"error": processer_info.Err_msg,
			})
		}
	}

	if len(processer_infos) == 1{
		rm.http_service.GET(api_path, get_process)	
	}else if len(processer_infos) == 2{	//mid
		mid_processer_info := processer_infos[1].(Route_Mid_Processer_Info)
		rm.http_service.GET(api_path, Process_Route_Middleware_Module(mid_processer_info.Process, mid_processer_info.Get_headers, mid_processer_info.Err_msg), get_process)	
	}else{
		public.DBG_ERR("param error")
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
	go route_manager.Init(bind_addr)
}