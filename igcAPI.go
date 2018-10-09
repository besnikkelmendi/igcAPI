package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/marni/goigc"
)

var urlArray []string

type UrlForm struct {
	URL string `jason:"url"`
}

func searchArray(x []string, y string) bool {
	for _, i := range x {
		if i == y {
			//found
			return false
		}
	}
	return true
}

func getIndex(x []string, y string) int {
	for i, j := range x {
		if j == y {
			//found
			return i
		}
	}
	return -1
}

func postHandler(w http.ResponseWriter, r *http.Request) {

	//igcMap := make(map[int64]igc.Track)

	if r.Method == "POST" {

		Url := &UrlForm{}
		//	Url.URL = r.FormValue("url")
		err := json.NewDecoder(r.Body).Decode(&Url)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		if getIndex(urlArray, Url.URL) == -1 {
			urlArray = append(urlArray, Url.URL)
		}

		track, _ := igc.ParseLocation(Url.URL)

		track.UniqueID = "besi"

		uID := len(urlArray)

		track.UniqueID = strconv.Itoa(uID)

		resp := "{\n\"id\": " + "\"" + track.UniqueID + "\"\n}"

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, resp)
	}
	if r.Method == "GET" {

		w.Header().Set("Content-Type", "application/json")

		resp := "["
		for i, _ := range urlArray {

			resp += strconv.Itoa(i + 1)
			if i+1 == len(urlArray) {
				break
			}
			resp += ","
		}
		resp += "]"

		fmt.Fprint(w, resp)

	}

}

var start = time.Now()

func IGChandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	end := time.Now()
	duration := end.Sub(start)
	fmt.Fprint(w, "{\n\"uptime\": \""+duration.String()+"\",\n\"info\": \"Service for IGC tracks.\",\n \"version0\": \"v1\"}")

}

func main() {

	http.HandleFunc("/igcinfo/api/", IGChandler)
	http.HandleFunc("/igcinfo/api/igc/", postHandler)
	//http.HandleFunc("/igcinfo/api/igc/")
	http.ListenAndServe(":8080", nil)

	s := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	track, err := igc.ParseLocation(s)
	if err != nil {
		fmt.Errorf("Problem reading the track", err)
	}

	fmt.Printf("Pilot: %s, gliderType: %s, date: %s",
		track.Pilot, track.GliderType, track.Date.String())
}
