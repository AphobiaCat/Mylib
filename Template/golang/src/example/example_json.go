package example

import(
	"mylib/src/public"
)

func Example_Json_Op(){
	obj := map[string]interface{}{
		"test1" : 1,
		"test2" : 2,
	}

	ret := public.Build_Net_Json(obj)
	public.DBG_LOG(ret)

	type Test_Data struct{
		Test1	int `json:"test"`
		Test2	int	`json:"test2"`
	}

	var res Test_Data

	public.Parser_Json("{\"test\":1, \"test2\":2}", &res)

	
	public.DBG_LOG(res)
	public.DBG_LOG(res.Test1)
	public.DBG_LOG(res.Test2)	
}