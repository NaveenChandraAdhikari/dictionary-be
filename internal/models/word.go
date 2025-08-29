package models

import "time"

type Word struct {
	ID   int    `json:"id,omitempty" db:"id,omitempty"`
	Word string `json:"word,omitempty" db:"word,omitempty"`
	//Meaning   string `json:"meaning,omitempty" db:"meaning,omitempty"`
	CreatedBy int `json:"created_by,omitempty" db:"created_by,omitempty"`

	Phonetic string   `json:"phonetic,omitempty" db:"phonetic,omitempty"`
	Origin   string   `json:"origin,omitempty" db:"origin,omitempty"`
	Meanings []string `json:"meanings,omitempty" db:"meanings,omitempty"`

	//OPTIONAL EXEC  ID
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}
