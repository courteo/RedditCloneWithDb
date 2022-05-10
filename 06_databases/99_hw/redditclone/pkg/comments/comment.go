package comments

import (
	"redditclone/pkg/forms"
)

type Comment struct {
	ID          string         `json:"id"`
	Description string         `json:"body"`
	CurrentTime string         `json:"created"`
	CreatedBy   forms.UserForm `json:"author"`
	PostID      string         `json:"postId"`
}
