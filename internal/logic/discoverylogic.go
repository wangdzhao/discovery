package logic

import (
	"context"
	"fmt"

	"github.com/wangdzhao/community/discovery/internal/svc"
	"github.com/wangdzhao/community/discovery/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DiscoveryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoveryLogic {
	return &DiscoveryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoveryLogic) Discovery(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return &types.Response{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}
