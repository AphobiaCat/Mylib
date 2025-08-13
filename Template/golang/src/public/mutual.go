package public

type callback func(interface{}) (interface{}, bool)

var callbacks map[string]callback

func Reg(key string, user_callback callback){
	callbacks[key] = user_callback
}

func Call(key string, payload interface{})(interface{}, bool){
	if call, exist := callbacks[key]; exist{
		return call(payload)
	}else{

		DBG_ERR("now support[", callbacks, "]")
	
		return "callback[" + key + "] no exist.", false
	}
}

func init(){
	callbacks = make(map[string]callback)
}

