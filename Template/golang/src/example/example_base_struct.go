package example

import(
	"mylib/src/public"
	"mylib/src/module/base_struct"
	
)

type Test_Thread_Map struct{
	Name	string	`json:"name"`
	Age		int		`json:"age"`
}

func test_thread_map2(tmp *base_struct.Thread_Map[Test_Thread_Map]){

	for i:= 0; i < 100; i++{
		tmp.Ready_Set("111")

		val, exist := tmp.Get("111")
		
		public.DBG_LOG("val[", val, "] exist[", exist, "]")

		val.Age += 1
		
		tmp.Set("111", val)
	}
}

func Example_thread_map(){
	
	tmp := base_struct.New_Thread_Map[Test_Thread_Map]("hello_thread_map")

	var tmp2 Test_Thread_Map

	tmp2.Name = "Dunty"
	tmp2.Age = 0

	tmp.Ready_Set("111")
	tmp.Set("111", tmp2)
	
	go test_thread_map2(&tmp)
	test_thread_map2(&tmp)
}

func notify_room(notify_who string, who_active string, msg base_struct.Notify){
	public.DBG_LOG("user[", notify_who, "] recive other user[", who_active, "] msg[", msg, "]")
}

func Example_room(){
	
	rm := base_struct.New_Room("test_room", notify_room)

	rm.Create_Room("Dunty")
	room_id, succ := rm.Create_Room("Dunty")
	public.DBG_LOG("new room[", room_id, "] succ[", succ, "]")
	rm.Join_Room("Mila", room_id)
	rm.Join_Room("Duduo", room_id)

	rm.Do_Sth("Mila", "speak", "hello room")
	rm.Do_Sth("Dunty", "speak", "hello room")
	rm.Do_Sth("Duduo", "speak", "hello room")
	
	rm.Exit_Room("Mila")

	list := rm.List_Room("10", "0")

	public.DBG_LOG(list)
	
	rm.Exit_Room("Dunty")
}


