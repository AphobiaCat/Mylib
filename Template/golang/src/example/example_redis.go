package example


import (
	"time"

	"mylib/src/module/redis_manager"
	"mylib/src/public"
)

type Test_Redis_Data struct{
	Data		int		`json:"data"`
	OtherData	string	`json:"other_data"`
}

func test2(){

	for i := 0; i < 200; i++{
		ret := redis_manager.Get_Value("test_key")
		public.DBG_LOG(ret);
		time.Sleep(20 * time.Millisecond)
	}
}

func test1(){

	var value Test_Redis_Data

	for i := 0; i < 100; i++{
		ret := redis_manager.Borrow_Value("test_key")
		
		public.Parser_Json(ret.(string), &value)

		value.Data 		+= 1
		value.OtherData += " h"

		return_val := public.Build_Net_Json(value)
		
		time.Sleep(20 * time.Millisecond)
		redis_manager.Return_Value("test_key", return_val.String())
		time.Sleep(20 * time.Millisecond)
	}
}

func example_redis_test1(){
	var value Test_Redis_Data
	value.Data 		= 0
	value.OtherData	= "hello wolrd"
	
	result := public.Build_Net_Json(value)

	redis_manager.Set_Value("test_key", result.String())

	go test1()
	time.Sleep(10 * time.Millisecond)
	go test2()

	time.Sleep(10 * time.Second)

	public.DBG_LOG("example_redis_test1 close");
}


func example_redis_test2(){

	for i := 0; i < 120; i++{
		count := redis_manager.Timer_Count("test_redis_timer_count", 60, 5)
		public.DBG_LOG("timer count[", count, "]")
		public.Sleep(200)
	}
}

func Example_Redis_Manager(){
	example_redis_test1()
	example_redis_test2()

	//redis_manager.Close_Redis()
}


