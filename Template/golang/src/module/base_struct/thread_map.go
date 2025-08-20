package base_struct

import (
	"mylib/src/public"
	"sync"
	redis "mylib/src/module/redis_manager"
)

type Thread_Map[Val any] struct{
	val_lock		[]sync.Mutex
	lock_index		map[string]int
	index_lock		sync.Mutex

	redis_key		string
}

func (this *Thread_Map[Val]) init(map_redis_key string){

	all_info := redis.HGetAll(map_redis_key)

	this.val_lock = make([]sync.Mutex, len(all_info))
	this.lock_index = make(map[string]int)
	
	i := 0

	for key, _ := range all_info{
		this.lock_index[key] = i
		i += 1
	}

	this.redis_key = map_redis_key
}

func (this *Thread_Map[Val]) Get(key string) (ret Val){
	result := redis.HGet(this.redis_key, key)

	public.Parser_Json(result, &ret)
	return
}

func (this *Thread_Map[Val]) Get_All() (ret map[string]Val){
	all_info := redis.HGetAll(map_redis_key)

	ret = make(map[string]Val)

	for key, json_info := range all_info{
		var tmp Val
		public.Parser_Json(json_info, &tmp)
	
		ret[key] = tmp
	}

	return 
}


func (this *Thread_Map[Val]) Exist(key string) (bool){
	return redis.HExist(this.redis_key, key)
}

func (this *Thread_Map[Val]) Ready_Set(key string){
	this.index_lock.Lock()
	index, exist := this.lock_index[key]
	if !exist{
		this.val_lock = append(this.val_lock, sync.Mutex{})
		index = len(this.val_lock) - 1
		this.lock_index[key] = index
	}
	this.index_lock.Unlock()
	
	this.val_lock[index].Lock()
	
	return
}

func (this *Thread_Map[Val]) Set(key string, new_val Val){
	this.index_lock.Lock()
	index, exist := this.lock_index[key]
	this.index_lock.Unlock()
	if !exist{
		return
	}
	
	redis.HSet(this.redis_key, key, new_val)

	this.val_lock[index].Unlock()
	return
}

func (this *Thread_Map[Val]) Del(key string){
	redis.HDel(this.redis_key, key)
}

func New_Thread_Map[Val any](map_redis_key string)Thread_Map[Val]{
	ret := Thread_Map[Val]{}

	ret.init(map_redis_key)

	return ret
}

