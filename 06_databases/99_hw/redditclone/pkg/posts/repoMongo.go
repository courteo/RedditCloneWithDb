package posts

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
	"log"
	"redditclone/pkg/forms"
	"sort"
)

var (
	ErrNoPost = errors.New(" No post found")
)

type PostMemoryRepository struct {
	data   *mongo.Collection
	LastId int
}

func NewMemoryRepo(collection *mongo.Collection) *PostMemoryRepository {
	return &PostMemoryRepository{
		data:   collection,
		LastId: 1,
	}
}

// ВСЕ ГЕТТЕРЫ

func (repo *PostMemoryRepository) GetAll() (res []*Post, err error) {
	post1 := []*Post{}

	cur, err := repo.data.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Post
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		post1 = append(post1, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	sort.SliceStable(post1, func(i, j int) bool {
		return post1[i].Score >= post1[i].Score
	})
	return post1, nil

}

func (repo *PostMemoryRepository) GetPostsByUser(author forms.UserForm) (res []*Post, err error) {
	post1 := []*Post{}

	cur, err := repo.data.Find(context.TODO(), bson.M{"author": author})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Post
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		post1 = append(post1, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())
	sort.SliceStable(post1, func(i, j int) bool {
		return post1[i].Score >= post1[i].Score
	})
	return post1, nil
}

func (repo *PostMemoryRepository) GetPostsCategory(category string) (res []*Post, err error) {
	post1 := []*Post{}

	cur, err := repo.data.Find(context.TODO(), bson.M{"category": category})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Post
		err := cur.Decode(&elem)
		if err != nil {
			return nil, fmt.Errorf("no user")
		}

		post1 = append(post1, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("no user")
	}

	// Close the cursor once finished
	cur.Close(context.TODO())
	sort.SliceStable(post1, func(i, j int) bool {
		return post1[i].Score >= post1[i].Score
	})
	return post1, nil
}

func (repo *PostMemoryRepository) GetByID(id string) (*Post, error) {
	posts := &Post{}

	pos := repo.data.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&posts)
	if pos != nil {
		return nil, fmt.Errorf("no user")
	}
	return posts, nil
}

// ДОБАВЛЕНИЕ И УДАЛЕНИЕ ПОСТА

func (repo *PostMemoryRepository) Update(post *Post) (*Post, error) {
	//posts := &posts.Post{}
	update := bson.D{{"$set", bson.M{"comments": post.Comments, "score": post.Score, "upvotePercentage": post.UpVotedPercentage}}}
	_, pos := repo.data.UpdateOne(context.TODO(), bson.M{"_id": post.ID}, update)
	if pos != nil {
		return nil, fmt.Errorf("no user")
	}
	return post, nil
}

func (repo *PostMemoryRepository) GetLastId() int {
	return repo.LastId
}

func (repo *PostMemoryRepository) Add(post *Post) error {

	//post.ID = strconv.Itoa(repo.GetLastId())
	//post.Score = 1
	//post.UpVotedPercentage = 100
	//post.Views = 0
	//post.Comments = []comments.Comment{}
	//post.Votes = make([]*forms.VoteForm, 0, 10)
	//post.Votes = append(post.Votes, &forms.VoteForm{Vote: 1, ID: post.CreatedBy.ID})
	repo.LastId++
	_, err := repo.data.InsertOne(context.TODO(), post)
	if err != nil {
		return fmt.Errorf("no user")
	}
	return nil
}

func (repo *PostMemoryRepository) Delete(id string) bool {
	_, err := repo.data.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err == mgo.ErrNotFound {
		return false
	}
	return true

}

// ДРУГИЕ

func (repo *PostMemoryRepository) IncreaseViews(newPost *Post) {
	newPost.Views++
}

func (repo *PostMemoryRepository) IncreaseVote(fd *forms.VoteForm, post *Post) *Post {
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

	repo.UpvotePercentage(post)
	repo.Score(post)
	return post
}

func (repo *PostMemoryRepository) DecreaseVote(fd *forms.VoteForm, post *Post) *Post {
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

func (repo *PostMemoryRepository) UpvotePercentage(post *Post) {

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

func (repo *PostMemoryRepository) Score(post *Post) {
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
