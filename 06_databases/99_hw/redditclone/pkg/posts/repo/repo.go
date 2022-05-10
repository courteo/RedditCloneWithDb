package repo

import (
	"fmt"
	"redditclone/pkg/comments"
	"redditclone/pkg/forms"
	"redditclone/pkg/posts"
	"redditclone/pkg/posts/mocks"
	"strconv"
)

type MyRepo struct {
	Db posts.PostRepo
}

func (d *MyRepo) GetAll() (res []*posts.Post, err error) {
	post1, err := d.Db.GetAll()
	if err != nil {
		return nil, fmt.Errorf("no user")
	}

	return post1, nil
}

func (d *MyRepo) GetPostsByUser(author forms.UserForm) (res []*posts.Post, err error) {
	post1, err := d.Db.GetPostsByUser(author)
	if err != nil {
		return nil, fmt.Errorf("no user")
	}

	return post1, nil
}

func (d *MyRepo) GetPostsCategory(category string) (res []*posts.Post, err error) {
	post1, err := d.Db.GetPostsCategory(category)
	if err != nil {
		return nil, fmt.Errorf("no user")
	}
	return post1, nil
}

func (d *MyRepo) GetByID(id string) (*posts.Post, error) {
	post1, err := d.Db.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("no user")
	}
	return post1, nil
}

func (d *MyRepo) Update(post *posts.Post) (*posts.Post, error) {
	res, err := d.Db.Update(post)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *MyRepo) Add(post *posts.Post) error {
	post.ID = strconv.Itoa(d.Db.GetLastId())
	post.Score = 1
	post.UpVotedPercentage = 100
	post.Views = 0
	post.Comments = []comments.Comment{}
	post.Votes = make([]*forms.VoteForm, 0, 10)
	post.Votes = append(post.Votes, &forms.VoteForm{Vote: 1, ID: post.CreatedBy.ID})

	err := d.Db.Add(post)
	if err != nil {
		return fmt.Errorf("no user")
	}
	return nil
}

func (d *MyRepo) Delete(id string) bool {
	res := d.Db.Delete(id)
	return res
}

func (d *MyRepo) IncreaseViews(newPost *posts.Post) {
	newPost.Views++
}

func (d *MyRepo) IncreaseVote(fd *forms.VoteForm, post *posts.Post) *posts.Post {
	allnum := len(post.Votes)
	if allnum == 0 {
		return post
	}

	i := -1
	for idx, vote := range post.Votes {
		if vote.ID != fd.ID {
			continue
		}
		i = idx
	}

	if i != -1 {
		if i < len(post.Votes)-1 {
			copy(post.Votes[i:], post.Votes[i+1:])
		}
		post.Votes[len(post.Votes)-1] = nil
		post.Votes = post.Votes[:len(post.Votes)-1]
	}

	post.Votes = append(post.Votes, fd)

	d.UpvotePercentage(post)
	d.Score(post)
	return post
}

func (repo *MyRepo) DecreaseVote(fd *forms.VoteForm, post *posts.Post) *posts.Post {
	allnum := len(post.Votes)
	if allnum == 0 {
		return post
	}

	i := -1
	for idx, vote := range post.Votes {
		if vote.ID != fd.ID {
			continue
		}
		i = idx
	}

	if i < 0 {
		return post
	}

	if i < len(post.Votes)-1 {
		copy(post.Votes[i:], post.Votes[i+1:])
	}
	post.Votes[len(post.Votes)-1] = nil
	post.Votes = post.Votes[:len(post.Votes)-1]
	repo.UpvotePercentage(post)
	repo.Score(post)
	return post
}

func (d *MyRepo) UpvotePercentage(post *posts.Post) {

	allnum := len(post.Votes)
	if allnum == 0 {
		post.UpVotedPercentage = 0
		return
	}

	var currNums int
	for _, vote := range post.Votes {
		if vote.Vote == 1 {
			currNums++
		}
	}
	post.UpVotedPercentage = uint32(currNums * 100 / allnum)
}

func (d *MyRepo) Score(post *posts.Post) {
	allnum := len(post.Votes)
	if allnum == 0 {
		post.Score = 0
		return
	}

	var currNums int
	for _, vote := range post.Votes {
		currNums += vote.Vote
	}
	post.Score = currNums
}

func InitMyRepoTest() MyRepo {
	return MyRepo{Db: &mocks.PostRepo{}}
}
