package cachesql_manager

import (
	"context"
    "github.com/redis/go-redis/v9"
    "time"

    "mylib/src/public"
    "mylib/src/module/app"
)	

var cache_sql_manager Cache_Sql_Manager

type Cache_Sql_Manager struct {
	rdb 	*redis.Client
	ctx 	context.Context
}

type Standard_CSM_Cache struct {
	D 		string	`json:"d"` //user data
	L 		int64	`json:"l"` //last update time
	LW 		int64	`json:"lw"` //last work time
	W 		bool	`json:"w"` //is working update status
	Wait 	bool	`json:"wait"` //wait frist
}

type New_Cache_Func func() interface{}

func (csm *Cache_Sql_Manager) Init(server_ip string, password string, DB int) {

	csm.ctx = context.Background()

	csm.rdb = redis.NewClient(&redis.Options{
		Addr: server_ip,
		Password: password,
		DB: DB,
	})

	_, err := csm.rdb.Ping(csm.ctx).Result()

	if err != nil {
		public.DBG_ERR("unable connet Redis:", err)
	}else{
		public.DBG_LOG("connect redis server succ")
	}
}

func (csm *Cache_Sql_Manager) Set_Cache(key string, value interface{}, config_time ...int64) {
	max_alive_time 		:= int64(1000 * 1000 * 1000 * 60 * 2)

	if len(config_time) == 1{
		max_alive_time	= config_time[0] * 1000 * 1000 * 1000
	}

	now_time := public.Now_Time_S()

	var new_cache_data Standard_CSM_Cache
	
	new_cache_data.D	= public.Build_Json(value)
	new_cache_data.L	= now_time + max_alive_time	//keep no update by get cache default force update time
	new_cache_data.LW	= 0
	new_cache_data.W	= false
	new_cache_data.Wait	= false

	err := csm.rdb.Set(csm.ctx, key, public.Build_Json(new_cache_data), time.Duration(max_alive_time)).Err()
	if err != nil {
		public.DBG_ERR("set value failed", err)
	}
}

