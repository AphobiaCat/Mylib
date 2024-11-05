package http_manager

import(
	"io/ioutil"
	"net/http"
	"net/url"

	"mylib/src/public"
)

var default_headers map[string]string
var default_headers_have_init bool = false

func Post(base_url string, body_data interface{}, header_map... map[string]string) string{

	req, err := http.NewRequest("POST", base_url, public.Build_Net_Jason(body_data))
	if err != nil {
		public.DBG_ERR("Error creating request:", err)
		return ""
	}

	for key, val := range default_headers{
		req.Header.Set(key, val)
	}

	if len(header_map) == 1{
		for key, val := range header_map[0]{
			req.Header.Set(key, val)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		public.DBG_ERR("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		public.DBG_ERR("Error reading response:", err)
		return ""
	}
	
	return string(body)
}

func Get(base_url string, params... map[string]string) string{
	if len(params) > 2{
		public.DBG_ERR("one is params, two is headers")
	}

	req_params := url.Values{}

	params_headers_flag := 0
	switch len(params){
		case 1: //only params
			params_headers_flag = 1
		case 2: //params and headers
			params_headers_flag = 3
	}

	if params_headers_flag & 0x1 == 1{
		for key, val := range params[0]{
			req_params.Add(key, val)
		}
	}
	
	req, err := http.NewRequest("GET", base_url + "?" + req_params.Encode(), nil)
	if err != nil {
		public.DBG_ERR("Error creating request:", err)
		return ""
	}

	for key, val := range default_headers{
		req.Header.Set(key, val)
	}

	if params_headers_flag & 0x2 == 0x2{
		for key, val := range params[1]{	//usr headers
			req.Header.Set(key, val)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		public.DBG_ERR("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		public.DBG_ERR("Error reading response:", err)
		return ""
	}
	
	return string(body)
}

func Set_Default_Headers(header_map map[string]string){
	if !default_headers_have_init{
		default_headers = make(map[string]string)
	}

	//clear(default_headers) go 1.21 or later

	for key, _ := range default_headers{
		delete(default_headers, key)
	}

	for key, val := range header_map{
		default_headers[key] = val
	}
}

