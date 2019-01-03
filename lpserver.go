package lpserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/influxdata/influxdb/models"

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

func processData(db string, p string, rp string, consistency string, data []byte) error {
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

// StartServer to start the line protocol server
func (lp *LPServer) StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/query", queryHandler)
	r.HandleFunc("/write", writeHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", lp.host, lp.port),
	}

	log.Fatal(srv.ListenAndServe())
}
