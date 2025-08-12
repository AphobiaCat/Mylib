package app_b

import (
	"mylib/src/public"
)

func process(payload string)(string, bool){
	return "B:" + payload, true
}

func Test_App_B() {
	public.DBG_LOG("hello app_b package")
	ret, _ := public.Callbacks["A"]("hello")

	public.DBG_LOG("B recv callback:", ret)
}

func init(){
	public.Reg("B", process)
}


