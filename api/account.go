package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	db "github.com/beabear/simplebank/db/sqlc"
	"github.com/beabear/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	result, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if meErr, ok := err.(*mysql.MySQLError); ok {
			switch meErr.Number {
			case 1452, 1062:
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

	accountId, err := result.LastInsertId()
	if err != nil {
		if meErr, ok := err.(*mysql.MySQLError); ok {
			log.Printf("Error %d (%s): %s", meErr.Number, meErr.SQLState, meErr.Message)
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if meErr, ok := err.(*mysql.MySQLError); ok {
			log.Printf("Error %d (%s): %s", meErr.Number, meErr.SQLState, meErr.Message)
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
