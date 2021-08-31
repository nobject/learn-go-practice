package pointer

import (
	"fmt"
	"unsafe"
)

func slicePointer() {
	//slice的结构
	//type slice struct {
	//	array unsafe.Pointer
	//	len   int
	//	cap   int
	//}
	slice := make([]int, 1, 10)
	sliceLen := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&slice)) + uintptr(8)))
	fmt.Printf("slice len:%d, slice cacl len:%d\n", len(slice), sliceLen)
	sliceCap := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&slice)) + uintptr(16)))
	fmt.Printf("slice cap:%d, slice cacl cap:%d\n", cap(slice), sliceCap)
}

func mapPointer() {
}

func structPointer() {

}
