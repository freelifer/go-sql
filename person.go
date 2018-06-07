package main

import (
	"fmt"
	"github.com/freelifer/gosql"
	"reflect"
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
}

type PersonDaoImpl struct {
	gosql.Engine
}

func (p *PersonDaoImpl) Table() string {
	return "1111"
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

func main() {
	// IPersonDao
	var personDao IPersonDao
	personDao = gosql.GetBean("IPersonDao").(*PersonDaoImpl)
	fmt.Println(personDao.GetPersonCount())

	// p := &PersonDaoImpl{}
	// mtV := reflect.ValueOf(p)
	// fmt.Println(mtV.MethodByName("Table").Call(nil)[0])
}
