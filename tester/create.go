package tester

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/brams/drive"
)

type DataTest struct {
	One   string  `json:"one"`
	Two   int     `json:"two"`
	More  float32 `json:"more"`
	Bytes []byte  `json:"bytes"`
	Bool  bool    `json:"bool"`
}

func CreateDb() {
	drive.Drop("test")
	err := drive.CreateDb("test", "one", "two", "more", "bytes", "bool")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	err = drive.CreateDbInMemory("mema", "one", "two", "more", "bytes", "bool")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	err = drive.CreateDbInMemory("memb", "one", "two", "more", "bytes", "bool")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	err = drive.CreateDbInMemory("memd", "one", "two", "more", "bytes", "bool")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	err = drive.CreateDbInMemory("memc", "one", "two", "more", "bytes", "bool")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	db, err := drive.Open("test")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	dbm, err := drive.Open("mema")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	dbb, err := drive.Open("memb")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	dbd, err := drive.Open("memd")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	dbc, err := drive.Open("memc")
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	defer db.Close()
	defer dbm.Close()
	defer dbb.Close()
	defer dbd.Close()
	defer dbc.Close()
	go pusher(dbm, 10.1, rand.Intn(30000))
	go pusher(dbb, 12345.5555, rand.Intn(30000))
	go pusher(dbc, 22675.5555, rand.Intn(30000))
	go pusher(dbd, 11322.5555, rand.Intn(30000))

	var t DataTest
	more := 10.4
	two := 1234
	i := 0
	for ; i < 40; i++ {

		t.More = float32(more)
		t.Two = two
		t.One = fmt.Sprintf("%d", i%10)
		t.Bool = rand.Intn(2)%2 == 0
		t.Bytes = []byte(t.One)

		buf, _ := json.Marshal(t)
		err = db.WriteRecord(buf)
		if err != nil {
			logger.Error.Println(err.Error())
			return
		}
		more += 10.23
		two += 7
	}
}

func pusher(db *drive.Db, more float32, two int) {
	for {
		var t DataTest
		i := 0
		mmore := more
		mtwo := two
		start := time.Now()
		for ; i < 1000000; i++ {

			t.More = float32(mmore)
			t.Two = mtwo
			t.One = fmt.Sprintf("%d", i%10)
			t.Bool = rand.Intn(2)%2 == 0
			t.Bytes = []byte(t.One)

			buf, _ := json.Marshal(t)
			err := db.WriteRecord(buf)
			if err != nil {
				logger.Error.Println(err.Error())
				return
			}
			mmore += 10.23
			mtwo += 7
		}
		logger.Info.Printf("push %s\t%d\t%v", db.Name, i, time.Since(start))
		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)
	}
}
