package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/berachain/stargazer/wasp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GetTransactionByHash(c *gin.Context) {
	ctx := context.Background()
	hash := c.Param("hash")
	// address := common.HexToAddress(hash)
	txnModel, err := h.queryClient.GetTransactionByHash(ctx, hash)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Response(c, http.StatusNotFound, ERROR_TXN_NOT_FOUND, nil)
		return
	} else if err != nil {
		Response(c, http.StatusInternalServerError, ERROR, nil)
		return
	}
	Response(c, http.StatusOK, SUCCESS, txnModel)
}
func (h *Handler) GetLatestTransactions(c *gin.Context) {
	ctx := context.Background()
	pagination := wasp.GeneratePaginationFromRequest(c)

	txns, err := h.queryClient.GetLatestTransactions(ctx, pagination)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Response(c, http.StatusNotFound, ERROR_TXN_NOT_FOUND, nil)
		return
	} else if err != nil {
		Response(c, http.StatusInternalServerError, ERROR, nil)
		return
	}
	Response(c, http.StatusOK, SUCCESS, txns)
}
func (h *Handler) GetTransactionsByBlockNumber(c *gin.Context) {
	ctx := context.Background()
	pagination := wasp.GeneratePaginationFromRequest(c)
	height := c.Param("height")

	txns, err := h.queryClient.GetTransactionsByBlockNumber(ctx, height, pagination)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Response(c, http.StatusNotFound, ERROR_TXN_NOT_FOUND, nil)
		return
	} else if err != nil {
		Response(c, http.StatusInternalServerError, ERROR, nil)
		return
	}
	Response(c, http.StatusOK, SUCCESS, txns)
}
func (h *Handler) GetTransactionsByBlockHash(c *gin.Context) {
	ctx := context.Background()
	pagination := wasp.GeneratePaginationFromRequest(c)
	hash := c.Param("hash")

	txns, err := h.queryClient.GetTransactionsByBlockHash(ctx, hash, pagination)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Response(c, http.StatusNotFound, ERROR_TXN_NOT_FOUND, nil)
		return
	} else if err != nil {
		Response(c, http.StatusInternalServerError, ERROR, nil)
		return
	}
	Response(c, http.StatusOK, SUCCESS, txns)
}
func (h *Handler) GetTransactionCount(c *gin.Context)   {}
func (h *Handler) GetTransactionReceipt(c *gin.Context) {}
