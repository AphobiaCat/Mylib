package main

var database_name string = "test"
var user_name string = "root"
var user_passwd string = "123123"
var database_ip string = "127.0.0.1:3306"
var dsn_str string = user_name + ":" + user_passwd + "@tcp(" + database_ip + ")/" + database_name + "?charset=utf8mb4&parseTime=True&loc=Local"

var redis_server_ip		string = "127.0.0.1:6379"
var redis_server_passwd	string = "123123"
var redis_DB			int    = 0



