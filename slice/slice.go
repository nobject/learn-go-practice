package slice

import "fmt"

// 初始化
func InitSlice() {
	// nil slice
	var sliceInt []int

	// empty slice
	var sliceInt2 = make([]int, 0)

	// always true
	if sliceInt == nil {
		fmt.Printf("sliceInt is nil,slice data:%p,len:%d,cap:%d\n", sliceInt, len(sliceInt), cap(sliceInt))
	}

	// always false
	if sliceInt2 == nil {
		fmt.Println("sliceInt2 is nil")
	}

	fmt.Printf("sliceInt2 data:%p,len:%d,cap:%d\n", sliceInt2, len(sliceInt2), cap(sliceInt2))
	var sliceInt3 = make([]int, 0)
	fmt.Printf("sliceInt3 data:%p,len:%d,cap:%d\n", sliceInt3, len(sliceInt3), cap(sliceInt3))
}

// 容量与长度
func CapAndLen() {
	// 使用make初始化切片时，默认的容量与长度一致，cap代表slice的底层数组的长度，len代表slice本身的长度
	slice1 := make([]int, 5)
	fmt.Printf("slice1:%d,len:%d,cap:%d\n", slice1, len(slice1), cap(slice1))
	slice2 := make([]int, 5, 10)
	fmt.Printf("slice2:%d,len:%d,cap:%d\n", slice2, len(slice2), cap(slice2))
}

// 切片
func ArraySlice() {
	// 无论切几次，如果在未扩容前，底层数组还是同一个，容量也都是底层数组的容量
	arr := [5]int{1, 2, 3, 4, 5}
	slice := arr[0:4]
	slice2 := slice[1:2]
	fmt.Printf("arr:%p,arr:%d\n", &arr, arr)
	fmt.Printf("slice:%p,slice:%d,len:%d,cap,%d\n", slice, slice, len(slice), cap(slice))
	fmt.Printf("slice2:%p,slice2:%d,len:%d,cap,%d\n", slice2, slice2, len(slice2), cap(slice2))
}

// 扩容
func Append() {
	//slice := []int{1, 2, 3, 4, 5}
	//slice2 := slice[0:3]
	//// 发生扩容后，如超过底层数组的容量，会新建一个新的切片，引用的是新的底层数组
	//// 当原容量不够时，新的切片的容量会按照之前容量的2倍扩容
	//slice3 := append(slice2, []int{6, 7, 8}...)
	//fmt.Printf("slice:%p,slice:%d,len:%d,cap,%d\n", slice, slice, len(slice), cap(slice))
	//fmt.Printf("slice2:%p,slice2:%d,len:%d,cap,%d\n", slice2, slice2, len(slice2), cap(slice2))
	//fmt.Printf("slice3:%p,slice3:%d,len:%d,cap,%d\n", slice3, slice3, len(slice3), cap(slice3))
	//// 当原来的容量为1024后，新的切片容量会按1.25倍扩容
	//slice4 := make([]int, 1024)
	//fmt.Printf("slice4:%p,len:%d,cap,%d\n", slice4, len(slice4), cap(slice4))
	//slice4 = append(slice4, []int{1}...)
	//fmt.Printf("after append,slice4:%p,len:%d,cap,%d\n", slice4, len(slice4), cap(slice4))

	slice5 := make([]int, 1280)
	fmt.Printf("slice5:%p,len:%d,cap,%d\n", slice5, len(slice5), cap(slice5))
	slice5 = append(slice5, []int{1}...)
	fmt.Printf("after append,slice5:%p,len:%d,cap,%d\n", slice5, len(slice5), cap(slice5))

	// 当容量一口气增加到原容量2倍以上,新容量以会当前容量为基础
	slice6 := make([]int, 5)
	fmt.Printf("slice6:%p,len:%d,cap,%d\n", slice6, len(slice6), cap(slice6))
	slice6 = append(slice6, make([]int, 12)...)
	fmt.Printf("after append,slice6:%p,len:%d,cap,%d\n", slice6, len(slice6), cap(slice6))

	// 以上只能确定切片的大致容量，还需要根据切片中的元素大小对齐内存，
	// 当数组中元素所占的字节大小为 1、8 或者 2 的倍数时，运行时会使用如下所示的代码对齐内存，所以并不是绝对的2倍或1.25倍
}

// 复制
func Copy() {
	// 等号赋值操作与copy操作的不同
	// copy复制为值复制，改变原切片的值不会影响新切片。而等号复制为指针复制，改变原切片或新切片都会对另一个产生影响。
	slice := []int{1, 2, 3, 4, 5}
	slice2 := slice
	slice3 := make([]int, 5)
	copy(slice3, slice)
	fmt.Printf("slice:%p,slice:%d,len:%d,cap,%d\n", slice, slice, len(slice), cap(slice))
	fmt.Printf("slice2:%p,slice2:%d,len:%d,cap,%d\n", slice2, slice2, len(slice2), cap(slice2))
	fmt.Printf("slice3:%p,slice3:%d,len:%d,cap,%d\n", slice3, slice3, len(slice3), cap(slice3))
	slice2[3] = 8
	slice3[3] = 9
	fmt.Println("after change value")
	fmt.Printf("slice:%p,slice:%d,len:%d,cap,%d\n", slice, slice, len(slice), cap(slice))
	fmt.Printf("slice2:%p,slice2:%d,len:%d,cap,%d\n", slice2, slice2, len(slice2), cap(slice2))
	fmt.Printf("slice3:%p,slice3:%d,len:%d,cap,%d\n", slice3, slice3, len(slice3), cap(slice3))
}

// 参数传递
func ParamPass() {
	slice := []int{1, 2, 3, 4, 5}
	fmt.Printf("slice:%p,slice:%d,len:%d,cap,%d\n", slice, slice, len(slice), cap(slice))
	param(slice)
	fmt.Println("after call param")
	fmt.Printf("slice:%p,slice:%d,len:%d,cap,%d\n", slice, slice, len(slice), cap(slice))
	param2(slice)
	fmt.Println("after call param2")
	fmt.Printf("slice:%p,slice:%d,len:%d,cap,%d\n", slice, slice, len(slice), cap(slice))
}
func param(b []int) {
	b[1] = 10
}

func param2(b []int) {
	b = append(b, 20)
}
