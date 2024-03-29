你好，我叫蒋贻华，毕业于浙江师范大学计算机专业。应聘的岗位是golang工程师，golang方面的经验是2年左右。
按最近几家公司工作及项目介绍一下：
最近一家公司商米科技有限公司，一家以智能终端设备为主营业务的公司，在公司里，主要负责的项目有三个。

1. 远程控制，控制设备的屏幕，用于解决设备出现故障等问题。
   方案选择1：七牛云的推流，将屏幕以rtmp推流的方式推给控制端，效果不理想，延迟2秒以上。
   rtmp 基于tcp，所以他的延时会比较高，延迟5秒内，优化的好还得延迟至少2秒
   webrtc主要是基于udp，控制屏幕相关的可以考虑丢帧的情况，webrtc的延时
   后面调研使用了webrtc，webrtc天然兼容浏览器，也提供了丰富的api，延迟可以在500ms左右
   webrtc通信前的信令交互, 我们建立两条socket连接，分别与前端交互，app交互，传递信令(web端告诉服务端有哪些能力，支不支持webrtc，支不支持音视频等)信息，鉴权，有没有登录连接等
   webrtc服务器的搭建，
   ● 负责云端socket与设备端socket长连接的业务代码编写与维护
   ● 负责云端数据统计api的展示
   ● 对原有服务进行解耦，重构，项目间采用grpc进行通信

2. 远程设备管理（mdm） 
   MDM功能是客户可以通过云端操纵机具的开关机，app卸载，文件分发，日志获取，adb调试，清除锁屏密码等功能，相较于远程协助，更便于用户的操作。
   协助制定云端的接口规范，负责与云端交互的api接口开发
   
   为什么不使用socket.io而使用mqtt协议
   MDM，远程设备管理（控制设备的开关机）, 选用mqtt协议的原因
      1.使用发布/订阅消息模式，提供一对多的消息发布，解除应用程序耦合。
      2.对负载内容屏蔽的消息传输。
      3.使用 TCP/IP 提供网络连接。
      4.有三种消息发布服务质量：
         "至多一次"，消息发布完全依赖底层 TCP/IP 网络。会发生消息丢失或重复。这一级别可用于如下情况，环境传感器数据，丢失一次读记录无所谓，因为不久后还会有第二次发送。
         "至少一次"，确保消息到达，但消息重复可能会发生。
         "只有一次"，确保消息到达一次。这一级别可用于如下情况，在计费系统中，消息重复或丢失会导致不正确的结果。
      5.小型传输，开销很小（固定长度的头部是 2 字节），协议交换最小化，以降低网络流量。

   ● 负责开发地理围栏，文件分发，日志获取，定时开关机，定制桌面，远程锁机等模块
   ● 调研emq作为broker，搭建集群及压测，提高服务的稳定性，目前长连接可维持在百万级别
   ● 基于mqtt协议的topic收发业务代码的实现
   ● 对外提供openapi的接口实现
   ● mdm收费功能的数据库设计，代码开发
   ● 迁移原项目的虚机部署方式 ，改由docker + k8s 部署
   ● 将原先的api接口从php用golang重构
   
下面说说项目相关：
1. 推送项目
   推送服务原先采用小米推送服务，该版在兼容小米推送的基础上，增加基于mqtt长连接的推送服务，小米推送逐渐从公司的新设备中移除，
   而该版将承接大部分公司的推送服务。 使用rocketMQ作为消息队列服务，服务间通信采用grpc的方式，redis用于保存每次推送的任务的到达数等。
   涉及的技术：rocketMQ, golang, mqtt, redis, grpc
   ● 使用redis维护每台设备的在线状态并记录日活
   ● 基于小米推送提供的接口实现基于小米推送的golang版本sdk
   ● push service云端接口开发，使用grpc + grpc gateway的方式提供给业务端接口
   ● 推送回执及离线推送等业务场景实现

   redis：   1.上下线的状态记录,使用的普通的string类型，
             2.记录设备的日活用hyperloglog
             3.任务对应的推送设备，使用set数据结构，推送成功即删除
   rocketmq: 设备上线通知，另一个服务收到设备上线通知，会查找该设备因离线而未推送的任务，然后进行推送
             设备推送成功上报回执通知，上报推送结果，将redis中对应任务的设备sn从set中删除
   难点：设备因为内存，版本，网络环境等各种问题，可能会导致设备上下线通知，有的设备可能一天出现这种情况5，6W次，异常频繁，导致服务可能会处理很多无效功，
   对于每日设备上下线，目前会用redis记录每个设备上次上线的时间与下线的时间与当天上线次数，如果与前一次上线时间很接近，则抛弃，如果当天上线次数大于200次，可认为设备异常，将不处理之后的逻辑，
