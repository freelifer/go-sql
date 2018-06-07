package main

import (
	"github.com/freelifer/go-sql"
	"reflect"
)

func init() {
	gosql.GoSqlApp.ReflectTypeMap["IPersonDao"] = reflect.TypeOf(PersonDaoImpl)
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

func (p *PersonDaoImpl) GetPersonName(id int64) (string, error) {
	return "", nil
}

func (p *PersonDaoImpl) AddPerson(person Person) {

}

func (p *PersonDaoImpl) GetPersonCount() int64 {
	return 1
}

func (p *PersonDaoImpl) ListPersons() *Person {
	return nil
}

func main() {
	// IPersonDao
	personDao := gosql.GetBean("IPersonDao")
	fmt.Println(personDao)
}
