package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID       string
	Cmd      string
	Args     []string
	Interval int
}

// Task map
type TaskMap struct {
	m map[string]chan<- int
	sync.Mutex
}

func NewTaskMap() *TaskMap {
	return &TaskMap{
		m: make(map[string]chan<- int),
	}
}

func (tm *TaskMap) Add(t Task) error {
	tm.Lock()
	defer tm.Unlock()
	_, ok := tm.m[t.ID]
	if ok {
		return errors.New("The task print-time already exists.")
	}
	done := make(chan int, 1)
	tm.m[t.ID] = done

	go func(ms int) {
		ticker := time.NewTicker(time.Millisecond * time.Duration(ms))
		for {
			select {
			case <-done:
				log.Println("crontask cancel!")
				return
			case <-ticker.C:
				Args := ""
				for _, v := range t.Args {
					Args = Args + v + " "
				}
				Args = strings.Trim(Args, " ")
				cmd := exec.Command(t.Cmd, Args)
				log.Println("Running command and waiting for it to finish... [ID]", t.ID)
				err := cmd.Run()
				log.Printf("Command finished with error: %v", err)
			}
		}
	}(t.Interval)
	return nil
}

func (tm *TaskMap) Del(ID string) error {
	tm.Lock()
	defer tm.Unlock()
	_, ok := tm.m[ID]
	if !ok {
		return errors.New("The task print-time is not found.")
	}
	done := tm.m[ID]
	delete(tm.m, ID)

	done <- 0
	return nil
}

type Reply struct {
	Ok    string
	Id    string
	Error string
}

var gTaskMap *TaskMap

func main() {
	gTaskMap = NewTaskMap()

	/*
		arg := flag.Int("port", 4567, "port for gocron")
		flag.Parse()
		port := ":"+strconv.Itoa(*arg)
	*/
	port := ":4567"
	if len(os.Args) == 2 {
		strPort := os.Args[1]
		p, err := strconv.Atoi(strPort)
		if err != nil {
			log.Fatal("invalid args")
		} else {
			if p > 1024 && p < 65536 {
				port = strconv.Itoa(p)
			} else {
				log.Fatal("invalid args")
			}
		}
	}

	http.HandleFunc("/", handleFunc)
	log.Println("GoCron listening on ", port)
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	rsp := Reply{
		Ok:    "true",
		Id:    "",
		Error: "",
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("get body fail")
		rsp.Ok = "false"
		rsp.Error = err.Error()
		reply(w, http.StatusBadRequest, rsp)
		return
	}

	var t Task
	if err := json.Unmarshal([]byte(body), &t); err != nil {
		rsp.Ok = "false"
		rsp.Error = err.Error()
		reply(w, http.StatusBadRequest, rsp)
		return
	}
	log.Println(t.ID, t.Cmd, t.Args[0], t.Interval, r.Method)

	switch r.Method {
	case "POST":
		if err := gTaskMap.Add(t); err != nil {
			rsp.Ok = "false"
			rsp.Error = err.Error()
			reply(w, http.StatusConflict, rsp)
			return
		} else {
		}
	case "DELETE":
		if err := gTaskMap.Del(t.ID); err != nil {
			rsp.Ok = "false"
			rsp.Error = err.Error()
			reply(w, http.StatusNotFound, rsp)
			return
		} else {
		}
	}

	rsp.Id = t.ID
	reply(w, http.StatusOK, rsp)
}

func reply(w http.ResponseWriter, status int, reply Reply) {
	w.WriteHeader(status)
	b, err := json.Marshal(reply)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprintf(w, "%s", b)
}
