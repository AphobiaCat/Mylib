package app_a

import (
	"mylib/src/public"
)

func process(payload string)(string, bool){
	return "A:" + payload, true
}

func Test_App_A() {
	public.DBG_LOG("hello app_a package")
	ret, _ := public.Callbacks["B"]("hello")

	public.DBG_LOG("A recv callback:", ret)
}

func init(){
	public.Reg("A", process)
}


