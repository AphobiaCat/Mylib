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

		val := tmp.Get("111")
		
		public.DBG_LOG("val[", val, "]")

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

