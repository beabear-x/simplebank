package gapi

import (
	"fmt"

	db "github.com/beabear/simplebank/db/sqlc"
	"github.com/beabear/simplebank/pb"
	"github.com/beabear/simplebank/token"
	"github.com/beabear/simplebank/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimplebankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
