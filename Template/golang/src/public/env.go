package public
//test
var Database_Name string = "test"
var User_Name string = "Dunty"
var User_Passwd string = "123123Dunty"
var Database_Ip string = "192.168.136.130:3306"
var Dsn_Str string = User_Name + ":" + User_Passwd + "@tcp(" + Database_Ip + ")/" + Database_Name + "?charset=utf8mb4&parseTime=True&loc=Local"

var Redis_Server_Ip		string = "192.168.136.130:6379"
var Redis_Server_Passwd	string = "123123Dunty"
var Redis_DB			int    = 1


