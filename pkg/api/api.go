package api

import (
	"fmt"
	"gclustercheck/pkg/env"
	"gclustercheck/pkg/nodestate"
	output "gclustercheck/pkg/output"
	"log"
	"time"

	air "github.com/aofei/air"
)

// XtraDBClient embed nodestate.Client
type XtraDBClient struct {
	*nodestate.Cfg
}

// Init api for requests
func Init() {
	mysqlenv, apienv := env.Get()
	xtradb := XtraDBClient{
		&nodestate.Cfg{
			MysqlHostname:       mysqlenv.MysqlHostname,
			MysqlUsername:       mysqlenv.MysqlUsername,
			MysqlPassword:       mysqlenv.MysqlPassword,
			MysqlPort:           mysqlenv.MysqlPort,
			ReadOnlyIsAvailable: apienv.AvailableWhenReadOnly,
			DonorIsAvailable:    apienv.AvailableWhenDonor,
		},
	}
	air.Default.Address = fmt.Sprintf("0.0.0.0:%v", apienv.APIPort)
	air.Default.AppName = "gclustercheck"
	air.Default.GET("/", xtradb.clustercheck)
}

func executionTime(start time.Time, functionName string) {
	elapsedTime := time.Since(start)
	log.Printf("%s took %s", functionName, elapsedTime)
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

func (xtradb *XtraDBClient) clustercheck(req *air.Request, res *air.Response) error {
	defer executionTime(time.Now(), "clustercheck")
	responses := output.Messages()
	sqlStates, err := xtradb.Check()
	fullStatusMsg := fmt.Sprintf(responses.FullStatus, sqlStates.Offline, sqlStates.Synced, sqlStates.ReadOnly, sqlStates.Donor, sqlStates.Error)
	if err != nil {
		log.Printf("err = %s", err)
		offline(res, responses.Error, fullStatusMsg)
	}
	if sqlStates.Synced && !sqlStates.Offline {
		if sqlStates.ReadOnly {
			online(res, responses.ReadOnly, fullStatusMsg)
		} else {
			online(res, responses.Synced, fullStatusMsg)
		}
	} else {
		if sqlStates.Donor {
			offline(res, responses.Donor, fullStatusMsg)
		} else if sqlStates.ReadOnly {
			offline(res, responses.ReadOnly, fullStatusMsg)
		} else {
			offline(res, responses.Error, fullStatusMsg)
		}
	}
	return nil
}
