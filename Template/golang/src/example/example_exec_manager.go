package example

import(
	"mylib/src/public"
	"mylib/src/module/exec_manager"
)

func Example_Exec(){
	_, succ := exec_manager.Exec(true, "ls")

	public.DBG_LOG("run ", succ)
}

