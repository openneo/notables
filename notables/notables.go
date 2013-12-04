package notables

import "time"

type Notable struct {
	PetName   string    `gorethink:"petName"`
	ImageHash string    `gorethink:"imageHash"`
	Observed  time.Time `gorethink:"observed"`
}
