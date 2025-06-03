package redis_manager

import (
    "context"
    "github.com/redis/go-redis/v9"
    "sync"

    "mylib/src/public"
    "mylib/src/module/app"
)


var redis_manager Redis_Manager

type Redis_Manager struct{
	rdb						*redis.Client
	ctx						context.Context

	value_lock				[]sync.Mutex
	value_lock_index		map[string]int
	value_lock_index_lock	sync.Mutex
}

func (rm *Redis_Manager) Set_Value(value_key string, value interface{}){

	rm.value_lock_index_lock.Lock()

	val, exist := rm.value_lock_index[value_key]

	if !exist{
		rm.value_lock = append(rm.value_lock, sync.Mutex{})
		rm.value_lock_index[value_key] = len(rm.value_lock) - 1
		val = rm.value_lock_index[value_key]
	}
	
	rm.value_lock_index_lock.Unlock()

	rm.value_lock[val].Lock()

	err := rm.rdb.Set(rm.ctx, value_key, value, 0).Err()
	if err != nil {
		rm.value_lock[val].Unlock()
		public.DBG_ERR("set value failed", err)
		return
	}
	
	rm.value_lock[val].Unlock()
}

func (rm *Redis_Manager) Get_Value(value_key string) interface{}{

	rm.value_lock_index_lock.Lock()

	val, exist := rm.value_lock_index[value_key]

	if !exist{
		rm.value_lock = append(rm.value_lock, sync.Mutex{})
		rm.value_lock_index[value_key] = len(rm.value_lock) - 1
		val = rm.value_lock_index[value_key]
	}
	
	rm.value_lock_index_lock.Unlock()

	rm.value_lock[val].Lock()

	ret_val, err := rm.rdb.Get(rm.ctx, value_key).Result()
	if err != nil {
		rm.value_lock[val].Unlock()
		public.DBG_ERR("get value failed", err)
		return ret_val
	}
	//public.DBG_LOG("key value:", val)
	
	rm.value_lock[val].Unlock()

	return ret_val
}

func (rm *Redis_Manager) Return_Value(value_key string, value interface{}){
	rm.value_lock_index_lock.Lock()

	val, exist := rm.value_lock_index[value_key]

	if !exist{
		rm.value_lock_index_lock.Unlock()
		public.DBG_ERR("Return_Value value failed, this value no Borrow")
		return
	}
	
	rm.value_lock_index_lock.Unlock()

	

	err := rm.rdb.Set(rm.ctx, value_key, value, 0).Err()
	if err != nil {
		rm.value_lock[val].Unlock()
		public.DBG_ERR("set value failed", err)
		return
	}
	
	rm.value_lock[val].Unlock()
}

func (rm *Redis_Manager) Borrow_Value(value_key string) interface{}{
	rm.value_lock_index_lock.Lock()

	val, exist := rm.value_lock_index[value_key]

	if !exist{
		rm.value_lock = append(rm.value_lock, sync.Mutex{})
		rm.value_lock_index[value_key] = len(rm.value_lock) - 1
		val = rm.value_lock_index[value_key]
	}
	
	rm.value_lock_index_lock.Unlock()

	rm.value_lock[val].Lock()

	ret_val, err := rm.rdb.Get(rm.ctx, value_key).Result()
	if err != nil {
		rm.value_lock[val].Unlock()
		public.DBG_ERR("get value failed", err)
		return ret_val
	}

	return ret_val
}


func (rm *Redis_Manager) Queue_Set(redis_key string, data interface{}){
	err := rm.rdb.LPush(rm.ctx, redis_key, public.Build_Json(data)).Err()

	if err != nil{
		public.DBG_ERR("queue set value failed", err)
	}
}

func (rm *Redis_Manager) Queue_Get(redis_key string)(string, bool){
	task, err := rm.rdb.RPop(rm.ctx, redis_key).Result()
	 
	if err != nil {
		if err != redis.Nil {
			public.DBG_ERR("queue get value failed", err)
		}
		return "", false
	}
	return task, true
}

func (rm *Redis_Manager) Stack_Set(redis_key string, data interface{}){
	err := rm.rdb.LPush(rm.ctx, redis_key, public.Build_Json(data)).Err()

	if err != nil{
		public.DBG_ERR("stack set value failed", err)
	}
}

func (rm *Redis_Manager) Stack_Get(redis_key string)(string, bool){
	task, err := rm.rdb.LPop(rm.ctx, redis_key).Result()

	if err != nil {
		if err != redis.Nil {
			public.DBG_ERR("stack get value failed", err)
		}
		return "", false
	}
	return task, true
}



func Set_Value(value_key string, value interface{}){
	redis_manager.Set_Value(value_key, value)
}

func Return_Value(value_key string, value interface{}){
	redis_manager.Return_Value(value_key, value)
}

func Get_Value(value_key string) interface{}{
	return redis_manager.Get_Value(value_key)
}

func Borrow_Value(value_key string) interface{}{
	return redis_manager.Borrow_Value(value_key)
}

func Queue_Set(redis_key string, data interface{}){
	redis_manager.Queue_Set(redis_key, data)
}

func Queue_Get(redis_key string)(string, bool){
	return redis_manager.Queue_Get(redis_key)
}

func Stack_Set(redis_key string, data interface{}){
	redis_manager.Stack_Set(redis_key, data)
}

func Stack_Get(redis_key string)(string, bool){
	return redis_manager.Stack_Get(redis_key)
}



func init(){

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

	redis_manager.value_lock_index = make(map[string]int)

	redis_manager.ctx = context.Background()

	redis_manager.rdb = redis.NewClient(&redis.Options{
		Addr:     redis_ip,
		Password: redis_passwd,
		DB:       int(redis_db),    
	})
	
	_, err := redis_manager.rdb.Ping(redis_manager.ctx).Result()
	if err != nil {
		public.DBG_ERR("unable connet Redis:", err)
		panic(err)
	}
	public.DBG_LOG("connect redis server succ")
	

	//rdb := redis_manager.rdb
	//public.DBG_LOG_VAR(rdb)
}

func Close_Redis(){
	redis_manager.rdb.Close()
}

