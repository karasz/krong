package krong

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	u := User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(&u)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
		return
	}

	_, err = c.WriteOne(fmt.Sprintf("INSERT OR REPLACE INTO %s (name,email,type,address) VALUES('%s','%s','%s','%s');", USERS_TABLE_NAME, u.Name, u.Email, u.Type, u.Address))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error writting user to database", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
func getUser(w http.ResponseWriter, r *http.Request) {
	uid, e := strconv.Atoi(mux.Vars(r)["id"])
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
	}
	wr, err := c.QueryOne(fmt.Sprintf("SELECT * FROM %s WHERE _rowid_=%d", USERS_TABLE_NAME, uid))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	u := User{}
	for wr.Next() {
		wr.Scan(u)
	}
	if err := json.NewEncoder(w).Encode(u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}

}

func updateUser(w http.ResponseWriter, r *http.Request) {
	uid, e := strconv.Atoi(mux.Vars(r)["id"])
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
	}
	u := User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}
	_, err := c.WriteOne(fmt.Sprintf("REPLACE INTO %s (name,email,type,address) VALUES('%s','%s','%s','%s') WHERE _rowid_=%d;", USERS_TABLE_NAME, u.Name, u.Email, u.Type, u.Address, uid))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error writting user to database", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	uid, e := strconv.Atoi(mux.Vars(r)["id"])
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
	}
	_, err := c.WriteOne(fmt.Sprintf("DELETE FROM %s WHERE _rowid_=%d", USERS_TABLE_NAME, uid))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
}

func createJob(w http.ResponseWriter, r *http.Request) {
	j := Job{}
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(&j)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
		return
	}

	_, err = c.WriteOne(fmt.Sprintf("INSERT OR REPLACE INTO %s (name,displayname,timezone,schedule,successcount,errorcount,lastsuccess,lasterror,disabled,retries,concurrency,status,next,ephemeral,expiresat,webhook,agent,owner,type) VALUES ('%s','%s','%s','%s',%d,%d,'%s','%s',%d,%d,'%s','%s','%s',%d,'%s','%s','%s',%d,'%s')", JOBS_TABLE_NAME, j.Name, j.DisplayName, j.Timezone, j.Schedule, j.SuccessCount, j.ErrorCount, j.LastSuccess, j.LastError, j.Disabled.IsSetBool2Int(), j.Retries, j.Concurrency, j.Status, j.Next, j.Ephemeral.IsSetBool2Int(), j.ExpiresAt, j.WebHook, j.Agent, j.Owner, j.Type))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error writting user to database", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
func getJob(w http.ResponseWriter, r *http.Request) {
	jid, e := strconv.Atoi(mux.Vars(r)["id"])
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
		return
	}
	wr, err := c.QueryOne(fmt.Sprintf("SELECT * FROM %s WHERE _rowid_=%d", JOBS_TABLE_NAME, jid))
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	j := Job{}
	for wr.Next() {
		wr.Scan(j)
	}
	if err := json.NewEncoder(w).Encode(j); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}

}
func updateJob(w http.ResponseWriter, r *http.Request) {
	jid, e := strconv.Atoi(mux.Vars(r)["id"])
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
	}
	j := Job{}
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}
	_, err := c.WriteOne(fmt.Sprintf("REPLACE INTO %s (name,displayname,timezone,schedule,successcount,errorcount,lastsuccess,lasterror,disabled,retries,concurrency,status,next,ephemeral,expiresat,webhook,agent,owner) VALUES ('%s','%s','%s','%s',%d,%d,'%s','%s',%d,%d,'%s','%s','%s',%d,'%s','%s','%s',%d) WHERE _rowid_=%d", JOBS_TABLE_NAME, j.Name, j.DisplayName, j.Timezone, j.Schedule, j.SuccessCount, j.ErrorCount, j.LastSuccess, j.LastError, j.Disabled.IsSetBool2Int(), j.Retries, j.Concurrency, j.Status, j.Next, j.Ephemeral.IsSetBool2Int(), j.ExpiresAt, j.WebHook, j.Agent, j.Owner, jid))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error writting user to database", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&j); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

func deleteJob(w http.ResponseWriter, r *http.Request) {
	jid, e := strconv.Atoi(mux.Vars(r)["id"])
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
	}
	_, err := c.WriteOne(fmt.Sprintf("DELETE FROM %s WHERE _rowid_=%d", JOBS_TABLE_NAME, jid))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
}
