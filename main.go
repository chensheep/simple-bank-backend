package main

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/chensheep/simple-bank-backend/api"
	db "github.com/chensheep/simple-bank-backend/db/sqlc"
	"github.com/chensheep/simple-bank-backend/gapi"
	"github.com/chensheep/simple-bank-backend/pb"
	"github.com/chensheep/simple-bank-backend/util"
	"github.com/chensheep/simple-bank-backend/worker"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/lib/pq"

	_ "embed"
)

//go:embed doc/swagger/*
var swaggerFS embed.FS

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	store := db.NewSQLStore(conn)

	redisClientOpt := asynq.RedisClientOpt{Addr: config.RedisServerAddress}
	taskDistributor := worker.NewRedisDistrubuter(redisClientOpt)
	go runTaskProcessor(redisClientOpt, store)

	go createGatewayServer(config, store, taskDistributor)
	createGRPCServer(config, store, taskDistributor)
	log.Info().Msg("main existed")
}

func createGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}

func runTaskProcessor(redisClientOpt asynq.RedisClientOpt, store db.Store) {
	processor := worker.NewRedisTaskProcessor(redisClientOpt, store)
	log.Info().Msg("start task processor")
	err := processor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start task processor")
	}
	log.Info().Msg("task processor existed")
}

func createGRPCServer(config util.Config, store db.Store, taskDistributor worker.TaskDistrubutor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	log.Info().Msgf("start gRPC server on %s", config.GRPCServerAddress)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal().Err(err).Msg("grpc server failed to serve")
	}
}

func createGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistrubutor) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jsonOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOpts)
	err = pb.RegisterSimpleBankServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	subFS, err := fs.Sub(swaggerFS, "doc/swagger")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load swagger files")
	}
	fs := http.FileServer(http.FS(subFS))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server on %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}

	log.Info().Msg("gRPC gateway existed")
}
