package main

func test_timer(){
	DBG_LOG("hello wolrd")
}

func test_timer2(){
	DBG_LOG("hello wolrd2")
}

func test_timer3(){
	DBG_LOG("hello wolrd3")
}

func example_timer_manager(){
	Init_Timer();
	
	Reg_Timer(1, test_timer)
	Reg_Timer(2, test_timer2)
	Reg_Timer(3, test_timer3)
}

