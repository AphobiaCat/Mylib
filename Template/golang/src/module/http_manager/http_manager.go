package http_manager

import(
	"io/ioutil"
	"net/http"
	"net/url"

	"mylib/src/public"
	"time"
)

var default_headers map[string]string
var default_headers_have_init bool = false

var req_timeout time.Duration


type HTTP_Manager struct{
	url		string
	headers map[string]string
	params	url.Values
	body_s	interface{}
	timeout	time.Duration
	is_get	bool
}

func Post(base_url string) *HTTP_Manager{
	ret := &HTTP_Manager{
		url:     base_url,
		timeout: req_timeout,
		is_get:  false,
		params:  url.Values{},
	}

	return ret
}

func Get(base_url string) *HTTP_Manager{
	ret := &HTTP_Manager{
		url:     base_url,
		timeout: req_timeout,
		is_get:  true,
		params:  url.Values{}, // 同上
	}
	return ret
}

func (hm *HTTP_Manager) Header(key string, val string) *HTTP_Manager{
	if hm.headers == nil{
		hm.headers = make(map[string]string)
	}

	hm.headers[key] = val

	return hm
}

func (hm *HTTP_Manager) Body(val interface{}) *HTTP_Manager{
	hm.body_s = val

	return hm
}

func (hm *HTTP_Manager) Param(key string, val string) *HTTP_Manager{
	hm.params.Add(key, val)
	
	return hm
}

func (hm *HTTP_Manager) Timeout(timeout_sec int) *HTTP_Manager{
	hm.timeout = time.Duration(timeout_sec * 1000 * 1000 * 1000)

	return hm
}

func (hm *HTTP_Manager) Send() string{

	var req *http.Request
	var err error

	if hm.is_get{
		req, err = http.NewRequest("GET", hm.url + "?" + hm.params.Encode(), nil)
		if err != nil {
			public.DBG_ERR("Error creating request:", err)
			return ""
		}
	}else{
		req, err = http.NewRequest("POST", hm.url, public.Build_Net_Json(hm.body_s))
		if err != nil {
			public.DBG_ERR("Error creating request:", err)
			return ""
		}
	}

	for key, val := range default_headers{
		req.Header.Set(key, val)
	}


	for key, val := range hm.headers{
		req.Header.Set(key, val)
	}

	client := &http.Client{
		Timeout: hm.timeout,
	}
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

func Set_Default_Headers(header_map map[string]string, keep_old ...bool)(old_config map[string]string){
	if !default_headers_have_init{
		default_headers = make(map[string]string)
	}

	//clear(default_headers) go 1.21 or later

	old_config = default_headers

	if len(keep_old) == 0{
		for key, _ := range default_headers{
			delete(default_headers, key)
		}
	}

	for key, val := range header_map{
		default_headers[key] = val
	}

	return old_config
}

func Set_Default_Timeout(timeout_sec int){
	req_timeout = time.Duration(timeout_sec * 1000 * 1000 * 1000)
}

func init(){
	req_timeout = time.Duration(30 * 1000 * 1000 * 1000)

	new_header := make(map[string]string)
	new_header["Accept"]		= "*/*"
	new_header["Content-Type"]	= "application/json"

	Set_Default_Headers(new_header)
}

