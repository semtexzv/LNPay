package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lightningnetwork/lnd/lnrpc"
	"lnpay/kube"
	"lnpay/lnet"
	"net/http"
	"os"
)

func GetMetadata(c *gin.Context) {
	info, err := lnet.LnClient.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Could not connect to lightning node"})
		return
	}

	c.JSON(http.StatusOK, &struct {
		PubKey  string `json:"pub_key"`
		IP      string `json:"ip"`
		Network string `json:"network"`
	}{
		PubKey:  info.IdentityPubkey,
		IP:      kube.GetSvcEndpoint("lnd-external"),
		Network: os.Getenv("NETWORK"),
	})

}
