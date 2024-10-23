package main

import (
	"time"
)

type Test_Redis_Data struct{
	Data		int		`json:"data"`
	OtherData	string	`json:"other_data"`
}

func test2(){

	for i := 0; i < 200; i++{
		ret := Get_Value("test_key")
		DBG_LOG(ret);
		time.Sleep(20 * time.Millisecond)
	}
}

func test1(){

	var value Test_Redis_Data

	for i := 0; i < 100; i++{
		ret := Borrow_Value("test_key")
		
		Parser_Jason(ret.(string), &value)

		value.Data 		+= 1
		value.OtherData += " h"

		return_val := Build_Net_Jason(value)
		
		time.Sleep(20 * time.Millisecond)
		Return_Value("test_key", return_val.String())
		time.Sleep(20 * time.Millisecond)
	}
}

func example_redis_manager(){
	Init_Redis(redis_server_ip, redis_server_passwd, redis_DB)

	var value Test_Redis_Data
	value.Data 		= 0
	value.OtherData	= "hello wolrd"
	
	result := Build_Net_Jason(value)

	Set_Value("test_key", result.String())

	go test1()
	time.Sleep(10 * time.Millisecond)
	go test2()

	time.Sleep(10 * time.Second)

	Close_Redis()
	DBG_LOG("close");
}


