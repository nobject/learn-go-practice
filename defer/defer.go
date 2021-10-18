package main

func main() {
	// defer执行的时候i变量已经变为5
	//var whatever [6]struct{}
	//for i := range whatever {
	//	defer func() {
	//		fmt.Println(i)
	//	}()
	//}

	// defer是按照先进后出输出的
	//var whatever [6]struct{}
	//for i := range whatever {
	//	defer func(i int) {
	//		fmt.Println(i)
	//	}(i)
	//}

	// defer与return返回之间的关系，return的返回不是原子性的，先会对值进行赋值，然后真正返回，然而defer的操作是是两者之间进行的
	println(f1())
	println(f2())
	println(f3())
}

func f1() (r int) {
	defer func() {
		r++
	}()
	return 0
}

func f2() (r int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

func f3() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}
