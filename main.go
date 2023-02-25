package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/beabear/simplebank/api"
	db "github.com/beabear/simplebank/db/sqlc"
	"github.com/beabear/simplebank/gapi"
	"github.com/beabear/simplebank/pb"
	"github.com/beabear/simplebank/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGatewayServer(config, store)
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
		log.Fatal("cannot create listenner:", err)
	}

	log.Printf("start gRPC server at %s", listenner.Addr().String())
	err = grpcServer.Serve(listenner)
	if err != nil {
		log.Fatal("cannot start gRPC server: ", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimplebankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listenner, err := net.Listen("tcp", config.HTTPServerAddresss)
	if err != nil {
		log.Fatal("cannot create listenner: ", err)
	}

	log.Printf("start HTTP gateway server at %s", listenner.Addr().String())
	err = http.Serve(listenner, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server: ", err)
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
