package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"myAutoGoRPC/ws"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
)

const (
	rpcPort = ":18080"
	wsPort  = ":18081"
)

func startRPCServer(wg *sync.WaitGroup) {
	defer wg.Done()

	rpc.Register(new(Server))
	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		log.Fatalln("fatal error: ", err)
	}
	fmt.Println("RPC server started on ", rpcPort)

	for {
		conn, err := lis.Accept() // 接收客户端连接请求
		if err != nil {
			log.Println("Error accepting connection:", err)
			return
		}
		go jsonrpc.ServeConn(conn) // 启动 Goroutine 处理连接
	}
}

func startWebSocketServer(wg *sync.WaitGroup) {
	defer wg.Done()
	r := mux.NewRouter()

	r.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		param := vars["id"] // 从 URL 中获取名为 "id" 的参数
		ws.HandleConnection(w, r, ws.CM, param)
	})

	fmt.Println("WebSocket server started on", wsPort)
	if err := http.ListenAndServe(wsPort, r); err != nil {
		fmt.Println("Error starting server:", err)
	}

}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go startRPCServer(&wg)       // 启动 RPC 服务器
	go startWebSocketServer(&wg) // 启动 WebSocket 服务器

	wg.Wait() // 等待所有 Goroutine 完成
}
