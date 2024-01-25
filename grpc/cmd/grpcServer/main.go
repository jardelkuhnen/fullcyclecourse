package main

import (
	"database/sql"
	"github.com/jardelkuhnen/fullcyclecourse/grpc/internal/database"
	"github.com/jardelkuhnen/fullcyclecourse/grpc/internal/pb"
	"github.com/jardelkuhnen/fullcyclecourse/grpc/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDb := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDb)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}