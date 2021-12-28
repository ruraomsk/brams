package tester

import (
	"math/rand"
	"time"

	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/brams/drive"
)

func Update(name string, first string, last bool) {
	db, err := drive.Open(name)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	defer db.Close()
	for {
		var list []string
		start := time.Now()
		list, err = db.ReadListKeys(0, first, nil, nil, nil, last)
		if err != nil {
			logger.Error.Println(err.Error())
			return
		}
		for _, key := range list {
			buf, err := db.ReadRecordFromList(key)
			if err != nil {
				logger.Error.Println(err.Error())
				return
			}
			err = db.WriteRecord(buf)
			if err != nil {
				logger.Error.Println(err.Error())
				return
			}
		}
		logger.Info.Printf("get %s\t: %d\t%v", name, len(list), time.Since(start))
		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)
	}
}
