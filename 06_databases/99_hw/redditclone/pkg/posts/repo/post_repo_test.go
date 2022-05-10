package repo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"redditclone/pkg/comments"
	"redditclone/pkg/forms"
	"redditclone/pkg/posts"
	"redditclone/pkg/posts/mocks"
	"testing"
)

func GetPost() ([]*posts.Post, *posts.Post) {
	expectedPosts := []*posts.Post{
		{
			ID:                "1",
			Title:             "qwe",
			Text:              "privet",
			Category:          "music",
			CreatedBy:         forms.UserForm{ID: "1", Login: "ata"},
			Score:             0,
			UpVotedPercentage: 100,
			Type:              "text",
			CurrentTime:       "1",
			Votes:             make([]*forms.VoteForm, 0, 10),
			Comments:          []comments.Comment{},
		},
	}

	p := &posts.Post{
		ID:                "1",
		Title:             "qwe",
		Text:              "privet",
		Category:          "music",
		CreatedBy:         forms.UserForm{ID: "1", Login: "ata"},
		Score:             0,
		UpVotedPercentage: 100,
		Type:              "text",
		CurrentTime:       "1",
		Votes: []*forms.VoteForm{
			{
				ID:   "1",
				Vote: 1,
			},
			{
				ID:   "2",
				Vote: 1,
			},
			{
				ID:   "3",
				Vote: -1,
			},
			{
				ID:   "4",
				Vote: 1,
			},
		},
		Comments: []comments.Comment{},
	}
	return expectedPosts, p
}

func TestBadGetAll(t *testing.T) {
	dBase := InitMyRepoTest()
	//expectedPosts, p := GetPost()
	dBase.Db.(*mocks.PostRepo).On("GetAll").Return(nil, fmt.Errorf("no user"))
	_, err1 := dBase.GetAll()
	assert.Equal(t, err1, fmt.Errorf("no user"))
}

func TestBadGetCategory(t *testing.T) {
	dBase := InitMyRepoTest()
	//expectedPosts, p := GetPost()
	dBase.Db.(*mocks.PostRepo).On("GetPostsCategory", "q").Return(nil, fmt.Errorf("no user"))
	_, err1 := dBase.GetPostsCategory("q")
	assert.Equal(t, err1, fmt.Errorf("no user"))
}

func TestBadGetById(t *testing.T) {
	dBase := InitMyRepoTest()
	//expectedPosts, p := GetPost()
	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(nil, fmt.Errorf("no user"))
	_, err1 := dBase.GetByID("1")
	assert.Equal(t, err1, fmt.Errorf("no user"))
}

func TestBadUpdate(t *testing.T) {
	dBase := InitMyRepoTest()
	_, p := GetPost()
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(nil, fmt.Errorf("no user"))
	_, err1 := dBase.Update(p)
	assert.Equal(t, err1, fmt.Errorf("no user"))
}

func TestBadGetByUser(t *testing.T) {
	dBase := InitMyRepoTest()
	fd := forms.UserForm{ID: "1", Login: "ata"}
	dBase.Db.(*mocks.PostRepo).On("GetPostsByUser", fd).Return(nil, fmt.Errorf("no user"))
	_, err1 := dBase.GetPostsByUser(fd)
	assert.Equal(t, err1, fmt.Errorf("no user"))
}

func TestBadAdd(t *testing.T) {
	dBase := InitMyRepoTest()
	_, p := GetPost()
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(fmt.Errorf("no user"))
	err1 := dBase.Add(p)
	assert.Equal(t, err1, fmt.Errorf("no user"))
}

func TestGetAllAdd(t *testing.T) {
	dBase := InitMyRepoTest()
	expectedPosts, p := GetPost()

	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	dBase.Db.(*mocks.PostRepo).On("GetAll").Return(expectedPosts, nil)

	err := dBase.Add(p)
	assert.Empty(t, err)
	_, err1 := dBase.GetAll()
	assert.Empty(t, err1)
	dBase.Score(p)

	dBase.Db.(*mocks.PostRepo).On("Delete", "1").Return(true)
	dBase.Delete("1")
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(p, nil)
	dBase.Update(p)
}

