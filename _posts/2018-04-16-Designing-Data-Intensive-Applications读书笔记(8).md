---
layout: post
title:  "Designing Data-Intensive Applications读书笔记(8)"
categories: DDIA
tags:  读书笔记 分布式
author: Borui
---

* content
{:toc}

# The Trouble with Distributed Systems
分布式系统和单机系统有很大的差异,分布式系统会面临很多新的故障问题.

## 故障与部分失败
单机的程序,如果硬件没有故障,那么相同操作下产生相同的结果.如果硬件出现了故障,那么程序将无法运作.因此对于运行编写良好程序的单机应用,要么完好的工作,要么彻底无法提供服务,不会存在中间的状态.在计算机的设计中有一个深思熟虑的设计:如果内部错误发生了,我们宁愿死机也不能返回错误的结果,因为处理起来,错误的结果很难处理也很困扰人.计算机将底层不确定都屏蔽掉,始终对应用层提供确定性.

在分布式系统里,消息在网络中传播所花费的时间是不确定的,而多个节点故障发生也是不确定的,因此这些不确定性加上部分失效导致分布式系统很难去处理这些问题.

### 云计算与超级计算
有两种理念来建立大规模的计算系统.
1. 一种理念是高性能计算(high-performance computing (HPC)).超级计算机拥有成千上万个cpu被用于计算敏感的科学任务,例如天气预报或者分子动态.
2. 一种是云计算,涉及到多用户的数据中心,通过以太网或者ip连接的商品电脑,弹性的资源分配以及计量计费.
3. 传统企业的数据中心,介于这两种之间.

在超级计算机上,一个作业不时的将当前的计算状态持久化到存储系统上,如果一个节点失败了,通常会将整个集群停下,在故障节点恢复后,计算工作通过最近的checkpoint继续进行.因此超级计算机更像是一个单机计算机,将部分失效演变成完全失效.但是我们更关注互联网服务,这跟超级计算机区别很大:
1. 许多互联网相关的应用是在线的,这意味着他们需要时刻和用户做低延迟的交互.将整个集群停下来不对外提供服务这是不可接受的,但是离线任务则完全没有问题.
2. 超级计算机通常运行在特殊的硬件上,通常每个节点相当可靠,节点间通过共享内存和远程直接内存访问设备(RDMA)来通信.而云服务里的节点都是商品机器,通过规模经济来使用较少的花销获取同等的性能,但是错误率很高.
3. 大型数据中心的网络通常预计于ip或者以太网,通过Clos拓扑结构提供高速双工带宽.超级计算机通常使用特殊的网络拓扑,例如多维网格或者环状,在已知的通信模式下为HPC提供更好的性能.
4. 系统越大,某个组件就更有可能产生故障,随着时间流逝,故障被修复新故障产生.但是在一个上千节点的系统里,假定故障随时发生是合理的.如果错误处理策略都只是简单的放弃处理,那么大型系统将花费大量的时间来故障恢复而不是做有用的事情.
5. 如果一个系统能够容错失效节点并且作为一个整体对外坚持提供服务,对于运维来说这是很有用的特征.
6. 物理分布式部署环境下(将数据放置在离用户更近的地方来降低访问延迟),通信通常都是经过因特网环境,对比局域网会更慢更不可靠.而超级计算机所有节点都被放置在一起.

如果我们想使得我们的分布式系统工作,我们必须接受**部分失效**和**建立容错机制**.换句话说,我们需要在不可靠的组件上建立可靠的系统.(这个可靠是有限度的,不可能十分完美的可靠).直观上来说,一个系统的可靠程度,取决于系统中最不可靠的组件,事实并非如此:
1. 纠错码允许数字信号通过一个偶尔会导致部分位错误的通信信道来传输准确的信息.例如无线网络中存在无线电干扰.
2. IP协议是不可靠的,会丢包延迟重复和重排序.但是tcp协议在ip基础上提供了更可靠的传输层.它能保证丢包被重发,重复消息被消除,接受到的消息顺序和发送一直.

尽管系统可以比它所依赖的组件更可靠,但是可靠性是有限度的.纠错码只能纠正少量的单字节错误,TCP也没有办法解决延迟问题.这些系统并不完美,但是通过屏蔽底层的一些错误,遗留的错误很容易被定位和解决.

## 不可靠网络
我们重点是关注无共享分布式系统,系统内的机器通过网络互联,网络也是系统间通信的唯一手段,每个机器拥有自己的内存和硬盘,但是机器之间无法相互访问内存和硬盘,只能通过网络访问彼此的服务接口.无共享不是构建系统的唯一方式,但这已经成为构建互联网服务的主导方式:
1. 因为不需要特殊的硬件,因此相对便宜,能够利用商业级别的云计算服务.
2. 通过多个地理隔离的分布式数据中心进行冗余,能获取高性能.

