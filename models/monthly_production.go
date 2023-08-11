package models

import "time"

type MonthlyProduction struct {
	Choke          float64   `json:"choke"`
	Co2Injection   float64   `json:"co2Injection"`
	CreatedAt      time.Time `json:"createdAt"`
	Date           time.Time `json:"date"`
	CustomNumber0  int       `json:"customNumber0"`
	CustomNumber1  int       `json:"customNumber1"`
	CustomNumber2  int       `json:"customNumber2"`
	CustomNumber3  int       `json:"customNumber3"`
	CustomNumber4  int       `json:"customNumber4"`
	DaysOn         int       `json:"daysOn"`
	Gas            float64   `json:"gas"`
	GasInjection   float64   `json:"gasInjection"`
	Ngl            float64   `json:"ngl"`
	Oil            float64   `json:"oil"`
	OperationalTag string    `json:"operationalTag"`
	SteamInjection float64   `json:"steamInjection"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Water          float64   `json:"water"`
	WaterInjection float64   `json:"waterInjection"`
	Well           string    `json:"well"`
}
