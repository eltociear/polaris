package types

import (
	"github.com/berachain/stargazer/wasp/proto"
)

type GrpcServer interface {
	proto.BlockServiceServer
}
