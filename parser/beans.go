package parser

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type Beans struct {
	XMLName xml.Name     `xml:"beans"`
	Bean    []BeanObject `xml:"bean"`
}

type BeanObject struct {
	Name     string         `xml:"name,attr"`
	Property []BeanProperty `xml:"property"`
}

type BeanProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// exported interfaces
func ParseFile(path string) (*Beans, error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	beans := Beans{}
	err = xml.Unmarshal(data, &beans)
	if err != nil {
		return nil, err
	}

	return &beans, nil
}

func ParseBeans(beans *Beans) map[string]map[string]BeanProperty {
	var data map[string]map[string]BeanProperty
	data = make(map[string]map[string]BeanProperty)
	for _, value := range beans.Bean {
		beanName := value.Name
		beanProerty := value.Property
		m1 := make(map[string]BeanProperty)
		for _, v := range beanProerty {
			m1[v.Name] = v
		}
		data[beanName] = m1
	}
	return data
}

func GetPropertyValue(propertyMap map[string]BeanProperty, name string, defValue string) string {
	if v, ok := propertyMap[name]; ok {
		return v.Value
	} else {
		return defValue
	}
}
