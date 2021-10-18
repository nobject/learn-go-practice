### linux常用命令
```text
ls cp rm mkdir touch vim chmod chgrp chroot ss, tcp ping netstat top等
```

### vim 常用命令
```text
0 跳至行首，不管有无缩进，就是跳到第0个字符 (常用)
$ 跳至行尾 (常用)
gg 跳至文首 (常用)
G 调至文尾 (常用)
dd 删除光标所在行 (常用)
n+[Enter]	n 为数字。光标向下移动 n 行(常用)
:1,$s/word1/word2/g	从第一行到最后一行寻找 word1 字符串，并将该字符串取代为 word2 ！(常用)
dw 删除一个字(word)
𝑦𝑦复制一行
/pattern 向后搜索字符串pattern
n 下一个匹配(如果是/搜索，则是向下的下一个，?搜索则是向上的下一个)
N 上一个匹配(同上)
:w 将缓冲区写入文件，即保存修改
:wq 保存修改并退出
:x 保存修改并退出
:q 退出，如果对缓冲区进行过修改，则会提示
:q! 强制退出，放弃修改
:set nu 显示行号
:set nonu	与 set nu 相反，为取消行号！
i	从目前光标所在处输入
[Esc]	退出编辑模式，回到一般模式中(常用)
u	复原前一个动作。(常用)
[Ctrl]+r	重做上一个动作。(常用)
```
### awk
```text
基本用法：
awk 动作 文件名
awk '{print $0}' demo.txt
$后面表示第几个字段，$0表示整行，$1表示第1个字段，NF表示最后一个字段，NR表示当前的行数

echo 'this is a test' | awk '{print $0}'
cat /etc/hosts | awk '{print $NF}'
cat /etc/hosts | awk '{print $1, $NF}'
cat /etc/hosts | awk '{print NR, $1, $NF}'

awk '条件 动作' 文件名
awk -F ':' '{if ($1 > "m") print $1}'
cat /etc/hosts | awk '{if (NR %2 == 1) print $1; else print $2}'
```
### sed
1、简明教程：

https://coolshell.cn/articles/9104.html#%E7%94%A8s%E5%91%BD%E4%BB%A4%E6%9B%BF%E6%8D%A2

http://sed.sourceforge.net/sed1line_zh-CN.html

2、mac上使用遇到的问题
使用如下命令报错问题
```
sed -i "s/my/john/g" example.txt
```
```
sed: 1: "example.txt": extra characters at the end of h command
```
需要改写下如下
```
sed -i "" "s/my/john/g" example.txt
```
只对-i有影响，对其他的-n等无影响

sed -i 后面的双引号中可写任意字符串或者为空，含义是用于生成源文件的备份文件的文件名。

### grep
1、搜索某个文件里面是否包含字符串
```
grep xxx file.log
```
2、在多个文件中检索某个字符串
```
grep xxx file1.log file2.log file3.log
```
还可以使用通配符
```
grep xxx *.log
```
3、显示所检索内容在文件中的行数，可以使用参数-n
```
grep -n xxx file.log
```
4、检索时需要忽略大小写问题，可以使用参数“-i”
```
grep -i xxx file.log
```
5、从文件内容查找不匹配指定字符串的行
```
grep -v xxx file.go
```
6、搜索、查找匹配的行数：
```
grep -c xxx file.go
```
```
grep xxx file.go | wc -l
```
7、在命令中添加-A，-B，-C参数，可分别获取某关键词出现位置后面、前面、前后n行的内容：
```
grep -A|B|C[num] xxx file.log