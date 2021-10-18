package memory

import "fmt"

type People struct {
	Name string
	Age  int
}

func NewPeople(name string, age int) *People{
	p := new(People)

	p.Name = name
	p.Age = age
	return p
}

func BigArr()  {
	arr := make([]int,8193)
	for i :=0;i <=1000;i++ {
		arr[i] = i
	}
}

func DynamicType()  {
	var a interface{}
	a = 5
	fmt.Println(a)
}

func StoreSlicePointer()  {
	slice := make([]*int,1)
	a := 5
	slice[0] = &a
}