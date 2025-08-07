package entities

import (
	"time"
)

type Post struct {
	ID        string
	Title     string
	Content   string
	AuthorID  string
	Tags      []string
	ViewCount int
	Likes     []string
	Dislikes  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// content should collections of block
/*
	should content r
		* Section
		  - Heading
		  -

*/
