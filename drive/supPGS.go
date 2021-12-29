package drive

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/brams/setup"
)

var createTable = `create table _table (
	 	name text not null,
		keys jsonb not null, 
		CONSTRAINT _table_pkey PRIMARY KEY (name)  
	);`
var err error

func InitPGS(set setup.PgDB) error {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		set.Host, set.User, set.Password, set.DBname)
	mdb, err = sql.Open("postgres", dbinfo)
	if err != nil {
		return fmt.Errorf("запрос на открытие %s %s", dbinfo, err.Error())
	}
	_, err = mdb.Exec("select * from _table;")
	if err != nil {
		logger.Info.Printf("Главной таблицы не существует создаем...")
		mdb.Exec(createTable)
	}
	return nil
}

func GetListDBs() []string {
	rows, err := mdb.Query("select name from _table;")
	if err != nil {
		return make([]string, 0)
	}
	defer rows.Close()
	result := make([]string, 0)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		result = append(result, name)
	}
	return result
}

func GetKeys(name string) []string {
	rows, err := mdb.Query(fmt.Sprintf("select keys from _table where name='%s';", name))
	if err != nil {
		return make([]string, 0)
	}
	defer rows.Close()
	for rows.Next() {
		var k []byte
		var df []string
		rows.Scan(&k)
		err = json.Unmarshal(k, &df)
		if err != nil {
			logger.Error.Print(err.Error())
			return make([]string, 0)
		}
		return df
	}
	return make([]string, 0)
}
func CreateDbPGS(name string, defkeys []string) error {
	return nil
}
