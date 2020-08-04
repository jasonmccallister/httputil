package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func redisHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_URL"), os.Getenv("REDIS_PORT"))

	db := 0
	if os.Getenv("REDIS_DB") != "" {
		e, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("REDIS_DB must be a valid integer"))
			return
		}

		db = e
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("unable to ping the REDIS_URL, err: " + err.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(pong))
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	env := map[string]string{}

	for _, keyval := range os.Environ() {
		keyval := strings.SplitN(keyval, "=", 2)
		env[keyval[0]] = keyval[1]
	}

	bytes, err := json.Marshal(env)
	if err != nil {
		w.Write([]byte("{}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(bytes))
}

func postgresHandler(w http.ResponseWriter, r *http.Request) {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_SERVER"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	fmt.Println(conn)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to open the postgres database err: " + err.Error()))
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to ping the postgres database, err: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PONG"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	fmt.Println("listening on port", port)

	http.HandleFunc("/", envHandler)
	http.HandleFunc("/redis", redisHandler)
	http.HandleFunc("/postgres", postgresHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
