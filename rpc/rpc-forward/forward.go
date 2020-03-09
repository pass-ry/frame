package forward

import (
	"encoding/json"

	"gitlab.ifchange.com/data/cordwood/log"
	curl "gitlab.ifchange.com/data/cordwood/rpc/rpc-curl"
	handler "gitlab.ifchange.com/data/cordwood/rpc/rpc-handler"
)

func Request(serverRequest *handler.Request) (clientRequest *curl.Request) {
	clientRequest = curl.NewRequest()

	clientRequest.SetLogID(serverRequest.GetLogID())
	clientRequest.SetC(serverRequest.GetC())
	clientRequest.SetM(serverRequest.GetM())

	paramsJson := serverRequest.GetP()
	var params interface{}

	err := json.Unmarshal(paramsJson, &params)
	if err != nil {
		log.Errorf("Cordwood/rpc/rpc-forward-request %v %v",
			err, serverRequest)
		return
	}
	clientRequest.SetP(params)
	return
}

func Response(clientResponse *curl.Response) (serverResponse *handler.Response) {
	serverResponse = handler.NewResponse()

	serverResponse.SetLogID(clientResponse.GetLogID())
	serverResponse.SetErrNo(clientResponse.GetErrNo())
	serverResponse.SetErrMsg(clientResponse.GetErrMsg())

	resultsJson := clientResponse.GetResults()
	var results interface{}

	err := json.Unmarshal(resultsJson, &results)
	if err != nil {
		log.Errorf("Cordwood/rpc/rpc-forward-response %v %v",
			err, clientResponse)
		return
	}
	serverResponse.SetResults(results)
	return
}
