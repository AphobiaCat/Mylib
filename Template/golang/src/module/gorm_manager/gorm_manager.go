package gorm_manager

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"


    "mylib/src/public"
)

var gorm_manager Gorm_Manager

type Gorm_Manager struct{
	dsn		string	
	db		*gorm.DB
}


func (gm *Gorm_Manager) Init(dsn string, models ...interface{}){

	gm.dsn = dsn

	var err error

	//	connect
	gm.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Warn, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      false,       // Don't include params in the SQL log
				Colorful:                  true,        // Enable color
			},
		),
	})
	if err != nil {
		panic("failed to connect database")
	}
	// gm.db = gm.db.Debug()
	// aoto migrate, according struct create/update tables

	//gm.db.Migrator().DropTable(models...)	//for test

    if err := gm.db.AutoMigrate(models...); err != nil {
        panic(err)
    }

	//DBG_LOG_VAR(gm.db)
	
}

func (gm *Gorm_Manager) Sellect_All(all_data interface{}, conds ...interface{}) error{
	result := gm.db.Find(all_data, conds...)
	if result.Error != nil {
		public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) Create(new_item interface{}) error{
	//user := User{Name: "John Doe", Email: "john@example.com", Password: "secret"}

	result := gm.db.Create(new_item)
	if result.Error != nil {
		public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) Fetch(fetched_data interface{}, key interface{}) error{
    result := gm.db.First(fetched_data, key)
    if result.Error != nil {
		public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) Foreign(fetched_data interface{}, key interface{}, foreign_volume string) error{
	result := gm.db.Preload(foreign_volume).First(fetched_data, key)

	if result.Error != nil {
		public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) Fetch_Where(fetched_data interface{}, where_data interface{}) error{
    result := gm.db.Where(where_data).Find(fetched_data)

    if result.Error != nil {
		public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) Foreign_Where(fetched_data interface{}, where_data interface{}, foreign_volume []string, conds ...interface{}) error{
	query := gm.db
	for _, volume := range foreign_volume {
	    query = query.Preload(volume)
	}
	result := query.Where(where_data).Find(fetched_data, conds...)

	if result.Error != nil {
		public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) Update(data interface{}, volume string, new_data interface{}) error{
    result := gm.db.Model(data).Update(volume, new_data)

    if result.Error != nil {
		public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) Updates(data interface{}, new_data interface{}) error{
	result := gm.db.Model(data).Select("*").Updates(new_data)

	if result.Error != nil {
	    public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) Delete(data interface{}) error{
    result := gm.db.Delete(data)

    if result.Error != nil {
		public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) SQL_Query(fetched_data interface{}, query string, datas ...interface{}) error{
    result := gm.db.Raw(query, datas...).Scan(fetched_data)

    if result.Error != nil {
		public.DBG_ERR("Error:", result.Error)
	}
	return result.Error
}

func (gm *Gorm_Manager) Exec(sql string, args ...interface{}) error {
	result := gm.db.Exec(sql, args...)
	if result.Error != nil {
		public.DBG_ERR("ExecRaw Error:", result.Error)
	}
	return result.Error
}

//------------------------------API---------------------------------

func Init_Gorm(v ...interface{}){
	dsn, exist := public.Global[string]("dsn")
	
	if exist{
		gorm_manager.Init(dsn, v...)
	}else{
		panic("database no init")
	}
}

func Gorm_Sellect_All(all_data interface{}, conds ...interface{}) error{
	return gorm_manager.Sellect_All(all_data, conds...)
}

func Gorm_Create(new_item interface{}) error{
	return gorm_manager.Create(new_item)
}

func Gorm_Fetch(fetched_data interface{}, key interface{}) error{
	return gorm_manager.Fetch(fetched_data, key)
}

func Gorm_Foreign(fetched_data interface{}, key interface{}, foreign_volume string) error{
	return gorm_manager.Foreign(fetched_data, key, foreign_volume)
}

func Gorm_Fetch_Where(fetched_data interface{}, where_data interface{}) error{
	return gorm_manager.Fetch_Where(fetched_data, where_data)
}

func Gorm_Foreign_Where(fetched_data interface{}, where_data interface{}, foreign_volume []string, conds ...interface{}) error{
	return gorm_manager.Foreign_Where(fetched_data, where_data, foreign_volume, conds...)
}

func Gorm_Update(data interface{}, volume string, new_data interface{}) error{
	return gorm_manager.Update(data, volume, new_data)
}

func Gorm_Updates(data interface{}, new_data interface{}) error{
    return gorm_manager.Updates(data, new_data)
}

func Gorm_Delete(data interface{}) error{
	return gorm_manager.Delete(data)
}

func Gorm_SQL_Query(fetched_data interface{}, query string, datas ...interface{}) error{
	return gorm_manager.SQL_Query(fetched_data, query, datas...)
}

func Gorm_SQL_Exec(sql string, args ...interface{}) error {
	return gorm_manager.Exec(sql, args...)
}

func Get_DB() *gorm.DB {
	return gorm_manager.db
}

func init(){

	dsn := public.Config.Database.User + ":" + public.Config.Database.Password + "@tcp(" + public.Config.Database.Ip + ")/" + public.Config.Database.Name + "?charset=utf8mb4&parseTime=True&loc=Local"

	//public.DBG_LOG("database dsn:", dsn)

	public.Set_Global("dsn", dsn)
}

