package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"todoService/cmd/server/controller"
	todov1 "todoService/gen/go/todo/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// loggingInterceptor は各RPC呼び出し前後にログを出力するインターセプターです。
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	// リクエスト元の情報を取得
	p, _ := peer.FromContext(ctx)
	log.Printf("Request - Method:%s Peer:%s StartTime:%s Request:%v", info.FullMethod, p.Addr, start.Format(time.RFC3339), req)

	// ハンドラーを呼び出す（実際のAPI処理）
	resp, err = handler(ctx, req)

	// 処理完了後のログ
	log.Printf("Response - Method:%s ElapsedTime:%s Error:%v Response:%v", info.FullMethod, time.Since(start), err, resp)

	return resp, err
}


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
	s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	todov1.RegisterTodoServiceServer(s, grpcTaskController)
	// 3. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 4.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}