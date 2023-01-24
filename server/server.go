package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/OpenFilWallet/OpenFilWallet/chain"
	"github.com/OpenFilWallet/OpenFilWallet/client"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/juju/ratelimit"
	"sync"
	"time"
)

var log = logging.Logger("msgboat")

type Boat struct {
	lk             sync.Mutex
	nodes          map[string]string
	apis           map[string]client.LotusClient
	limiterBuckets *ratelimit.Bucket
}

func NewBoat(nodes map[string]string) (*Boat, error) {
	bucket := ratelimit.NewBucketWithQuantum(
		time.Second,
		10,
		10,
	)

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
		nodes:          nodes,
		apis:           apis,
		limiterBuckets: bucket,
	}, nil
}

// Send Post
func (b *Boat) Send(c *gin.Context) {
	b.lk.Lock()
	defer b.lk.Unlock()

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
		cid, err := api.Api.MpoolPush(ctx, signedMsg)
		if err != nil {
			log.Warnw("MpoolPush fail", "name", name, "err", err)
			continue
		}

		msg, _ := json.Marshal(signedMsg)
		log.Infow("Send", "name", name, "cid", cid.String(), "msg", string(msg))

		return cid, nil
	}

	return cid.Cid{}, errors.New("")
}

func (b *Boat) Status(c *gin.Context) {
	b.lk.Lock()
	defer b.lk.Unlock()

	i := 0
	for _, api := range b.apis {
		_, err := api.Api.ChainHead(context.Background())
		if err != nil {
			log.Warnw("ChainHead", "err", err)
			continue
		}
		i++
	}

	if i == 0 {
		ReturnError(c, NewError(500, "no nodes available"))
		return
	}

	ReturnOk(c, client.Response{
		Code:    200,
		Message: "Good",
	})
}