互联网和大部分数据中心的内部网络(通常是以太网),都是异步包网络(asynchronous packet networks).这种网络环境下,网络不提供保证数据包什么时间到达目的地,数据包是否能到达,如果你发送了请求然后等待响应,可能会有很多情况发生:
1. 你的请求丢失了(由于忘记插网线了)
2. 你的请求正在排队,然后晚些时候会被发送(获取是因为网络负载比较重)
3. 远程节点失效了
4. 远程节点临时停止响应了(由于长时间的GC),但后续会恢复工作
5. 远程节点响应了, 但是响应结果丢失在网络里了(或许由于交换机配置错误)
6. 远程节点响应了,但是响应结果被延迟了,晚些会发送(由于负载较高导致)

发送方甚至不清楚消息到底发送出去没.通常通过**超时**机制来解决.但实际上即使是超时了,你也不清楚到底发生了什么.

### 实际中的网络问题
网络故障时刻存在,即使是在公司维护的数据中心里,人为配置错误是主要的诱因.而共有网络里,可能交换机软件升级导致了数据延迟,鲨鱼咬断了海底光缆,甚至出现接收不到数据,但是能发送数据,因为网络链路不保证双工运行.

**网络分区**:当网络的一部分和另一部分无法连接,这被称为网络分区或者网路断裂.我们统称为网络故障.

处理网络故障并不意味着要容忍故障:如果你的网络相当可靠,那么处理方式可以是简单的将错误信息反馈给用户.但是你需要了解系统面临网络故障时的反应,以及需要确保系统能够从网络故障里恢复.

### 检测故障
需要系统需要自动的检测故障节点.
1. 负载均衡器需要停止将请求发送到僵死的节点.
2. 单主复制的分布式系统里,如果主节点失败了,其中一个从节点需要成为新的主节点.

然而网络使得检测故障节点很困难,在一些特定的场景下,可能会得到一些直白的反馈来告诉你节点是否还在工作:
1. 如果你能连接节点所在的机器,但是目标端口没有在监听(这意味着进程挂掉了),那么操作系统能够通过回复RST或者FIN包来关闭或者复用TCP连接.但是如果处理流程中进程崩溃了,你并不知道到底处理了多少数据.
2. 如果节点进程崩溃了(被管理员杀死进程了),但是节点所在的操作系统还在运行,然后通过脚本通知其他节点消息,使得其他节点不必等待超时.HBase就是这么处理的.
3. 如果你有权限连接到数据中心交换机的管理界面,就可以在硬件层面获取链接错误信息.但是如果你通过互联网连接,或者是共享的数据中心,没有访问权限,或者由于网络问题你无法访问管理界面,这就不可行了.
4. 如果一个路由器确定你想访问的ip地址是不可达的,会返回给你一个ICMP的目标不可达的包反馈.但是路由器也没有魔法能力来进行故障检测,它和其他网络组件服从一样的限制.

快速获悉远程节点宕机是很有用的,但你不能指望它.通常情况下,当发生故障时,你会什么响应也没有.你可以重试若干次等到超时,然后猜测远程节点已经宕机了.

### 超时和无限时延
如果超时是我们唯一检测故障的手段,那么超时时间应该设置多久呢?很不幸没有简单的答案.
1. 设置过长,节点被认为已经挂掉前需要一直等待,用户要么等要么看到错误信息.
2. 设置过短,检测速度很快,但是有很大概率会错误的认为节点已经宕机,也许只是偶尔那个时刻机器变慢了.

过早的声明一个节点失效了是有问题的,这会造成同样的操作被执行两次.而且当节点被声明失效了,那么其他节点将会接替该节点的工作,这会给其他节点增加负担.如果当前系统正在经历高负荷,而过早地声明一个失效了会加重这种境地.通常情况下节点都只是变慢了而没有完全是失效,这样一来会造成级联故障,所有节点都发觉其他节点因为变慢而声明他们都宕机了,于是整个系统都不工作了.

