package main

import (
	"fmt"
	"io"
	"log"
	"os"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
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

func genDumpFile(databaseName string) (*os.File, error) {
	dumpPath := fmt.Sprintf("%s.sql", databaseName)
	return os.Create(dumpPath)
}

func backupMySQL(mysqlconf mysqlConfig, w io.Writer) error {
	bckCmd := exec.Command("/usr/local/bin/mysqldump", "-h"+mysqlconf.host, "-P"+mysqlconf.port, "-u"+mysqlconf.username, "-p"+mysqlconf.password, mysqlconf.database)

	//bckCmd.Stdout = w http.ResponseWriter, r *http.Request
	bckCmd.Stdout = w
	return bckCmd.Run()
}

func main() {
	configPath := "/Users/joaocarreira/go/src/github.com/johnnybus/backup-mysql-s3/backup.conf"
	mysqlconf, s3conf := getConfig(configPath)
	w, err := genDumpFile(mysqlconf.database)
	check(err)
	defer w.Close()

	err = backupMySQL(mysqlconf, w)
	check(err)

	fmt.Printf("%v - %v\n", mysqlconf, s3conf)
}
