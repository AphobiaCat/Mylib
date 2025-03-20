package route_manager

import(
	"mylib/src/public"

	"github.com/golang-jwt/jwt/v5"
	
	"time"
)

var jwtKey = []byte("my_secret_key")

func Route_Generate_Jwt(data interface{}, expiration_time_s_option ...int)(string, bool){
	expiration_time := 60 * 60

	if len(expiration_time_s_option) != 0 && expiration_time_s_option[0] > 0{
		expiration_time = expiration_time_s_option[0]
	}

	claims := jwt.MapClaims{
		"data"	:	public.Build_Json(data),
		"exp"	:   time.Now().Add(time.Duration(expiration_time) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ret, err := token.SignedString(jwtKey)

	if err != nil{
		public.DBG_ERR("generate jwt failed. user data:", data)
		return ret, false
	}
	
	return ret, true
}

func Route_Parser_Jwt(tokenString string)(string, bool){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		public.DBG_ERR("parser jwt failed. err:", err)
	
		return "", false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		ret, exist := claims["data"]

		if !exist{
			public.DBG_ERR("jwt haven't data")
			return "", false
		}
	
		return ret.(string), true
	}

	public.DBG_ERR("parser jwt failed. err:", err)

	return "", false
}

func Route_Get_Jwt_Mid(headers map[string]string)(map[string]string, bool){

	token := headers["token"]
	data, succ := Route_Parser_Jwt(token)

	if succ{
		ret_data := make(map[string]string)
		ret_data["data"] = data
		return ret_data, true
	}

	return map[string]string{}, false
}


