package main

import (
    "context"
    "github.com/redis/go-redis/v9"
    "sync"
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
		DBG_ERR("set value failed", err)
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
		DBG_ERR("get value failed", err)
		return ret_val
	}
	//DBG_LOG("key value:", val)
	
	rm.value_lock[val].Unlock()

	return ret_val
}

func (rm *Redis_Manager) Return_Value(value_key string, value interface{}){
	rm.value_lock_index_lock.Lock()

	val, exist := rm.value_lock_index[value_key]

	if !exist{
		rm.value_lock_index_lock.Unlock()
		DBG_ERR("Return_Value value failed, this value no Borrow")
		return
	}
	
	rm.value_lock_index_lock.Unlock()

	

	err := rm.rdb.Set(rm.ctx, value_key, value, 0).Err()
	if err != nil {
		rm.value_lock[val].Unlock()
		DBG_ERR("set value failed", err)
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
		DBG_ERR("get value failed", err)
		return ret_val
	}

	return ret_val
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

func Init_Redis(server_ip string, password string, DB int){

	redis_manager.value_lock_index = make(map[string]int)

	redis_manager.ctx = context.Background()

	redis_manager.rdb = redis.NewClient(&redis.Options{
		Addr:     server_ip,
		Password: password,
		DB:       DB,    
	})
	
	_, err := redis_manager.rdb.Ping(redis_manager.ctx).Result()
	if err != nil {
		DBG_ERR("unable connet Redis:", err)
		panic(err)
	}
	DBG_LOG("connect redis server succ")
	

	//rdb := redis_manager.rdb
	//DBG_LOG_VAR(rdb)
}

func Close_Redis(){
	redis_manager.rdb.Close()
}

