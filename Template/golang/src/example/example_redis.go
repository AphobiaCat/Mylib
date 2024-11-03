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
		
		public.Parser_Jason(ret.(string), &value)

		value.Data 		+= 1
		value.OtherData += " h"

		return_val := public.Build_Net_Jason(value)
		
		time.Sleep(20 * time.Millisecond)
		redis_manager.Return_Value("test_key", return_val.String())
		time.Sleep(20 * time.Millisecond)
	}
}

func Example_Redis_Manager(){

	var value Test_Redis_Data
	value.Data 		= 0
	value.OtherData	= "hello wolrd"
	
	result := public.Build_Net_Jason(value)

	redis_manager.Set_Value("test_key", result.String())

	go test1()
	time.Sleep(10 * time.Millisecond)
	go test2()

	time.Sleep(10 * time.Second)

	redis_manager.Close_Redis()
	public.DBG_LOG("close");
}