当网络延迟时间是有保证的时候,我们可以设定一个合理的超时时间.然而目前的异步网络遵循的是尽快的发包,但不保证发送时间的上限.
#### 网络拥塞和排队
1. 如果一些不同的节点同时往一个目标节点发包,网络交换机必须对他们进行排队,然后一个一个把它们推送到目标节点.当网络链路繁忙时,包必须等待,我们称之为网络拥塞.如果排队的包太多了,后来的包就会被丢弃需要被重发,尽管此时网络其实是好的.
2. 当包到达了目标机器,但是当前所有的cpu都在忙碌,那么到达的请求就需要在操作系统里排队,等待应用程序来处理.这取决于机器的负载程度,等多久是没有限制的.
3. 在虚拟化环境里,一个运行的os通常会暂停几十毫秒来等待其他虚拟机运行在cpu上.在这个阶段里,os没有办法消费任何网络请求,到达的数据统统都被虚拟机排队,这进一步增加了延迟.
4. TCP执行**流量控制**,也叫作**拥塞避免或者背压**.为了避免加重网络链路或者目标节点的负载,会限制发送节点的速度.这意味着发送时也需要额外的排队.

TCP会在一定超时范围内,将尚未确认的包认为已经丢失,超时时间通过之前观察到的往返之间来估算.丢失的包被自动重新发送,尽管应用层看不到丢包和重发,但是能看到的直观的结果就是时延.

当系统接近最大容量时,排队时延会非常大.拥有大量备用容量的系统可以很快的消耗队列,然而在高度被利用的系统中,长队列很快就会排起来.

在公有云和多租户机房里,资源是被许多消费者共享的.批量任务例如MapReduce很容易就把网络链路吃满.由于无法控制共享资源中其他用户的资源使用,如果你周围的节点使用了很多资源,那么网络延迟时间更加波动.

**在这种环境下,只能摸索着选择超时时间.通过较长时间内,跨多台机器之间来观察网络往返的时间分布,最终来决定延迟的可预期变化.然后再考虑应用程序的特征,在故障检测延迟和过早超时风险之间进行权衡.更好的方法就是不去配置固定超时,而是根据响应时间动态的调整超时.TCP就是这么做的.**

####TCP和UDP
一些对延迟敏感的应用会采用UDP,而不是TCP.这其实是在可靠性和延迟波动间做出的权衡.因为udp不会进行流量控制也不会重发丢失的包,因此避免了一部分网络延迟.udp适用于丢包延迟数据无所谓的场景,例如ip电话.

### 同步vs异步网络
如果分布式系统是部署在**最大时延固定**并且**不丢包**的网络环境下,问题将简单的多.那么为什么我们不部署在这种硬件环境呢?这里需要对比数据中心的网络和传统固定线路的电话网络.

当我们通过电话网络打电话时,首先需要建立一个**电路**:一个固定的,保证足够带宽,贯穿两端的线路,这个电路直到通话结束才关闭.这种网络是同步的,传输中的数据不需要排队,带宽固定因此时延固定,我们称之为**有界时延**.而传统tcp不同,一个电路具有固定带宽则独立使用,而tcp则是尽可能的使用网络带宽,当进行传输大小变化的数据块时,会尽量在最短的时间发送出去,并且当tcp连接空闲时是不占用带宽的.以太网和IP协议都是包交换协议,并没有电路的概念.

那么为什么会这么选择呢,这是为了**突发流量**做出的优化,因为打电话期间每秒的流量是固定的,而对于互联网来说,时刻发送的数据都是不固定的,因此需要尽可能的发送出去.如果我们通过电路传输数据,就需要提前申请带宽:如果申请小了,那么传送的速度太慢了,导致网络能力被浪费.如果申请的多了,电路建立不起来(因为没有足够的带宽支持).因此采用电路传送突发的流量,会导致网络能力没有被充分利用从而速度很慢.像TCP会动态的调整传送速度来利用网络的能力.

也有混合两种网络的网络,叫做ATM.无线带宽技术有些类似:它实现了链路层上端到端的流量控制,从而减少了网络排队的需求,但是仍然会经历链路拥塞的延迟.利用Qos(quality of service,包的优先级和调度)和接纳控制(admission control, 限速的发送方),能够在包网络中模拟电路交换,从而提供有界的延迟.但是在多租户数据中心和公有云,或者互联网通信场景下,这种服务质量还没有启用.

#### 延迟和资源利用
更一般的讲,我们可以吧变动的延迟看作是动态资源分区的结果.电路交换是一种**静态分配**方式,就算只有一通电话在拨打,其他的带宽也没有被利用.包交换网络是一种**动态分配**方式,数据被尽可能快的发送出去,这种方式有排队的缺点,但是能够最大化利用带宽.

在cpu利用上有相同的境遇.如果在多线程中动态的共享cpu核,那么线程之间需要互相等待和排队,因此挂起时间无法预测.但是硬件的利用率要比你给线程分配固定数量的cpu周期要好.更好的硬件利用率是采用虚拟机的动机.

