package _struct

import (
	"fmt"
	"reflect"
)

func CompareStruct() {
	struct1 := struct {
		Name string
		Age  int
	}{
		"jiang", 20,
	}
	struct2 := struct {
		Name string
		Age  int
	}{
		"jiang", 20,
	}
	if struct1 == struct2 {
		fmt.Println("属性相同顺序的结构体，可以比较")
	}
	struct3 := struct {
		Name string
		m    map[int]string
	}{
		"jiang", map[int]string{20: "jiang"},
	}
	struct4 := struct {
		Name string
		m    map[int]string
	}{
		"jiang", map[int]string{20: "jiang"},
	}
	// error
	//if struct3 == struct4 {
	//	fmt.Println("属性相同顺序的结构体，可以比较")
	//}

	if reflect.DeepEqual(struct3, struct4) {
		fmt.Println("struct3 == struct4")
	} else {
		fmt.Println("struct3 != struct4")
	}
}
