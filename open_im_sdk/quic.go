package open_im_sdk

import (
	"context"
	"fmt"
	"github.com/quic-go/quic-go/http3"
	"io"
	"net/http"
	"time"
)

func init() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		cli := http.Client{
			Transport: &http3.Transport{},
		}
		const http3url = "https://dns.alidns.com/resolve?name=www.taobao.com.&type=1"
		fmt.Println("http3 test:", http3url)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, http3url, nil)
		if err != nil {
			fmt.Println("http.NewRequestWithContext error", err)
			return
		}
		resp, err := cli.Do(req)
		if err != nil {
			fmt.Println("cli.Do error", err)
			return
		}
		defer resp.Body.Close()
		fmt.Println("response status code:", resp.Status)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("io.ReadAll error", err)
			return
		}
		fmt.Println("http version:", resp.Proto)
		fmt.Println("response body:", string(body))
	}()
}