在静态资源分配下,延迟是有保证的.但是也降低了利用率,换言之就是代价太大.而动态资源分配下,资源利用率提高了,但是延迟也变得不可控.

所以,变化的延迟不是自然产生的,而是代价/收益权衡的结果.
## 不可靠的时钟
对时间的利用主要是利用**时间段**和**时间点**.网络中的每台机器都有自己的时钟,实际上是一种硬件设备,通常是石英晶体振荡器.这种设备并不是完全精确的,因此每台机器都有自己的时间概念,可能快也可能慢.可以在某种程度上同步这个时间,目前最常使用的机制就是**NTP(Network Time Protocol)**:该协议通过一组服务集群提供的时间来调整计算机的时间,而这些集群则通过更准确的时间源——GPS.
### 单调性与时刻时钟
现代计算机至少拥有两种不同的时钟:时刻时钟和单调时钟.
#### 时刻时钟
时刻时钟和直观上我们认知的时钟一直,返回当前的日期和时间.时刻时钟通常通过NTP进行同步,这意味着两台机器上的时间戳是一致的.但是时刻时钟也有很古怪的地方,例如当前时钟速度快于NTP服务器,某个时刻会被NTP服务器纠正回到之前的某个时刻.这种跳跃和忽略闰秒使得时刻时钟不能用来评估时间流逝.
#### 单调时钟
单调时钟很适合评估时间段,例如超时时间和响应时间.单调时钟顾名思义总能保证时间是向前的,而不像时刻时钟会跳回过去某个时刻.单调时钟能用来评估时间的流逝.但是单调时钟的绝对值是没有意义的,不同计算机之间的单调时钟值没有任何比较的意义.

在多cpu套接字的服务器上,每个cpu有自己的计时器,并且彼此之间不用同步.操作系统会补偿这种差异,并试图赋予应用线程单调的时钟视图,尽管这些线程会调度到不同的cpu上.但是尽量不要太依赖这种单调性保证.

NTP会调整单调时钟的移动速率,但是不会影响单调性,而且单调性时钟的精度通常很高,能够精确到微秒甚至更小.在分布式系统里,单调时钟适用于度量时间的流逝,因为不需要考虑不同节点时钟的同步.并且对微小的差异不敏感.

### 时钟同步与精度
单调时钟不需要同步,但是时刻时钟需要借助外部时间源来同步,但也没有我们想象的那么精确:
1. 计算机上的石英钟不是非常精确,经常漂移.而且漂移的变化取决于计算机的温度.Google认为会有百万分之一的误差.如果30s同步一次,误差为6ms,如果每天同步一次,误差为17s.这个漂移限制了我们能获取的最好精度.
2. 如果计算机时钟和NTP服务器相差太大,可能会拒绝同步或者强制被重置.任何在这个时刻附近的应用都会看到时间回去了或者跳跃向前了.
3. 如果一个节点意外地屏蔽了NTP服务器,这种错误的配置可能会被忽略一段时间.
4. NTP同步只能和网络延迟一样好,因此在一个拥挤的网络上,它的准确性是有限的.通过互联网同步最小35ms的错误是可以实现的,而且在网络峰值时错误会在1s左右.这取决于配置,较大的网络延迟会导致NTP客户端完全放弃同步.
5. 一些NTP服务器可能被错误配置导致数小时的差异.NTP客户端相当鲁棒,因为会请求若干个服务器并忽略异常值.
6. 闰秒会导致一分钟内包含59s或者61s,这会打乱没有考虑闰秒的系统针对时间的假定.最好的办法是让NTP服务器在一天内逐渐完成对闰秒的调整.
7. 在虚拟机上,硬件时钟也是虚拟的,当虚拟机暂停运行时,从应用层的观点来看时钟突然跳到前面去了.
8. 如果运行软件的设备并不完全收到你控制(手机或者嵌入式设备),那么就更不能信任硬件时钟了.一些用户可以故意将硬件时钟设置成不正确的日期和时间.

获取高精度时钟是可能的,需要投入大量的资源.

### 依靠同步时钟
时钟问题有部分原因是因为很容易被忽略:如果cpu有缺陷或者网络被错误配置,机器不大可能正常工作,因此会很快被发现和修复.但是如果石英钟有缺陷或者NTP客户端配置错误,大部分事情看起来都是好的,尽管时钟在逐渐的漂移离实际时间越来越远.如果有软件依赖精确的时钟同步,那么结果可能是沉默或者微秒的数据丢失,而不是系统崩溃.

