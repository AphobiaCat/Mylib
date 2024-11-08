package example

import(
	"mylib/src/module/gorm_manager"
	"mylib/src/module/cachesql_manager"
	"mylib/src/public"
)

type UsrInfo struct{
	Name	string	`json:"name" gorm:"primaryKey"`
	Age		int		`json:"age" gorm:"age"`
	Email	string	`json:"email" gorm:"email"`
}

func Example_Cachesql(){

	gorm_manager.Init_Gorm(public.Dsn_Str, &UsrInfo{})

	gorm_manager.Gorm_Create(&UsrInfo{Name:"Dunty", Age:25, Email:"Dunty@gmail.com"})


	key := "UsrInfo_Dunty"

	usr_data := cachesql_manager.Get_Cache(key, func()interface{}{
		var sql_data UsrInfo
		
		gorm_manager.Gorm_Fetch_Where(&sql_data, &UsrInfo{Name:"Dunty"})

		return sql_data
	}, 10, 60, 20)

	var usr_info UsrInfo
	public.Parser_Json(usr_data, &usr_info)
	public.DBG_LOG(usr_info)




	usr_info.Age = 26
	cachesql_manager.Set_Cache(key, usr_info, 10, 60, 20)	
	usr_data = cachesql_manager.Get_Cache(key, func()interface{}{
		var sql_data UsrInfo
		
		gorm_manager.Gorm_Fetch_Where(&sql_data, &UsrInfo{Name:"Dunty"})

		return sql_data
	}, 10, 60, 20)

	public.Parser_Json(usr_data, &usr_info)
	public.DBG_LOG(usr_info)



	for ;;{
		usr_data := cachesql_manager.Get_Cache(key, func()interface{}{
			var sql_data UsrInfo
			
			gorm_manager.Gorm_Fetch_Where(&sql_data, &UsrInfo{Name:"Dunty"})

			return sql_data
		}, 10, 60, 20)

		var usr_info UsrInfo
		public.Parser_Json(usr_data, &usr_info)
		public.DBG_LOG(usr_info)

		public.Sleep(1000)
	}
}

