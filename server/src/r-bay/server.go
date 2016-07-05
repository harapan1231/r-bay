package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"

	"database/sql"
	_ "github.com/lib/pq"
)

const (
	ServerHost = "******"
	ServerPort = "******"
	DbHost     = "******"
	DbPort     = "******"
	DbUser     = "******"
	DbPwd      = "******"
)

type MUser struct {
	UserId   string
	UserName string
	JobType  string
	TeamId   string
}

func getUser(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", DbUser, DbPwd, DbUser, DbHost, DbPort))
	if err != nil {
		log.Fatal(err)
		return
	}

	rows, err := db.Query("SELECT user_id, user_name, job_type, team_id FROM m_user")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	i := 0
	ret := make([]MUser, 0)
	for rows.Next() {

		var r MUser

		if err = rows.Scan(&r.UserId, &r.UserName, &r.JobType, &r.TeamId); err != nil {
			log.Fatal(err)
			return
		}

		ret = append(ret, r)
		i = i + 1
	}

	encoder := json.NewEncoder(res)
	encoder.Encode(ret)
}

func putUserState(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	userId := pat.Param(ctx, "uid")
	date := pat.Param(ctx, "date")
	time := pat.Param(ctx, "time")
	state := pat.Param(ctx, "state")

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", DbUser, DbPwd, DbUser, DbHost, DbPort))
  defer db.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	rows, err := db.Query("SELECT 1 FROM t_user_state WHERE user_id = $1 AND date = $2 AND time = $3", userId, date, time)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	var query string
	if rows.Next() {
		query = "UPDATE t_user_state state = $4 WHERE user_id = $1 AND date = $2 AND time = $3"
	} else {
		query = "INSERT INTO t_user_state (user_id, date, time, state) VALUES ($1, $2, $3, $4)"
	}
	_, err = db.Exec(query, userId, date, time, state)
	if err != nil {
		log.Fatal(err)
		return
	}
  //TODO return response code
}

func main() {
	mux := goji.NewMux()
	mux.HandleFuncC(pat.Post("/api/get-user"), getUser)
  mux.HandleFuncC(pat.Post("/api/put-user-state/:uid/:date/:time/:state"), putUserState)

	http.ListenAndServe(fmt.Sprintf("%s:%d", ServerHost, ServerPort), mux)
}
