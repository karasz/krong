package krong

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/amoghe/distillog"
	"github.com/gorilla/mux"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/robfig/cron/v3"
	"github.com/rqlite/gorqlite"
	"gopkg.in/natefinch/lumberjack.v2"
)

type KronG struct {
	database gorqlite.Connection
	logger   distillog.Logger
	pidFile  string
	*cron.Cron
	config Config
	http.Server
}

var (
	ErrProcessRunning = errors.New("process is running")
	ErrFileStale      = errors.New("pidfile exists but process is not running")
	ErrFileInvalid    = errors.New("pidfile has invalid contents")
)

func InitKronG() *KronG {
	var c Config

	err := hclsimple.DecodeFile("config.hcl", nil, &c)
	if err != nil {
		fmt.Println(err)
	}

	lumberjackHandle := &lumberjack.Logger{
		Filename:   c.LogFile,
		MaxSize:    500,
		MaxAge:     28,
		MaxBackups: 3,
		LocalTime:  false,
		Compress:   false,
	}

	_, err = lumberjackHandle.Write([]byte(""))
	if err != nil {
		fmt.Println(err)
	}
	klogger := distillog.NewStreamLogger("krong", lumberjackHandle)

	data, err := initializeDatabase(c.DBURL)
	if err != nil {
		klogger.Errorf("Error connecting to database. %s\n", err.Error())
	}

	klogger.Infoln("database successfully connected.")

	k := &KronG{
		database: data,
		pidFile:  c.PidFile,
		logger:   klogger,
		Cron:     cron.New(cron.WithSeconds()),
		config:   c,
		Server:   http.Server{},
	}
	return k

}

func (k *KronG) StartServer() {
	defer k.Shutdown()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go k.SignalHandler(signals)

	if err := k.writePID(); err != nil {
		log.Printf("Can't create pid file! %v \n", err)
		os.Exit(1)
	}
	defer k.removePID()
	k.Start()
	go k.startAPI()
	js, err := getAllJobs()
	if err != nil {
		k.logger.Errorln("Cannot load Jobs from database")
		k.Shutdown()
	}
	for _, j := range js {
		k.Cron.AddJob(j.Schedule, j)
	}

	<-signals

}

func (k *KronG) startAPI() {
	k.logger.Infoln("Started API server")
	r := mux.NewRouter()

	setHandlers(r)
	// We might want to use a different port, but we do not want to expose this
	k.logger.Errorln(http.ListenAndServe("127.0.0.1:8080", r))
}

func setHandlers(r *mux.Router) {
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	r.HandleFunc("/jobs", createUser).Methods("POST")
	r.HandleFunc("/jobs/{id}", getUser).Methods("GET")
	r.HandleFunc("/jobs/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/jobs/{id}", deleteUser).Methods("DELETE")
}

func (k *KronG) Shutdown() {
	k.Stop()
	dbname := fmt.Sprintf("krong-%s.sql", time.Now().Format("2006-01-02-15-04-05"))
	err := backupDB(k.config.DBURL, dbname, true)
	if err != nil {
		k.logger.Warningln(err)
	}

	k.removePID()
	k.database.Close()
	k.logger.Infoln("database successfully closed.")
	err = k.Server.Shutdown(context.Background())
	if err != nil {
		k.logger.Infoln("error shutting down the api server :%s", err.Error())
	}

	k.logger.Close()
	os.Exit(0)
}

func (k *KronG) SignalHandler(signals chan os.Signal) {
	for {
		sig := <-signals
		//k.Logger.Infof("King KronG.  Received a signal: %v\n", sig.String())
		switch sig {
		//case syscall.SIGUSR1:
		//case syscall.SIGUSR2:
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGTERM:
			k.logger.Infoln("King KronG. Shutdown Server")
			k.Shutdown()
		default:
			k.logger.Infof("King KronG. Received a UNKNOWN signal: %v\n", sig.String())
		}
	}
}

func (k *KronG) writePID() error {
	return write(k.pidFile, os.Getpid())
}

func write(filename string, pid int) error {
	content, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		content = []byte("0")
	} else if err != nil {
		return err
	}

	oldpid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return ErrFileInvalid
	}

	if err == nil {
		if pidIsRunning(oldpid) {
			return ErrProcessRunning
		}
	}

	return os.WriteFile(filename, []byte(fmt.Sprintf("%d\n", pid)), 0644)
}

func (k *KronG) removePID() error {
	return os.RemoveAll(k.pidFile)
}

func pidIsRunning(pid int) bool {
	process, err := os.FindProcess(pid)

	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil && err.Error() == "os: process not initialized" {
		return false
	}

	if err != nil && err.Error() == "os: process already finished" {
		return false
	}

	return true
}
