package example



import (
	"time"

	"mylib/src/module/gorm_manager"
	"mylib/src/public"
)


type Test_GORM_Data struct {
	ID       	uint   `gorm:"primaryKey;autoIncrement"`
	Name     	string `gorm:"size:100"`
	Email    	string `gorm:"size:100;uniqueIndex"`
	Password 	string `gorm:"size:100"`
	CreatedAt	time.Time 	// auto add create time
	UpdatedAt	time.Time	// auto add update time
}

func example_gorm(){
	gorm_manager.Init_Gorm(public.Dsn_Str, &Test_GORM_Data{})

	tgd		:= Test_GORM_Data{Name:"Dunty", Email:"Dunty@gmail.com", Password:"123"}
	tgd2	:= Test_GORM_Data{Name:"Cirila", Email:"Cirila@Gmail.com", Password:"123"}
	

	gorm_manager.Gorm_Create(&tgd)
	gorm_manager.Gorm_Create(&tgd2)


	var tgd3 Test_GORM_Data

	gorm_manager.Gorm_Fetch(&tgd3, 2)

	public.DBG_LOG(tgd3)
	
	gorm_manager.Gorm_Update(&tgd3, "Email", "123123123")

	tgd3.Name = "Dunty"

	gorm_manager.Gorm_Updates(&tgd3, &tgd3)
	
	gorm_manager.Gorm_Delete(&tgd)
}


//========================================================================================


type User struct {
    ID       string    `gorm:"primaryKey"`
    Name     string    `gorm:"size:100"`
    Email    string    `gorm:"size:100;uniqueIndex"`
    Comments []Comment `gorm:"foreignKey:UserID"`

    CreatedAt	time.Time 	// auto add create time
	UpdatedAt	time.Time	// auto add update time
}

type Comment struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"`
    Content   string    `gorm:"type:text"`
    UserID    string	`gorm:"index"` 					// foreignKey, connect user
    ParentID  *uint     `gorm:"index"` 					// self index, use *uint because this can be null
    Mint      string    `gorm:"size:100;index"` 		// use to query
    Replies   []Comment `gorm:"foreignKey:ParentID"` 	// sub comment self connect

    CreatedAt	time.Time 	// auto add create time
	UpdatedAt	time.Time	// auto add update time
}


func Example_2_Gorm(){
	//if wanna retest , should add "gm.db.Migrator().DropTable(models...)" to gorm_manager.go before gm.db.AutoMigrate(models...)

	gorm_manager.Init_Gorm(public.Dsn_Str, &User{}, &Comment{})

	user1 := User{ID:"123456", Name:"Cirila", Email:"Cirila@gmail.com"}	
	user2 := User{ID:"789101", Name:"Dunty"	, Email:"Dunty@gmail.com"}
	user3 := User{ID:"112131", Name:"World"	, Email:"World@gmail.com"}
	
	gorm_manager.Gorm_Create(&user1)
	gorm_manager.Gorm_Create(&user2)
	gorm_manager.Gorm_Create(&user3)

	var tmp_num uint
	tmp_num = 1
	var tmp_num2 uint
	tmp_num2 = 5
	var tmp_num3 uint
	tmp_num3 = 7
	

	public.DBG_LOG("create comment")

	comment  := Comment{Content:"Hello World1", UserID:"123456", Mint:"0x1"}
	comment1 := Comment{Content:"Hello World2", UserID:"112131", Mint:"0x1"		, ParentID:&tmp_num}
	comment2 := Comment{Content:"Hello World3", UserID:"112131", Mint:"0x1"}
	comment3 := Comment{Content:"Hello World4", UserID:"789101", Mint:"0x1"}
	comment4 := Comment{Content:"Hello World5", UserID:"123456", Mint:"0x1"		, ParentID:&tmp_num}
	comment5 := Comment{Content:"Hello World6", UserID:"789101", Mint:"0x1"}
	comment6 := Comment{Content:"Hello World7", UserID:"112131", Mint:"0x1"		, ParentID:&tmp_num2}
	comment7 := Comment{Content:"Hello World8", UserID:"123456", Mint:"0x1"}
	comment8 := Comment{Content:"Hello World9", UserID:"789101", Mint:"0x1"}
	comment9 := Comment{Content:"Hello World10", UserID:"112131", Mint:"0x1"		, ParentID:&tmp_num3}

	gorm_manager.Gorm_Create(&comment)	
	gorm_manager.Gorm_Create(&comment1)
	gorm_manager.Gorm_Create(&comment2)
	gorm_manager.Gorm_Create(&comment3)
	gorm_manager.Gorm_Create(&comment4)
	gorm_manager.Gorm_Create(&comment5)
	gorm_manager.Gorm_Create(&comment6)
	gorm_manager.Gorm_Create(&comment7)
	gorm_manager.Gorm_Create(&comment8)
	gorm_manager.Gorm_Create(&comment9)

	public.DBG_LOG("query")

	var comment_of_mint []Comment
	gorm_manager.Gorm_Fetch_Where(&comment_of_mint, &Comment{Mint:"0x1"})
	public.DBG_LOG(comment_of_mint)

	var comment_of_user User
	gorm_manager.Gorm_Foreign_Where(&comment_of_user, &User{Name:"Dunty"}, []string{"Comments"})
	public.DBG_LOG(comment_of_user)

	var foregin_use User
	gorm_manager.Gorm_Foreign(&foregin_use, "789101", "Comments")
	public.DBG_LOG(foregin_use)

	var comment_and_sub_comment_of_mint []Comment
	gorm_manager.Gorm_Foreign_Where(&comment_and_sub_comment_of_mint, &Comment{Mint:"0x1"}, []string{"Replies"})
	public.DBG_LOG(comment_and_sub_comment_of_mint)
}


