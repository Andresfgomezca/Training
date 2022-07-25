package models

import "time"

//"time"
type Employee struct{
	ID int
	Name string
	DOJ time.Time
	Skillset []string
}