func TestWithVotes(t *testing.T) {
	dBase := InitMyRepoTest()
	_, p := GetPost()

	ansP := p
	ansP.Votes[1].Vote = -1
	fd := forms.VoteForm{
		ID:   "2",
		Vote: -1,
	}
	fd2 := forms.VoteForm{
		ID:   "3",
		Vote: 0,
	}
	ans2 := p
	ans2.Votes[2].Vote = 0
	dBase.Db.(*mocks.PostRepo).On("Score", p).Return(nil)
	dBase.Db.(*mocks.PostRepo).On("IncreaseVote", fd, p).Return(ansP)
	dBase.Db.(*mocks.PostRepo).On("DecreaseVote", fd2, p).Return(ans2)
	dBase.Db.(*mocks.PostRepo).On("UpvotePercentage", p).Return(nil)
	dBase.Db.(*mocks.PostRepo).On("IncreaseViews", p).Return(nil)
	post := dBase.IncreaseVote(&fd, p)
	assert.Equal(t, ansP, post)
	post = dBase.DecreaseVote(&fd2, p)
	assert.Equal(t, ans2, post)
	post = dBase.DecreaseVote(&forms.VoteForm{
		ID:   "42",
		Vote: 0,
	}, p)
	assert.Equal(t, p, post)

	dBase.Score(p)
	dBase.UpvotePercentage(p)
	dBase.IncreaseViews(p)
	p.Votes = nil
	post = dBase.IncreaseVote(&fd, p)
	assert.Equal(t, p, post)
	post = dBase.DecreaseVote(&fd2, p)
	assert.Equal(t, p, post)
	dBase.Score(p)
	dBase.UpvotePercentage(p)
}

func TestPostByCategory(t *testing.T) {
	dBase := InitMyRepoTest()

	expectedPosts := []*posts.Post{
		{
			ID:                "1",
			Title:             "qwe",
			Text:              "privet",
			Category:          "music",
			CreatedBy:         forms.UserForm{ID: "1", Login: "ata"},
			Score:             0,
			UpVotedPercentage: 100,
			Type:              "text",
			CurrentTime:       "1",
			Votes:             make([]*forms.VoteForm, 0, 10),
			Comments:          []comments.Comment{},
		},
	}

	p := &posts.Post{
		ID:                "1",
		Title:             "qwe",
		Text:              "privet",
		Category:          "music",
		CreatedBy:         forms.UserForm{ID: "1", Login: "ata"},
		Score:             0,
		UpVotedPercentage: 100,
		Type:              "text",
		CurrentTime:       "1",
		Votes:             make([]*forms.VoteForm, 0, 10),
		Comments:          []comments.Comment{},
	}

	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)

	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	dBase.Db.(*mocks.PostRepo).On("GetPostsCategory", "music").Return(expectedPosts, nil)

	p2 := p
	p2.Category = "ewq"
	err := dBase.Add(p)
	assert.Empty(t, err)
	err = dBase.Add(p2)
	assert.Empty(t, err)
	ans, err := dBase.GetPostsCategory("music")
	assert.Equal(t, ans, expectedPosts)
}

func TestPostById(t *testing.T) {
	dBase := InitMyRepoTest()

	p := &posts.Post{
		ID:                "1",
		Title:             "qwe",
		Text:              "privet",
		Category:          "music",
		CreatedBy:         forms.UserForm{ID: "1", Login: "ata"},
		Score:             0,
		UpVotedPercentage: 100,
		Type:              "text",
		CurrentTime:       "1",
		Votes:             make([]*forms.VoteForm, 0, 10),
		Comments:          []comments.Comment{},
	}
	p2 := p
	p2.ID = "2"
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)

	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	dBase.Db.(*mocks.PostRepo).On("GetByID", "2").Return(p2, nil)

	err := dBase.Add(p)
	assert.Empty(t, err)
	err = dBase.Add(p2)
	assert.Empty(t, err)
	ans, err := dBase.GetByID("2")
	assert.Equal(t, ans, p2)
}

func TestPostByUser(t *testing.T) {
	dBase := InitMyRepoTest()

	expectedPosts := []*posts.Post{
		{
			ID:                "1",
			Title:             "qwe",
			Text:              "privet",
			Category:          "music",
			CreatedBy:         forms.UserForm{ID: "1", Login: "ata"},
			Score:             0,
			UpVotedPercentage: 100,
			Type:              "text",
			CurrentTime:       "1",
			Votes:             make([]*forms.VoteForm, 0, 10),
			Comments:          []comments.Comment{},
		},
	}

	p := &posts.Post{
		ID:                "1",
		Title:             "qwe",
		Text:              "privet",
		Category:          "music",
		CreatedBy:         forms.UserForm{ID: "1", Login: "ata"},
		Score:             0,
		UpVotedPercentage: 100,
		Type:              "text",
		CurrentTime:       "1",
		Votes:             make([]*forms.VoteForm, 0, 10),
		Comments:          []comments.Comment{},
	}
	fd := forms.UserForm{ID: "1", Login: "ata"}

	p2 := p
	p2.CreatedBy.ID = "2"
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)

	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	dBase.Db.(*mocks.PostRepo).On("GetPostsByUser", fd).Return(expectedPosts, nil)

	err := dBase.Add(p)
	assert.Empty(t, err)
	err = dBase.Add(p2)
	assert.Empty(t, err)
	ans, err := dBase.GetPostsByUser(fd)
	assert.Equal(t, ans, expectedPosts)
}
