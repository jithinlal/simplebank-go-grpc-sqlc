package gapi

import (
	"fmt"

	"github.com/jithinlal/simplebank/worker"

	db "github.com/jithinlal/simplebank/db/sqlc"
	"github.com/jithinlal/simplebank/pb"
	"github.com/jithinlal/simplebank/token"
	"github.com/jithinlal/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
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
