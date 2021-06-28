package fbmsg

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/olvrng/rbot/be/pkg/l"
	"github.com/olvrng/rbot/be/pkg/xerrors"
)

var ll = l.New()
var ls = ll.Sugar()

const EndpointURL = "https://graph.facebook.com/v2.6/me/messages"

type Config struct {
	VerifyToken     string
	PageAccessToken string
}

func (c *Config) Verify() error {
	if c.VerifyToken == "" {
		return xerrors.Errorf(xerrors.Internal, nil, "no verify token")
	}
	if c.PageAccessToken == "" {
		return xerrors.Errorf(xerrors.Internal, nil, "no page access token")
	}
	return nil
}

type Client struct {
	Config
	HTTP *resty.Client
}

func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Verify(); err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	restyClient := resty.NewWithClient(httpClient)
	c := &Client{
		Config: cfg,
		HTTP:   restyClient,
	}
	return c, nil
}

func (c *Client) callAPI(ctx context.Context, body interface{}) (*resty.Response, error) {

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	r := c.HTTP.R().SetContext(ctx)
	r.SetHeader("Content-Type", "application/json")
	r.SetQueryParam("access_token", c.PageAccessToken)
	r.SetBody(data)
	resp, err := r.Post(EndpointURL)
	if err != nil {
		ll.Error("messenger: can not call api", l.String("url", EndpointURL), l.Error(err))
		return nil, err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		ll.Debug("messenger: call api successfully", l.String("url", EndpointURL))
		return resp, nil
	}

	ls.Errorf("messenger: call api error, code=%v request=%s response=%s", resp.StatusCode(), data, resp.String())
	err = xerrors.Errorf(xerrors.Internal, nil, "messenger: response %v", resp.StatusCode())
	return resp, err
}

func (c *Client) CallSendAPI(ctx context.Context, req *SendRequest) error {
	_, err := c.callAPI(ctx, req)
	return err
}
