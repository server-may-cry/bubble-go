package main

func main() {}

/*
import (
	"fmt"
	"time"

	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
)

type pgUser struct {
	tableName               struct{} `sql:"users"`
	Id                      int64
	SysId                   int   // vk
	ExtId                   int64 // 1234
	ReachedStage01          int32
	ReachedSubStage01       int32
	IgnoreSavePointBlock    int
	InifinityExtra00        int
	InifinityExtra01        int
	InifinityExtra02        int
	InifinityExtra03        int
	InifinityExtra04        int
	InifinityExtra05        int
	InifinityExtra06        int
	InifinityExtra07        int
	InifinityExtra08        int
	InifinityExtra09        int
	RemainingTries          int
	RestoreTriesAt          int64 // timestamp
	Credits                 int16
	FriendsBonusCreditsTime int64       // timestamp
	ProgressStandart        interface{} // json
}

type moUser struct {
	Id                      uint64
	SysId                   uint8  // 1 = vk
	ExtId                   uint64 // 1234
	ReachedStage01          uint8
	ReachedSubStage01       uint8
	IgnoreSavePointBlock    uint8 // as int
	InfinityExtra00         uint8 // as int
	InfinityExtra01         uint8 // as int
	InfinityExtra02         uint8 // as int
	InfinityExtra03         uint8 // as int
	InfinityExtra04         uint8 // as int
	InfinityExtra05         uint8 // as int
	InfinityExtra06         uint8 // as int
	InfinityExtra07         uint8 // as int
	InfinityExtra08         uint8 // as int
	InfinityExtra09         uint8 // as int
	RemainingTries          uint8
	RestoreTriesAt          time.Time
	Credits                 uint16
	FriendsBonusCreditsTime time.Time
	ProgressStandart        interface{} // [][]int8 // json -1 0 1 2 3 stars count
}

func main() {
	session, err := mgo.Dial("url")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	err = session.DB("dbname").DropDatabase()
	if err != nil {
		panic(err)
	}

	// Collection users
	users := session.DB("heroku_cqp2d1s6").C("users")

	// Index
	index := mgo.Index{
		Key:        []string{"sysId", "extId"},
		Unique:     true,
		DropDups:   true,
		Background: false,
		Sparse:     true,
	}

	err = users.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	db := pg.Connect(&pg.Options{
		Addr:     "",
		User:     "",
		Password: "",
		Database: "",
	})

	res := []pgUser{}
	before := time.Now()
	err = db.Model(&pgUser{}).Select(&res)
	pass := time.Since(before)
	fmt.Println(pass)
	if err != nil {
		panic(err)
	}
	before2 := time.Now()
	for _, usr := range res {
		// fmt.Printf("%+v \n", usr)
		newUser := moUser{
			Id:                      uint64(usr.Id),
			SysId:                   uint8(usr.SysId),
			ExtId:                   uint64(usr.ExtId),
			ReachedStage01:          uint8(usr.ReachedStage01),
			ReachedSubStage01:       uint8(usr.ReachedSubStage01),
			IgnoreSavePointBlock:    uint8(usr.IgnoreSavePointBlock),
			InfinityExtra00:         uint8(usr.InifinityExtra00),
			InfinityExtra01:         uint8(usr.InifinityExtra01),
			InfinityExtra02:         uint8(usr.InifinityExtra02),
			InfinityExtra03:         uint8(usr.InifinityExtra03),
			InfinityExtra04:         uint8(usr.InifinityExtra04),
			InfinityExtra05:         uint8(usr.InifinityExtra05),
			InfinityExtra06:         uint8(usr.InifinityExtra06),
			InfinityExtra07:         uint8(usr.InifinityExtra07),
			InfinityExtra08:         uint8(usr.InifinityExtra08),
			InfinityExtra09:         uint8(usr.InifinityExtra09),
			RemainingTries:          uint8(usr.RemainingTries),
			RestoreTriesAt:          time.Unix(usr.RestoreTriesAt, 0),
			Credits:                 uint16(usr.Credits),
			FriendsBonusCreditsTime: time.Unix(usr.FriendsBonusCreditsTime, 0),
			ProgressStandart:        usr.ProgressStandart,
		}
		err = users.Insert(newUser)
		if err != nil {
			panic(err)
		}
	}
	count, err := users.Count()
	if err != nil {
		panic(err)
	}
	fmt.Println("Count", count)

	fmt.Println("ok")
	pass2 := time.Since(before2)
	fmt.Println(pass2)
}
*/
