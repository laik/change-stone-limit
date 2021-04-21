package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// "path":"a.b.c"
// data = {"a":{"b":{"c":123}}}
// Set(data,"a.b.c",123)
func Set(data map[interface{}]interface{}, path string, value interface{}) {
	head, remain := shift(path)
	_, exist := from(data)[head]
	if !exist {
		data[head] = make(map[interface{}]interface{})
	}
	if remain == "" {
		data[head] = value
		return
	}
	Set(data[head].(map[interface{}]interface{}), remain, value)
}

// data = {"a":{"b":{"c":123}}}
// Get(data,"a.b.c") = 123
func Get(data map[string]interface{}, path string) (value interface{}) {
	head, remain := shift(path)
	_, exist := data[head]
	if exist {
		if remain == "" {
			return data[head]
		}
		switch data[head].(type) {
		case map[string]interface{}:
			return Get(data[head].(map[string]interface{}), remain)
		}
	}
	return nil
}

// data = {"a":{"b":{"c":123}}}
// Delete(data,"a.b.c") = {"a":{"b":""}}
func Delete(data map[string]interface{}, path string) {
	head, remain := shift(path)
	_, exist := data[head]
	if exist {
		if remain == "" {
			delete(data, head)
			return
		}
		switch data[head].(type) {
		case map[string]interface{}:
			Delete(data[head].(map[string]interface{}), remain)
			return
		}
	}
}

func shift(path string) (head string, remain string) {
	slice := strings.Split(path, ".")
	if len(slice) < 1 {
		return "", ""
	}
	if len(slice) < 2 {
		remain = ""
		head = slice[0]
		return
	}
	return slice[0], strings.Join(slice[1:], ".")
}

var object = make(map[string]interface{})

func from(src map[interface{}]interface{}) map[string]interface{} {
	dest := make(map[string]interface{})
	for k, v := range src {
		dest[k.(string)] = v
	}
	return dest
}

//func copyFromSlice(src []interface{}, dest []interface{}){}

func main() {
	text, err := ioutil.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal([]byte(text), &object); err != nil {
		panic(err)
	}
	for _, item := range Get(object, "items").([]interface{}) {
		stone := from(item.(map[interface{}]interface{}))
		stoneSpec := Get(stone, "spec")
		stoneSpecObject := from(stoneSpec.(map[interface{}]interface{}))
		stoneSpecTemplate := from(Get(stoneSpecObject, "template").(map[interface{}]interface{}))
		stoneContainerSpec := from(Get(stoneSpecTemplate, "spec").(map[interface{}]interface{}))

		for _, c := range Get(stoneContainerSpec, "containers").([]interface{}) {
			Set(c.(map[interface{}]interface{}), "resources.limits.cpu", "300m")
			Set(c.(map[interface{}]interface{}), "resources.limits.memory", "1024M")
			Set(c.(map[interface{}]interface{}), "resources.requests.cpu", "100m")
			Set(c.(map[interface{}]interface{}), "resources.requests.memory", "300M")

		}
		out, err := yaml.Marshal(item)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n---\n", out)
	}

}
