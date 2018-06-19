package gosql

import (
	"fmt"
	"github.com/freelifer/gosql"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"testing"
)

func init() {
	// typ := reflect.TypeOf(&Object{}).Elem()
	// value := reflect.New(typ).Interface().(*Object)
	// t.Log(value)
	// gosql.GoSqlApp.ReflectTypeMap["IPersonDao"] = reflect.TypeOf(PersonDaoImpl{})
	gosql.AddReflectType("IPersonDao", reflect.TypeOf(PersonDaoImpl{}))
}

type Person struct {
}

// 数据接口
type IPersonDao interface {
	GetPersonName(id int64) (name string, err error)
	AddPerson(person Person)
	GetPersonCount() int64
	ListPersons() *Person
	AddTest() (int64, error)
}

type PersonDaoImpl struct {
	gosql.Engine
}

func (p *PersonDaoImpl) Table() string {
	return `
	   CREATE TABLE IF NOT EXISTS userinfo(
	       uid INTEGER PRIMARY KEY AUTOINCREMENT,
	       username VARCHAR(64) NULL,
	       department VARCHAR(64) NULL,
	       created DATE NULL
	   );
	   `
}

func (p *PersonDaoImpl) AddTest() (int64, error) {
	//插入数据
	stmt, err := p.GetSql().Prepare("INSERT INTO userinfo(username, department, created) values(?,?,?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *PersonDaoImpl) GetPersonName(id int64) (string, error) {
	return "", nil
}

func (p *PersonDaoImpl) AddPerson(person Person) {

}

func (p *PersonDaoImpl) GetPersonCount() int64 {
	return p.Count()
}

func (p *PersonDaoImpl) ListPersons() *Person {
	return nil
}

// go test *.go -v
func Test_GetWxUser(t *testing.T) {
	// IPersonDao
	var personDao IPersonDao
	personDao = gosql.GetBean("IPersonDao").(*PersonDaoImpl)
	fmt.Println(personDao.GetPersonCount())
	_, err := personDao.AddTest()
	if err != nil {
		t.Error(err)
		return
	}
	// 	fmt.Println(id)
	// p := &PersonDaoImpl{}
	// mtV := reflect.ValueOf(p)
	// fmt.Println(mtV.MethodByName("Table").Call(nil)[0])
}