因此当你的应用比较依赖同步的时钟时,必须要监控所有机器上时间的偏差.任何偏差较大的节点都需要从集群摘除,这样能避免损失.
#### 有序事件的时间戳
LWW的根本问题:
+ 数据库写入神秘地消失:一个慢时钟的节点无法复写之前快时钟节点写入的数据.这导致数据毫无征兆的消失了.
+ LWW无法区分连续发生的写入和并发的写入.为了防止违反因果关系,额外的因果追踪机制需要被采用,例如版本向量.
+ 两个节点可能同时产生时间戳一致的写入操作,需要额外的决定参数(通常是一个大的随机数)来解决冲突,但这也会导致因果关系颠倒.

我们无法准确的定义最新的数据,而NTP的同步精度也受到网络延迟的限制.因此一般基于自增长计数的**逻辑时钟**被采用.
#### 时钟读数有一个置信区间
公网上通过NTP服务同步时间,精度可能是几十ms,而当网络拥堵时可能会高达100ms.因此时钟的读数往往不意味着某个时刻,更应该被认为是一个时间区间.这个区间的边界取决于你的时间源,如果是GPS接收器或者原子时钟,那么生产商会提供.如果是通过NTP服务,那么就需要将自上次同步后的石英漂移+NTP服务器不确定+网络往返时间.不过大部分系统并不提供一个不确定的区间给你参考.

谷歌的Spanner里包含TrueTime的API,明白地提供本地时钟的置信区间,一般返回[earliest,latest],分别表示最早的可能时间戳和最晚的可能时间戳.

#### 用于全局快照的同步时间
快照隔离性中,最常见是通过单调的递增事务ID.在单机系统里,简单的计数器就足够了.但是分布式数据库中,一个全局的单调的递增事务ID就很难产生了,因为需要协调.由于许多小而且快速的事务,因此分布式系统产生事务ID成为了系统瓶颈.

Spanner通过TrueTime返回始终置信区间来实现快照隔离性.对于A = [A_earliest, A_latest] and B = [B_earliest, B_latest],如果两个区间没有重叠,例如A_earliest <= A_latest < B_earliest <= B_latest,那么就可以说B清楚地发生在A之后.为了确保时间戳能反映出因果关系,Spanner会在提交读写事务前有意等待置信区间的长度.那么就需要确保置信区间的长度足够短,Google在每个数据中心部署了GPS接收器或者原子时钟,使得时钟被同步在7ms以内.

### 进程暂停
有这样一种场危险的时钟使用场景在分布式系统中,单主架构里,leader如何确保自己还是leader呢?一种选择是从其他节点获取租约,流程大体如下:
```java
while (true) {
    request = getIncomingRequest();
    // Ensure that the lease always has at least 10 seconds remaining
    if (lease.expiryTimeMillis - System.currentTimeMillis() < 10000) {
        lease = lease.renew();
    }
    if (lease.isValid()) {
        process(request);
    }
}
```
这里是有问题的,如果使用时刻时钟,那么超时时间戳和本地时间戳是不同机器给出的机器,如果时间同步有异常,这种策略是有问题的.

及时我们采用本地单调时钟,会有另外一个问题.因为获取当前时间戳和处理请求之间我们一般认为会运行的很快,因此设定10s的过期时间是足够的.但是如果程序执行中出现了进程暂停呢?比如在处理请求时暂停了15s,而恢复运行后线程并不知道被暂停了那么久,而此时处理请求其实是不安全的.线程暂停的因素很多:
1. 许多编程语言运行时会有GC,gc过程会stw,有时甚至会高达几分钟无响应.即使是并发收集器如CMS,仍然存在stw的操作.尽管可以通过调节gc参数来降低暂停,但是为了鲁棒性保证我们必须做最坏的打算.
2. 在虚拟环境下,虚拟机会被暂停(停止所有线程,然后保存内存内容到磁盘)和恢复(恢复内存内容然后继续运行).这个过程会发生在进程执行的任意时刻并持续任意时长.这种特性有时会被用来做从一台主机到另一台主机的无重启实时迁移,至于暂停多久,取决于恢复数据到内存有多快.
3. 使用笔记本的用户关上电脑盖子时,执行也会在任意时刻被暂停或继续执行.
4. 当操作系统进行上下文切换时,当虚拟机切换不同虚拟机时,执行线程会在任意时刻被暂停.在其他虚拟机上运行的cpu时间被称为**窃取时间**,如果负载很重或者排队线程很多,那么线程恢复运行要等待一段时间.
5. 如果线程执行同步磁盘操作,线程会被挂起等待I/O操作完成.很多语言里及时没有显示的磁盘操作也会发生磁盘操作,例如类加载器的延迟加载.I/O阻塞和gc的暂停会加剧暂停问题.如果磁盘是网络文件系统或者网络块设备,那么I/O的时延还要受制于网络的时延.
6. 如果操作系统被设置为允许交换到磁盘,内存访问引发页错误会使得页被从磁盘加载到内存,I/O操作发生会导致线程暂停.如果内存压力过大,也会使得不同的页从内存交换到磁盘上.极端情况下,系统将花大量的时间在交换页上,因此一般情况下系统会禁止这个操作.
7. unix进程会被SIGSTOP中断.

