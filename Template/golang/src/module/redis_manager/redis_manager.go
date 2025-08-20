package redis_manager

import (
    "context"
    "github.com/redis/go-redis/v9"
    "sync"
	"crypto/tls"

    "mylib/src/public"
    "time"
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

func (rm *Redis_Manager) LPUSH(redis_key string, data string){
	err := rm.rdb.LPush(rm.ctx, redis_key, data).Err()

	if err != nil{
		public.DBG_ERR("queue set value failed", err)
	}
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

func (rm *Redis_Manager) List_Range(redis_key string, start_pos int64, end_pos int64)([]string, bool){
	values, err := rm.rdb.LRange(rm.ctx, redis_key, start_pos, end_pos).Result()
    if err != nil {
    	public.DBG_ERR("redis list range err:", err)
        return values, false
    }

    return values, true
}

func (rm *Redis_Manager) Add_Num(redis_key string, num int64, timeout_s ...int64)(int64, bool){
	values, err := rm.rdb.IncrBy(rm.ctx, redis_key, num).Result()
    if err != nil {
    	public.DBG_ERR("redis incr by int err:", err)
        return values, false
    }

	if len(timeout_s) > 0{
		rm.rdb.Expire(rm.ctx, redis_key, time.Duration(timeout_s[0] * 1000 * 1000 * 1000))
	}

    return values, true
}


func (rm *Redis_Manager) Add_Float_Num(redis_key string, num float64, timeout_s ...int64)(float64, bool){
	values, err := rm.rdb.IncrByFloat(rm.ctx, redis_key, num).Result()
    if err != nil {
    	public.DBG_ERR("redis incr by float err:", err)
        return values, false
    }

	if len(timeout_s) > 0{
		rm.rdb.Expire(rm.ctx, redis_key, time.Duration(timeout_s[0] * 1000 * 1000 * 1000))
	}

    return values, true
}

func (rm *Redis_Manager) Get(redis_key string)(string, bool){
	values, err := rm.rdb.Get(rm.ctx, redis_key).Result()
    if err != nil {
    	public.DBG_ERR("get err:", err)
        return values, false
    }

    return values, true
}


func (rm *Redis_Manager) Delete(redis_key string){
	err := rm.rdb.Del(rm.ctx, redis_key)
	if err != nil {
		public.DBG_ERR("del value failed", err)
	}
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

func LPUSH(redis_key string, data string){
	redis_manager.LPUSH(redis_key, data)
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

func List_Range(redis_key string, start_pos int64, end_pos int64)([]string, bool){
	return redis_manager.List_Range(redis_key, start_pos, end_pos)
}

func Add_Num(key string, num int64, timeout_s ...int64) (int64, bool){
	return redis_manager.Add_Num(key, num, timeout_s...)
}

func Add_Float_Num(key string, num float64, timeout_s ...int64) (float64, bool){
	return redis_manager.Add_Float_Num(key, num, timeout_s...)
}

func Get(key string)(string, bool) {
	return redis_manager.Get(key)
}

func Delete(key string) {
	redis_manager.Delete(key)
}


func init(){

	redis_manager.value_lock_index = make(map[string]int)

	redis_manager.ctx = context.Background()

	if public.Config.Redis.EnableTls{
		redis_manager.rdb = redis.NewClient(&redis.Options{
			Addr:     public.Config.Redis.Ip,
			Password: public.Config.Redis.Password,
			DB:       public.Config.Redis.DB,
			TLSConfig: &tls.Config{},
		})
	}else{
		redis_manager.rdb = redis.NewClient(&redis.Options{
			Addr:     public.Config.Redis.Ip,
			Password: public.Config.Redis.Password,
			DB:       public.Config.Redis.DB,
		})
	}
	
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

