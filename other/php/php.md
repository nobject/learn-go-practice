### zval结构
zval由三个联合体组成，总共占16个字字，zend_value联合体用于存储真正的value值，u1用于存储value的类型，u2...
```c
typedef union _zend_value {
	zend_long         lval;				/* long value */
	double            dval;				/* double value */
	zend_refcounted  *counted;
	zend_string      *str;
	zend_array       *arr;
	zend_object      *obj;
	zend_resource    *res;
	zend_reference   *ref;
	zend_ast_ref     *ast;
	zval             *zv;
	void             *ptr;
	zend_class_entry *ce;
	zend_function    *func;
	struct {
		uint32_t w1;
		uint32_t w2;
	} ww;
} zend_value;

struct _zval_struct {
	zend_value        value;			/* value */
	union {
		uint32_t type_info;
		struct {
			ZEND_ENDIAN_LOHI_3(
				zend_uchar    type,			/* active type */
				zend_uchar    type_flags,
				union {
					uint16_t  extra;        /* not further specified */
				} u)
		} v;
	} u1;
	union {
		uint32_t     next;                 /* hash collision chain */
		uint32_t     cache_slot;           /* cache slot (for RECV_INIT) */
		uint32_t     opline_num;           /* opline number (for FAST_CALL) */
		uint32_t     lineno;               /* line number (for ast nodes) */
		uint32_t     num_args;             /* arguments number for EX(This) */
		uint32_t     fe_pos;               /* foreach position */
		uint32_t     fe_iter_idx;          /* foreach iterator index */
		uint32_t     property_guard;       /* single property guard */
		uint32_t     constant_flags;       /* constant flags */
		uint32_t     extra;                /* not further specified */
	} u2;
};

typedef struct _zend_refcounted_h {
	uint32_t         refcount;			/* reference counter 32-bit */
	union {
		uint32_t type_info;
	} u;
} zend_refcounted_h;

struct _zend_refcounted {
	zend_refcounted_h gc;
};
```
### 引用记数
https://github.com/pangudashu/php7-internal/blob/master/2/zval.md
变量是简单类型（true/false/double/long/null）时直接拷贝值，不需要引用计数。
变量是临时字串，在赋值时会用到引用计数，但如果变量是字符常量(字串字面量)，则不会用到。
变量是对象(zval.u1.v.type=IS_OBJECT), 资源(zval.u1.v.type=IS_RESOURCE), 引用(zval.u1.v.type=IS_REFERENCE)类型时， 赋值一定会用到引用计数。
变量是普通的数组， 赋值时会用到引用计数，变量是IS_ARRAY_IMMUTABLE时，赋值不使用引用计数。
### 垃圾回收
PHP变量的回收是根据refcount实现的，当unset、return时会将变量的引用计数减掉，如果refcount减到0则直接释放value，这是变量的简单gc


#### 循环引用
数组或对象时，容易产生循环引用，即使unset掉refcount也是大于0的
```php
$a = [1];
$a[] = &$a;

unset($a);
```
所以refcount > 0，无法通过简单的gc机制回收，这种变量就是垃圾，垃圾回收器要处理的就是这种情况，目前垃圾只会出现在array、object两种类型中，所以只会针对这两种情况作特殊处理：
当销毁一个变量时，如果发现减掉refcount后仍然大于0，且类型是IS_ARRAY、IS_OBJECT则将此value放入gc可能垃圾双向链表中，等这个链表达到一定数量后启动检查程序将所有变量检查一遍，如果确定是垃圾则销毁释放。


