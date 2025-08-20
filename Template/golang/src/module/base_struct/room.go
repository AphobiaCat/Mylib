package base_struct

import(
	"mylib/src/public"
	"sync.Mutex"
)


type Room struct{
	Creator				string			`json:"creator"`
	Members				[]string		`json:"members"`
}

type Room_Manager struct{
	room			Thread_Map[Room]	// room_id index
	room_index		map[string]string	// creator -> index
	room_index_lock sync.Mutex
}

func (o *Room_Manager) Init(redis_key string){
	o.room = New_Thread_Map[Room](redis_key)
	o.room_index = make(map[string]string)

	room_info := o.room.Get_All()

	for room_id, room := range room_info{
		o.room_index[room.Creator] = room_id
	}
}

func (o *Room_Manager) Create_Room(creator string){

	o.room_index_lock.Lock()
	room_id, exist := o.room_index[creator]
	o.room_index_lock.Unlock()

	if exist{
		o.Close_Room(room_id)
	}

	for{
		rand_id_str := public.ConvertNumToStr(int64(public.Rand_U64() % 1000000))

		if len(rand_id_str) != 6{
			fill_zero := "000000"
		
			rand_id_str += fill_zero[:6 - len(rand_id_str)]
		}
	
		if !o.room.Exist(rand_id_str){

			o.room.Ready_Set(rand_id_str)
			o.room.Set(rand_id_str, Room{Creator: creator})

			o.room_index_lock.Lock()
			o.room_index[creator] = rand_id_str
			o.room_index_lock.Unlock()
		
			return
		}
	}
}

func (o *Room_Manager) Join_Room(who string, room_id string){
	o.room.Ready_Set(room_id)
	
	room := o.room.Get(room_id)

	room.Members = append(room.Members, who)
	
	o.room.Set(room_id, Room{Creator: creator})
}

func (o *Room_Manager) Close_Room(room_id string){

	room := o.room.Get(room_id)

	o.room_index_lock.Lock()
	delete(o.room_index, room.Creator)
	o.room_index_lock.Unlock()

	o.room.Del(room_id)
}



func Init_Room(redis_key string)Room_Manager{
	room := &Room_Manager{}
	room.Init(redis_key)
	

	return room
}

