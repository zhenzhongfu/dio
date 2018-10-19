package main

import (
	"fmt"
	"reflect"
)

type ee struct {
	goroot string `env:"GOROOT"`
	gopath string `env:"GOPATH"`
}

func main() {
	e := ee{}
	// 获取struct tag
	t := reflect.ValueOf(e)
	for i := 0; i < t.NumField(); i++ {
		field := t.Type().Field(i)
		fmt.Println(field.Tag.Get("env"))
	}
}
