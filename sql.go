package gosql

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/freelifer/gosql/parser"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"reflect"
	"time"
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
	var engine *sql.DB
	var err error
	db := sqlXml.Info.Db
	dbName := sqlXml.Info.Dbname
	dbUname := sqlXml.Info.Dbuname
	dbpassw := sqlXml.Info.Dbpassw
	fmt.Printf("connDb: %s %s %s %s", db, dbName, dbUname, dbpassw)

	if "mysql" == db {
		conn := fmt.Sprintf("%s:%s@/%s?charset=utf8", dbUname, dbpassw, dbName)
		engine, err = sql.Open("mysql", conn) //第一个参数为驱动名
	} else {
		engine, err = sql.Open("sqlite3", dbName)
	}
	return engine, err
}

func createTables(sqlXml *Sql) {
	// for index, value := range sqlXml.Info.DbTables {
	// 	fmt.Printf("arr[%d]=%d \n", index, value)
	// 	db.Exec(value.Created)
	// }
}

func Exec(key string, args ...interface{}) {

}

func Insert(key string, args ...interface{}) {

}

func Update(key string, args ...interface{}) {

}

var sqlStatements map[string]SqlStatement

func main() {
	v, err := Parse("sql.config") // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println(v)

	// Conn(v)
	// connDb(v)

	fmt.Println("=------------")
	dao, e := parser.ParseFile("sql-dao.xml")
	if e != nil {
		fmt.Printf("error: %v", e)
		return
	}
	mm := parser.ParseBeans(dao)

	fmt.Println(mm)
}

type BeanFactory struct {
	database map[string]parser.BeanProperty
}

type Engine struct {
	Database map[string]parser.BeanProperty
	Factory  *BeanFactory
}

func (engine *Engine) Table() string {
	return ""
}

func (engine *Engine) Count() int64 {
	return 8
}

func GetBean(beanName string) interface{} {
	if v, ok := GoSqlApp.ReflectObjectMap[beanName]; ok {
		return v
	} else {
		fmt.Println("Key Not Found")
		return nil
	}
}

func GetFactory() *BeanFactory {
	return GoSqlApp.Factory
}

var (
	// GoSqlApp is an application instance
	GoSqlApp *App
)

type App struct {
	Engine           *sql.DB
	Factory          *BeanFactory
	ReflectObjectMap map[string]interface{}
}

func init() {
	// create gosql application
	GoSqlApp = NewApp()
}

func NewApp() *App {
	var engine *sql.DB
	beans, e := parser.ParseFile("sql-dao.xml")
	if e != nil {
		fmt.Printf("error: %v", e)
		return nil
	}
	beanMap := parser.ParseBeans(beans)
	if v, ok := beanMap["database"]; ok {
		engine, _ = ConnDB(v)
		// error
	}

	objects := make(map[string]interface{})
	app := &App{Engine: engine, ReflectObjectMap: objects}
	return app
}

func ConnDB(beanPropertyMap map[string]BeanProperty) (*sql.DB, error) {
	var engine *sql.DB
	var err error

	db := GetPropertyValue(beanPropertyMap, "db", "sqlite3")
	dbname := GetPropertyValue(beanPropertyMap, "dbname", "test")
	dbuname := GetPropertyValue(beanPropertyMap, "dbuname", "root")
	dbpassw := GetPropertyValue(beanPropertyMap, "dbpassw", "123456")

	fmt.Printf("connDb: %s %s %s %s", db, dbName, dbUname, dbpassw)

	if "mysql" == db {
		conn := fmt.Sprintf("%s:%s@/%s?charset=utf8", dbUname, dbpassw, dbName)
		engine, err = sql.Open("mysql", conn) //第一个参数为驱动名
	} else {
		engine, err = sql.Open("sqlite3", dbName)
	}
	return engine, err
}

func AddReflectType(name string, t reflect.Type) {
	var object interface{}
	object = reflect.New(t).Interface()

	mtV := reflect.ValueOf(object)
	CreateTable(mtV.MethodByName("Table").Call(nil)[0].Interface().(string))

	GoSqlApp.ReflectObjectMap[name] = object
}

// Method ‘Table’ 创建数据库表
func CreateTable(table string) {
	if table != "" {

	}
}
