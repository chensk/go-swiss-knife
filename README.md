# go-swiss-knife

面向golang语言的工具包，小巧灵活而强大，就像瑞士军刀一样

# 背景

使用go语言开发时，经常面临一些标准库无法完美解决的问题，例如分布式锁等，为了避免重复开发，将这些通用场景的解决方案抽象在这里，尽量做到高度抽象、开箱即用，并且灵活可拓展。 各包的名字可顾名思义地对应各场景的问题，总结如下：

* lock：分布式锁问题
* sessionlib: 会话实现，包括基于redis的分布式会话

# 安装

`go get github.com/chensk/go-swiss-knife@1.0.3`

# 示例

## 分布式锁

分布式锁是常见的问题，通常的解决办法是使用redis或者数据库来实现。不管是底层使用什么技术实现，我们希望对用户是透明的，因此将分布式锁和底层存储的实现分享，用户可通过指定provider来切换存储实现。 分布式锁的接口定义为：

```go
type Lock interface {
    Unlock() error
}
```

通过加锁接口获取一个锁实例，要释放锁时调用`Unlock`方法即可。加锁的方式分为阻塞式和非阻塞式，分别对应LockTask和TryLockTask：

```go
// 阻塞式加锁，如果锁被占用，阻塞至其被释放；name用于区分锁
func LockTask(name string, options ...Options) (Lock, error)
// 非阻塞式加锁，如果锁被占用，立即返回，此时error非空
func TryLockTask(name string, options ...Options) (Lock, error)
```

支持的选项及其含义如下：

```go
// 指定持锁最长时间，可用于避免机器宕机导致锁无法释放等问题
func WithLockAtMost(atMost time.Duration) Options
// 指定阻塞式加锁的等待超时时间，如果timeout>0，阻塞时最多等待timeout时间，超过时将直接返回error
func WithLockTimeout(timeout time.Duration) Options
// 指定分布式锁的存储实现，例如指定基于gorm的mysql实现，可通过WithProvider(provider.NewMysqlLockProvider(db))来指定
func WithProvider(provider LockProvider) Options
```

如果使用基于mysql的分布式锁，需要先在数据库中创建一张shedlock表：

```mysql
CREATE TABLE `shedlock`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name`       varchar(64)         NOT NULL,
    `lock_until` timestamp(3)        NULL DEFAULT NULL,
    `locked_at`  timestamp(3)        NULL DEFAULT NULL,
    `locked_by`  varchar(255)             DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8
```

一个使用gorm实现的分布式锁的例子如下：

```go
db, _ := gorm.Open(mysql.Open("root:root1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"),
&gorm.Config{SkipDefaultTransaction: true})
...
// 非阻塞式
loc, err := TryLockTask("test-lock", []Options{
    WithProvider(provider.NewMysqlLockProvider(db)),
    WithLockAtMost(10 * time.Millisecond),
}...)
// 阻塞式
loc, err = LockTask("test-lock", []Options{
    WithProvider(provider.NewMysqlLockProvider(db)),
    WithLockAtMost(10 * time.Millisecond),
}...)
// 业务逻辑
...
_ = loc.Unlock()
```

## 会话

session在web开发中非常常见，常用于处理用户登录认证的问题。一个常见的模式是将sessionId保存在客户端的cookie中，而服务端保存sessionId和用户登录code的映射关系，
用户登录成功后服务端会创建一个新的会话，并将会话id保存在客户端cookie中，下次请求时服务端会读取客户端cookie中的sessionId，继而查找登录code，如果成功找到则视为用户已登录成功。

go语言没有session的默认实现，这里提供了一个会话的基本实现，支持：

* 在本地会话和分布式会话之间切换
* 可设置会话有效期，过期后会话会自动清除

用户可通过CreateSession函数创建一个新的会话：

```go
func CreateSession(sessionIdGetter SessionIdGetter, sessionIdSetter SessionIdSetter, options []SessionOptions) (Session, error)

type SessionIdGetter func () string

type SessionIdSetter func (string)
```

sessionIdGetter和SessionIdSetter抽象了sessionId的读写过程，例如sessionId保存在请求cookie中，或者保存在请求header中，用户可根据自己的情况实现SessionIdGetter。
如果会话不存在，会创建一个随机字符串，并调用SessionIdSetter将sessionId写回用户指定的位置，例如cookie中或者header。

支持的选项包括：

```go
// 指定session过期时间，默认24h
func WithExpiration(expiration time.Duration) SessionOptions
// 如果需要使用redis保存会话，可指定redis集群的ip:port列表
func WithRedisClusters(clusters []string) SessionOptions
// 可选：指定redis的请求超时时间，默认5s
func WithRedisTimeout(timeout time.Duration) SessionOptions
// 可自己实现会话的保存方式，例如通过db等
func WithSessionStore(store SessionStore) SessionOptions
```

如果没有指定redis cluster，则使用本地内存来保存会话，注意这种方式在分布式场景下可能有会话状态不一致的问题，因此分布式场景下建议使用redis来保存。

使用session的例子如下：

```go
var sid string
// 第一次会话
s, _ := CreateSession(func() string {
    return sid
}, func (s string) {
    sid = s
}, []SessionOptions{
    // 使用本地内存保存session
    WithExpiration(2 * time.Second),
})

_ = s.Set("test_key", "test_value")
...
// 保存会话
_ := s.Save(context.Background())
...
// 第二次会话
s, _ := CreateSession(func () string {
    return sid
}, func (s string) {
    sid = s
}, []SessionOptions{
    WithExpiration(2 * time.Second),
})
// 读取第一次会话存储的数据
v, ok = s.Get("test_key")
```

# 更改日志

* v1.0.3 实现基于db的分布式锁、会话实现等。