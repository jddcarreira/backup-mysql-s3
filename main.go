package main

import (
	"fmt"
	"log"
	"os/exec"

	"gopkg.in/ini.v1"
)

type mysqlConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

type s3Config struct {
	keyID      string
	secreteKey string
	bucket     string
}

func getConfig(configPath string) (mysqlconf mysqlConfig, s3conf s3Config) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatal(4, "Fail to parse 'conf/app.ini': %v", err)
	}

	mysqlCfg := cfg.Section("mysql")
	mysqlconf = mysqlConfig{
		host:     mysqlCfg.Key("host").String(),
		port:     mysqlCfg.Key("port").String(),
		username: mysqlCfg.Key("username").String(),
		password: mysqlCfg.Key("password").String(),
		database: mysqlCfg.Key("database").String(),
	}

	s3Cfg := cfg.Section("s3")
	s3conf = s3Config{
		keyID:      s3Cfg.Key("access_key_id").String(),
		secreteKey: s3Cfg.Key("secret_access_key").String(),
		bucket:     s3Cfg.Key("bucket").String(),
	}

	return mysqlconf, s3conf
}

func backupMySQL(mysqlconf mysqlConfig) {
	backupCmd := "mysqldump -h " + mysqlconf.host + " -P " + mysqlconf.port + " -u " + mysqlconf.username + " -p" + mysqlconf.password + " " + mysqlconf.database
	exec.Command(backupCmd)
}

func main() {
	configPath := "/Users/joaocarreira/go/src/github.com/johnnybus/backup-mysq-s3/backup.conf"
	mysqlconf, s3conf := getConfig(configPath)
	fmt.Printf("%v - %v\n", mysqlconf, s3conf)
}
