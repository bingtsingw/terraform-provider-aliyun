package aliyun

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/fc-go-sdk"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ApiVersion string

const (
	ApiVersion20160815 = ApiVersion("2016-08-15")
)

const DefaultClientRetryCountSmall = 5

type Config struct {
	AccessKey  string
	SecretKey  string
	Region     Region
	RegionId   string
}

func (c *Config) Client() Client {
	fcconn, _ := c.newFcClient()

	client := Client{
		fcconn: fcconn,
	}

	return client
}

func (c *Config) getSdkConfig() *sdk.Config {
	return sdk.NewConfig().
		WithMaxRetryTime(DefaultClientRetryCountSmall).
		WithTimeout(time.Duration(30) * time.Second).
		WithEnableAsync(true).
		WithGoRoutinePoolSize(100).
		WithMaxTaskQueueSize(10000).
		WithDebug(false).
		WithHttpTransport(c.getTransport()).
		WithScheme("HTTPS")
}

func (c *Config) getTransport() *http.Transport {
	handshakeTimeout, err := strconv.Atoi(os.Getenv("TLSHandshakeTimeout"))
	if err != nil {
		handshakeTimeout = 120
	}
	transport := &http.Transport{}
	transport.TLSHandshakeTimeout = time.Duration(handshakeTimeout) * time.Second

	return transport
}

func (c *Config) newFcClient() (*fc.Client, error) {
	endpoint := fmt.Sprintf("https://%s.fc.aliyuncs.com", c.RegionId)

	fcconn, err := fc.NewClient(endpoint, string(ApiVersion20160815), c.AccessKey, c.SecretKey)

	return fcconn, err
}
