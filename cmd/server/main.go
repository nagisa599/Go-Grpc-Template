package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"todoService/cmd/server/controller"
	todov1 "todoService/gen/go/todo/v1"

	"google.golang.org/grpc"
)

// 自作サービス構造体のコンストラクタを定義

func main() {
	// 1. 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	grpcTaskController := controller.NewGrpcController()
	// 2. gRPCサーバーを作成
	s := grpc.NewServer()
	todov1.RegisterTodoServiceServer(s, grpcTaskController)
	// 3. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// 4.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}