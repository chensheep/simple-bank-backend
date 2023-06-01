package gapi

import (
	"fmt"

	db "github.com/chensheep/simple-bank-backend/db/sqlc"
	"github.com/chensheep/simple-bank-backend/pb"
	"github.com/chensheep/simple-bank-backend/util"
	"github.com/chensheep/simple-bank-backend/worker"

	"github.com/chensheep/simple-bank-backend/token"
)

type Server struct {
	pb.UnimplementedSimpleBankServiceServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistrubutor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistrubutor) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
