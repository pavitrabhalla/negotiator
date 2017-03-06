package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	grant := &Grant{
		Shares:      500000,
		StrikePrice: 2,
		Commencement: Date{
			time.Now(),
		},
		VestingCliff:      12,
		VestingPeriod:     48,
		VestingSchedule:   "monthly",
		OutstandingShares: 2000000,
		CompanyValuation:  1000000,
		NumRounds:         2,
		RoundDetails: []Round{
			Round{
				Valuation:    2000000,
				AmountRaised: 18000000,
				FundingDate: Date{
					time.Now().AddDate(0, 18, 0),
				},
			},
			Round{
				Valuation:    250000000,
				AmountRaised: 750000000,
				FundingDate: Date{
					time.Now().AddDate(0, 30, 0),
				},
			},
		},
		ExitDate: Date{
			time.Now().AddDate(0, 40, 0),
		},
		ExitValuation: 30000000000,
		ExitAmount:    50000000000,
	}

	graphURL := Generate(grant)
	fmt.Printf("\n %s \n", graphURL)
}
