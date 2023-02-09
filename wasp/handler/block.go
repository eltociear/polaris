package handler

import (
	// "github.com/astaxie/beego/validation"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GetLatestBlock(c *gin.Context) {
	ctx := context.Background()
	blockModel, err := h.queryClient.GetLatestBlock(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Response(c, http.StatusNotFound, ERROR_BLOCK_NOT_FOUND, nil)
		return
	} else if err != nil {
		Response(c, http.StatusInternalServerError, ERROR, nil)
		return
	}
	fmt.Println(blockModel)
	Response(c, http.StatusOK, SUCCESS, blockModel)
}
func (h *Handler) GetBlocks(c *gin.Context)            {}
func (h *Handler) GetBlock(c *gin.Context)             {}
func (h *Handler) GetBlockTransactions(c *gin.Context) {}
