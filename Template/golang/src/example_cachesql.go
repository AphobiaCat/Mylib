package main


type UsrInfo_GORM_Data struct {
	Name	string	`gorm:"primaryKey"`
	Age		int		`gorm:"age"`
	Email	string	`gorm:"email"`
}

type UsrInfo struct{
	Name	string	`json:"name"`
	Age		int		`json:"age"`
	Email	string	`json:"email"`
}

func example_cachesql(){

	Init_Gorm(dsn_str, &UsrInfo_GORM_Data{})

	Gorm_Create(&UsrInfo_GORM_Data{Name:"Dunty", Age:25, Email:"Dunty@gmail.com"})


	Init_Cache_Sql()

	key := "UsrInfo_Dunty"

	usr_data := Get_Cache(key, func()interface{}{
		var sql_data UsrInfo_GORM_Data
		
		Gorm_Fetch_Where(&sql_data, &UsrInfo_GORM_Data{Name:"Dunty"})

		return sql_data
	}, 10, 60, 20)

	var usr_info UsrInfo
	Parser_Jason(usr_data, &usr_info)
	DBG_LOG(usr_info)




	usr_info.Age = 26
	Set_Cache(key, usr_info, 10, 60, 20)	
	usr_data = Get_Cache(key, func()interface{}{
		var sql_data UsrInfo_GORM_Data
		
		Gorm_Fetch_Where(&sql_data, &UsrInfo_GORM_Data{Name:"Dunty"})

		return sql_data
	}, 10, 60, 20)

	Parser_Jason(usr_data, &usr_info)
	DBG_LOG(usr_info)



	for ;;{
		usr_data := Get_Cache(key, func()interface{}{
			var sql_data UsrInfo_GORM_Data
			
			Gorm_Fetch_Where(&sql_data, &UsrInfo_GORM_Data{Name:"Dunty"})

			return sql_data
		}, 10, 60, 20)

		var usr_info UsrInfo
		Parser_Jason(usr_data, &usr_info)
		DBG_LOG(usr_info)

		sleep(1000)
	}
}

