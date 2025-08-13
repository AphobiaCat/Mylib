package app_b

import (
	"mylib/src/public"
)

func process(payload interface{})(interface{}, bool){
	return "B:" + payload.(string), true
}

func Test_App_B() {
	public.DBG_LOG("hello app_b package")
	ret, _ := public.Call("A", "hello")

	public.DBG_LOG("B recv callback:", ret)
}

func init(){
	public.Reg("B", process)
}


