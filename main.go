package main

import (
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/gorilla/schema"
	"github.com/plotly/golang-api/plotly"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/graph", PostHandler)

	log.Println("Listening...")
	http.ListenAndServe(":8000", nil)
}

// Grant: Struct to store form input values
type Grant struct {
	Shares            int     `schema:"shares"`
	StrikePrice       int     `schema:"strikePrice"`
	Commencement      Date    `schema:"commencement"`
	VestingCliff      int     `schema:"vestingCliff"`
	VestingPeriod     int     `schema:"vestingPeriod"`
	VestingSchedule   string  `schema:"vestingSchedule"`
	OutstandingShares int     `schema:"outstandingShares"`
	CompanyValuation  int     `schema:"companyValuation"`
	NumRounds         int     `schema:"numRounds"`
	RoundDetails      []Round `schema:"roundDetails"`
	ExitDate          Date    `schema:"exitDate"`
	ExitValuation     Date    `schema:"exitValuation"`
}

// Round - Struct to store funding round information
type Round struct {
	Valuation    int  `schema:"valuation"`
	AmountRaised int  `schema:"amountRaised"`
	FundingDate  Date `schema:"fundingDate"`
}

// Date - Custom object to handle HTML date objects as time.Time in Go
type Date struct {
	time.Time
}

// UnmarshalText - Implements date parsing from HTML form to time.Time in Go
func (d *Date) UnmarshalText(text []byte) (err error) {
	d.Time, err = time.Parse("2006-01-02", string(text))
	return
}

// Generate takes a Grant object as input, and returns a URL to the graph
// displaying FD% over 4 years from commencement, or until Exit date.
func Generate(grant *Grant) string {
	var outDates plotly.Array
	var outValues plotly.Array
	//var exitDate time.Time

	thisDate := grant.Commencement.UTC()
	cliffDate := thisDate.AddDate(0, grant.VestingCliff, 0)
	lastDate := thisDate.AddDate(0, 48, 0)

	// if grant.ExitDate.Time.IsZero() != true {
	// 	exitDate = grant.ExitDate.UTC()
	// 	lastDate = grant.ExitDate.UTC()
	// }

	vestingInterval := 1 //month
	dilution := float64(1)

	if grant.VestingSchedule == "quarterly" {
		vestingInterval = 3
	}

	thisFDPercent := float64(0)
	thisMonth := 0
	doPlot := true

	// if grant.ExitDate.Time.IsZero() != true {
	// 	duration = grant.ExitDate.Time.Sub(grant.Commencement.Time)
	// }

	for doPlot == true {
		if thisDate.After(lastDate) {
			break
		}

		outDates = append(outDates, thisDate)

		if len(grant.RoundDetails) != 0 {
			if thisDate.After(grant.RoundDetails[0].FundingDate.Time) {
				dilution = dilution * (float64(grant.RoundDetails[0].Valuation) / (float64(grant.RoundDetails[0].Valuation) + float64(grant.RoundDetails[0].AmountRaised)))
				grant.RoundDetails = grant.RoundDetails[1:]
			}
		}

		if thisDate.After(cliffDate) {
			thisFDPercent = ((float64(grant.Shares) / float64(48) * float64(thisMonth)) / float64(grant.OutstandingShares)) * dilution
		}

		outValues = append(outValues, thisFDPercent)
		thisDate = thisDate.AddDate(0, vestingInterval, 0)
		thisMonth = thisMonth + vestingInterval
	}

	f := plotly.Figure{
		Data: []plotly.Trace{
			plotly.Trace{
				Type: "scatter",
				X:    outDates,
				Y:    outValues,
			},
		},
	}

	result, err := f.Save("Negotiator")
	if err != nil {
		fmt.Printf("Error happened %v", err)
		return err.Error()
	}

	return result.Url
}

// PostHandler serves to POST /graph
func PostHandler(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		return
	}

	grant := new(Grant)
	decoder := schema.NewDecoder()
	decoder.Decode(grant, req.PostForm)

	graphURL := Generate(grant)

	http.Redirect(res, req, graphURL, 301)
}
