package http_m

import(
	"encoding/json"
	"bytes"
	"net/http"
	"time"
    "io/ioutil"
	//"reflect"
)

func POST(url string, body_data interface{}, header_map map[string]string) bool{

	jsonData, err := json.Marshal(body_data)
	if err != nil {
		DBG_ERR("Error marshalling JSON:", err)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		DBG_ERR("Error creating request:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	
	for key, value := range header_map{	
		req.Header.Set(key, value)
	}

	client := &http.Client{
    	Timeout: time.Second * 10,
    }

    resp, err := client.Do(req)
    if err != nil {
        DBG_ERR("Error sending request:", err)
        return false
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        DBG_LOG("Request successful!")
    } else {
        DBG_ERR("Request with status code:", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        DBG_ERR("Error reading response body:", err)
        return false
    }

	//DBG_LOG("Response Body[", string(body), "]")

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
        return false
	}

	err_code := dat["code"]
	err_code_ := err_code.(float64)

	if err_code_ != 0 {
		return false
	}

	return true
}

