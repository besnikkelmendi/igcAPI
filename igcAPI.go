package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/marni/goigc"
)

var urlArray []string                //used to store all the urls that we recive
var igcMap = make(map[int]igc.Track) //stores all the tracks that have been requested to us
var finalID int
var finalIDstr string
var field []string
var id string
var rspField string
var path string //used to check if the path and the handler corespond

type UrlForm struct {
	URL string `jason:"url"`
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
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

//returns the index of the string we are searching for and it returns -1 if i doesn't find it
func getIndex(x []string, y string) int {
	for i, j := range x {
		if j == y {
			//found
			return i
		}
	}
	return -1
}

//Calculates the total track lenght
func trackLength(track igc.Track) float64 {

	totalDistance := 0.0

	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return totalDistance
}

//This is the main handler
func postHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") //it's used to make the response show up as json

	//checking if the request method is post
	if r.Method == "POST" {
		//checking if the path that will accept the request is correct
		path = strings.TrimPrefix(r.URL.Path, "/igcinfo/api/igc/")
		if path != "" {
			http.Error(w, "", 400)
			return
		}
		Url := &UrlForm{}

		err := json.NewDecoder(r.Body).Decode(&Url) //obtaining the url recived from the request's body
		if err != nil {
			http.Error(w, err.Error(), 400) //checking for errors in the process and returning bad request if so
			return
		} else {

			track, err1 := igc.ParseLocation(Url.URL) //Used for parsing the obtained url

			if err1 != nil {
				http.Error(w, http.StatusText(400), http.StatusBadRequest) //checking for errors in the process and returning bad request if so
			} else {

				//checking if we haven't recived the same url before
				if getIndex(urlArray, Url.URL) == -1 {

					urlArray = append(urlArray, Url.URL)
					track, _ := igc.ParseLocation(Url.URL)
					igcMap[len(urlArray)-1] = track
				}

				uID := len(urlArray)

				track.UniqueID = strconv.Itoa(uID - 1) //I decided to use the array index as UniqueID

				resp := "{\n\"id\": " + "\"" + track.UniqueID + "\"\n}" //formating the response in json format

				w.Header().Set("Content-Type", "application/json")

				fmt.Fprint(w, resp)
			}
		}
	} else if r.Method == "GET" {

		id = strings.TrimPrefix(r.URL.Path, "/igcinfo/api/igc/") //the start of path falidation in this case checking if there's an id
		if id != "" {
			field = strings.Split(id, "/")
			//this one check if there's an field as well
			if field[1] != "" {
				finalIDstr = field[0]
				finalIDstr = strings.Replace(finalIDstr, "/", "", -1)
				finalID0, err := strconv.Atoi(finalIDstr)
				if err != nil {
					http.Error(w, "", http.StatusBadRequest)
				} else {
					finalID = finalID0
					if finalID > len(igcMap)-1 {
						http.Error(w, "", 404) //Returns an error if the id request isn't in our system
					} else {

						rspField = field[1]
						rspField = strings.Replace(rspField, "/", "", -1)

						//checking for the requested field and returning the value if it's found
						switch rspField {
						case "pilot":
							fmt.Fprintf(w, "%s", igcMap[finalID].Pilot)
							break
						case "glider":
							fmt.Fprintf(w, "%s", igcMap[finalID].GliderType)
							break
						case "glider_id":
							fmt.Fprintf(w, "%s", igcMap[finalID].GliderID)
							break
						case "track_length":
							fmt.Fprintf(w, "%f", trackLength(igcMap[finalID]))
							break
						case "H_date":
							fmt.Fprintf(w, "%s", igcMap[finalID].Date.String())
							break
						default:
							http.Error(w, "", http.StatusNotFound)
						}
					}
				}

			} else {

				id = strings.Replace(id, "/", "", -1)
				finalID0, err := strconv.Atoi(id)
				if err != nil {
					fmt.Fprint(w, "The id you wrote is not an integer")
					return
				} else {
					finalID = finalID0
					if finalID > len(igcMap)-1 {
						http.Error(w, "", 404)
					} else {
						resp := "{\n"
						resp += "  \"H_date\": " + "\"" + igcMap[finalID].Date.String() + "\",\n"
						resp += "  \"pilot\": " + "\"" + igcMap[finalID].Pilot + "\",\n"
						resp += "  \"glider\": " + "\"" + igcMap[finalID].GliderType + "\",\n"
						resp += "  \"glider_id\": " + "\"" + igcMap[finalID].GliderID + "\",\n"
						resp += "  \"track_lenght\": " + "\"" + fmt.Sprintf("%f", trackLength(igcMap[finalID])) + "\"\n"
						resp += "}"
						fmt.Fprint(w, resp)
					}
				}
			}
		} else {

			w.Header().Set("Content-Type", "application/json")

			resp := "["
			for i, _ := range urlArray {

				resp += strconv.Itoa(i)
				if i+1 == len(urlArray) {
					break
				}
				resp += ","
			}
			resp += "]"

			fmt.Fprint(w, resp)
		}

	} else {

		w.WriteHeader(http.StatusNotFound)
	}
}

var start = time.Now()

func IGChandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	end := time.Now()
	duration := end.Sub(start)
	fmt.Fprint(w, "{\n\"uptime\": \""+duration.String()+"\",\n\"info\": \"Service for IGC tracks.\",\n \"version0\": \"v1.4\"}")

}

func main() {
	http.HandleFunc("/igcinfo/api/", IGChandler)
	http.HandleFunc("/igcinfo/api/igc/", postHandler)
	//	http.ListenAndServe(":8080", nil)
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
