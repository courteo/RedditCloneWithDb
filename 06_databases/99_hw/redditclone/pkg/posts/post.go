package posts

import (
	"redditclone/pkg/comments"
	"redditclone/pkg/forms"
)

type Post struct {
	ID                string             `json:"id" bson:"_id"`
	Title             string             `json:"title" bson:"title"`
	URL               string             `json:"url" bson:"url"`
	Text              string             `json:"text" bson:"text"`
	Category          string             `json:"category" bson:"category"`
	Comments          []comments.Comment `json:"comments" bson:"comments"`
	CreatedBy         forms.UserForm     `json:"author" bson:"author"`
	Score             int                `json:"score" bson:"score"`
	UpVotedPercentage uint32             `json:"upvotePercentage" bson:"upvotePercentage"`
	Views             uint32             `json:"views" bson:"views"`
	Type              string             `json:"type" bson:"type"`
	CurrentTime       string             `json:"created" bson:"created"`
	Votes             []*forms.VoteForm  `json:"votes" bson:"votes"`
	ComCount          int                `json:"count" bson:"count"`
}

type PostRepo interface {
	IncreaseViews(newPost *Post)
	GetPostsCategory(category string) ([]*Post, error)
	GetPostsByUser(author forms.UserForm) ([]*Post, error)
	GetAll() ([]*Post, error)
	GetByID(id string) (*Post, error)
	Add(post *Post) error
	Update(post *Post) (*Post, error)
	GetLastId() int
	IncreaseVote(fd *forms.VoteForm, post *Post) *Post
	DecreaseVote(fd *forms.VoteForm, post *Post) *Post
	Delete(id string) bool
}
