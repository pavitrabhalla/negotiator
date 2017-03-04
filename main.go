package main

import (
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/plotly/golang-api/plotly"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/graph", postHandler)

	log.Println("Listening...")
	http.ListenAndServe(":8000", nil)
}

type grant struct {
	shares            int
	strikePrice       int
	commencement      time.Time
	vestingCliff      int
	vestingPeriod     int
	vestingSchedule   string
	outstandingShares int
	companyValuation  int
	numRounds         int
	roundDetails      []round
}

type round struct {
	valuation    int
	amountRaised int
}

func postHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	log.Println(req.Form)
	f := plotly.Figure{
		Data: []plotly.Trace{
			plotly.Trace{
				Type: "scatter",
				X: plotly.Array{
					4.54, 3, 34, 35, 362,
				},
				Y: plotly.Array{
					1, 2, 3, 4, 5,
				},
			},
		},
	}

	result, err := f.Save("new golang file")
	if err != nil {
		fmt.Printf("Error happened %v", err)
		return
	}

	http.Redirect(res, req, result.Url, 301)
}
