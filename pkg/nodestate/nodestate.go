package nodestate

import (
	"database/sql"
	"fmt"
	env "gclustercheck/pkg/env"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Status stores sql member status
type Status struct {
	Synced   bool
	ReadOnly bool
	Offline  bool
	Error    string
}

func sqlClient(env *env.MysqlCfg) (db *sql.DB, err error) {
	sqlURI := fmt.Sprintf("%s:%s@tcp(%s:%v)/", env.MysqlUsername, env.MysqlPassword, env.MysqlHostname, env.MysqlPort)
	db, err = sql.Open("mysql", sqlURI)
	if err != nil {
		log.Println(err.Error())
	}
	return db, err
}

func sqlPing(sqlX *sql.DB) (bool, error) {

	nodeOffline := true
	err := sqlX.Ping()
	if err != nil {
		err = fmt.Errorf("Failed to ping database %s", err)
		return nodeOffline, err
	}
	nodeOffline = false
	return nodeOffline, err
}

func sqlShowStatus(sqlX *sql.DB, sqlStmt string) (string, error) {
	selectOut, err := sqlX.Query(sqlStmt)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	var result string
	var variableName string
	for selectOut.Next() {
		err = selectOut.Scan(&variableName, &result)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return result, err
}

func getWsrepLocalState(sqlX *sql.DB) (bool, error) {
	var err error
	sqlStmt := "SHOW STATUS LIKE '%wsrep_local_state';"
	result, err := sqlShowStatus(sqlX, sqlStmt)
	if err != nil {
		log.Println(err)
	}

	//fmt.Println("result = ", result)
	var synced bool = false
	if result == "" {
		err := fmt.Errorf("No result returned for: %s", sqlStmt)
		return synced, err
	}

	if result == "4" {
		synced = true
	}
	return synced, err
}

// Check returns struct of current states of sql
func Check(env *env.MysqlCfg) (*Status, error) {
	sqlX, err := sqlClient(env)
	defer sqlX.Close()
	// Set default status states
	status := &Status{
		Synced:   false,
		ReadOnly: false,
		Offline:  true,
		Error:    "None",
	}

	if err != nil {
		log.Println(err.Error())
		status.Error = err.Error()
		status.Offline = true
		return status, err
	}

	sqlPingSuccess, err := sqlPing(sqlX)
	if err != nil {
		status.Offline = true
		status.Error = err.Error()
		return status, err
	}

	syncedLocalState, err := getWsrepLocalState(sqlX)
	if err != nil {
		status.Error = err.Error()
		return status, err
	}
	status.Synced = syncedLocalState
	status.Offline = sqlPingSuccess
	status.Error = "None"
	return status, err
}
