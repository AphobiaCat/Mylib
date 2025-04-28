package example

import(
	"mylib/src/public"
	"mylib/src/module/thread_map"
	
)

func test_thread_map2(tmp *thread_map.Thread_Map[string, uint32]){
	ret, exist := tmp.Get_Val("Dunty")

	if exist{
		public.DBG_LOG("tmp val:", ret)
	}else{
		public.DBG_ERR("tmp val no exist")
	}
}

func Example_thread_map(){
	
	tmp := thread_map.New_Thread_Map[string, uint32]("hello", 10, 60)

	tmp.New_Or_Update("Dunty", 1)

	ret, exist := tmp.Get_Val("Dunty")

	if exist{
		public.DBG_LOG("tmp val:", ret)
	}else{
		public.DBG_ERR("tmp val no exist")
	}

	tmp.Borrow_Val("Dunty")
	
	go test_thread_map2(&tmp)

	public.Sleep(1000)

	tmp.Return_Val("Dunty", 2)
}

