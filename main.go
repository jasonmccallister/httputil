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
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type response struct {
	Environment map[string]string `json:"environment"`
	Headers     map[string]string `json:"headers"`
	Status      int               `json:"status"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("listening on port", port)

	http.HandleFunc("/", envHandler)
	http.HandleFunc("/redis", redisHandler)
	http.HandleFunc("/postgres", postgresHandler)
	http.HandleFunc("/mysql", mysqlHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func redisHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_URL"), os.Getenv("REDIS_PORT"))

	db := 0
	if os.Getenv("REDIS_DB") != "" {
		e, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("REDIS_DB must be a valid integer"))
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
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("unable to ping the REDIS_URL, err: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(fmt.Sprintf(`{"ping":"%s"}`, pong))); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := response{}

	// get the environment variables
	env := map[string]string{}
	for _, v := range os.Environ() {
		v := strings.SplitN(v, "=", 2)
		if strings.Contains(v[0], "DB_DSN") || strings.Contains(v[0], "DB_PASSWORD") || strings.Contains(v[0], "DB_SERVER") {
			env[v[0]] = "************"
			fmt.Println("Found", v[0], "with the value", v[1])
		} else {
			env[v[0]] = v[1]
		}

	}
	resp.Environment = env

	// get the request headers
	headers := map[string]string{}
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			headers[name] = value
		}
	}
	resp.Headers = headers
	resp.Status = http.StatusOK

	b, err := json.Marshal(resp)
	if err != nil {
		_, _ = w.Write([]byte("{}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func postgresHandler(w http.ResponseWriter, r *http.Request) {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_SERVER"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("pgqsl", conn)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("unable to open the database err: " + err.Error()))
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("unable to ping the database, err: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	pong := fmt.Sprintf(`{"ping":"%s"}`, "PONG")
	if _, err := w.Write([]byte(pong)); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}

func mysqlHandler(w http.ResponseWriter, r *http.Request) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("unable to open the database err: " + err.Error()))
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("unable to ping the database, err: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	pong := fmt.Sprintf(`{"ping":"%s"}`, "PONG")
	if _, err := w.Write([]byte(pong)); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}
