package handlers

import (

	"log"
	"net/http"


	"github.com/openfaas/faas/gateway/metrics"
	"github.com/openfaas/faas/gateway/types"
	"io/ioutil"
	"encoding/json"
	"github.com/openfaas/faas/gateway/requests"
	"bytes"
)


// MakeDeleteFunctionProxyHandler create a handler which delete function
func MakeDeleteFunctionProxyHandler(proxy *types.HTTPClientReverseProxy, metrics metrics.MetricOptions, baseURLResolver BaseURLResolver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := baseURLResolver.Resolve(r)
		body, _ := ioutil.ReadAll(r.Body)
		req:=requests.DeleteFunctionRequest{}
		json.Unmarshal(body, &req)
		delete(proxy.Client,baseURL,req.FunctionName,metrics)
	}
}

func delete(proxyClient *http.Client,baseURL string,functionName string,metrics metrics.MetricOptions){
	str :=make(map[string]string)
	str["functionName"]=functionName
	reqBytes, _ := json.Marshal(&str)
	reader := bytes.NewReader(reqBytes)
	request, _ := http.NewRequest(http.MethodDelete, baseURL+"/system/functions",reader)
	request.Header.Set("Content-Type", "application/json")
	_,err:=proxyClient.Do(request)
	if err==nil{
		log.Printf("delete functon [%s] is ok",functionName)
		trackFunctionStop(metrics,functionName)
	}
}


