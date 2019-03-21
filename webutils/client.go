package webutils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	lbc = make(map[string]*fasthttp.LBClient)
)

func AddLBClient(name string, servers []string) {
	cli, ok := lbc[name]
	if !ok {
		cli = &fasthttp.LBClient{}
		lbc[name] = cli
	}
	for _, addr := range servers {
		c := &fasthttp.HostClient{
			Addr: addr,
		}
		cli.Clients = append(cli.Clients, c)
	}
}

func NewRequest(host, method, requestURI string, args map[string]string, body []byte, cookies map[string]interface{}) *fasthttp.Request {
	req := &fasthttp.Request{}
	req.Header.SetHost(host)
	req.Header.SetMethod(method)
	if args != nil {
		req.SetRequestURI(requestURI + "?" + GenRequestArgs(args))
	} else {
		req.SetRequestURI(requestURI)
	}
	if len(body) > 0 {
		io.Copy(req.BodyWriter(), bytes.NewReader(body))
	}
	if cookies != nil {
		for key, v := range cookies {
			req.Header.SetCookie(key, v.(string))
		}
	}
	return req
}

func NewRequestWithSign(host, method, requestURI, secret string, args map[string]string, body []byte, cookies map[string]interface{}) *fasthttp.Request {
	req := &fasthttp.Request{}
	req.Header.SetHost(host)
	req.Header.SetMethod(method)
	if args != nil {
		// 如果没有签名参数则添加
		_, ok := args[PARAM_SIGN]
		if !ok {
			args[PARAM_SIGN] = GenSignature(method, requestURI, args, secret)
		}
		req.SetRequestURI(requestURI + "?" + GenSortRequestArgs(args))
	} else {
		req.SetRequestURI(requestURI)
	}
	if len(body) > 0 {
		io.Copy(req.BodyWriter(), bytes.NewReader(body))
	}
	if cookies != nil {
		for key, v := range cookies {
			req.Header.SetCookie(key, v.(string))
		}
	}
	return req
}

func NewResponse() *fasthttp.Response {
	return &fasthttp.Response{}
}

func LBCDoTimeout(cliName string, req *fasthttp.Request, resp *fasthttp.Response, timeout int) ([]byte, int, error) {
	cli, ok := lbc[cliName]
	if !ok {
		return []byte{}, http.StatusNotFound, fmt.Errorf("not found client by:%s", cliName)
	}
	err := cli.DoTimeout(req, resp, time.Duration(timeout)*time.Second)
	if err != nil {
		return []byte{}, http.StatusInternalServerError, err
	}
	return resp.Body(), resp.StatusCode(), err
}

func HttpForData(host, method, url string, args map[string]string, body []byte, cookies map[string]interface{}, timeout int) ([]byte, error) {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}
	req.Header.SetHost(host)
	req.Header.SetMethod(method)

	if args != nil {
		var params []string
		for k, v := range args {
			p := fmt.Sprintf("%s=%s", k, v)
			params = append(params, p)
		}
		req.SetRequestURI(url + "?" + strings.Join(params, "&"))
	} else {
		req.SetRequestURI(url)
	}
	if len(body) > 0 {
		io.Copy(req.BodyWriter(), bytes.NewReader(body))
	}
	if cookies != nil {
		for key, v := range cookies {
			req.Header.SetCookie(key, v.(string))
		}
	}
	err := fasthttp.DoTimeout(req, resp, time.Duration(timeout)*time.Second)
	if err != nil {
		return []byte{}, err
	}
	return resp.Body(), nil
}

func HttpForResponse(host, method, url string, args map[string]string, body []byte, cookies map[string]interface{}, timeout int) (*fasthttp.Response, error) {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}
	req.Header.SetHost(host)
	req.Header.SetMethod(method)

	if args != nil {
		var params []string
		for k, v := range args {
			p := fmt.Sprintf("%s=%s", k, v)
			params = append(params, p)
		}
		req.SetRequestURI(url + "?" + strings.Join(params, "&"))
	} else {
		req.SetRequestURI(url)
	}
	if len(body) > 0 {
		io.Copy(req.BodyWriter(), bytes.NewReader(body))
	}
	if cookies != nil {
		for key, v := range cookies {
			req.Header.SetCookie(key, v.(string))
		}
	}
	err := fasthttp.DoTimeout(req, resp, time.Duration(timeout)*time.Second)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func CheckSignByRequest(secret string, req *http.Request) bool {
	_sign := ""
	sign := ""

	argsParam := make(map[string]string)
	for k, vs := range req.URL.Query() {
		argsParam[k] = vs[0]
	}

	_sign = argsParam[PARAM_SIGN]
	if len(argsParam) > 0 {
		delete(argsParam, PARAM_SIGN)
		sign = GenSignature(req.Method, req.URL.Path, argsParam, secret)
	}
	if _sign == "" || sign == "" {
		return false
	}
	return _sign == sign
}

func CheckSignByFasthttp(secret string, ctx *fasthttp.RequestCtx) bool {
	_sign := ""
	sign := ""

	args := ctx.Request.URI().QueryArgs()
	if args != nil {
		argsParam := make(map[string]string)
		args.VisitAll(func(key, value []byte) {
			argsParam[string(key)] = string(value)
		})
		_sign = argsParam[PARAM_SIGN]
		delete(argsParam, PARAM_SIGN)

		sign = GenSignature(string(ctx.Method()), string(ctx.URI().Path()), argsParam, secret)
	}
	if _sign == "" || sign == "" {
		return false
	}
	return _sign == sign
}

// return: isHttps/host/uri/error
// eg: https://www.baidu.com/api/v1/name,return: true,www.baidu.com,/api/v1/name,nil
func ParseUrl(url string) (bool, string, string) {
	var beginPos, endPos int
	var host, uri string
	isHttps := strings.HasPrefix(url, "https://")
	if isHttps {
		beginPos += 8
	}
	if strings.HasPrefix(url, "http://") {
		beginPos += 7
	}
	endPos = strings.Index(url[beginPos:], "/")
	if endPos > 0 {
		host = url[beginPos : beginPos+endPos]
		uriPos := strings.Index(url[beginPos+endPos:], "?")
		if uriPos > 0 {
			uri = url[beginPos+endPos : beginPos+endPos+uriPos]
		} else {
			uri = url[beginPos+endPos:]
		}
	} else {
		host = url[beginPos:]
		uri = "/"
	}
	return isHttps, host, uri
}

func NewMutilpartBuf(fieldname, filename string) (*bytes.Buffer, string, error) {
	if filename == "" {
		return nil, "", fmt.Errorf("filename is empty")
	}
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, "", err
	}
	if !fi.Mode().IsRegular() {
		return nil, "", fmt.Errorf(filename + " is not a regular file")
	}
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	defer w.Close()

	pkg, err := w.CreateFormFile(fieldname, filepath.Base(filename))
	if err != nil {
		return nil, "", err
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	_, err = io.Copy(pkg, f)
	if err != nil {
		return nil, "", err
	}
	return buf, w.FormDataContentType(), nil
}
