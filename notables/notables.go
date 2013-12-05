package notables

import "time"

type Notable struct {
	PetName   string    `gorethink:"petName" json:"petName"`
	ImageHash string    `gorethink:"imageHash" json:"imageHash"`
	Observed  time.Time `gorethink:"observed" json:"observed"`
}
