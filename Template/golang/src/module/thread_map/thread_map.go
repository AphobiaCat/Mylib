package thread_map

import (
	"mylib/src/public"
	"sync"
	"time"
	cache "mylib/src/module/cachesql_manager"
)


type Thread_Map[Key comparable, Val any] struct{
	total_lock		sync.Mutex				`json:"-"`

	map_lock		sync.Mutex				`json:"-"`
	Map				map[Key]uint32			`json:"map"`
	
	val_lock		[]sync.Mutex			`json:"-"`
	Exist			[]bool					`json:"exist"`
	Val_Array		[]Val					`json:"val_array"`
	
	Array_Len		uint32					`json:"array_len"`
	Now_Index		uint32					`json:"now_index"`
					
	Is_Delete		bool					`json:"is_delete"`
}

func (this *Thread_Map[Key, Val]) init(map_redis_key string, save_interval uint64, save_time int64){

	this.total_lock.Lock()
	defer this.total_lock.Unlock()

	this.Map = make(map[Key]uint32)
	
	tmp_info_str := cache.Get_Cache(map_redis_key, func()interface{}{
		return ""
	})

	public.Parser_Json(tmp_info_str, &this)

	if this.Array_Len != 0{
		this.val_lock = make([]sync.Mutex, this.Array_Len)
	}
	
	ticker := time.NewTicker(1 * time.Second)

	save_data_point := &this
	
	go func(){
		defer ticker.Stop()
		now_time := uint64(0)
		for !this.Is_Delete {
			select{
				case _ = <-ticker.C:
					now_time++
			}
			
			if now_time % save_interval == 0{
				cache.Set_Cache(map_redis_key, *save_data_point, save_time)
			}
		}
    }()
}

func (this *Thread_Map[Key, Val]) New(new_key Key, new_val Val)bool{

	this.map_lock.Lock()
	_, exist := this.Map[new_key]
	if exist{
		this.map_lock.Unlock()
		public.DBG_ERR("this key id[", new_key, "] exist")
		return false
	}
	this.map_lock.Unlock()

	map_index	:= uint32(0xFFFFFFFF)

	for i := uint32(0); i < this.Array_Len; i++{

		this.val_lock[i].Lock()
	
		if !this.Exist[(i + this.Now_Index) % this.Array_Len]{
			map_index = (i + this.Now_Index) % this.Array_Len
			this.Exist[map_index] = true
			
			this.val_lock[i].Unlock()
			break
		}
		this.val_lock[i].Unlock()
	}

	this.total_lock.Lock()

	if map_index == 0xFFFFFFFF{	//map full
		tmp_val_array	:= this.Val_Array
		tmp_exist_array	:= this.Exist

		map_index		= this.Array_Len

		this.Array_Len	+= 100

		this.Val_Array	= make([]Val, this.Array_Len)
		this.Exist		= make([]bool, this.Array_Len)
		this.val_lock	= make([]sync.Mutex, this.Array_Len)
		
		for index, val	:= range tmp_val_array{
			this.Val_Array[index]	= val
			this.Exist[index]		= tmp_exist_array[index]
		}
	}

	this.total_lock.Unlock()

	this.val_lock[map_index].Lock()
	this.Val_Array[map_index]	= new_val
	this.Exist[map_index]		= true
	this.val_lock[map_index].Unlock()
	this.Now_Index = map_index + 1

	this.map_lock.Lock()
	this.Map[new_key] = map_index
	this.map_lock.Unlock()

	return true

}

func (this *Thread_Map[Key, Val]) New_Or_Update(new_key Key, new_val Val)(last Val, is_update bool){
	this.map_lock.Lock()	
	map_index, exist := this.Map[new_key]
	if exist{
		this.map_lock.Unlock()
		this.val_lock[map_index].Lock()
		this.total_lock.Lock()
		last = this.Val_Array[map_index]
		this.Val_Array[map_index] = new_val
		this.total_lock.Unlock()
		this.val_lock[map_index].Unlock()

		is_update = true
		return last, is_update
	}
	this.map_lock.Unlock()

	if succ := this.New(new_key, new_val); !succ{
		public.DBG_ERR("key[", new_key, "] val[", new_val, "] create error")
	}

	last = new_val
	return last, is_update
}

func (this *Thread_Map[Key, Val]) Delete(delete_key Key)bool{
	this.map_lock.Lock()	
	map_index, exist := this.Map[delete_key]
	if exist{
		delete(this.Map, delete_key)
		this.map_lock.Unlock()
		this.val_lock[map_index].Lock()
		this.total_lock.Lock()
		this.Exist[map_index] = false
		this.total_lock.Unlock()
		this.val_lock[map_index].Unlock()

		return true
	}
	this.map_lock.Unlock()
	return false
}

func (this *Thread_Map[Key, Val]) Get_Val(key Key)(Val, bool){
	this.map_lock.Lock()
	map_index, exist := this.Map[key]
	this.map_lock.Unlock()

	if exist{

		this.val_lock[map_index].Lock()
		this.total_lock.Lock()
		ret := this.Val_Array[map_index]
		this.total_lock.Unlock()
		this.val_lock[map_index].Unlock()
		
		return ret, true
	}

	var zero Val
	return zero, false
}

func (this *Thread_Map[Key, Val]) Borrow_Val(key Key)(Val, bool){
	this.map_lock.Lock()
	map_index, exist := this.Map[key]
	this.map_lock.Unlock()

	if exist{

		this.val_lock[map_index].Lock()
		this.total_lock.Lock()
		ret := this.Val_Array[map_index]
		this.total_lock.Unlock()
		
		return ret, true
	}

	var zero Val
	return zero, false
}

func (this *Thread_Map[Key, Val]) Return_Val(key Key, val Val)bool{
	this.map_lock.Lock()
	map_index, exist := this.Map[key]
	this.map_lock.Unlock()

	if exist{

		this.total_lock.Lock()
		this.Val_Array[map_index] = val
		this.total_lock.Unlock()
		this.val_lock[map_index].Unlock()
		
		return true
	}

	return false
}


func (this *Thread_Map[Key, Val]) Delete_Map(){
	this.Is_Delete = true
}

func New_Thread_Map[Key comparable, Val any](map_redis_key string, save_interval uint64, save_time int64)Thread_Map[Key, Val]{
	ret := Thread_Map[Key, Val]{}

	ret.init(map_redis_key, save_interval, save_time)

	return ret
}

