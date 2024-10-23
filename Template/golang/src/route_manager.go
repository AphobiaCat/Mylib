package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(auth_password string) gin.HandlerFunc {
    return func(c *gin.Context) {

        authHeader := c.GetHeader("auth")

        if authHeader != auth_password {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort() // 阻止请求继续传递
            return
        }

        c.Next()
    }
}


type Login_Data struct {
    Uid			string `json:"uid"`
    Timestamp	uint32 `json:"timestamp"`
    Invite_from	string `json:"invite_from"`
    Name		string `json:"name"`
    Is_premium	bool   `json:"is_premium"`
}

func Test_Post(context *gin.Context) {

	var input Login_Data
	
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":-1,
		})
		return
	}

	DBG_LOG(input)

	context.JSON(http.StatusOK, gin.H{
		"code":0,
	})
}

func Test_Get(context *gin.Context) {

	if key, exists := context.GetQuery("key"); exists {

		if key == "1007036321_32722527"{
			DBG_LOG("haha")
		}
	
	    context.JSON(http.StatusOK, gin.H{
	        "key": key,
	    })
	} else {
	    context.JSON(http.StatusOK, gin.H{
	        "error": "Parameter 'key' not found",
	    })
	}
}



func Start_Http_Route() {
	r := gin.New()
	
	//r.Use(authMiddleware(token))

	r.POST("/test_post", AuthMiddleware("hello_world"), Test_Post)
	r.GET("/test_get", Test_Get)

	if err := r.Run("0.0.0.0:" + "7992"); err != nil {
		panic(err)
	}
}

func Init_Htpp_Server(){
	go Start_Http_Route()
}


