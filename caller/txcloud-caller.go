package txycaller

import (
	"github.com/rosbit/txcloud-sign"
	"github.com/rosbit/gnet"
	"txcloud-api-caller/conf"
	"net/http"
	"fmt"
	"os"
)

type CallerParams struct {
	Method  string
	Service string
	Version string
	Endpoint string
}

func CallTxCloud(action string, region string, callerParams *CallerParams, signedHeaders map[string]string, body []byte, res interface{}, debug bool) error {
	ak := &conf.ServiceConf.AccessKey
	headers := txcsign.MakeTxCloudSignV30Headers(ak.SecretId, ak.SecretKey, callerParams.Method, callerParams.Service, action, region, callerParams.Version, "/", signedHeaders, body)

	o := func(*gnet.Options){}
	if debug {
		o = gnet.BodyLogger(os.Stderr)
		fmt.Fprintf(os.Stderr, "headers: %#v\n", headers)
	}
	status, err := gnet.HttpCallJ(callerParams.Endpoint, res, gnet.M(callerParams.Method), gnet.Params(body), gnet.Headers(headers), o)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("status: %d\n", status)
	}
	return nil
}

