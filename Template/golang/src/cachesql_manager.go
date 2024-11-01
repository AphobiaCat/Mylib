package main

import (
	"context"
    "github.com/redis/go-redis/v9"
    "time"
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
		DBG_ERR("unable connet Redis:", err)
	}

	DBG_LOG("connect redis server succ")

}

func (csm *Cache_Sql_Manager) Set_Cache(key string, value interface{}, config_time ...int64) {
	max_alive_time 		:= int64(60 * 10)

	if len(config_time) == 1{
		max_alive_time	= config_time[2]
	}

	now_time := now_time_s()

	var new_cache_data Standard_CSM_Cache
	
	new_cache_data.D	= Build_Jason(value)
	new_cache_data.L	= now_time
	new_cache_data.LW	= 0
	new_cache_data.W	= false
	new_cache_data.Wait	= false

	err := csm.rdb.Set(csm.ctx, key, Build_Jason(new_cache_data), time.Duration(max_alive_time)).Err()
	if err != nil {
		DBG_ERR("set value failed", err)
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

	now_time := now_time_s()

	ret_val, err := csm.rdb.Get(csm.ctx, key).Result()

	if err != nil {
		if err == redis.Nil {	
			ret_val = ""
		}else{
			DBG_ERR("get value failed", err)
			return ""
		}
	}

	//DBG_LOG("ret_val : ", ret_val)

	var cache_data Standard_CSM_Cache
	Parser_Jason(ret_val, &cache_data)

	if cache_data.Wait == true {
		
		//DBG_LOG("PAHT CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")
		for ;;{
			ret_val, err := csm.rdb.Get(csm.ctx, key).Result()

			if err != nil {
				DBG_ERR("get value failed", err)
				return ret_val
			}

			Parser_Jason(ret_val, &cache_data)

			if cache_data.D != "" {
				break
			}

			if cache_data.Wait == false {
				break
			}

			sleep(3000)

			_now_time := now_time_s()

			if _now_time-now_time > 10000 {
				csm.rdb.Del(csm.ctx, key)
				break
			}
		}
	}

	if cache_data.D == "" {
		//DBG_LOG("PAHT AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAa")

	
		var new_cache_data Standard_CSM_Cache
		new_cache_data.Wait = true
		err := csm.rdb.Set(csm.ctx, key, Build_Jason(new_cache_data), time.Duration(max_alive_time)).Err()
		if err != nil {
			DBG_ERR("set value failed", err)
			return ""
		}

		defer func() {
			if r := recover(); r != nil {
				DBG_ERR("err:", r)
				csm.rdb.Del(csm.ctx, key)
			}
		}()

		new_data := new_cache_func()

		new_data_str := Build_Jason(new_data)

		now_time = now_time_s()

		new_cache_data.D	= new_data_str
		new_cache_data.L	= now_time
		new_cache_data.LW	= 0
		new_cache_data.W	= false
		new_cache_data.Wait	= false

		err = csm.rdb.Set(csm.ctx, key, Build_Jason(new_cache_data), time.Duration(max_alive_time)).Err()
		if err != nil {
			DBG_ERR("set value failed", err)
			return ""
		}

		return new_data_str

	} else if ((now_time - cache_data.L >= force_update_time) && !cache_data.W) || ((cache_data.LW != 0) && (now_time - cache_data.LW >= max_work_time)) {		
		//DBG_LOG("PAHT BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")

		var new_cache_data Standard_CSM_Cache
		new_cache_data.D	= cache_data.D
		new_cache_data.L	= now_time
		new_cache_data.LW	= now_time
		new_cache_data.W	= true
		new_cache_data.Wait	= false

		err := csm.rdb.Set(csm.ctx, key, Build_Jason(new_cache_data), time.Duration(max_alive_time)).Err()
		if err != nil {
			DBG_ERR("set value failed", err)
			return ""
		}

		defer func() {
			if r := recover(); r != nil {
				DBG_ERR("err:", r)
				csm.rdb.Del(csm.ctx, key)
			}
		}()

		new_data := new_cache_func()

		new_data_str := Build_Jason(new_data)

		now_time = now_time_s()

		new_cache_data.D	= new_data_str
		new_cache_data.L	= now_time
		new_cache_data.LW	= 0
		new_cache_data.W	= false
		new_cache_data.Wait	= false

		err = csm.rdb.Set(csm.ctx, key, Build_Jason(new_cache_data), time.Duration(max_alive_time)).Err()
		if err != nil {
			DBG_ERR("set value failed", err)
			return new_data_str
		}

		return new_data_str
	}

	//DBG_LOG("now_time - cache_data.L: ", (now_time - cache_data.L), "    force_update_time:", force_update_time)

	return cache_data.D
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

func Init_Cache_Sql() {
	cache_sql_manager.Init(redis_server_ip, redis_server_passwd, redis_DB)
}


