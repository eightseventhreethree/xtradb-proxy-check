package env

import (
	viper "github.com/spf13/viper"
)

// MysqlCfg is struct to store Env for MySQL
type MysqlCfg struct {
	MysqlUsername string
	MysqlPassword string
	MysqlHostname string
	MysqlPort     int
	Verbose       bool
}

// APICfg is struct to store Env for API
type APICfg struct {
	APIPort               int
	AvailableWhenDonor    bool
	AvailableWhenReadOnly bool
	Verbose               bool
}

// Get returns EnvCfg with environment config
func Get() (*MysqlCfg, *APICfg) {
	vcfg := viper.New()
	vcfg.SetEnvPrefix("CLUSTERCHK")
	vcfg.AutomaticEnv()
	vcfg.SetDefault("mysql_username", "clustercheckuser")
	vcfg.SetDefault("mysql_password", "clustercheckpassword123")
	vcfg.SetDefault("mysql_hostname", "127.0.0.1")
	vcfg.SetDefault("mysql_port", 3306)
	vcfg.SetDefault("api_port", 9200)
	vcfg.SetDefault("available_when_donor", true)
	vcfg.SetDefault("available_when_read_only", false)
	vcfg.SetDefault("verbose", false)

	return &MysqlCfg{
			MysqlUsername: vcfg.GetString("mysql_username"),
			MysqlPassword: vcfg.GetString("mysql_password"),
			MysqlHostname: vcfg.GetString("mysql_hostname"),
			MysqlPort:     vcfg.GetInt("mysql_port"),
			Verbose:       vcfg.GetBool("verbose"),
		}, &APICfg{
			APIPort:               vcfg.GetInt("api_port"),
			AvailableWhenDonor:    vcfg.GetBool("available_when_donor"),
			AvailableWhenReadOnly: vcfg.GetBool("available_when_read_only"),
			Verbose:               vcfg.GetBool("verbose"),
		}
}
