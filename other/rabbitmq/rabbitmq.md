https://juejin.cn/post/6844903805453074446
### 一些概念
ConnectionFactory（连接管理器）：应用程序与Rabbit之间建立连接的管理器，程序代码中使用。

Channel（信道）：消息推送使用的通道，在客户端的每个连接里，可建立多个channel。

Exchange（交换器）：消息交换机，指定消息按什么规则，路由到哪个队列。

Queue（队列）：消息的载体，每个消息都会被投到一个或多个队列。

RoutingKey（路由键）：exchange根据这个关键字进行消息投递。

Binding（绑定键）：把exchange和queue按照路由规则绑定起来

Consumer:消息消费者，接受消息的程序

Producer:消息生产者，投递消息的程序

#### exchange
direct exchange:直接匹配，通过exchange名称 + routingKey来发送与接收消息
fanout exchange:广播订阅，向所有的消费者发布消息，但是只有消费者将队列绑定到该路由器才能收到消息
topic exchange：主题匹配订阅，这里的主题指的是routingKey,routingKey可以采用通配符
默认的exchange:如果使用空字符串去声明一个exchange，系统就会使用amq.direct这个exchange，我们创建一个queue时，默认的都会有一个和新建queue同名的routingKey绑定到这个默认的exchange上去

### rabbitmq 中 vhost 的作用是什么？
vhost 可以理解为虚拟 broker ，即 mini-RabbitMQ server。
其内部均含有独立的 queue、exchange 和 binding 等，但最最重要的是，其拥有独立的权限系统，可以做到 vhost 范围的用户控制。
当然，从 RabbitMQ 的全局角度，vhost 可以作为不同权限隔离的手段（一个典型的例子就是不同的应用可以跑在不同的 vhost 中）。

### rabbitmq 怎么实现延迟消息队列？
通过消息过期后进入死信交换器，再由交换器转发到延迟消费队列，实现延迟功能；
使用 RabbitMQ-delayed-message-exchange 插件实现延迟功能。

### 消息确认
消息确认
消费者应用（Consumer applications） - 用来接受和处理消息的应用 - 在处理消息的时候偶尔会失败或者有时会直接崩溃掉。而且网络原因也有可能引起各种问题。这就给我们出了个难题，AMQP代理在什么时候删除消息才是正确的？AMQP 0-9-1 规范给我们两种建议：

当消息代理（broker）将消息发送给应用后立即删除。（使用AMQP方法：basic.deliver或basic.get-ok）
待应用（application）发送一个确认回执（acknowledgement）后再删除消息。（使用AMQP方法：basic.ack）
前者被称作自动确认模式（automatic acknowledgement model），后者被称作显式确认模式（explicit acknowledgement model）。在显式模式下，由消费者应用来选择什么时候发送确认回执（acknowledgement）。应用可以在收到消息后立即发送，或将未处理的消息存储后发送，或等到消息被处理完毕后再发送确认回执（例如，成功获取一个网页内容并将其存储之后）。


### 集群
- 镜像模式
  在其他的机器上也会保留当前机器上queue的数据，即使有一台宕机了，也可以使用另一台机器的slave去读取
  脑裂: 脑裂问题是分布式系统中最常见的问题，指在一个高可用（HA）系统中，当联系着的两个节点断开联系时，本来为一个整体的系统，分裂为两个独立节点，这时两个节点开始争抢共享资源，结果会导致系统混乱，数据损坏。对于无状态服务的HA，无所谓脑裂不脑裂；但对有状态服务，数据相关服务(比如MySQL，消息队列)的HA，必须要严格防止脑裂。（但有些生产环境下的系统按照无状态服务HA的那一套去配置有状态服务，结果可想而知...），有一些存储系统像数据库，kv存储都已经有很好的一致性协议解决了raft paxos协议解决了，这里我们需要格外注意这里的脑裂处理流程。脑裂带来的最大问题就是分区问题，分区在rabbitmq中有三种配置模式ignore(默认方式),
  ignore的方式会导致脑裂，因为有问题的机器并不会从集群中剔除
  ignore: 假设你的集群运行在网络非常可靠的情况，所有的节点都是在相同交换机下，然后交换机在将流量路由到外部。如果任何其他群集发生故障（或者有一个双节点群集，不希望运行任何群集关闭的任何风险。
  pause_minority: 假设你的网络不太可靠，你的节点跨域了通地域多个数据中心，然后数据中心可能会异常，你希望到集群某个中心异常的时候，其他两个数据中心服务继续工作，当数据中心恢复后，节点能自动增加到集群中。就像阿里云的可用区一样。
  autoheal: 假设你的网络可能不可靠，你更关注服务的连续性而不是数据完整性，这个时候可能有一个双节点群集。
  
- 非镜像普通模式，在其他机器上的只会保留当前机器上queue的元数据，所以如果有一台宕机了，并不能从其他机器上读取

### rabbitmq为啥不用了
- 我觉得rabbitmq的语言，生态现在都不行，erlang太小众
- 吞吐量不够，即使搭建了集群，但是同一个queue的数据我记得是只在同一台机器上，同一台机器上再怎么样吞吐能力都是不够的
- 消费积压如果严重了，也可能导致消费能力有限
- 搭建的集群容易出问题，之前自建的集群，因为少配置或者操作问题，脑裂了几次
- 目前我们公司rabbitmq上的只有老服务的延时队列，新的基本上都使用rocketMQ或kafka
- 官方的示例demo实际在生产环境上直接使用会有问题，connection,channel，rabbitmq崩了都会有问题

### demo
```go
package main

import (
        "log"
        "os"
        "strings"

        "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
        if err != nil {
                log.Fatalf("%s: %s", msg, err)
        }
}

func main() {
        conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()

        ch, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer ch.Close()

        err = ch.ExchangeDeclare(
                "logs_topic", // name
                "topic",      // type
                true,         // durable
                false,        // auto-deleted
                false,        // internal
                false,        // no-wait
                nil,          // arguments
        )
        failOnError(err, "Failed to declare an exchange")

        body := bodyFrom(os.Args)
        err = ch.Publish(
                "logs_topic",          // exchange
                severityFrom(os.Args), // routing key
                false, // mandatory
                false, // immediate
                amqp.Publishing{
                        ContentType: "text/plain",
                        Body:        []byte(body),
                })
        failOnError(err, "Failed to publish a message")

        log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
        var s string
        if (len(args) < 3) || os.Args[2] == "" {
                s = "hello"
        } else {
                s = strings.Join(args[2:], " ")
        }
        return s
}

func severityFrom(args []string) string {
        var s string
        if (len(args) < 2) || os.Args[1] == "" {
                s = "anonymous.info"
        } else {
                s = os.Args[1]
        }
        return s
}



package main

import (
"log"
"os"

"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_topic", s)
		err = ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"logs_topic", // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
```




