### docker的实现原理
https://coolshell.cn/articles/17010.html
>一个“容器”，实际上是一个由 Linux Namespace、 Linux Cgroups 和 rootfs 三种技术构建出来的进程的隔离环境。
> Namespace 的作用是“隔离”，隔离进程，网络，挂载等 , Cgroups 的作用是“限制”, 限制容器的资源。所以容器，其实是一种特殊的进程而已
```text
1. 一组联合挂载在 /var/lib/docker/aufs/mnt 上的 rootfs，这一部分我们称为“容器镜
     像”（Container Image），是容器的静态视图；
     
2. 一个由 Namespace+Cgroups 构成的隔离环境，这一部分我们称为“容器运行
   时”（Container Runtime），是容器的动态视图。     
```
深入：
```text
在使用Docker 的时候，并没有一个真正的“Docker 容器”运行在宿主机里面。Docker 项目帮助用户启动的，还是原来的应用进程，只不过在创建这些进程时，Docker 为它们加上了各
种各样的 Namespace 参数。

cgroup主要的作用：
就是限制一个进程组能够使用的资源上限，包括 CPU、内存、磁盘、网络带宽等等。

Mount Namespace 跟其他 Namespace 的使用略有不同的地方：它对容器进程视图的改变，一定是伴随着挂载操作（mount）才能生效

rootfs 只是一个操作系统所包含的文件、配置和目录，并不包括操
作系统内核。在 Linux 操作系统中，这两部分是分开存放的，操作系统只有在开机启动时
才会加载指定版本的内核镜像。
```


### docker与虚拟机的区别
```text
虚拟机：通过硬件虚拟化功能，模拟出了运行一个操作系统需要的各种硬件，比如CPU、内存、I/O 设备等等。然后，它在这些虚拟的硬件上安装了一个新的操作系统，即
Guest OS。


```
### docker的常用命令
### dockerfile的常用命令
