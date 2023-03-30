package sslcert

import (
	txycaller "txcloud-api-caller/caller"
	"fmt"
	"os"
	"flag"
	"encoding/json"
)

var callerParmas = txycaller.CallerParams {
	Method: "POST",
	Service: "ssl",
	Version: "2019-12-05",
	Endpoint: "https://ssl.tencentcloudapi.com/",
}

// UploadCertificate: 上传证书
// 参考文档: https://cloud.tencent.com/document/api/400/41665
func UploadSslCerts() {
	f := flag.NewFlagSet("upload-ssl-certs", flag.ExitOnError)
	certName := f.String("cert-name", "", "specify cert name")
	region := f.String("region", "ssl.tencentcloudapi.com", "speicfy region name")
	debug := f.Bool("debug", false, "debug mode or not")
	certFileName := f.String("cert-file-name", "", "specify cert file name")
	keyFileName := f.String("key-file-name", "", "specify key file name")
	if err := f.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}
	if len(*certName) == 0 ||  len(*region) == 0 ||  len(*certFileName) == 0 || len(*keyFileName) == 0 {
		fmt.Fprintf(os.Stderr, "usage: %s upload-ssl-certs -certName=xxx -region=xxx -cert-file-name=xxx -key-file-name=xxx -debug=true|false\n", os.Args[0])
		os.Exit(4)
	}
	cert, err := os.ReadFile(*certFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(5)
	}
	key, err := os.ReadFile(*keyFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(6)
	}

	action := "UploadCertificate"
	body, _ := json.Marshal(map[string]interface{}{
		"CertificatePublicKey": string(cert),
		"CertificatePrivateKey": string(key),
		"CertificateType": "SVR",
		"CertificateUse": "LIVE",
		"Alias": *certName,
		"Repeatable": false,
	})
	signedHeaders := map[string]string{
		"Host": *region,
		"Content-Type": "application/json; charset=utf-8",
	}
	var res struct {
		Response struct {
			CertificateId string
			RepeatCertId string
			Error *struct {
				Code string
				Message string
			} `json:"Error,omitempty"`
			RequestId string
		}
	}
	if err := txycaller.CallTxCloud(action, "", &callerParmas, signedHeaders, body, &res, *debug); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(5)
	}
	fmt.Printf("%s\n", res.Response.CertificateId)
}