func (csm *Cache_Sql_Manager) Get_Cache(key string, new_cache_func New_Cache_Func, config_time ...int64) string {
	force_update_time 	:= int64(60 * 2)
	max_work_time 		:= int64(60 * 5)
	max_alive_time 		:= int64(1000 * 1000 * 1000 * 60 * 10)	//ns -> us -> ms -> s

	switch len(config_time) {
	case 0:
		//default config
	case 1:
		force_update_time	= config_time[0]
	case 2:
		force_update_time	= config_time[0]
		max_work_time 		= config_time[1]
	case 3:
		force_update_time	= config_time[0]
		max_work_time 		= config_time[1]
		max_alive_time		= config_time[2] * 1000 * 1000 * 1000
	}

	now_time := public.Now_Time_S()

	ret_val, err := csm.rdb.Get(csm.ctx, key).Result()

	if err != nil {
		if err == redis.Nil {	
			ret_val = ""
		}else{
			public.DBG_ERR("get value failed", err)
			return ""
		}
	}

	//public.DBG_LOG("ret_val : ", ret_val)

	var cache_data Standard_CSM_Cache
	public.Parser_Json(ret_val, &cache_data)

	if cache_data.Wait == true {
		
		for ;;{
			ret_val, err := csm.rdb.Get(csm.ctx, key).Result()

			if err != nil {
				public.DBG_ERR("get value failed", err)
				return ret_val
			}

			public.Parser_Json(ret_val, &cache_data)

			if cache_data.D != "" {
				break
			}

			if cache_data.Wait == false {
				break
			}

			public.Sleep(3000)

			_now_time := public.Now_Time_S()

			if _now_time-now_time > 10000 {
				csm.rdb.Del(csm.ctx, key)
				break
			}
		}
	}

	if cache_data.D == "" {
	
		var new_cache_data Standard_CSM_Cache
		new_cache_data.Wait = true
		err := csm.rdb.Set(csm.ctx, key, public.Build_Json(new_cache_data), time.Duration(max_alive_time)).Err()
		if err != nil {
			public.DBG_ERR("set value failed", err)
			return ""
		}

		defer func() {
			if r := recover(); r != nil {
				public.DBG_ERR("err:", r)
				csm.rdb.Del(csm.ctx, key)
			}
		}()

		new_data := new_cache_func()

		new_data_str := public.Build_Json(new_data)

		now_time = public.Now_Time_S()

		new_cache_data.D	= new_data_str
		new_cache_data.L	= now_time
		new_cache_data.LW	= 0
		new_cache_data.W	= false
		new_cache_data.Wait	= false

		err = csm.rdb.Set(csm.ctx, key, public.Build_Json(new_cache_data), time.Duration(max_alive_time)).Err()
		if err != nil {
			public.DBG_ERR("set value failed", err)
			return ""
		}

		return new_data_str

	} else if ((now_time - cache_data.L >= force_update_time) && !cache_data.W) || ((cache_data.LW != 0) && (now_time - cache_data.LW >= max_work_time)) {		

		go func(){
			var new_cache_data Standard_CSM_Cache
			new_cache_data.D	= cache_data.D
			new_cache_data.L	= now_time
			new_cache_data.LW	= now_time
			new_cache_data.W	= true
			new_cache_data.Wait = false
	
			err := csm.rdb.Set(csm.ctx, key, public.Build_Json(new_cache_data), time.Duration(max_alive_time)).Err()
			if err != nil {
				public.DBG_ERR("set value failed", err)
			}
	
			defer func() {
				if r := recover(); r != nil {
					public.DBG_ERR("err:", r)
					csm.rdb.Del(csm.ctx, key)
				}
			}()
	
			new_data := new_cache_func()
	
			new_data_str := public.Build_Json(new_data)
	
			now_time = public.Now_Time_S()
	
			new_cache_data.D	= new_data_str
			new_cache_data.L	= now_time
			new_cache_data.LW	= 0
			new_cache_data.W	= false
			new_cache_data.Wait = false
	
			err = csm.rdb.Set(csm.ctx, key, public.Build_Json(new_cache_data), time.Duration(max_alive_time)).Err()
			if err != nil {
				public.DBG_ERR("set value failed", err)

			}
		}()

		return cache_data.D		
	}

	//public.DBG_LOG("now_time - cache_data.L: ", (now_time - cache_data.L), "    force_update_time:", force_update_time)

	return cache_data.D
}

func (csm *Cache_Sql_Manager) Del_Cache(key string) {
	err := csm.rdb.Del(csm.ctx, key)
	if err != nil{
		public.DBG_ERR("del cache err:", err)
	}
}


func Set_Cache(key string, value interface{}, config_time ...int64){
	//config_time[0]	force_update_time
	//config_time[1]	max_work_time
	//config_time[2]	max_alive_time

	cache_sql_manager.Set_Cache(key, value, config_time...)
}

func Get_Cache(key string, new_cache_func New_Cache_Func, config_time ...int64) string {
	//config_time[0]	force_update_time
	//config_time[1]	max_work_time
	//config_time[2]	max_alive_time
	
	return cache_sql_manager.Get_Cache(key, new_cache_func, config_time...)
}

func Del_Cache(key string) {
	cache_sql_manager.Del_Cache(key)
}


func init() {

	redis_global_param_exist := true

	redis_ip, e		:= app.Global[string]("redis_ip")
	redis_global_param_exist = redis_global_param_exist && e
	redis_passwd, e	:= app.Global[string]("redis_passwd")
	redis_global_param_exist = redis_global_param_exist && e
	redis_db, e		:= app.Global[float64]("redis_db")
	redis_global_param_exist = redis_global_param_exist && e

	if !redis_global_param_exist{
		panic("redis no config")
	}

	cache_sql_manager.Init(redis_ip, redis_passwd, int(redis_db))
}


