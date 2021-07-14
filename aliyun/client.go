package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/aliyun/fc-go-sdk"
)

type Client struct {
	fcconn   *fc.Client
	crconn   *cr.Client
	dcdnconn *dcdn.Client
}
