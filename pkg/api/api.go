package api

import (
	"fmt"
	"gclustercheck/pkg/env"
	"gclustercheck/pkg/nodestate"
	output "gclustercheck/pkg/output"

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
	return res.WriteString(fullResp)
}

func offline(res *air.Response, response string, fullstatus string) error {
	res.Status = 503
	fullResp := response + fullstatus
	return res.WriteString(fullResp)
}

func clustercheckHandler(req *air.Request, res *air.Response) error {
	responses := output.Messages()
	mysqlenv, _ := env.Get()
	sqlStates := nodestate.Check(mysqlenv)
	fullStatusMsg := fmt.Sprintf(responses.FullStatus, sqlStates.Offline, sqlStates.Synced, sqlStates.ReadOnly)
	if sqlStates.Synced && !sqlStates.Offline {
		online(res, responses.Synced, fullStatusMsg)
	}
	offline(res, responses.Unsynced, fullStatusMsg)

	return nil
}
