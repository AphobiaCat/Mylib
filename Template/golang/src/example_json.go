package main



func example_json_op(){
	obj := map[string]interface{}{
		"test1" : 1,
		"test2" : 2,
	}

	ret := Build_Net_Jason(obj)
	DBG_LOG(ret)

	type Test_Data struct{
		Test1	int `json:"test"`
		Test2	int	`json:"test2"`
	}

	var res Test_Data

	Parser_Jason("{\"test\":1, \"test2\":2}", &res)

	
	DBG_LOG(res)
	DBG_LOG(res.Test1)
	DBG_LOG(res.Test2)	
}