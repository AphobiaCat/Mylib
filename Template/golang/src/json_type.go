package main

import(
	"encoding/json"
	"bytes"
	//"reflect"
)

/*
	obj := map[string]interface{}{
		"test1" : 1,
		"test2" : 2,
	}
	ret := Build_Net_Jason(obj)
*/

func Build_Jason(obj interface{}) string{

    jsonData, err := json.Marshal(obj)
	if err != nil {
		DBG_ERR("Error marshalling JSON:", err)
	}
	
	return string(jsonData)
}


func Build_Net_Jason(obj interface{}) *bytes.Buffer{

    jsonData, err := json.Marshal(obj)
	if err != nil {
		DBG_ERR("Error marshalling JSON:", err)
	}

	//DBG_LOG(jsonData)

	return bytes.NewBuffer(jsonData)
}

func Parser_Jason(message string, res interface{}) {

	/*
	type Parser_Jason_Test_Data struct{
		Code	int		`json:"code"`
		Encode	uint32	`json:"encode"`
	}
	*/
		
	json.Unmarshal([]byte(message), res)
}