如果当变量的refcount减少后大于0，PHP并不会立即进行对这个变量进行垃圾鉴定，而是放入一个缓冲buffer中，等这个buffer满了以后(10000个值)再统一进行处理，加入buffer的是变量zend_value的zend_refcounted_h:
```text
typedef struct _zend_refcounted_h {
    uint32_t         refcount; //记录zend_value的引用数
    union {
        struct {
            zend_uchar    type,  //zend_value的类型,与zval.u1.type一致
            zend_uchar    flags,
            uint16_t      gc_info //GC信息，垃圾回收的过程会用到
        } v;
        uint32_t type_info;
    } u;
} zend_refcounted_h;
```
一个变量只能加入一次buffer，为了防止重复加入，变量加入后会把zend_refcounted_h.gc_info置为GC_PURPLE，即标为紫色，下次refcount减少时如果发现已经加入过了则不再重复插入。垃圾缓存区是一个双向链表，等到缓存区满了以后则启动垃圾检查过程：遍历缓存区，再对当前变量的所有成员进行遍历，然后把成员的refcount减1(如果成员还包含子成员则也进行递归遍历，其实就是深度优先的遍历)，最后再检查当前变量的引用，如果减为了0则为垃圾。这个算法的原理很简单，垃圾是由于成员引用自身导致的，那么就对所有的成员减一遍引用，结果如果发现变量本身refcount变为了0则就表明其引用全部来自自身成员。具体的过程如下：

(1) 从buffer链表的roots开始遍历，把当前value标为灰色(zend_refcounted_h.gc_info置为GC_GREY)，然后对当前value的成员进行深度优先遍历，把成员value的refcount减1，并且也标为灰色；

(2) 重复遍历buffer链表，检查当前value引用是否为0，为0则表示确实是垃圾，把它标为白色(GC_WHITE)，如果不为0则排除了引用全部来自自身成员的可能，表示还有外部的引用，并不是垃圾，这时候因为步骤(1)对成员进行了refcount减1操作，需要再还原回去，对所有成员进行深度遍历，把成员refcount加1，同时标为黑色；

