package lpserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

// LPServer is Line Protocol Server
type LPServer struct {
	host string
	port int
}

// NewLPServer to create new line protocol server
func NewLPServer(host string, port int) *LPServer {
	return &LPServer{
		host: host,
		port: port,
	}
}

type lpData struct {
	sync.Mutex
	name string
	file *os.File
}

var state *lpData

func processData(db string, p string, rp string, consistency string, data []byte) error {
	fmt.Printf("db: %s\n", db)
	/*
		if points, err := models.ParsePointsWithPrecision(data, time.Now().UTC(), p); err == nil {
			for _, point := range points {
				name := point.Name()
				fmt.Println(string(name))
				tags := point.Tags()
				for i, tag := range tags {
					fmt.Printf("%d [%s = %s]\n", i, tag.Key, tag.Value)
				}
				if fields, err := point.Fields(); err == nil {
					for k := range fields {
						fmt.Printf("[%s=%s]\n", k, fields[k])
					}
				} else {
					return err
				}
			}
		} else {
			return err
		}
	*/
	return nil
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer r.Body.Close()

	db := r.FormValue("db")
	p := r.FormValue("precision")
	rp := r.FormValue("rp")
	consistency := r.FormValue("consistency")

	processData(db, p, rp, consistency, data)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

type body struct {
	Name string `json:"name"`
}

func getStoreName(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	body := &body{}
	if err := decoder.Decode(body); err != nil {
		return "", err

	}
	return body.Name, nil
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["op"]
	switch action {
	case "open":
		if name, err := getStoreName(r); err == nil {
			state.Lock()
			if state.name == "" {
				if state.file, err = os.Create(name); err == nil {
					state.name = name
					log.Printf("open: %s\n", name)
				} else {
					log.Printf("%v", err)
				}
			}
			state.Unlock()
		}
	case "close":
		if name, err := getStoreName(r); err == nil {
			state.Lock()
			state.Unlock()
			if state.name != "" {
				if err = state.file.Close(); err == nil {
					state.file = nil
					state.name = ""
					log.Printf("close: %s\n", name)
				} else {
					log.Printf("%v", err)
				}
			}

		}
	}
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func lpInit() {
	state = &lpData{}
}

// StartServer to start the line protocol server
func (lp *LPServer) StartServer() {
	lpInit()
	r := mux.NewRouter()
	r.HandleFunc("/query", queryHandler)
	r.HandleFunc("/write", writeHandler)
	r.HandleFunc("/store/{op}", storeHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", lp.host, lp.port),
	}
	log.Fatal(srv.ListenAndServe())
}
