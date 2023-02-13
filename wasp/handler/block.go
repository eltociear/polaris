package handler

import (
	// "github.com/astaxie/beego/validation"
	"context"
	"errors"
	"net/http"

	"github.com/berachain/stargazer/wasp"
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
	Response(c, http.StatusOK, SUCCESS, blockModel)
}
func (h *Handler) GetBlock(c *gin.Context) {
	ctx := context.Background()
	height := c.Param("height")

	blockModel, err := h.queryClient.GetBlock(ctx, height)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Response(c, http.StatusNotFound, ERROR_BLOCK_NOT_FOUND, nil)
		return
	} else if err != nil {
		Response(c, http.StatusInternalServerError, ERROR, nil)
		return
	}
	Response(c, http.StatusOK, SUCCESS, blockModel)
}

// localhost:10081/blocks/?page=1&limit=4
func (h *Handler) GetBlocks(c *gin.Context) {
	ctx := context.Background()
	pagination := wasp.GeneratePaginationFromRequest(c)
	blockModel, err := h.queryClient.GetBlocks(ctx, pagination)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Response(c, http.StatusNotFound, ERROR_BLOCK_NOT_FOUND, nil)
		return
	} else if err != nil {
		Response(c, http.StatusInternalServerError, ERROR, nil)
		return
	}
	Response(c, http.StatusOK, SUCCESS, blockModel)
}
