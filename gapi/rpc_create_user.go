package gapi

import (
	"context"

	db "github.com/beabear/simplebank/db/sqlc"
	"github.com/beabear/simplebank/pb"
	"github.com/beabear/simplebank/util"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		if meErr, ok := err.(*mysql.MySQLError); ok {
			switch meErr.Number {
			case 1062:
				// 1452: foreign_key_violation, 1062: unique_violation
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			default:
				return nil, status.Errorf(codes.Internal, "Error %d (%s): %s", meErr.Number, meErr.SQLState, meErr.Message)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	user, err := server.store.GetUser(ctx, arg.Username)
	if err != nil {
		if meErr, ok := err.(*mysql.MySQLError); ok {
			return nil, status.Errorf(codes.Internal, "Error %d (%s): %s", meErr.Number, meErr.SQLState, meErr.Message)
		}
		return nil, status.Errorf(codes.NotFound, "failed to find the created user")
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}
