package notables

import "time"

type Notable struct {
    PetName   string
    ImageHash string
    Observed  time.Time
}