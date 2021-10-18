package main

import (
	"fmt"
)

type A struct {
	name string
}
func testForRange()  {
	slice := make([]*A,2)
	slice[0] = &A{"hello0"}
	slice[1] = &A{"hello1"}
	for _,v:=range slice{
		go func(test *A) {
			fmt.Printf("%#v",test)
		}(v)
	}
	for {}
}

func printName(a *A)  {

}
func main()  {
	testForRange()
}