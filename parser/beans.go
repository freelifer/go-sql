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