(3) 再次遍历buffer链表，将非GC_WHITE的节点从roots链表中删除，最终roots链表中全部为真正的垃圾，最后将这些垃圾清除。
### php-fpm
### hash表实现
https://mp.weixin.qq.com/s/J8eICn4BSvwmAyrbeG1xKA
```c
typedef struct _Bucket {
    //始终为zval，PHP7将zval嵌入到bucket中，每一个zval只有16字节，当zval是IS_PTR类型是，
    //可以通过zval.value.ptr指向任何类型数据
    zval              val;  
    // h值，表示数字key或者字符串key的h值 
    zend_ulong        h;      
    //字符串key,区别于PHP5，不再是char *类型的指针，而是一个指向zend_string的指针
    //zend_string是带有字符串长度，h值，gc信息的字符串数组，可以大幅度提升性能和空间效率        
    zend_string      *key;              
} Bucket;

typedef struct _zend_array HashTable;

struct _zend_array {
    //引用计数，引用计数不是zval的字段而是被设计在zval所在value字段所指向的结构体中
    zend_refcounted_h gc;
    //总共4字节 可以存储一个uint32_t类型的flags，也可以存储由4个unsigned char组成的结构体v
    union {
        struct {
            //兼容不同操作系统的大小端
            ZEND_ENDIAN_LOHI_4(
                //用各个bit表达HashTable的各种标记，总共6中flag，对应flags的第1位到第6位
//在zend_hash.h文件中38~43行
//#define HASH_FLAG_PERSISTENT       (1<<0)  是否使用持久化内存，不使用内存池
//#define HASH_FLAG_APPLY_PROTECTION (1<<1)  是否开启递归遍历保护
//#define HASH_FLAG_PACKED           (1<<2) 是否是packed array
//#define HASH_FLAG_INITIALIZED      (1<<3) 是否已经初始化
//#define HASH_FLAG_STATIC_KEYS      (1<<4)  标记key为数字key或者字符串key
//#define HASH_FLAG_HAS_EMPTY_IND    (1<<5) 是否存在空的间接val
                zend_uchar    flags, 
                 //递归遍历数，为了解决循环引用导致死循环问题
                zend_uchar    nApplyCount,
                 //迭代器计数，php中每一个foreach语句都会在全局遍历EG中创建一个迭代器，
                //包含正在遍历的HashTable和游标信息。该字段记录了当前runtime正在迭代当前
                //HashTable的迭代器的数量
                zend_uchar    nIteratorsCount,
                  //调试目的用
//#define HT_OK                 0x00  正常状态 各种数据完全一致
//#define HT_IS_DESTROYING      0x40 正在删除所有内容，包括arBuckets本身
//#define HT_DESTROYED          0x80  已删除，包括arBuckets本身
//#define HT_CLEANING           0xc0 正在清除所有的arBuckets指向的内容，但不包括arBuckets本身
                zend_uchar    consistency)
        } v;
        uint32_t flags;
    } u;  
    //掩码。一般为-nTableSize，区别于PHP5,PHP7中的掩码始终为负数
    uint32_t          nTableMask;
    //实际的存储容器，通过指针指向一段连续的内存，存储着bucket数组
    Bucket           *arData;
    //所有已使用bucket的数量，包括有效bucket和无效bucket数量，在bucket数组中，下标从
    //0~(nNumUsed-1)的bucket都属于已使用bucket，而下标从nNumUsed~(nTableSize-1)的bucket
    //都属于未使用bucket
    uint32_t          nNumUsed;
   //有效的bucket数量，总是<=nNumUsed
    uint32_t          nNumOfElements;
    //HashTable的大小，表示arData所指向的bucket数组大小，即所有bucket数量，取值始终是2^n,
    //最小是8，最大32位系统是（2^30）,64位系统(2^31)
    uint32_t          nTableSize;
    //HashTable的全局默认游标，在PHP7中reset/key/current/next/prev/end等等函数都与该字段
    //有关，是一个有符号整型，由于所有bucket内存是连续的，不需要再根据指针维护正在遍历的
    //bucket，而是只维护正在遍历的bucket所在数组中的下标就可以
    uint32_t          nInternalPointer;
    //HashTable的自然key，即数组插入元素无须指定key，key会以nNextFreeElement的值为准。
    //该值初始化为0，比如$a[] = 1, 实际上插入到key等于0的bucket上，然后nNextFreeElement
    //会变为1，代表下一个自然插入的元素的key为1
    zend_long         nNextFreeElement;
   //析构函数，当bucket元素被更新或者删除时，会对bucket的value调用该函数，如果value是引用计数
   //的类型，那么会对value引用计数-1，引发可能的GC
    dtor_func_t       pDestructor;
};
```
```text
1) key: 通过key可以快速找到value。一般可以为数字或者字符串。
2) value: 值可以为复杂结构
3) bucket: 桶，HashTable存储数据的单元，用来存储key/value以及其他辅助信息的容器
4) slot:槽，HashTable有多少个槽，一个bucket必须属于具体的某一个slot，一个slot可以有多个
   bucket
5) 哈希函数:一般都是自己实现（time33），在存储的时候，会将key通过hash函数确定所在的slot
6) 哈希冲突: 当多个key经过哈希计算后，得到的slot位置是同一个，那么就会冲突，一般这个时候会有
   2种解决办法:链地址法和开放地址法。其中php采用的是链地址法，即将同一个slot中的bucket通过
   链表连接起来
```
#### 哈希表中h的作用
```text
1) HashTable中的key可能是数字也可能是字符串，所以在设计bucket的key时，分为字符串key和数字
   key，在上图中的bucket中，“h”代表数字key，“key”代表字符串key，对于数字key，hash1函数并没
   有做任何事情，h值就是数字key
2) 每个字符串key，经过hash1函数都会计算一个h值。可以加快字符串之间的比较速度。如果要比较2个
   字符串是否相等，首先比较这2个key的h值是否相等，如果相等再比对2个key的长度和内容。否则可以
   判定不相等。这样可以提高HashTable插入，查找速度
```
#### packed array与hash array的区别：
```text
packed array
1） key全是数字key
2） key按插入顺序排列，必须是递增的
3）每个key-value对的存储位置都是确定的，都存储在bucket数组的第key个元素上
packed array不需要slot索引数组，而hash array需要slot指向对应的bucket。


```


#### hash实现原理
php7哈希掩码的值是负数的，其实是用于slot的范围。

当一个key进行插入的时候，如果是个数字索引，则不需要进行hash操作，直接往对应的空余的bucket中插入，插入的bucket的索引值记录在slot索引数组中。如果key为字符串，还需要先进行hash处理，转换成数字，然后与哈希掩码作或运算，找到具体的slot位置，用来记录bucket的索引值
#### 哈希冲突解决方案
当两个key对应的是同一个slot下的索引，那么就产生哈希冲突，由于每个slot索引数组存的是一个int32的值，那hash冲突后，php7的处理方案是用新的bucket的索引位置替换掉旧的索引位置，而新的bucket存储的val中的next会记录被顶替掉的那个索引位置。这样即使冲突后，也可以通过next的指针去寻找之前的值

