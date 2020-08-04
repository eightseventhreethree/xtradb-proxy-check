package nodestate

import (
	"database/sql"
	"fmt"
	env "gclustercheck/pkg/env"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	*env.MysqlCfg
}

// Status stores sql member status
type Status struct {
	Synced   bool
	ReadOnly bool
	Offline  bool
}

func sqlClient(env *env.MysqlCfg) (db *sql.DB) {
	sqlURI := fmt.Sprintf("%s:%s@tcp(%s:%v)/", env.MysqlUsername, env.MysqlPassword, env.MysqlHostname, env.MysqlPort)
	fmt.Println("sqlURI", sqlURI)
	db, err := sql.Open("mysql", sqlURI)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func sqlPing(sqlX *sql.DB) bool {

	nodeOffline := true
	err := sqlX.Ping()
	if err != nil {
		err = fmt.Errorf("Failed to ping database %s", err)
		return nodeOffline

	}
	nodeOffline = false
	fmt.Println("nodeOnline = ", nodeOffline)
	//defer sqlX.Close()
	return nodeOffline
}

func sqlShowStatus(sqlX *sql.DB, sqlStmt string) string {
	selectOut, err := sqlX.Query(sqlStmt)
	if err != nil {
		panic(err.Error())
	}
	var result string
	var variableName string
	for selectOut.Next() {
		err = selectOut.Scan(&variableName, &result)
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println("Variable name = ", variableName)
	return result
}

func getWsrepLocalState(sqlX *sql.DB) bool {
	result := sqlShowStatus(sqlX, "SHOW STATUS LIKE '%wsrep_local_state';")

	fmt.Println("result = ", result)
	var synced bool = false
	err := fmt.Errorf("No error")
	if err != nil {
		log.Printf("getWsrepLocalState err %s\n", err)
	}
	if result == "4" {
		synced = true
	}
	return synced
}

// Check returns struct of current states of sql
func Check(env *env.MysqlCfg) *Status {
	sqlX := sqlClient(env)

	defer sqlX.Close()
	return &Status{
		Synced:   getWsrepLocalState(sqlX),
		ReadOnly: false,
		Offline:  sqlPing(sqlX),
	}
}
