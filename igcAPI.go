package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/marni/goigc"
)

var urlArray []string
var igcMap = make(map[int]igc.Track)
var finalID int
var finalIDstr string
var field []string
var id string

type UrlForm struct {
	URL string `jason:"url"`
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
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

func trackLength(track igc.Track) float64 {

	totalDistance := 0.0

	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return totalDistance
}

func getHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")

	}

}

func postHandler(w http.ResponseWriter, r *http.Request) {

	//igcMap := make(map[int64]igc.Track)
	w.Header().Set("Content-Type", "application/json")

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
			track, _ := igc.ParseLocation(Url.URL)
			igcMap[len(urlArray)-1] = track
		}

		track, _ := igc.ParseLocation(Url.URL)

		//track.UniqueID = "besi"

		uID := len(urlArray)

		track.UniqueID = strconv.Itoa(uID - 1)

		resp := "{\n\"id\": " + "\"" + track.UniqueID + "\"\n}"

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, resp)
	}
	if r.Method == "GET" {

		id = strings.TrimPrefix(r.URL.Path, "/igcinfo/api/igc/")
		if id != "" {
			field = strings.Split(id, "/")
			if field[1] != "" {
				finalIDstr = field[0]
				finalIDstr = strings.Replace(finalIDstr, "/", "", -1)
				finalID0, err := strconv.Atoi(finalIDstr)

				if err != nil {
					fmt.Fprint(w, "The id you wrote is not an integer")
				} else {
					finalID = finalID0
				}

			} else {

				id = strings.Replace(id, "/", "", -1)
				finalID0, err := strconv.Atoi(id)
				if err != nil {
					fmt.Fprint(w, "The id you wrote is not an integer")
				} else {
					finalID = finalID0
				}

			}

			resp := "{"
			resp += "\"H_date\": " + "\"" + igcMap[finalID].Date.String() + "\","
			resp += "\"pilot\": " + "\"" + igcMap[finalID].Pilot + "\","
			resp += "\"glider\": " + "\"" + igcMap[finalID].GliderType + "\","
			resp += "\"glider_id\": " + "\"" + igcMap[finalID].GliderID + "\","
			resp += "\"track_lenght\": " + "\"" + fmt.Sprintf("%f", trackLength(igcMap[finalID])) + "\""
			resp += "}"

			fmt.Fprint(w, resp)
		} else {

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

	/*s := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	track, err := igc.ParseLocation(s)
	if err != nil {
		fmt.Errorf("Problem reading the track", err)
	}

	fmt.Printf("Pilot: %s, gliderType: %s, date: %s",
		track.Pilot, track.GliderType, track.Date.String())*/
}
