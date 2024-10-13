package main

import (
	"errors"
	"fmt"
	"github.com/Dasongzi1366/AutoGo/accessibility"
	"github.com/Dasongzi1366/AutoGo/app"
	"github.com/Dasongzi1366/AutoGo/images"
	"myAutoGoRPC/ws"
	"strconv"
	"time"
)

type Arg struct {
	Parameter []interface{} `json:"parameter"`
	Id        int           `json:"id"`
}

// rpc服务结构体
type Server struct{}

var A11Y = accessibility.New()

// 处理字符串
func (this *Server) ParseStr(req Arg, res *interface{}) (err error) {

	wsId := strconv.Itoa(req.Id)
	// 发送ws消息：“收到来自服务端的请求”
	err = ws.CM.SendMessage(wsId, fmt.Sprintf("reqId_%v:收到来自服务端的请求", wsId)) // Id 应该是连接的唯一标识符
	if err != nil {
		return
	}

	if len(req.Parameter) < 2 {
		err = errors.New("参数数量错误")
		return
	}
	str1, ok1 := req.Parameter[0].(string)
	str2, ok2 := req.Parameter[1].(string)
	if !ok1 || !ok2 {
		err = errors.New("参数类型错误")
	}

	// 发送ws消息：“处理来自服务端的请求”
	err = ws.CM.SendMessage(wsId, fmt.Sprintf("reqId_%v:处理来自服务端的请求", wsId))
	if err != nil {
		return err
	}

	// 模拟消息处理时间
	time.Sleep(3 * time.Second)

	*res = fmt.Sprintf("卜: %s 命: %s", str1, str2)

	// 发送ws消息：“处理完成”
	err = ws.CM.SendMessage(wsId, fmt.Sprintf("reqId_%v:处理完成", wsId))
	time.Sleep(1 * time.Second)
	return err
}

func (this *Server) ClickScreen(req Arg, res *interface{}) (err error) {
	wsId := strconv.Itoa(req.Id)
	// 发送ws消息：“收到来自服务端的请求”
	err = ws.CM.SendMessage(wsId, fmt.Sprintf("reqId_%v:收到来自服务端的请求", wsId)) // Id 应该是连接的唯一标识符
	if err != nil {
		return
	}

	A11Y.Click("QQ")

	return
}

func (this *Server) RunPredict(req Arg, res *interface{}) (err error) {
	wsId := strconv.Itoa(req.Id)
	// 发送ws消息：“收到来自服务端的请求”
	err = ws.CM.SendMessage(wsId, fmt.Sprintf("reqId_%v:收到来自服务端的请求", wsId)) // Id 应该是连接的唯一标识符
	if err != nil {
		return
	}

	appPkg := "com.yongxianghui.yyhl"
	if !app.Launch(appPkg) {
		sendWsMsg(wsId, fmt.Sprint("启动失败"))
		err = errors.New("启动失败")
		return
	}
	time.Sleep(time.Second)
	err = findUI2Click(wsId, "六爻")
	if err != nil {
		return
	}
	time.Sleep(time.Second)
	err = findUI2Click(wsId, "随机卦")
	if err != nil {
		return
	}
	time.Sleep(time.Second)
	err = findUI2Click(wsId, "开始排盘")

	time.Sleep(time.Second)
	screen := images.CaptureScreen()
	*res = images.ToBytes(screen, "png", 0)
	time.Sleep(time.Second)
	app.ForceStop(appPkg)
	return
}

func findUI2Click(wsId, txt string) (err error) {
	UI := A11Y.Text(txt).WaitFor(5000)
	if UI == nil {
		err = fmt.Errorf("未找到 \"%s\" ", txt)
		return
	}
	if !UI.Click() {
		err = fmt.Errorf("点击\"%s\"失败", txt)
		return
	}
	_ = ws.CM.SendMessage(wsId, fmt.Sprintf("reqId_%v:\"%s\"点击成功", wsId, txt))
	return
}

func sendWsMsg(wsId, msg string) {
	_ = ws.CM.SendMessage(wsId, msg)
}
