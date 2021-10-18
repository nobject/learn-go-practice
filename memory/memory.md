### 参考文章
https://zhuanlan.zhihu.com/p/59125443
https://draveness.me/golang/docs/part3-runtime/ch07-memory/golang-memory-allocator/
https://mp.weixin.qq.com/s/3gGbJaeuvx4klqcv34hmmw
https://www.cnblogs.com/zkweb/p/7880099.html
https://www.infoq.cn/article/IEhRLwmmIM7-11RYaLHR

### 逃逸分析
通过检查变量的作用域是否超出了它所在的栈来决定是否将它分配在堆上的技术， 其中变量的作用域超出了它所在的栈，这种行为称为逃逸。

好处：
1. 函数返回直接释放，不会引起垃圾回收，对性能没有影响
2. 减少内存碎片的产生
3. 减轻分配堆内存的开销，提高程序的运行速度

可能产生逃逸的情况
- 指针作为返回值（方法得非内联,如果被内联优化，有可能并未逃逸）
```go
type People struct {
    Name string
    Age  int
}
//go:noinline
func NewPeople(name string, age int) *People{
    p := new(People)

    p.Name = name
    p.Age = age
    return p
}
```
- 栈空间不足逃逸,比如大array或者slice
```go
func BigArr()  {
	arr := make([]int,8193)
	for i :=0;i <=1000;i++ {
		arr[i] = i
	}
}
```
- 动态类型逃逸
```go
func DynamicType()  {
	var a interface{}
	a = 5
	fmt.Println(a)
}
```
- 闭包引用对象逃逸
- 在一个切片上存储指针或者带指针的值, 该值会被放入堆中

函数传递指针真的比传值效率高吗？
```text
传递指针可以减少底层值的拷贝，可以提高效率，但是如果拷贝的数据量小，由于指针传递会产生逃逸，可能会使用堆，也可能会增加GC的负担，所以传递指针不一定是高效的,得看对象的大小，小对象可以使用值传递。
```
### 内存分配
> go内存分配借鉴TCMalloc(thread cache malloc)的思想， TCMalloc的思想：
> 1. 为了防止多线程内存竞争，为每个线程预分配内存，减少竞争时的加锁的可能性
> 2. 为线程预分配内存需要1次系统调用，后续分配就无需额外进行系统调用，在用户态即可处理

Go的内存管理：  
相关概念
- page：x64下每页大小为8KB
  
- span: 内存管理的基本单位，结构体名为mspan,一组连续的page
  
- mcache：保存各种大小的span，并按span class分类，小对象直接从mcache分配，可以无锁访问
  
- mcentral：所有线程共享的内存，需加锁访问，结构与mcache差不多，centralCache是每个级别span有1个链表，mcache每个级别span有两个链表
  
- mheap
  堆内存的抽象，把从os申请出的内存页组织成span，当mcentral的span不够时向mheap申请，mheap的span不够时向os申请，按页从os申请并组织成span供mcentral使用
  
- object size(size)：指申请内存的对象的大小
- size class(class): size级别
- span class：span的级别
- num of page(npage)：代表page的数量，span包含的页数

