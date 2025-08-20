package base_struct

import(
	"mylib/src/public"
	"github.com/redis/go-redis/v9"
	"sync"
)

type Room struct{
	Creator	string		`json:"creator"`
	Members	[]string	`json:"members"`
}

type All_Room struct{
	rooms			[]Room
	room_creator	map[string]string
}

var all_room			[]All_Room
var all_room_len		int
var all_room_index		map[string]int
var all_room_index_lock	sync.Mutex

func (rm *Redis_Manager) scan_room(redis_key string){

	index, exist := all_room_index[redis_key]

	if !exist{
		index = all_room_len
		all_room_index[redis_key] = all_room_len
		all_room_len += 1
	}

	values, err := rm.rdb.ZRange(rm.ctx, redis_key, 0, -1).WithScores().Result()
    if err != nil {
    	public.DBG_ERR("scan_room err:", err)
    	return
    }

	i := 0
	k := 1

    for ; k < len(values); i += 1, k += 1{
		tmp_room := Room{}
    
		public.Parser_Json(values[i], &tmp_room)
		all_room.room_creator[tmp_room.Creator] = values[k]
		all_room.rooms = append(all_room.rooms, tmp_room)
    }
}

func (rm *Redis_Manager) create_room(redis_key string, creator string)(int){
	if initial[redis_key] == false{
		initial[redis_key] = true
		rm.scan_room(redis_key)
	}

	rm.close_room(redis_key, creator)

	var room Room

	room.Creator = creator

	rand_id := int(public.Rand_U64() % 9000000) + 1000000

	err := rm.rdb.ZAdd(rm.ctx, redis_key, redis.Z{Score: rand_id, Member: public.Build_Json(room)}).Err()
	
    if err != nil {
    	public.DBG_ERR("create_room err[", err, "]")
        return -1
    }

    room_creator[creator] = rand_id

    return rand_id
}

func (rm *Redis_Manager) join_room(redis_key string, creator string){
	if initial[redis_key] == false{
		initial[redis_key] = true
		rm.scan_room(redis_key)
	}

	err := rm.rdb.ZRemRangeByScore(rm.ctx, redis_key, room_creator[creator], room_creator[creator]).Err()

	delete(room_creator, creator)

	if err != nil{
		public.DBG_ERR("close_room error ", err)
	}
}

func (rm *Redis_Manager) exit_room(redis_key string, creator string){
	if initial[redis_key] == false{
		initial[redis_key] = true
		rm.scan_room(redis_key)
	}

	err := rm.rdb.ZRemRangeByScore(rm.ctx, redis_key, room_creator[creator], room_creator[creator]).Err()

	delete(room_creator, creator)

	if err != nil{
		public.DBG_ERR("close_room error ", err)
	}
}

func (rm *Redis_Manager) close_room(redis_key string, creator string){
	if initial[redis_key] == false{
		initial[redis_key] = true
		rm.scan_room(redis_key)
	}

	err := rm.rdb.ZRemRangeByScore(rm.ctx, redis_key, room_creator[creator], room_creator[creator]).Err()

	delete(room_creator, creator)

	if err != nil{
		public.DBG_ERR("close_room error ", err)
	}
}


func Create_Room(redis_key string, reload_count int64, reset_time int64)(int64){
	return redis_manager.Timer_Count(redis_key, reload_count, reset_time)
}

func init(){
	all_room_index = make(map[string]int)
}

