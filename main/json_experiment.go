package main

import (
	"encoding/json"
	"fmt"
)

type Name struct {
	CommonName     string
	ScientificName string
}

func (n *Name) UnmarshalJSON(bytes []byte) error {
	var name string
	err := json.Unmarshal(bytes, &name)
	if err != nil {
		return err
	}
	n.CommonName = name
	n.ScientificName = ""
	return nil
}

type Elephant struct {
	Name Name `json:"name"`
	Age  int  `json:"age"`
}

func main() {
	alice := Elephant{}
	aliceJson := `{"name":"Alice","age":4}`
	err := json.Unmarshal([]byte(aliceJson), &alice)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", alice)
}

// main.YamlConfig{Name:main.Name{Path:"Alice", Url:""}, Age:2}
