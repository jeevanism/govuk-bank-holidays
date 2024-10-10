package main

import (
	"fmt"
	"govuk-bank-holidays/bankholidays"
	"log"
	"time"
)

func main() {
	// Fetch holidays from GOV.UK (or fallback to backup)
	holidays, err := bankholidays.FetchHolidays()
	if err != nil {
		log.Println("Error loading holidays:", err)
		return
	}

	// Initialize BankHolidays struct
	bh := bankholidays.BankHolidays{
		Holidays: holidays,
		Weekend:  map[time.Weekday]bool{time.Saturday: true, time.Sunday: true},
	}

	// Get today's date
	today := time.Now()

	// Check if today is a workday
	if bh.IsWorkDay(today, bankholidays.EnglandAndWales) {
		fmt.Println("Today is a workday")
	} else {
		fmt.Println("Today is not a workday")
	}

	// Get the next workday
	nextWorkday := bh.GetNextWorkDay(today, bankholidays.EnglandAndWales)
	fmt.Println("Next workday is:", nextWorkday)

	// Get the next holiday
	nextHoliday := bh.GetNextHoliday(today, bankholidays.EnglandAndWales)
	if nextHoliday != nil {
		fmt.Printf("Next holiday is %s on %s\n", nextHoliday.Title, nextHoliday.Date.Format("2006-01-02"))
	} else {
		fmt.Println("No upcoming holidays")
	}
}
