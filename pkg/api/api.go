package api

import (
	"fmt"
	"gclustercheck/pkg/env"
	"gclustercheck/pkg/nodestate"
	output "gclustercheck/pkg/output"
	"log"

	air "github.com/aofei/air"
)

// Init api for requests
func Init() {
	_, apienv := env.Get()
	air.Default.Address = fmt.Sprintf("0.0.0.0:%v", apienv.APIPort)
	air.Default.AppName = "gclustercheck"
	air.Default.GET("/", clustercheckHandler)
}

func online(res *air.Response, response string, fullstatus string) error {
	res.Status = 200
	fullResp := response + fullstatus
	log.Println("Online: ", fullResp)
	return res.WriteString(fullResp)
}

func offline(res *air.Response, response string, fullstatus string) error {
	res.Status = 503
	fullResp := response + fullstatus
	log.Println("Offline: ", fullResp)
	return res.WriteString(fullResp)
}

func clustercheckHandler(req *air.Request, res *air.Response) error {
	responses := output.Messages()
	mysqlenv, _ := env.Get()
	sqlStates, err := nodestate.Check(mysqlenv)
	fullStatusMsg := fmt.Sprintf(responses.FullStatus, sqlStates.Offline, sqlStates.Synced, sqlStates.ReadOnly, sqlStates.Error)
	if err != nil {
		offline(res, responses.Error, fullStatusMsg)
	}
	if sqlStates.Synced && !sqlStates.Offline {
		online(res, responses.Synced, fullStatusMsg)
	} else {
		offline(res, responses.Unsynced, fullStatusMsg)
	}
	return nil
}
