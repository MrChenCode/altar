package test

import (
	"altar/application/config"
	"altar/application/context"
	"altar/application/logger"
	"altar/application/router"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type Method string

const (
	Get     Method = "GET"
	Post    Method = "POST"
	Delete  Method = "DELETE"
	Head    Method = "HEAD"
	Put     Method = "PUT"
	Options Method = "OPTIONS"
	Trace   Method = "TRACE"
	Patch   Method = "PATCH"
	Move    Method = "MOVE"
	Copy    Method = "COPY"
	Link    Method = "LINK"
	Unlink  Method = "UNLINK"
	Wrapped Method = "WRAPPED"
	Connect Method = "CONNECT"

	ResultCustom string = "custom"
)

var (
	server *httptest.Server
	client *http.Client
)

type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

func init() {
	var err error
	server, err = NewHttpTest()
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ResponseHeaderTimeout: 10 * time.Second,
		DisableKeepAlives:     false,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   http.DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		TLSNextProto:          make(map[string]func(string, *tls.Conn) http.RoundTripper),
	}

	client = &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}

type Request struct {
	//请求的地址，最前面带/, 比如"/bookinfo"
	Path string
	//请求的方式Method
	Method Method
	//请求的参数数据, "a=b&c=d"
	Body []byte
	//发送的数据，urlencode/json 默认为urlencode
	RequestType string

	//期望返回的http code, 如果是非200的code，不再验证respCode和respResult
	//默认为200
	HttpCode int
	//期望返回值的code字段码，如果是非0，不再验证respResult的值
	//默认为0
	RespCode int

	//期望返回的result字段的值
	//验证具体返回值，需要多种情况
	//
	//第一： 不验证result返回数量，只验证code，此处r.RespResult为nil，则不判断返回的result
	//
	//第二： 目前提供的验证返回值的方式，不能满足更复杂的验证，则需要提供返回值调用方自行验证
	//		此处r.RespResult为ResultCustom，此函数返回值会返回result
	//
	//第三： 只验证返回值的字段是否齐全，不验证值内容，这时候r.RespResult是一个string的slice
	//		此时必须必须要求返回值的result字段是一个map[string]interface{}
	//
	//第四:  需要验证字段和具体的返回值是否相等, 这时候r.RespResult是一个map[string]interface{}
	//		也必须要求返回值的result字段是一个map[string]interface{}
	RespResult interface{}
}

func (r *Request) init() {
	if r.Method == "" {
		r.Method = Get
	}
	if r.HttpCode == 0 {
		r.HttpCode = 200
	}
}

func RunTestHttpCode(t *testing.T, path string, method Method, body string, httpcode int) {
	RunTestApi(t, &Request{
		Path:     path,
		Method:   method,
		Body:     []byte(body),
		HttpCode: httpcode,
	})
}

func RunTestResponseCode(t *testing.T, path string, method Method, body string, code int) {
	RunTestApi(t, &Request{
		Path:     path,
		Method:   method,
		Body:     []byte(body),
		HttpCode: 200,
		RespCode: code,
	})
}

//测试接口返回数据
func RunTestApi(t *testing.T, r *Request) interface{} {
	r.init()

	res, err := http.NewRequest(string(r.Method), server.URL+"/"+r.Path, bytes.NewBuffer(r.Body))
	require.Nil(t, err)
	res.Header = make(http.Header)
	switch r.RequestType {
	case "", "urlencode":
		res.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=utf-8")
	case "json":
		res.Header.Add("Content-type", "application/json;charset=utf-8")
	default:
		require.Fail(t, "invalid request type")
	}

	resp, err := client.Do(res)
	require.Nil(t, err)
	defer func() {
		require.Nil(t, resp.Body.Close())
	}()

	require.Equal(t, r.HttpCode, resp.StatusCode)

	//如果是验证非200的错误，到此处就结束了
	if r.HttpCode != 200 {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)

	//解析body
	response := &Response{}
	require.Nil(t, json.Unmarshal(body, response))

	//验证期望的response的返回code码，如果是非0的，则不再验证response的result
	require.Equal(t, r.RespCode, response.Code)
	if r.RespCode != 0 {
		return nil
	}

	//验证具体返回值，需要多种情况
	//
	//第一： 不验证result返回数量，只验证code，此处r.RespResult为nil，则不判断返回的result
	//
	//第二： 目前提供的验证返回值的方式，不能满足更复杂的验证，则需要提供返回值调用方自行验证
	//		此处r.RespResult为ResultCustom，此函数返回值会返回result
	//
	//第三： 只验证返回值的字段是否齐全，不验证值内容，这时候r.RespResult是一个string的slice
	//		此时必须必须要求返回值的result字段是一个map[string]interface{}
	//
	//第四:  需要验证字段和具体的返回值是否相等, 这时候r.RespResult是一个map[string]interface{}
	//		也必须要求返回值的result字段是一个map[string]interface{}
	switch v := r.RespResult.(type) {
	case nil:
		return nil
	case string:
		if ResultCustom == v {
			return response.Result
		}
		require.Fail(t, "invalid RespResult")
	case []string:
		result, ok := response.Result.(map[string]interface{})
		require.True(t, ok)

		//字段数量必须匹配期望的返回值数量
		require.Equal(t, len(result), len(v))
		for _, x := range v {
			_, ok := result[x]
			require.True(t, ok)
			delete(result, x)
		}
	case map[string]interface{}:
		result, ok := response.Result.(map[string]interface{})
		require.True(t, ok)

		//字段数量必须匹配期望的返回值数量
		require.Equal(t, len(result), len(v))
		for kk, vv := range v {
			rvv, ok := result[kk]
			require.True(t, ok)
			require.Equal(t, vv, rvv)
			delete(result, kk)
		}
	default:
		require.Fail(t, "invalid RespResult")
	}
	return nil
}

func NewHttpTest() (*httptest.Server, error) {
	//初始化配置文件
	var configFile string
	for _, v := range []string{"../altar.ini", "./altar_test.ini", "../altar_default.ini", "/etc/altar.ini"} {
		_, err := os.Stat(v)
		if err == nil {
			configFile = v
		}
	}
	if configFile == "" {
		return nil, errors.New("invalid configuration file")
	}
	config.RootPath = "../"

	c, err := config.NewConfig(configFile)
	if err != nil {
		return nil, err
	}

	//创建logs
	loginfo := logger.NewTestWriteNull()
	logwf := logger.NewTestWriteNull()

	ctx, err := context.NewController(c, loginfo, logwf)
	if err != nil {
		return nil, err
	}
	r := router.NewRouter(ctx)

	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(router.Recovery(logwf))
	r.Router(e)
	e.Routes()

	return httptest.NewServer(e), nil
}
