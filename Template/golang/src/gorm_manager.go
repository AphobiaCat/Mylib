package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
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
    gm.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // aoto migrate, according struct create/update tables

	//gm.db.Migrator().DropTable(models...)	//for test

    if err := gm.db.AutoMigrate(models...); err != nil {
        panic(err)
    }

	//DBG_LOG_VAR(gm.db)
	
}

func (gm *Gorm_Manager) Create(new_item interface{}){
	//user := User{Name: "John Doe", Email: "john@example.com", Password: "secret"}

	result := gm.db.Create(new_item)
	if result.Error != nil {
		DBG_ERR("Error:", result.Error)
	}
}

func (gm *Gorm_Manager) Fetch(fetched_data interface{}, key interface{}){
    gm.db.First(fetched_data, key)
}

func (gm *Gorm_Manager) Foreign(fetched_data interface{}, key interface{}, foreign_volume string){
	gm.db.Preload(foreign_volume).First(fetched_data, key)
}

func (gm *Gorm_Manager) Fetch_Where(fetched_data interface{}, where_data interface{}){
    gm.db.Where(where_data).Find(fetched_data)
}

func (gm *Gorm_Manager) Foreign_Where(fetched_data interface{}, where_data interface{}, foreign_volume string){
	gm.db.Preload(foreign_volume).Where(where_data).Find(fetched_data)
}

func (gm *Gorm_Manager) Update(data interface{}, volume string, new_data interface{}){
    gm.db.Model(data).Update(volume, new_data)
}

func (gm *Gorm_Manager) Delete(data interface{}){
    gm.db.Delete(data)
}



//------------------------------API---------------------------------

func Init_Gorm(dsn string, v ...interface{}){
	gorm_manager.Init(dsn, v...)
}

func Gorm_Create(new_item interface{}){
	gorm_manager.Create(new_item)
}

func Gorm_Fetch(fetched_data interface{}, key interface{}){
	gorm_manager.Fetch(fetched_data, key)
}

func Gorm_Foreign(fetched_data interface{}, key interface{}, foreign_volume string){
	gorm_manager.Foreign(fetched_data, key, foreign_volume)
}

func Gorm_Fetch_Where(fetched_data interface{}, where_data interface{}){
	gorm_manager.Fetch_Where(fetched_data, where_data)
}

func Gorm_Foreign_Where(fetched_data interface{}, where_data interface{}, foreign_volume string){
	gorm_manager.Foreign_Where(fetched_data, where_data, foreign_volume)
}

func Gorm_Update(data interface{}, volume string, new_data interface{}){
	gorm_manager.Update(data, volume, new_data)
}

func Gorm_Delete(data interface{}){
	gorm_manager.Delete(data)
}