任意时间的发生都会抢占线程的运行,然后后续再恢复运行.在单机上可以通过同步等手段来协调多线程程序运行,但在分布式环境下,当一台机器上的进程暂停了,其他机器上的进程还在继续运行着,甚至认为该机器已经死机了.

#### 响应时间保证
存在一些系统需要保证某个时刻必须响应,如果不响应将会造成整个系统的失败,这种被称为**强实时(hard real-time)**系统.在嵌入式系统中,实时意味着系统被小心设计和测试,从而能够在任何场景下满足特殊的时间保证.而在web场景下,意味着服务端向客户端推送以及没有硬性响应时间约束的流处理过程.

提供实时保证需要各个级别的支持:实时操作系统,在某段时间内保证cpu时间的前提下,允许进程被调度.库函数需要说明最坏的执行时间.动态内存申请需要被限制甚至被取消.大量的测试和评估工作必须被执行来确保满足了条件.

实时系统的实现代价很大,往往因为实时性导致吞吐量较低.对于服务端数据处理系统,提供实时保证是不经济的.因此分布式系统依然要忍受随时的暂停和时钟不稳定.

#### 限制gc影响
不一定非要昂贵的实时调度保证,进程暂停的负面影响也可以减缓些.当编程语言运行时计划gc时,这个过程是可变动的.一种新型的想法是,把gc看做一个简单计划好的暂停节点动作.当一个节点将要gc时,运行时会发送请求给应用,应用层会停止发送新的请求到该节点上,该节点处理完尚未结束的请求,然后就可以在没有请求的环境下进行gc了.这种方案屏蔽了暂停并且降低了响应延迟.许多延迟敏感的财务交易系统就采用这种方式.

另外一个变种是,对短期对象(可快速回收)执行垃圾回收,然后在积累足够多的长期存活对象以至于需要full gc之前,定期进行系统重启.一次只有一个节点进行重启,重启时流量被切换到其他节点上,这类似轮转升级.

这些方案都不能避免gc停顿,但能够减少对应用的影响.

## 知识,事实和谎言
分布式系统不同于单机:没有共享的内存,通过没有固定延迟的不可靠网络通信,会遭受部分失效,不可靠的时钟和进程暂停.

网络中的节点不可能准确的知道一切事情,只能基于通过网络接收到的信息(或者接收不到信息)来进行猜测.节点只能通过在网络中交换信息来获悉其他节点的状态.如果远程节点没有回复,那么将无从知晓它的状态,因为没办法区分是网络问题还是节点自身的问题.

在分布式系统中,我们可以针对我们的系统模型做出假设,然后以满足这些假设的方式设计出实际的系统.即使底层不可靠的系统模型提供较少的保证,我们依然可以得到可靠的行为.尽管在一个不可靠的系统模型上使得我们的软件运行良好是可能的,但这并不容易做到.因此我们需要知道我们能做哪些假设,以及我们需要提供哪些保证.
### 事实被多数定义
1. 节点可以接受信息,但无法发送消息.很快节点被其他节点认为死机了,但该节点自行无法得知.
2. 对比场景1,该节点发现发送的节点都没有响应,因此意识到网络问题.但什么也做不了.
3. 节点经历长时间的GC,GC阶段其他节点得不到响应从而宣称该节点死机了.而gc结束后,该节点毫无意识这段时间自己没有响应.

因此一个节点一定不能相信自己对某个情况的判断.许多分布式算法都依赖于quorum机制,也就是节点投票:为了减少对任意一个节点的依赖,需要来自若干节点的最小投票数来做决定.

大部分情况下,quorum是超过半数节点的绝大多数,其他类型的quorum也有在用.
#### leader和锁
+ 只有一个节点被允许成为数据库分区的leader,从而避免脑裂.(分区的leader相信自己是leader)
+ 只有一个事务或者客户端被允许获取特定资源或者对象的锁,从而阻止并发写.(锁的拥有者相信自己是唯一的拥有者)
+ 只有一个用户被允许注册一个特定的名字,名字必须唯一标示用户.(请求的处理进程相信当前用户获取了特定的用户名)

