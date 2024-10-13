# myAutoRunRPCServ
被远程调用的go jsonrpc服务

没什么好说的，非常简单，这里用到了一个一个使用 Golang 构建的 Android 自动化测试框架（AutoGo）：
```json
https://github.com/Dasongzi1366/AutoGo
```
怎么用点进去有详细说明

--- 
然后本项目主要用到了一个websocket作为消息回调
然后用到了一个jsonrpc作为远程调用

如果你觉得一个Server不够用
可以再注册一个
```go
rpc.Register(new(Server))
```
在rpc.go里面实现 Server的方法就好，方法名大写，格式按照我的抄，这样的格式就是一个服务了
```go
func (this *Server) <你的方法名>(req Arg, res *interface{}) (err error){}
```
req里面有一个id，这个id是客户端传过来的，用来和客户端进行websocket通信的唯一标识
```go
type Arg struct {
	Parameter []interface{} `json:"parameter"`
	Id        int           `json:"id"`
}

```
你在连接websocket服务时的id就是传过来的id，这俩id必须一致
其他就没什么了，写个方法，方法里面写逻辑，写完在客户端直接调用方法名就可以在手机端执行了

编译的化还是按照AutoGo的方式编译成二进制文件
用adb投递到手机端执行就好
AutoGo 里面有详细的介绍