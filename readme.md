# GOV.UK Bank Holidays

This Go package provides a simple way to fetch and work with UK bank holidays using data from the GOV.UK API. The package supports holidays for different regions (England and Wales, Scotland, Northern Ireland) and includes features to check if a date is a holiday, find the next workday, and more.

## Installation

You can install this package using Go modules:

```bash
go get github.com/jeevanism/govuk-bank-holidays

```

## Sample Usage

```
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/yourusername/govuk-bank-holidays/bankholidays"
)

func main() {
    holidays, err := bankholidays.FetchHolidays()
    if err != nil {
        log.Fatalf("Failed to fetch holidays: %v", err)
    }

    bh := bankholidays.BankHolidays{
        Holidays: holidays,
        Weekend:  map[time.Weekday]bool{time.Saturday: true, time.Sunday: true},
    }

    today := time.Now()
    if bh.IsWorkDay(today, bankholidays.EnglandAndWales) {
        fmt.Println("Today is a workday")
    } else {
        fmt.Println("Today is not a workday")
    }

    nextWorkday := bh.GetNextWorkDay(today, bankholidays.EnglandAndWales)
    fmt.Println("Next workday is:", nextWorkday)
}

```

## License

This project is licensed under the MIT License