如果一个节点继续以自封的leader身份发送消息,在设计不良好的系统中其他节点认可了这个消息,就会造成系统运行不正常.Hbase中曾经出现过这种bug:
![Figure 8-4. Incorrect implementation of a distributed lock: client 1 believes that it still has a valid lease, even though it has expired, and thus corrupts a file in storage.](https://raw.githubusercontent.com/codeborui/codeborui.github.io/master/img/11.jpg)

这种bug是因为租约未到期前,发生了较长时间的gc停顿,然后租约过期了.gc结束后该节点恢复过来,认为租约未过期并且自己仍然有写权限,从而导致了异常.
#### fencing令牌
为了解决上述问题,提出了fencing令牌的方法:
![Figure 8-5. Making access to storage safe by allowing writes only in the order of increasing fencing tokens.](https://raw.githubusercontent.com/codeborui/codeborui.github.io/master/img/12.jpg)
当从锁服务器获取锁或者租约时,也会返回一个fencing令牌,这个令牌是自增长的.任何客户端发送写请求时必须包含该令牌.这样,当新的获得租约的客户端和旧租约客户端冲突时,较大令牌值的客户端胜出.zk使用zxid(事务id)或者cversion(节点版本)作为fencing令牌.

这个机制需要**资源自身**参与到解决令牌冲突的过程中.如果存储层不支持,那么就只能有限制的使用这种方法(比如在文件名中包含令牌id).对于服务端来说,增加校验令牌的工作不见得是一件坏事情.
### 拜占庭错误
fencing令牌可以避免非故意的错误.但是如果一个节点故意想要破坏系统的保障,比如使用一个伪造的fencing令牌.目前我们都认为节点是不可靠的,但是很诚实.如果分布式环境中,节点故意伪造信息那么分布式问题会变得更难.在一个不可信的环境里达成一致,称为拜占庭将军问题.

如果一个系统在即使节点出故障并且不遵守协议的情况下,甚至有人恶意攻击干扰网络,仍然能够正确运行,这个系统被称为**拜占庭容错**的:
+ 在航天环境中,计算机内存和cpu寄存器的数据由于受到辐射的影响而损坏,导致它会以任意不可预测的方式响应其他节点.而系统直接失败代价太大了(航天器故障会杀死机上的所有人),因此飞行控制系统必须容忍拜占庭错误.
+ 在多参与机构的系统中,一些参与机构可能会试着欺骗或者诈骗别的机构.在这种环境下,简单相信其他节点的信息是不安全的,甚至有可能是恶意的企图信息.例如在p2p网络里的比特币或者其他区块链能够不依赖中心授权使得互不信任的双方达成一致.

但是我们现在不考虑拜占庭错误.在我们自己的数据中心里,节点都是我们自己的机构控制,辐射也足够低因此也不会造成内存数据出错.拜占庭容错的协议极其复杂而且一来硬件设备,大部分服务端系统中,拜占庭容错的部署方案是不切实际不可行的.我们的系统应该预期来自终端用户的输入会有各种各样的异常,因此我们才会做输入校验,数据清理,输出转义来避免sql注入和跨网站脚本攻击.但是,我们并不适用拜占庭容错协议,只是简单的让服务器决定哪些行为允许哪些不允许.拜占庭容错更多是和p2p环境相关.

程序的bug可以被看做拜占庭错误,但是因为你部署在了所有的节点,所以拜占庭容错协议也毫无办法.大部分拜占庭容错协议需要超过2/3的大多数,也就是说之多1/4的节点出现bug.

如果一个协议能够保护我们避免恶意攻击那很好,但是这不现实.因为攻击者能够攻击一个节点就能攻击所有的节点.因此传统的机制(授权,访问控制,加密,防火墙等等)仍然是对抗攻击者的主要保护措施.
#### 弱化的谎言
尽管我们假定节点通常是诚实的,但还是要采取一些措施还对空弱化的谎言.比如,由于硬件故障导致的错误消息,软件bug和错误配置.这些措施不是完全拜占庭容错的:
+ 由于硬件和操作系统,驱动和路由器的bug,网络包有时候会被损坏.通常损坏的包有tcp和udp里的校验码捕捉到,但有时候仍会逃避检测.简单的方式就是在应用层对数据正确性做校验.
+ 公开访问权限的应用必须小心对待用户的输入.防火墙内的服务可以减少输入限制.
+ NTP客户端配置多台服务器的地址.当同步的时候,消息会发送给所有服务器.客户端会少数服从多数.从而规避某台NTP服务器配置不正确导致时间错误.

### 系统模型与现实
解决分布式问题的算法,需要在不过分依赖硬件和软件配置的基础上来实现.因此我们需要形式化我们在系统中可能遇到的错误,我们称之为**系统模型**,这是对我们算法可能遇到的问题进行的一种抽象.

关于时间议题,主要有三种模型:
1. **同步模型**:同步模型假设网络延迟是有界的,进程暂停是有边界的,时钟错误也是有边界的.这不是说时钟是完全同步的,然后网络零延迟.而是说所有这些异常都是有一个上限的.同步模型并不是实际中系统常用的模型,因为这些异常往往都是无边界的.
2. **部分同步模型**:这意味着大部分时候系统都像是同步的系统,但是有时候会超过上述异常的边界.实际中大部分系统都是这种模型.
3. **异步模型**:在这种模型里,算法不做任何时间假设.这种模型中的算法功能很有限.

关于节点议题,也有三种模型:
1. **Crash-stop faults**:这种模型中,算法假设节点只会产生一种异常,就是崩溃了.任意时刻节点停止响应然后不会再恢复回来了.
2. **Crash-recovery faults**:这种模型里,节点可能在任意时刻崩溃了,但是可能在任意时刻恢复回来.这种模型里假设节点有稳定的存储来避免崩溃,但是内存的数据会丢失.
3. **Byzantine (arbitrary) faults**:节点可能出现任何问题.

现实中,包含**Crash-recovery faults**的**部分同步**模型是最有用的.

#### 算法正确性
我们说算法是正确的这意味着什么呢?如果一个算法是正确的,它需要满足一些特质.例如我们生成fencing令牌,那么我们需要算法提供这些特质:
+ 唯一性:获取fencing令牌的请求不会返回同样的值.
+ 单调序列:x请求获取token_x,y请求获取token_y,如果x先与y,那么token_x < token_y
+ 可用性:请求的节点只要不崩溃,最终就一定能获取到响应.

因此如果在某些系统模型下,在这些模型限制下发生的所有异常场景下,算法都能满足这些特征,我们说算法就是正确的.但这毫无意义,因为一旦所有节点都崩溃了,算法什么也做不了了.

#### 安全性和活跃性
需要区分两类特征:**安全性**和**活跃性**,上述例子里,唯一性和单调序列是安全性,可用性是活跃性问题.一般带有**最终**字样的定义,都是活跃性问题.

安全性:坏的事情不会发生.

活跃性:正确的事情最终一定会发生.

+ 如果安全性被违反了,违反的结果无法被撤回,因此危害已经造成了.
+ 如果活跃性被违反了,下次可能不会违反.

区分安全性和活跃性能够帮助我们处理复杂的系统.对于分布式算法,安全性必须被保证.就算所有机器都宕机了,至少不会返回一个错误的结果.对于活跃性我们可以给出一些警告,只有当大部分节点都正常,并且网络最终会恢复的情况下,请求可以收到一个响应.

#### 映射系统模型到真实世界
在推理出分布式算法的正确性上,**安全性,活跃性和系统模型**很有用.但是系统模型毕竟还是现实世界的简化.

比如crash-recovery模型中假设通过稳定存储里的数据恢复,但是假如这些数据也丢失了呢?系统模型只是第一步,实际中产生的错误还是需要理论分析和海量的测试.

## 摘要
分布式环境下我们会遇到的问题:
+ 无论何时通过网络发送数据包,都有可能丢失或者任意的延迟.同样响应也可能丢失或者任意延迟,所以如果你接收不到响应,你也不知道消息发没发送到.
+ 一个节点的时钟可能和其他节点并不同步,尽管你尽最大的努力去搭建NTP.时间可能突然跳到后面或者前面,所以依赖时间是危险的,因为很大可能性你评估不出来时间差的区间.
+ 一个进程可能在执行过程中的任意时刻任意位置突然中断,然后就被其他节点认为已经宕机了,而进程恢复后并不知道这一切.

部分失效是分布式系统的典型特征.在分布式系统,必须建立部分失效容错的机制,使得系统作为整体还是可以继续工作的.

建立容错机制,第一步就是检测,但是这很难.大部分系统没有准确的机制来检测节点是否失败了,大部分分布式算法依赖超时机制.但是超时不能区分究竟是网络问题还是节点问题

一旦检测出故障,容忍故障也不容易:没有全局变量,也没有共享内存,机器之间也没有共享机器状态的能力.节点之间甚至都不能就时间达成一致,唯一的方式就将交互信息在不可靠的网络上.主要的决策不能依赖单个节点,需要多个节点共同决定.

分布式系统在扩展性,容错和低延迟上都比单节点有优势.