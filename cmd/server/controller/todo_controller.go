package controller

import (
	"context"
	todov1 "todoService/gen/go/todo/v1"
)

type grpcTodoController struct {
	todov1.UnimplementedTodoServiceServer
}

func NewGrpcController() *grpcTodoController {
	return &grpcTodoController{}
}

func (s *grpcTodoController) GetTodo(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.GetTodoResponse, error) {
	return &todov1.GetTodoResponse{
		Id: 1,
		Title: "title",
		Description: "description",
	}, nil
}