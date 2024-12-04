package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ysy try!")
	w.Write([]byte("Hello, World!"))
}

func main() {
	// 创建一个HTTP服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(helloHandler),
	}
	fmt.Println("Starting server at :8080")

	// 启动一个goroutine来运行服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// 创建一个通道来监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞，直到接收到一个信号
	<-quit
	log.Println("Shutting down server...")

	// 创建一个上下文，设置超时时间为5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅地关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown: %v", err)
	}
	log.Println("Server gracefully stopped")
}
