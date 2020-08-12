package nodestate

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Status stores sql member status
type Status struct {
	Synced   bool
	ReadOnly bool
	Offline  bool
	Donor    bool
	Error    string
}

// Cfg info to create Mysql client
type Cfg struct {
	MysqlUsername       string
	MysqlPassword       string
	MysqlHostname       string
	MysqlPort           int
	ReadOnlyIsAvailable bool
	DonorIsAvailable    bool
}

// DBClient embeds sql.DB
type DBClient struct {
	*sql.DB
}

func (c *Cfg) sqlClient() (db *sql.DB, err error) {
	sqlURI := fmt.Sprintf("%s:%s@tcp(%s:%v)/", c.MysqlUsername, c.MysqlPassword, c.MysqlHostname, c.MysqlPort)
	db, err = sql.Open("mysql", sqlURI)
	if err != nil {
		log.Println(err.Error())
	}
	return db, err
}

func (db *DBClient) ping() (bool, error) {
	nodeOffline := true
	err := db.Ping()
	if err != nil {
		err = fmt.Errorf("Failed to ping database %s", err)
		return nodeOffline, err
	}
	nodeOffline = false
	return nodeOffline, err
}

func (db *DBClient) getVariables(sqlStmt string) (string, error) {
	selectOut, err := db.Query(sqlStmt)
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

func (db *DBClient) getState(sqlStmt string, matchValue string) (bool, error) {
	var err error
	result, err := db.getVariables(sqlStmt)
	if err != nil {
		log.Println(err)
	}

	var currentValueMatches bool = false
	if result == "" {
		err := fmt.Errorf("No result returned for: %s", sqlStmt)
		return currentValueMatches, err
	}

	if result == matchValue {
		currentValueMatches = true
	}
	return currentValueMatches, err
}

func (db *DBClient) isSynced() (bool, error) {
	sqlStmt := "SHOW STATUS LIKE '%wsrep_local_state';"
	synced, err := db.getState(sqlStmt, "4")
	return synced, err
}

func (db *DBClient) isDonor() (bool, error) {
	sqlStmt := "SHOW STATUS LIKE '%wsrep_local_state';"
	donor, err := db.getState(sqlStmt, "2")
	return donor, err
}

func (db *DBClient) isReadOnly() (bool, error) {
	sqlStmt := "SHOW GLOBAL VARIABLES LIKE 'read_only';"
	readonly, err := db.getState(sqlStmt, "ON")
	return readonly, err
}

// Check returns struct of current states of sql
func (c *Cfg) Check() (*Status, error) {
	sqlX, err := c.sqlClient()
	db := DBClient{
		sqlX,
	}
	defer db.Close()
	// Set default status states
	status := &Status{
		Synced:   false,
		ReadOnly: true,
		Offline:  true,
		Donor:    true,
		Error:    "None",
	}

	if err != nil {
		log.Println(err.Error())
		status.Error = err.Error()
		status.Offline = true
		return status, err
	}

	offline, err := db.ping()
	if err != nil {
		status.Offline = true
		status.Error = err.Error()
		return status, err
	}

	synced, err := db.isSynced()
	if err != nil {
		status.Error = err.Error()
		return status, err
	}

	readonly, err := db.isReadOnly()
	if err != nil {
		status.Error = err.Error()
		return status, err
	}

	donor, err := db.isDonor()
	if err != nil {
		status.Error = err.Error()
		return status, err
	}
	// If either overrides for readonly or donor are set, take the node offline
	if !c.ReadOnlyIsAvailable && readonly {
		offline = true
	} else if !c.DonorIsAvailable && donor {
		offline = true
	}

	status.Synced = synced
	status.Offline = offline
	status.ReadOnly = readonly
	status.Donor = donor
	status.Error = "None"
	return status, err
}
