package main

var database_name string = "test"
var user_name string = "Dunty"
var user_passwd string = "123123Dunty"
var database_ip string = "192.168.136.130:3306"
var dsn_str string = user_name + ":" + user_passwd + "@tcp(" + database_ip + ")/" + database_name + "?charset=utf8mb4&parseTime=True&loc=Local"

var redis_server_ip		string = "192.168.136.130:6379"
var redis_server_passwd	string = "123123Dunty"
var redis_DB			int    = 0



