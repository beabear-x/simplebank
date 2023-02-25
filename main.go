package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/beabear/simplebank/api"
	db "github.com/beabear/simplebank/db/sqlc"
	"github.com/beabear/simplebank/gapi"
	"github.com/beabear/simplebank/pb"
	"github.com/beabear/simplebank/util"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listenner, err := net.Listen("tcp", config.GRPCServerAddresss)
	if err != nil {
		log.Fatal("cannot create listenner")
	}

	log.Printf("start gRPC server at %s", listenner.Addr().String())
	err = grpcServer.Serve(listenner)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddresss)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
