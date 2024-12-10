package example

import(
	"mylib/src/public"
	"mylib/src/module/file_manager"
)

func Example_file(){
	file_manager.File_Append("hello wolrd!!!", "tmp_file.txt", 512)

	ret := file_manager.File_Read("tmp_file.txt")

	public.DBG_LOG("read file ret:", ret)
}