#### 什么时候扩容
```text
1). Hash Array的容量是分配是固定的，初始化时每次申请是2^n的容量，最小为8，最大为0x80000000

2).当容量足够的时候直接执行插入操作。

3).当容量不够时(nNumUsed >= nTableSize),检查已删除元素所占的比例，假如达到阈值
(ht->nNumUsed - ht->nNumOfElements > (ht->nNumOfElements >> 5)),则将已删除元素从 HashTable中移除，并进行重建索引操作。如果未达到阈值，则进行扩容，新的容量扩大到当前的2倍(2 * nTableSize),将当前Bucket数组复制到新的空间，并进行重建索引操作。

4).扩容并重建完成后，如果有足够的空间则再执行插入操作。
```

#### rehash流程
```text
1). Rehash 对应的源码中的 zend_hash_rehash(ht)方法
2). Rehash 的主要功能就是把HashTable bucket数组中标识为IS_UNDEF的数据删除，把有效的数据重
    新聚合到bucket数组并更新插入索引表
3).整个Rehash不会重新申请内存，而是在原有结构上进行聚合调整
   具体步骤如下:
   (1). 重置所有nIndex数组为-1.
   (2). 初始化2个bucket类型的指针p，q,循环遍历bucket数组 
   (3). 每次循环，p++，遇到第一个IS_UNDEF时，q=p；继续遍历
   (4). 当再次遇到一个正常数据时，把正常数据拷贝到q指向的位置,q++
   (5). 直到遍历完数组，更新nNumUsed等计数
```
###PHP弱类型变量是如何实现
https://www.kancloud.cn/martist/be_new_friends/1739248

