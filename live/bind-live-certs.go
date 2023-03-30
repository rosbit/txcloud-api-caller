package live

import (
	txycaller "txcloud-api-caller/caller"
	"strings"
	"fmt"
	"os"
	"flag"
	"encoding/json"
)

var callerParams = txycaller.CallerParams {
	Method: "POST",
	Service: "live",
	Version: "2018-08-01",
	Endpoint: "https://live.tencentcloudapi.com/",
}

// ModifyLiveDomainCertBindings: 批量绑定证书对应的播放域名
// 参考文档：https://cloud.tencent.com/document/api/267/78655
func BindSSLCerts() {
	f := flag.NewFlagSet("bind-live-ssl-certs", flag.ExitOnError)
	certId := f.String("cert-id", "", "specify cert id")
	region := f.String("region", "live.tencentcloudapi.com", "speicfy region name")
	domains := f.String("domains", "", "specify domains")
	debug := f.Bool("debug", false, "debug mode or not")
	if err := f.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}
	if len(*certId) == 0 ||  len(*region) == 0 ||  len(*domains) == 0 {
		fmt.Fprintf(os.Stderr, "usage: %s bind-live-ssl-certs -cert-id=xxx -region=xxx -domains=xxx -debug=true|false\n", os.Args[0])
		os.Exit(4)
	}

	ds := strings.FieldsFunc(*domains, func(c rune)bool{
		return c == ' ' || c == '\t' || c == ',' || c == ';'
	})

	action := "ModifyLiveDomainCertBindings"
	body, _ := json.Marshal(map[string]interface{}{
		"DomainInfos": func() []interface{}{
			res := make([]interface{}, len(ds))
			for i, d := range ds {
				res[i] = map[string]interface{}{
					"DomainName": d,
					"Status": 1,
				}
			}
			return res
		}(),
		"CloudCertId": *certId,
	})
	signedHeaders := map[string]string{
		"Host": *region,
		"Content-Type": "application/json; charset=utf-8",
	}
	var res struct {
		Response struct {
			MismatchedDomainNames []string
			Error *struct {
				Code string
				Message string
			} `json:"Error,omitempty"`
			RequestId string
		}
	}
	if err := txycaller.CallTxCloud(action, "", &callerParams, signedHeaders, body, &res, *debug); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(5)
	}
	md := res.Response.MismatchedDomainNames
	if len(md) > 0 {
		fmt.Fprintf(os.Stderr, "MismatchedDomainNames: %v\n", md)
		os.Exit(6)
	}
}
