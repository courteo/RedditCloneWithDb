package comments

import (
	"fmt"
	"strconv"
	"sync"
)

type CommentMemoryRepository struct {
	lastID uint32
	data   []*Comment
	mu     *sync.RWMutex
}

func NewMemoryRepo() *CommentMemoryRepository {
	return &CommentMemoryRepository{
		data: make([]*Comment, 0, 10),
		mu:   &sync.RWMutex{},
	}
}

func (repo *CommentMemoryRepository) GetAllComments() []*Comment {
	return repo.data
}

func (repo *CommentMemoryRepository) Add(Comment *Comment) string {
	//repo.mu.RLock()
	repo.lastID++
	Comment.ID = strconv.Itoa(int(repo.lastID))
	repo.data = append(repo.data, Comment)
	//repo.mu.RUnlock()
	return Comment.ID
}

func Delete(data []Comment, id string) ([]Comment, bool) {
	i := -1
	for idx, comment := range data {
		if comment.ID != id {
			continue
		}
		i = idx
	}
	if i < 0 {
		return data, false
	}
	fmt.Println("i  -", i)
	if i < len(data)-1 {
		//repo.mu.RLock()
		copy(data[i:], data[i+1:])
		//repo.mu.RUnlock()
	}
	//repo.mu.RLock()
	data = data[:len(data)-1]
	//repo.mu.RUnlock()
	return data, true
}