### php的static什么情况下会用，好处是什么
```text
  答：静态的东西都是给类用的（包括类常量），非静态的都是给对象用的
    （1）静态方法可以直接被类访问，不需要实例化
    （2）函数执行完静态属性的值会一直都在
```
### php怎么实现常驻进程的，如何配置，如何监控
```text

```
### php-fpm
```text
PHP-FPM 即 PHP FastCGI 进程管理器，要了解 PHP-FPM ，首先要看看 CGI 与 FastCGI 的关系。

CGI 的英文全名是 Common Gateway Interface，即通用网关接口，是 Web 服务器调用外部程序时所使用的一种服务端应用的规范。

早期的 Web 通信只是按照客户端请求将保存在 Web 服务器硬盘中的数据转发过去而已，这种情况下客户端每次获取的信息也是同样的内容（即静态请求，比如图片、样式文件、HTML文档），而随着 Web 的发展，Web 所能呈现的内容更加丰富，与用户的交互日益频繁，比如博客、论坛、电商网站、社交网络等。

这个时候仅仅通过静态资源已经无法满足 Web 通信的需求，所以引入 CGI 以便客户端请求能够触发 Web 服务器运行另一个外部程序，客户端所输入的数据也会传给这个外部程序，该程序运行结束后会将生成的 HTML 和其他数据通过 Web 服务器再返回给客户端（即动态请求，比如基于 PHP、Python、Java 实现的应用）。利用 CGI 可以针对用户请求动态返回给客户端各种各样动态变化的信息。



FastCGI 顾名思义，是 CGI 的升级版本，为了提升 CGI 的性能而生，CGI 针对每个 HTTP 请求都会 fork 一个新进程来进行处理（解析配置文件、初始化执行环境、处理请求），然后把这个进程处理完的结果通过 Web 服务器转发给用户，刚刚 fork 的新进程也随之退出，如果下次用户再请求动态资源，那么 Web 服务器又再次 fork 一个新进程，如此周而复始循环往复。

而 FastCGI 则会先 fork 一个 master 进程，解析配置文件，初始化执行环境，然后再 fork 多个 worker 进程（与 Nginx 有点像），当 HTTP 请求过来时，master 进程将其会传递给一个 worker 进程，然后立即可以接受下一个请求，这样就避免了重复的初始化操作，效率自然也就提高了。而且当 worker 进程不够用时，master 进程还可以根据配置预先启动几个 worker 进程等着；当空闲 worker 进程太多时，也会关掉一些，这样不仅提高了性能，还节约了系统资源。

这样一来，PHP-FPM 就好理解了，FastCGI 只是一个协议规范，需要每个语言具体去实现，PHP-FPM 就是 PHP 版本的 FastCGI 协议实现，有了它，就是实现 PHP 脚本与 Web 服务器（通常是 Nginx）之间的通信，同时它也是一个 PHP SAPI，从而构建起 PHP 解释器与 Web 服务器之间的桥梁。



PHP-FPM 负责管理一个进程池来处理来自 Web 服务器的 HTTP 动态请求，在 PHP-FPM 中，master 进程负责与 Web 服务器进行通信，接收 HTTP 请求，再将请求转发给 worker 进程进行处理，worker 进程主要负责动态执行 PHP 代码，处理完成后，将处理结果返回给 Web 服务器，再由 Web 服务器将结果发送给客户端。这就是 PHP-FPM 的基本工作原理，

PS：最大请求数：最大处理请求数是指一个php-fpm的worker进程在处理多少个请求后就终止掉，master进程会重新respawn一个新的。

这个配置的主要目的是避免php解释器或程序引用的第三方库造成的内存泄露。

pm.max_requests = 10240
```
### php的自动加载
https://learnku.com/articles/4681/analysis-of-the-principle-of-php-automatic-loading-function
```text
通常 PHP5 在使用一个类时，如果发现这个类没有加载，就会自动运行 _autoload () 函数，这个函数是我们在程序中自定义的，在这个函数中我们可以加载需要使用的类。下面是个简单的示例：

        function __autoload($classname) {
           require_once ($classname . "class.php"); 
        }
Copy
在我们这个简单的例子中，我们直接将类名加上扩展名 ”.class.php” 构成了类文件名，然后使用 require_once 将其加载。从这个例子中，我们可以看出 autoload 至少要做三件事情：

根据类名确定类文件名；
确定类文件所在的磁盘路径 (在我们的例子是最简单的情况，类与调用它们的 PHP 程序文件在同一个文件夹下)；
将类从磁盘文件中加载到系统中。

```
### php运行流程
```text
php运行原理

1.转换为tokens语言片段
2.解析为表达式
3.编译为opcodes
4.zend引擎执行opcodes

fastcgi 进程管理器

php-fpm php进程管理

iis isapi
apache apache2handle + fastcgi + php-fpm
nginx fastcgi + php-fpm
```

### php 7 新特性
```text
<?php
// 1.太空船操作符
// echo 3<=>2; // 3>2 ,返回1
// echo false<=>0; // 返回0，false 0 null array()
// echo 'a'<=>'b'; // 97<98 ,返回-1 转ASCII码比较，ord('a')=97
// ord('{')=123

// 2.函数加变量声明与返回值
// 测试一，参数类型定义
// function test(int $a): int{
// 		return $a;
// }
// echo test(9); //9
// echo test('9'); //9
// echo test('a'); //报错
// echo test('0a'); //警告
// echo test('02'); //2

// 测试二,返回值类型与定义的不一致
// function test(int $a): int{
// 	return (array)$a;
// }
// echo test(9); //报错

// 类型说明
// '932' = 932 都会通过，整型与字符串好像没区别，会自动转换
// 'a',一定是字符串

// 3.三元运算
// $a=$c??$b; // 等同于 $a=isset($c)?$c:$b;
// $a=$c?:$b; // 等同于 $a=$c?$c:$b;
// 所以建议用??

// $name = $_GET['name'] ?? 'default_value';
// echo $name;
// // 说明
// // http://t.com/php/php7.php 			返回default_value
// // http://t.com/php/php7.php?name=		返回空,不会返回default_value,为了避免这种情况，要多加个是否为空判断
// // http://t.com/php/php7.php?name=wang 	返回wang


// 4.define定义数组常量
// define('STATUS',['未通过','审核中','已通过']);
// var_dump(STATUS);

// 5.命名空间导入多个类可以合并
// use some\namespace\{ClassA, ClassB, ClassC as C};

?>
```