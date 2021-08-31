package main

import "fmt"

type People struct {
}

func (p *People) String() string {
	return fmt.Sprintf("peopele: %v", p)
}

func main() {
	p := &People{}
	p.String()
	//c := make(chan int)
	//
	//go func() {
	//	c <- 1 // send to channel
	//}()
	//
	//x := <-c // recv from channel
	//
	//fmt.Println(x)
}
