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
	}

	rows, err := db.Query("SELECT user_id, user_name, job_type, team_id FROM m_user")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	ret := make([]MUser, 0)
	for rows.Next() {

		var r MUser

		if err = rows.Scan(&r.UserId, &r.UserName, &r.JobType, &r.TeamId); err != nil {
			log.Fatal(err)
		}

		ret = append(ret, r)
		i = i + 1
	}

	encoder := json.NewEncoder(res)
	encoder.Encode(ret)
}

func main() {
	mux := goji.NewMux()
	mux.HandleFuncC(pat.Post("/api/get-user"), getUser)

	http.ListenAndServe(fmt.Sprintf("%s:%d", ServerHost, ServerPort), mux)
}
