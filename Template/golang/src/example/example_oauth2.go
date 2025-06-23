package example


import(
	"mylib/src/public"
	"mylib/src/module/oauth2_manager"
)

func github_user_info(user_info string)string{
	public.DBG_LOG("get user_info[", user_info, "]")
	public.DBG_LOG("return a token to front web. you can encode user_info here")

	return user_info
}

func Example_OAuth2(){
	oauth2_server := oauth2_manager.New_OAuth2_App(
		"7001", 										// bind port
		"ClientID", 									// your client id
		"ClientSecret", 								// your client secret
		"http://localhost:7001/oauth2/callback", 		// login callback must same with your plateform config
		[]string{"user:email"}, 						// scops, detail in plateform
		"https://github.com/login/oauth/authorize", 	// detail in target plateform doc
		"https://github.com/login/oauth/access_token",	// detail in target plateform doc
		"https://api.github.com/user", 					// target plateform query user api
		github_user_info,								// process and return token to your front |
		"https://www.google.com",						// <--------------------------------------|
	)
	
	public.DBG_LOG("server[", oauth2_server, "]")
	public.DBG_LOG("use http://127.0.0.1:7001/oauth2/login to login")

	for ;;{
		public.Sleep(1000 * 60)
	}
}

