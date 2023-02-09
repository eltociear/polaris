package handler

import (
	"github.com/berachain/stargazer/wasp/queryClient"
)

type Handler struct {
	queryClient *queryClient.QueryClient
}

func NewHandler(qc *queryClient.QueryClient) *Handler {
	return &Handler{
		queryClient: qc,
	}
}
