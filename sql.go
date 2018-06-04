package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// root tree
type Sql struct {
	XMLName xml.Name      `xml:"sql"`
	Version string        `xml:"version,attr"`
	Info    SqlInfo       `xml:"info"`
	Sqls    SqlStatements `xml:"sqls"`
}

type SqlInfo struct {
	XMLName  xml.Name `xml:"info"`
	Db       string   `xml:"db"`
	Dbname   string   `xml:"dbname"`
	Dbuname  string   `xml:"dbuname"`
	Dbpassw  string   `xml:"dbpassw"`
	DbTables []Table  `xml:"table"`
}

type Table struct {
	Name    string `xml:"name"`
	Created string `xml:"create"`
}

type SqlStatements struct {
	S []SqlStatement `xml:"s"`
}

type SqlStatement struct {
	Name       string `xml:"name"`
	Code       string `xml:"code"`
	ResultName string `xml:"result,attr"`
}

func Conn(s *Sql) {
	// db, err := sql.Open("mysql", "jesse:jesse@tcp(127.0.0.1:3306)/?charset=utf8") //第一个参数为驱动名
	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

	//创建表
	// sql_table := `
	//    CREATE TABLE IF NOT EXISTS userinfo(
	//        uid INTEGER PRIMARY KEY AUTOINCREMENT,
	//        username VARCHAR(64) NULL,
	//        department VARCHAR(64) NULL,
	//        created DATE NULL
	//    );
	//    `
	for index, value := range s.Info.DbTables {
		fmt.Printf("arr[%d]=%d \n", index, value)
		db.Exec(value.Created)
	}

	//插入数据
	stmt, err := db.Prepare("INSERT INTO userinfo(username, department, created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	//更新数据
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created time.Time
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	//删除数据
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Parse(config string) (*Sql, error) {
	file, err := os.Open(config) // For read access.
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	v := Sql{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func connDb(sqlXml *Sql) (*sql.DB, error) {
	var engine *sql.Db
	var err error
	db := sqlXml.Info.Db
	dbName := sqlXml.Info.Dbname
	dbUname := sqlXml.Info.Dbuname
	dbpassw := sqlXml.Info.Dbpassw
	fmt.Printf("connDb: %s %s %s %s", db, dbName, dbUname, dbpassw)

	if "mysql" == db {
		conn := fmt.Sprintf("%s:%s@/%s?charset=utf8", dbUname, dbpassw, dbName)
		db, err = sql.Open("mysql", conn) //第一个参数为驱动名
	} else {
		db, err = sql.Open("sqlite3", dbName)
	}
	return engine, err
}

func createTables(sqlXml *Sql) {
	for index, value := range sqlXml.Info.DbTables {
		fmt.Printf("arr[%d]=%d \n", index, value)
		db.Exec(value.Created)
	}
}

func Exec(key string, args ...interface{}) {

}
func main() {
	v, err := Parse("sql.config") // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println(v)

	Conn(v)
	connDb(v)
}
