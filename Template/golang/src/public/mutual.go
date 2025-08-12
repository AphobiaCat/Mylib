package public

type callback func(string) (string, bool)

var Callbacks map[string]callback

func Reg(key string, user_callback callback){
	Callbacks[key] = user_callback
}

func init(){
	Callbacks = make(map[string]callback)
}

