package example

import(
	"mylib/src/module/timer_manager"
	"mylib/src/public"
)

func test_timer(){
	public.DBG_LOG("hello wolrd")
}

func test_timer2(){
	public.DBG_LOG("hello wolrd2")
}

func test_timer3(){
	public.DBG_LOG("hello wolrd3")
}

func Example_Timer_Manager(){

	timer_manager.Reg_Timer(1, test_timer)
	timer_manager.Reg_Timer(2, test_timer2)
	timer_manager.Reg_Timer(3, test_timer3)
}

