package base_struct

import(
	"mylib/src/public"
	cache "mylib/src/module/cachesql_manager"
	"sync"
)

const(
	ROOM_CREATE	= "create_room"
	ROOM_JOIN	= "join_room"
	ROOM_EXIT	= "exit_room"
)

type Room struct{
	Creator				string			`json:"creator"`
	Members				[]string		`json:"members"`
}

type Notify struct{
	Action				string			`json:"action"`
	Payload				string			`json:"payload"`
}

type Room_Manager struct{
	room					Thread_Map[Room]	// room_id as index
	room_index				map[string]string	// who -> room id
	room_index_lock 		sync.Mutex
	notify					func(target_user string, who string, payload interface{})

	room_speed_access		map[string]Room		// room id
	room_speed_access_lock 	sync.Mutex
	redis_key				string
}

func (o *Room_Manager) Init(redis_key string, notify func(target_user string, who string, payload interface{})){
	o.room				= New_Thread_Map[Room](redis_key)
	o.room_index		= make(map[string]string)
	o.notify			= notify
	o.room_speed_access	= make(map[string]Room)
	o.redis_key			= redis_key

	room_info := o.room.Get_All()

	for room_id, room := range room_info{
		o.room_index[room.Creator] = room_id
		o.room_speed_access[room_id] = room
	}
}

func (o *Room_Manager) Create_Room(creator string)(rand_room_id_str string, succ bool){

	o.Exit_Room(creator)

	for i:= 0; i < 5; i++{
		rand_room_id_str = public.ConvertNumToStr(int64(public.Rand_U64() % 1000000))

		if len(rand_room_id_str) != 6{
			fill_zero := "000000"
		
			rand_room_id_str += fill_zero[:6 - len(rand_room_id_str)]
		}
	
		if !o.room.HExist(rand_room_id_str){

			new_room := Room{Creator: creator}

			o.room.Ready_Set(rand_room_id_str)
			o.room.Set(rand_room_id_str, new_room)

			o.room_index_lock.Lock()
			o.room_index[creator] = rand_room_id_str
			o.room_index_lock.Unlock()
		
			succ = true

			o.notify(creator, creator, Notify{Action: ROOM_CREATE})

			o.room_speed_access_lock.Lock()
			o.room_speed_access[rand_room_id_str] = new_room
			o.room_speed_access_lock.Unlock()
			
			return
		}
	}

	public.DBG_ERR("create room error retry than 5")

	return
}

func (o *Room_Manager) Join_Room(who string, room_id string)(Room, bool){

	o.Exit_Room(who)

	o.room.Ready_Set(room_id)
	
	room, exist := o.room.Get(room_id)

	if !exist{
		return room, false
	}

	room.Members = append(room.Members, who)
	
	o.room.Set(room_id, room)

	for _, user_id := range room.Members{
		o.notify(user_id, who, Notify{Action: ROOM_JOIN})
	}
	o.notify(room.Creator, who, Notify{Action: ROOM_JOIN})

	o.room_speed_access_lock.Lock()
	o.room_speed_access[room_id] = room
	o.room_speed_access_lock.Unlock()

	o.room_index_lock.Lock()
	o.room_index[who] = room_id
	o.room_index_lock.Unlock()

	return room, true
}

func (o *Room_Manager) Exit_Room(who string){
	o.room_index_lock.Lock()
	room_id, exist := o.room_index[who]
	o.room_index_lock.Unlock()

	if exist{	
		o.room.Ready_Set(room_id)
	
		room, exist := o.room.Get(room_id)

		if !exist{
			return
		}

		if room.Creator == who{
			o.room_speed_access_lock.Lock()
			delete(o.room_speed_access, room_id)
			o.room_speed_access_lock.Unlock()

			o.room_index_lock.Lock()
			delete(o.room_index, room.Creator)
			o.room_index_lock.Unlock()

			for _, user_id := range room.Members{
				o.room_index_lock.Lock()
				delete(o.room_index, user_id)
				o.room_index_lock.Unlock()
			}

			o.room.Del(room_id)
		}else{

			o.room_index_lock.Lock()
			delete(o.room_index, who)
			o.room_index_lock.Unlock()
		
			for index, user_id := range room.Members{
				if user_id == who{
					room.Members = append(room.Members[0: index], room.Members[index + 1:]...)
					break
				}
			}
			o.room.Set(room_id, room)

			o.room_speed_access_lock.Lock()
			o.room_speed_access[room_id] = room
			o.room_speed_access_lock.Unlock()
		}

		for _, user_id := range room.Members{
			o.notify(user_id, who, Notify{Action: ROOM_EXIT})
		}
		o.notify(room.Creator, who, Notify{Action: ROOM_EXIT})
	}
}

func (o *Room_Manager) Do_Sth(who string, action string, payload string)bool{
	o.room_index_lock.Lock()
	room_id, exist := o.room_index[who]
	o.room_index_lock.Unlock()

	if exist{
		o.room_speed_access_lock.Lock()
		room, exist := o.room_speed_access[room_id]
		o.room_speed_access_lock.Unlock()

		if exist{
			for _, user_id := range room.Members{
				o.notify(user_id, who, Notify{Action: action, Payload: payload})
			}
			o.notify(room.Creator, who, Notify{Action: action, Payload: payload})

			return true
		}

		return false
	}

	return false
}

func (o *Room_Manager) List_Room(limit string, offset string) string{	

	o.room_speed_access_lock.Lock()
	total_map := o.room_speed_access
	o.room_speed_access_lock.Unlock()

	public.DBG_ERR(total_map)

	ret_info := cache.Get_Cache(o.redis_key + "_list:" + limit + ":" + offset, func()interface{}{
		
		var ret []Room

		i := int64(0)
		wait_stop := false
		offset_num := public.ConvertStrToNum(offset)
		limit_num := public.ConvertStrToNum(limit)

		for _, room_info := range total_map{
		
			if i >= offset_num{
				wait_stop = true
			}

			if wait_stop{
				if i < limit_num + offset_num{
					ret = append(ret, room_info)
				}else{
					break
				}
			}

			i += 1
		}

		return ret
	}, 30, 30, 30)

	return ret_info
}


func New_Room(redis_key string, notify func(target_user string, who string, payload interface{}))*Room_Manager{
	room := &Room_Manager{}
	room.Init(redis_key, notify)
	
	return room
}

