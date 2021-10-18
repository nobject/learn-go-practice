package  main

import (
	"fmt"
	"reflect"
)

type Myint int

type User struct {
	name string `json:"name"`
	age int
}

func (u *User)GetUserName() string{
	return u.name
}

func main() {
	u := &User{name:"lyd",age:24}
	v := reflect.ValueOf(u)       //{lyd 24}
	//ret := v.Method(0).Call([]reflect.Value{})
	ret := v.MethodByName("GetUserName").Call([]reflect.Value{})
	fmt.Printf("%v\n",ret)
}
