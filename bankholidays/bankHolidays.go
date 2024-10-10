package bankholidays

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Division constants
const (
	EnglandAndWales = "england-and-wales"
	Scotland        = "scotland"
	NorthernIreland = "northern-ireland"
)

var AllDivisions = []string{EnglandAndWales, Scotland, NorthernIreland}

// Holiday struct for storing holiday information
type Holiday struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Notes   string    `json:"notes,omitempty"`
	Bunting bool      `json:"bunting,omitempty"`
}

// BankHolidays struct for holding holidays and weekend information
type BankHolidays struct {
	Holidays map[string][]Holiday
	Weekend  map[time.Weekday]bool
}

const govUKURL = "https://www.gov.uk/bank-holidays.json"

// LoadBackupData loads the local backup JSON file in case the request fails
func LoadBackupData() (map[string][]Holiday, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Construct the absolute path to bank-holidays.json based on the current directory
	backupPath := filepath.Join(currentDir, "bankholidays", "bank-holidays.json")

	// Print the path for debugging purposes
	fmt.Println("Looking for bank-holidays.json at:", backupPath)

	// Open the backup file
	file, err := os.Open(backupPath)
	if err != nil {
		log.Printf("Error loading holidays: %v", err)
		return nil, err
	}
	defer file.Close()

	var rawData map[string]struct {
		Events []map[string]interface{} `json:"events"`
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&rawData)
	if err != nil {
		return nil, err
	}

	// Flatten and convert the structure into map[string][]Holiday
	holidays := make(map[string][]Holiday)
	for division, item := range rawData {
		var holidayList []Holiday
		for _, event := range item.Events {
			// Parse the date in "YYYY-MM-DD" format
			date, err := time.Parse("2006-01-02", event["date"].(string))
			if err != nil {
				log.Printf("Error parsing date for holiday: %v", err)
				continue
			}
			holidayList = append(holidayList, Holiday{
				Title:   event["title"].(string),
				Date:    date,
				Notes:   event["notes"].(string),
				Bunting: event["bunting"].(bool),
			})
		}
		holidays[division] = holidayList
	}

	return holidays, nil
}

// FetchHolidays fetches data from the GOV.UK API
func FetchHolidays() (map[string][]Holiday, error) {
	resp, err := http.Get(govUKURL)
	if err != nil {
		log.Println("Error fetching data, falling back to backup")
		return LoadBackupData()
	}
	defer resp.Body.Close()

	var rawData map[string]struct {
		Events []map[string]interface{} `json:"events"`
	}
	err = json.NewDecoder(resp.Body).Decode(&rawData)
	if err != nil {
		log.Println("Error decoding data, falling back to backup")
		return LoadBackupData()
	}

	// Flatten and convert the structure into map[string][]Holiday
	holidays := make(map[string][]Holiday)
	for division, item := range rawData {
		var holidayList []Holiday
		for _, event := range item.Events {
			// Parse the date in "YYYY-MM-DD" format
			date, err := time.Parse("2006-01-02", event["date"].(string))
			if err != nil {
				log.Printf("Error parsing date for holiday: %v", err)
				continue
			}
			holidayList = append(holidayList, Holiday{
				Title:   event["title"].(string),
				Date:    date,
				Notes:   event["notes"].(string),
				Bunting: event["bunting"].(bool),
			})
		}
		holidays[division] = holidayList
	}

	return holidays, nil
}

// GetHolidays returns a list of holidays for a given division and year
func (b *BankHolidays) GetHolidays(division string, year int) []Holiday {
	holidays, found := b.Holidays[division]
	if !found {
		return nil
	}

	var result []Holiday
	for _, holiday := range holidays {
		if holiday.Date.Year() == year || year == 0 {
			result = append(result, holiday)
		}
	}

	return result
}

// IsHoliday checks if the given date is a bank holiday in the specified division
func (b *BankHolidays) IsHoliday(date time.Time, division string) bool {
	holidays := b.GetHolidays(division, date.Year())
	for _, holiday := range holidays {
		if holiday.Date.Equal(date) {
			return true
		}
	}
	return false
}

// IsWorkDay checks if the given date is a workday (not a weekend or holiday)
func (b *BankHolidays) IsWorkDay(date time.Time, division string) bool {
	if b.Weekend[date.Weekday()] {
		return false
	}
	return !b.IsHoliday(date, division)
}

// GetNextWorkDay returns the next available workday after the given date
func (b *BankHolidays) GetNextWorkDay(date time.Time, division string) time.Time {
	oneDay := time.Hour * 24
	for {
		date = date.Add(oneDay)
		if b.IsWorkDay(date, division) {
			return date
		}
	}
}

// GetNextHoliday returns the next holiday after the given date
func (b *BankHolidays) GetNextHoliday(date time.Time, division string) *Holiday {
	holidays := b.GetHolidays(division, date.Year())
	for _, holiday := range holidays {
		if holiday.Date.After(date) {
			return &holiday
		}
	}
	return nil
}

// GetPreviousHoliday returns the last holiday before the given date
func (b *BankHolidays) GetPreviousHoliday(date time.Time, division string) *Holiday {
	holidays := b.GetHolidays(division, date.Year())
	for i := len(holidays) - 1; i >= 0; i-- {
		if holidays[i].Date.Before(date) {
			return &holidays[i]
		}
	}
	return nil
}
