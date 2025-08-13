package app_a

import (
	"mylib/src/public"
)

func process(payload interface{})(interface{}, bool){
	return "A:" + payload.(string), true
}

func Test_App_A() {
	public.DBG_LOG("hello app_a package")
	ret, _ := public.Call("B", "hello")

	public.DBG_LOG("A recv callback:", ret)
}

func init(){
	public.Reg("A", process)
}