生成对象流程
```text
获取当前goroutine的m的mcache
判断大小，如果小于32kb
    如果小于16byte，微小对象的分配，级别固定的span
    小对象，小于1kb, 以8字节为跨度分级别，0-32个级别
    小对象，大于1kb, 以128字节为跨度分级别，32-67级别
    根据级别获取到对应的对象大小（每个级别都会设定好对应的对象大小，因此对象其实是存在一定的内存浪费的）
    尝试快速的从mcache对应这个span中分配
大对象直接从mheap分配, 这里的s是一个特殊的span, 它的class是0    
```
源码解析
```go
// 获取当前goroutine对应的m
    mp := acquirem()
// 如果当前的m正在执行分配任务，则抛出错误
	if mp.mallocing != 0 {
		throw("malloc deadlock")
	}
	if mp.gsignal == getg() {
		throw("malloc during signal")
	}
// 锁住当前的m进行分配
	mp.mallocing = 1

	shouldhelpgc := false
	dataSize := size
    // 获取当前goroutine的m的mcache
	c := getMCache()
	if c == nil {
		throw("mallocgc called without a P or outside bootstrapping")
	}
	var span *mspan
	var x unsafe.Pointer
	noscan := typ == nil || typ.ptrdata == 0
	// In some cases block zeroing can profitably (for latency reduction purposes)
	// be delayed till preemption is possible; isZeroed tracks that state.
	isZeroed := true
	// 判断大小，如果小于32kb
	if size <= maxSmallSize {
		// 如果小于16byte，微小对象的分配
		if noscan && size < maxTinySize {
			off := c.tinyoffset
			// Align tiny pointer for required (conservative) alignment.
			// 如果是8字节
			if size&7 == 0 {
				off = alignUp(off, 8)
			} else if sys.PtrSize == 4 && size == 12 {
				// Conservatively align 12-byte objects to 8 bytes on 32-bit
				// systems so that objects whose first field is a 64-bit
				// value is aligned to 8 bytes and does not cause a fault on
				// atomic access. See issue 37262.
				// TODO(mknyszek): Remove this workaround if/when issue 36606
				// is resolved.
				off = alignUp(off, 8)
			} else if size&3 == 0 {
				off = alignUp(off, 4)
			} else if size&1 == 0 {
				off = alignUp(off, 2)
			}
			if off+size <= maxTinySize && c.tiny != 0 {
				// The object fits into existing tiny block.
				x = unsafe.Pointer(c.tiny + off)
				c.tinyoffset = off + size
				c.tinyAllocs++
				mp.mallocing = 0
				releasem(mp)
				return x
			}
			// Allocate a new maxTinySize block.
			span = c.alloc[tinySpanClass]
			v := nextFreeFast(span)
			if v == 0 {
				v, span, shouldhelpgc = c.nextFree(tinySpanClass)
			}
			x = unsafe.Pointer(v)
			(*[2]uint64)(x)[0] = 0
			(*[2]uint64)(x)[1] = 0
			// See if we need to replace the existing tiny block with the new one
			// based on amount of remaining free space.
			if !raceenabled && (size < c.tinyoffset || c.tiny == 0) {
				// Note: disabled when race detector is on, see comment near end of this function.
				c.tiny = uintptr(x)
				c.tinyoffset = size
			}
			size = maxTinySize
		} else {
			var sizeclass uint8
			// 小对象，小于1kb
			if size <= smallSizeMax-8 {
				// 以8字节为跨度分级别，0-32个级别
				sizeclass = size_to_class8[divRoundUp(size, smallSizeDiv)]
			} else {
				// 以128字节为跨度分级别，32-67级别
				sizeclass = size_to_class128[divRoundUp(size-smallSizeMax, largeSizeDiv)]
			}
			// 根据级别获取到对应的对象大小（每个级别都会设定好对应的对象大小，因此对象其实是存在一定的内存浪费的）
			size = uintptr(class_to_size[sizeclass])
			// sizeclass * 2 + (noscan ? 1 : 0)
			spc := makeSpanClass(sizeclass, noscan)
			span = c.alloc[spc]
            // 尝试快速的从这个span中分配
			v := nextFreeFast(span)
			if v == 0 {
                // 分配失败, 可能需要从mcentral或者mheap中获取
                // 如果从mcentral或者mheap获取了新的span, 则shouldhelpgc会等于true
                // shouldhelpgc会等于true时会在下面判断是否要触发GC
				v, span, shouldhelpgc = c.nextFree(spc)
			}
			x = unsafe.Pointer(v)
			if needzero && span.needzero != 0 {
				memclrNoHeapPointers(unsafe.Pointer(v), size)
			}
		}
	} else {
        // 大对象直接从mheap分配, 这里的s是一个特殊的span, 它的class是0
		shouldhelpgc = true
		// For large allocations, keep track of zeroed state so that
		// bulk zeroing can be happen later in a preemptible context.
		span, isZeroed = c.allocLarge(size, needzero && !noscan, noscan)
		span.freeindex = 1
		span.allocCount = 1
		x = unsafe.Pointer(span.base())
		size = span.elemsize
	}
```