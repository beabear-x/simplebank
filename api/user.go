package api

import (
	"log"
	"net/http"
	"time"

	db "github.com/beabear/simplebank/db/sqlc"
	"github.com/beabear/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		if meErr, ok := err.(*mysql.MySQLError); ok {
			switch meErr.Number {
			case 1062:
				// 1452: foreign_key_violation, 1062: unique_violation
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			default:
				log.Printf("Error %d (%s): %s", meErr.Number, meErr.SQLState, meErr.Message)
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, arg.Username)
	if err != nil {
		if meErr, ok := err.(*mysql.MySQLError); ok {
			log.Printf("Error %d (%s): %s", meErr.Number, meErr.SQLState, meErr.Message)
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
