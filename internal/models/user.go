package models

import (
	"time"

	"github.com/google/uuid"
)

// user model represents user table fields in db
type User struct {
	ID 		  uuid.UUID `db:"id" json:"id"`
	Username  string	`db:"username" json:"username"`
	Email 	  string	`db:"email" json:"email"`
	Password  []byte	`db:"password" json:"password"`
	Role 	  string	`db:"role" json:"role"`
	CreatedAt time.Time	`db:"created_at" json:"created_at"`
}
