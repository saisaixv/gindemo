package main

import (
	"fmt"
	"syscall"
	"time"

	"github.com/saisai/gindemo/api"
	"github.com/saisai/gindemo/models"

	"github.com/saisai/gindemo/utils/cache"

	"github.com/saisai/gindemo/utils/log"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"github.com/go-xorm/xorm"
	//	"gopkg.in/redsync.v1"
	_ "github.com/go-sql-driver/mysql"
)

const (
	Version = "1.0.0"
)

var (
	config      *ini.File
	apiAddr     string
	ConfigFiles []interface{} = []interface{}{"/opt/saisai/profile.ini"}
)

func initConfig() error {
	var err error
	profile, _ := ConfigFiles[0].(string)
	if err = syscall.Access(profile, syscall.F_OK); err != nil {
		return err
	}
	config, err = ini.LooseLoad(profile, ConfigFiles[0:]...)
	if err != nil {
		return err
	}
	return nil
}

func initRedis(cfg *ini.File) (err error) {
	sec, err := cfg.GetSection("redis")
	if err != nil {
		return err
	}
	url := sec.Key("url").String()
	log.Infof("[init redis] url:'%s'", url)
	//	auth := sec.Key("password").String()
	//	log.Infof("[init redis] pwd:'%s'", auth)

	cache.Init(url, "", 20, 20, 10)

	log.Info("[init redis success]")

	err = cache.Ping()
	if err != nil {
		return err
	}

	return
}

func initDB(cfg *ini.File) (err error) {
	//create database
	sec, err := cfg.GetSection("db")
	if err != nil {
		return err
	}

	driver := sec.Key("driver").String()
	source := sec.Key("source").String()
	showSql := sec.Key("show_sql").MustBool()
	utc := sec.Key("utc").MustBool(false)
	useCache := sec.Key("use_cache").MustBool(false)
	log.Infof("[init DB] url:'%s %s' show_sql:%t, utc:%t, cache:%t\n", driver, source, showSql, utc, useCache)

	db, err := xorm.NewEngine(driver, source)
	if err != nil {
		return
	}
	if utc {
		db.TZLocation = time.UTC
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(20)
	db.ShowSQL(showSql)

	models.InitDB(db)

	//sync table struct
	if err := models.SyncTables(); err != nil {
		return err
	}
	return
}

func initApi(cfg *ini.File) error {
	sec, err := cfg.GetSection("api")
	if err != nil {
		return err
	}

	apiAddr = sec.Key("addr").String()
	debug := sec.Key("debug").MustBool(false)

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	return nil
}

func initApplication() error {

	if err := initRedis(config); err != nil {
		fmt.Println("initRedis err")
		return err
	}

	if err := initDB(config); err != nil {
		fmt.Println("initDB err")
		return err
	}

	if err := initApi(config); err != nil {

	}

	return nil
}

func run() {
	log.Infof("http run %s\n", apiAddr)
	err := api.Engine().Run(apiAddr)
	log.Error(err)
}

func main() {
	log.Infof("Start ...\n")
	if err := initConfig(); err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	if err := initApplication(); err != nil {
		log.Fatal(err)
	}

	run()

}
