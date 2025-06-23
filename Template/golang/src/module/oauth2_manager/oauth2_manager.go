package oauth2_manager

import (
	"context"
	"io"
	"net/http"

	"golang.org/x/oauth2"

	"mylib/src/public"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID		: "YOUR_CLIENT_ID",
		ClientSecret	: "YOUR_CLIENT_SECRET",
		RedirectURL		: "http://localhost:8080/callback",		// must same with plateform config
		Scopes			: []string{"user:email"},				// according to plateform
		Endpoint		: oauth2.Endpoint{
			AuthURL		: "https://example.com/oauth/authorize",// target plateform
			TokenURL	: "https://example.com/oauth/token",	// target plateform
		},
	}

	oauthState			= "random-string"						// need a random string
)

type OAuth2_Manager struct{
	oauth_config	*oauth2.Config
	oauth_state		string
}

func New_OAuth2_App(bind_port string, client_id string, client_secret string, redirect_url string, scopes []string, plateform_auth_url string, plateform_token_url string, plateform_query_user_api string, callback func(string)string, redirect_to_your_web string)OAuth2_Manager{
	var tmp OAuth2_Manager

	tmp.oauth_config = &oauth2.Config{
		ClientID		: client_id,
		ClientSecret	: client_secret,
		RedirectURL		: redirect_url,			// must same with plateform config
		Scopes			: scopes,				// according to plateform
		Endpoint		: oauth2.Endpoint{
			AuthURL		: plateform_auth_url,	// target plateform
			TokenURL	: plateform_token_url,	// target plateform
		},
	}

	tmp.oauth_state = public.ConvertNumToHexStr(int64(public.Rand_U64()))

	go func(){
		handle_login	:= func(w http.ResponseWriter, r *http.Request) {
			url := tmp.oauth_config.AuthCodeURL(tmp.oauth_state)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}
	
		handle_callback := func(w http.ResponseWriter, r *http.Request){
			if r.FormValue("state") != tmp.oauth_state {
				http.Error(w, "State mismatch", http.StatusBadRequest)
				return
			}
	
			code := r.FormValue("code")
			token, err := tmp.oauth_config.Exchange(context.Background(), code)
			if err != nil {
				public.DBG_ERR(w, "Failed to exchange token: " + err.Error(), http.StatusInternalServerError)
				return
			}
	
			// use token query user info
			client := tmp.oauth_config.Client(context.Background(), token)
			resp, err := client.Get(plateform_query_user_api) // according plateform
			if err != nil {
				http.Error(w, "Failed to get user info: " + err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()
	
			body, _ := io.ReadAll(resp.Body)
			ret_token := callback(string(body))

			http.Redirect(w, r, redirect_to_your_web + "?token=" + ret_token, http.StatusTemporaryRedirect)
		}
	
		http.HandleFunc("/oauth2/login", handle_login)
		http.HandleFunc("/oauth2/callback", handle_callback)

		public.DBG_LOG("oauth2 app server start at 0.0.0.0:" + bind_port)
	
		http.ListenAndServe(":" + bind_port, nil)
	}()

	return tmp
}

