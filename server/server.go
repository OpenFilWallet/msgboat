package server

import (
	"context"
	"errors"
	"github.com/OpenFilWallet/OpenFilWallet/chain"
	"github.com/OpenFilWallet/OpenFilWallet/client"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"time"
)

var log = logging.Logger("wallet-server")

type Boat struct {
	apis map[string]client.LotusClient
}

func NewBoat(nodes map[string]string) (*Boat, error) {
	var apis = map[string]client.LotusClient{}
	for name, rpcAddr := range nodes {
		cli, err := client.NewLotusClient(rpcAddr, "")
		if err != nil {
			log.Warnw("NewLotusClient", "err", err)
			continue
		}
		_, err = cli.Api.ChainHead(context.Background())
		if err != nil {
			log.Warnw("ChainHead", "err", err)
			continue
		}

		apis[name] = *cli
	}

	if len(apis) == 0 {
		return nil, errors.New("no nodes available")
	}

	return &Boat{
		apis: apis,
	}, nil
}

// Send Post
func (b *Boat) Send(c *gin.Context) {
	param := chain.SignedMessage{}
	err := c.BindJSON(&param)
	if err != nil {
		ReturnError(c, ParamErr)
		return
	}

	signedMsg, err := chain.DecodeSignedMessage(&param)
	if err != nil {
		ReturnError(c, NewError(500, err.Error()))
		return
	}

	cid, err := b.mpoolPush(signedMsg)
	if err != nil {
		ReturnError(c, NewError(500, err.Error()))
		return
	}

	ReturnOk(c, client.Response{
		Code:    200,
		Message: cid.String(),
	})
}

func (b *Boat) mpoolPush(signedMsg *types.SignedMessage) (cid.Cid, error) {
	ctx := context.Background()

	for name, api := range b.apis {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		cid, err := api.Api.MpoolPush(ctxWithTimeout, signedMsg)
		cancel()
		if err != nil {
			log.Warnw("MpoolPush fail", "name", name, "err", err)
			continue
		}
		return cid, nil
	}

	return cid.Cid{}, errors.New("")
}
