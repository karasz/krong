package krong

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/karasz/krong/isset"
	"github.com/rqlite/gorqlite"
)

const USERS_TABLE_NAME = "users"
const JOBS_TABLE_NAME = "jobs"

var c gorqlite.Connection

func initializeDatabase(name string) (gorqlite.Connection, error) {

	conn, err := gorqlite.Open(name)
	if err != nil {
		return gorqlite.Connection{}, err
	}

	conn.SetConsistencyLevel("strong")
	c = conn
	createTables()
	return conn, nil
}

func createTables() {

	err := createUsersTable()
	if err != nil {
		fmt.Printf("There was an error creating Users table. The error was %s", err.Error())
	}
	err = createJobsTable()
	if err != nil {
		fmt.Printf("There was an error creating Jobs table. The error was %s", err.Error())
	}

}

func dropTables() {
	jobsstr := "DROP TABLE jobs"
	_, jobserr := c.WriteOne(jobsstr)
	if jobserr != nil {
		fmt.Println(jobserr)
	}
	usersstr := "DROP TABLE users"
	_, userserr := c.WriteOne(usersstr)
	if userserr != nil {
		fmt.Println(userserr)
	}

}

func createUsersTable() error {
	tableString := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER NOT NULL PRIMARY KEY,name TEXT, email TEXT, type TEXT, address TEXT)", USERS_TABLE_NAME)
	_, err := c.WriteOne(tableString)
	if err != nil {
		return err
	}
	return nil
}

func createJobsTable() error {
	tableString := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER NOT NULL PRIMARY KEY,name TEXT,displayname TEXT,timezone TEXT,schedule TEXT,successcount INTEGER,errorcount INTEGER,lastsuccess TEXT,lasterror TEXT,disabled INTEGER,retries INTEGER,concurrency TEXT,status TEXT,next TEXT,ephemeral INTEGER,expiresat TEXT,webhook TEXT,agent TEXT,owner INTEGER,type TEXT)", JOBS_TABLE_NAME)
	_, err := c.WriteOne(tableString)
	if err != nil {
		return err
	}
	return nil
}

func backupDB(db, filepath string, sql bool) error {
	db = db + "db/backup"
	if sql {
		db = db + "?fmt=sql"
	}

	resp, err := http.Get(db)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func restoreDB(db, filepath string) error {
	bak, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err.Error())

	}
	defer bak.Close()

	db = db + "db/load"

	req, err := http.NewRequest("POST", db, bak)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}
func getAllJobs() ([]*Job, error) {
	wr, err := c.QueryOne(fmt.Sprintf("SELECT * FROM %s;", JOBS_TABLE_NAME))
	if err != nil {
		return []*Job{}, err
	}
	allJobs := make([]*Job, 0)
	j := &Job{}
	for wr.Next() {
		mp, err := wr.Map()

		if err != nil {
			fmt.Println(err)
			return []*Job{}, err
		}

		j.ID = returnInt(mp["id"])
		j.Name = returnString(mp["name"], false)
		j.DisplayName = returnString(mp["displayname"], false)
		j.Timezone = returnString(mp["timezone"], false)
		j.Schedule = returnString(mp["schedule"], false)
		j.SuccessCount = returnInt(mp["successcount"])
		j.ErrorCount = returnInt(mp["errorcount"])
		j.LastSuccess = returnTime(mp["lastsuccess"])
		j.LastError = returnTime(mp["lasterror"])
		b := isset.Bool{}
		j.Disabled = b.IsSetInt2Bool(returnInt(mp["disabled"]))
		j.Retries = returnInt(mp["retries"])
		j.Concurrency = returnString(mp["concurency"], false)
		j.Status = returnString(mp["status"], false)
		j.Next = returnTime(mp["next"])
		b.Reset()
		j.Ephemeral = b.IsSetInt2Bool(returnInt(mp["ephemeral"]))
		j.ExpiresAt = returnTime(mp["expiresat"])
		o := &WebHook{}
		err = json.Unmarshal([]byte(returnString(mp["webhook"], true)), o)
		if err != nil {
			j.WebHook = WebHook{}
		} else {
			j.WebHook = *o
		}
		a := &Agent{}
		err = json.Unmarshal([]byte(returnString(mp["agent"], true)), a)
		if err != nil {
			j.Agent = Agent{}
		} else {
			j.Agent = *a
		}
		j.Owner = returnInt(mp["owner"])
		j.Type = returnString(mp["type"], false)

		allJobs = append(allJobs, j)
	}
	return allJobs, nil
}

func returnString(i interface{}, json bool) string {
	if i != nil {
		return i.(string)
	}
	if json {
		return "{}"
	}
	return ""
}

func returnInt(i interface{}) int {
	if i != nil {
		return int(i.(float64))
	}
	return 0
}

func returnTime(i interface{}) time.Time {
	const layout = "2006-01-02 15:04:05"
	if i != nil {
		if t, err := time.Parse(layout, returnString(i, false)); err == nil {
			return t
		}
		q, _ := time.Parse(time.RFC3339, returnString(i, false))
		return q
	}
	return time.Time{}
}
