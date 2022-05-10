package handlers

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http/httptest"
	"redditclone/pkg/comments"
	"redditclone/pkg/forms"
	"redditclone/pkg/posts"
	"redditclone/pkg/posts/mocks"
	"redditclone/pkg/posts/repo"
	"redditclone/pkg/session"
	"strconv"
	"strings"
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
			CurrentTime:       "2022-05-10T13:45:57.6132698Z",
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

func TestGetUserPost(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	expectedPosts, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	w := httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(1), "/api/user/1").Return(session.NewSession(uint32(2)), nil)

	ses.Create(w, uint32(1), "/api/user/1")
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	req1 := httptest.NewRequest("GET", `/api/user/1`, nil)
	a := mux.SetURLVars(req1, map[string]string{
		"USER_LOGIN": "ata",
	})
	userForm := forms.UserForm{
		ID:    "1",
		Login: "ata",
	}

	ses.EXPECT().Check(a).Return(session.NewSession(uint32(1)), nil)
	dBase.Db.(*mocks.PostRepo).On("GetPostsByUser", userForm).Return(expectedPosts, nil)
	w1 := httptest.NewRecorder()
	service.GetUserPost(w1, a)

	resp := w1.Result()
	if resp.StatusCode != 200 {
		t.Errorf("no text found")
		return
	}
}

func TestBadGetUserPost(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	w := httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(1), "/api/user/1").Return(session.NewSession(uint32(2)), nil)

	ses.Create(w, uint32(1), "/api/user/1")
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	req1 := httptest.NewRequest("GET", `/api/user/1`, nil)
	a := mux.SetURLVars(req1, map[string]string{
		"USER_LOGIN": "ata",
	})
	userForm := forms.UserForm{
		ID:    "1",
		Login: "ata",
	}

	ses.EXPECT().Check(a).Return(nil, session.ErrNoAuth)
	ses.EXPECT().Check(a).Return(session.NewSession(uint32(1)), nil)
	dBase.Db.(*mocks.PostRepo).On("GetPostsByUser", userForm).Return(nil, fmt.Errorf("no user"))
	w1 := httptest.NewRecorder()
	service.GetUserPost(w1, a)
	service.GetUserPost(w1, a)

	resp := w1.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestAddText(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	w := httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(1), "/api/user/1").Return(session.NewSession(uint32(2)), nil)

	ses.Create(w, uint32(1), "/api/user/1")
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	req1 := httptest.NewRequest("POST", `/api/post/1`, strings.NewReader(`{"category": "music", "text": "privet", "title": "qwe", "type": "text"}`))
	req1.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	userForm := forms.UserForm{
		ID:    "2",
		Login: "ayta",
	}
	p.Score = 1
	p.Votes = []*forms.VoteForm{
		{
			ID:   "2",
			Vote: 1,
		},
	}
	p.CreatedBy = userForm
	p.CurrentTime = "2"
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	w1 := httptest.NewRecorder()
	service.Add(w1, req1)

	resp := w1.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestAddUrl(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	w := httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(1), "/api/user/1").Return(session.NewSession(uint32(2)), nil)

	ses.Create(w, uint32(1), "/api/user/1")
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	userForm := forms.UserForm{
		ID:    "2",
		Login: "ayta",
	}
	p.Score = 1
	p.Votes = []*forms.VoteForm{
		{
			ID:   "2",
			Vote: 1,
		},
	}
	p.CreatedBy = userForm
	p.CurrentTime = "2"
	p.Text = ""
	p.Type = "url"
	p.URL = "privet"
	req2 := httptest.NewRequest("POST", `/api/post/1`, strings.NewReader(`{"category": "music", "url": "privet", "title": "qwe", "type": "url"}`))
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	w1 := httptest.NewRecorder()
	service.Add(w1, req2)

	resp := w1.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestBadAdd(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)

	userForm := forms.UserForm{
		ID:    "2",
		Login: "ayta",
	}
	p.Score = 1
	p.Votes = []*forms.VoteForm{
		{
			ID:   "2",
			Vote: 1,
		},
	}
	p.CreatedBy = userForm
	p.CurrentTime = "2"
	p.Text = ""
	p.Type = "url"
	p.URL = "privet"
	req2 := httptest.NewRequest("POST", `/api/post/1`, strings.NewReader(`{"category": "music", "url": "privet", "title": "qwe", "type": "url"}`))
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")

	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(fmt.Errorf("no user"))
	w1 := httptest.NewRecorder()
	service.Add(w1, req2)

	resp := w1.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

	req2 = httptest.NewRequest("POST", `/api/post/1`, strings.NewReader(`{category: "music", "url": "privet", "title": "qwe", "type": "url"}`))
	service.Add(w1, req2)

}

func TestDelete(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	ansP := p
	ansP.Views++
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	dBase.Db.(*mocks.PostRepo).On("Delete", "1").Return(true)

	req2 := httptest.NewRequest("DELETE", `/api/post/1`, nil)
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	w2 := httptest.NewRecorder()
	service.Delete(w2, b)

	resp := w2.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestBadDelete(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	ansP := p
	ansP.Views++
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	dBase.Db.(*mocks.PostRepo).On("Delete", "2").Return(false)

	req2 := httptest.NewRequest("DELETE", `/api/post/2`, nil)
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "2",
	})
	w2 := httptest.NewRecorder()
	service.Delete(w2, b)

	resp := w2.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestIncreaseVote(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	ansP := p
	ansP.Views++
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	fd := &forms.VoteForm{
		ID:   "1",
		Vote: 1,
	}

	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	dBase.Db.(*mocks.PostRepo).On("IncreaseVote", fd, p).Return(p)
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(ansP, nil)

	req2 := httptest.NewRequest("GET", `/api/post/1/upvote`, nil)
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	w2 := httptest.NewRecorder()
	service.Upvote(w2, b)

	resp := w2.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestBadIncreaseVote(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	ansP := p
	ansP.Views++
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	fd := &forms.VoteForm{
		ID:   "1",
		Vote: 1,
	}

	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	dBase.Db.(*mocks.PostRepo).On("IncreaseVote", fd, p).Return(p)
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(nil, fmt.Errorf("no user"))
	req2 := httptest.NewRequest("GET", `/api/post/1/upvote`, nil)
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	w2 := httptest.NewRecorder()
	service.Upvote(w2, b)

	resp := w2.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

	b = mux.SetURLVars(req2, map[string]string{
		"POST_ID": "2",
	})
	w2 = httptest.NewRecorder()
	dBase.Db.(*mocks.PostRepo).On("GetByID", "2").Return(nil, fmt.Errorf("no user"))
	service.Upvote(w2, b)

	resp = w2.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	title = "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestGet(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	service := &PostsHandler{
		PostRepo: dBase,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}
	_, p := GetPost()
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	dBase.Db.(*mocks.PostRepo).On("GetAll").Return(nil, fmt.Errorf("no user"))
	dBase.Db.(*mocks.PostRepo).On("GetByID", "4").Return(nil, fmt.Errorf("no user"))

	dBase.Db.(*mocks.PostRepo).On("GetPostsCategory", "music").Return(nil, fmt.Errorf("no user"))

	req := httptest.NewRequest("GET", "/api/posts/", nil)
	w := httptest.NewRecorder()
	service.GetAllPosts(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

	req1 := httptest.NewRequest("GET", `/api/posts/music`, nil)
	a := mux.SetURLVars(req1, map[string]string{
		"CATEGORY_NAME": "music",
	})
	w1 := httptest.NewRecorder()
	service.GetCategory(w1, a)

	resp = w1.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	title = "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

	req2 := httptest.NewRequest("GET", `/api/post/1`, nil)
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "4",
	})
	w2 := httptest.NewRecorder()
	service.GetPost(w2, b)

	resp = w2.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	title = "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	dBase.Db.(*mocks.PostRepo).On("IncreaseViews", p).Return(nil)
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(nil, fmt.Errorf("no user"))

	req3 := httptest.NewRequest("GET", `/api/post/1`, nil)
	c := mux.SetURLVars(req3, map[string]string{
		"POST_ID": "1",
	})
	w3 := httptest.NewRecorder()
	service.GetPost(w3, c)

	resp = w3.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	title = "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestDecreaseVote(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	ansP := p
	ansP.Views++
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	fd := &forms.VoteForm{
		ID:   "1",
		Vote: 1,
	}

	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	dBase.Db.(*mocks.PostRepo).On("IncreaseVote", fd, p).Return(p)
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(ansP, nil)

	req2 := httptest.NewRequest("GET", `/api/post/1/upvote`, nil)
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	w2 := httptest.NewRecorder()
	service.Downvote(w2, b)

	resp := w2.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestBadDecreaseVote(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	ansP := p
	ansP.Views++
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	fd := &forms.VoteForm{
		ID:   "1",
		Vote: 1,
	}

	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	dBase.Db.(*mocks.PostRepo).On("IncreaseVote", fd, p).Return(p)
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(nil, fmt.Errorf("no user"))
	req2 := httptest.NewRequest("GET", `/api/post/1/upvote`, nil)
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	w2 := httptest.NewRecorder()
	service.Downvote(w2, b)

	resp := w2.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

	b = mux.SetURLVars(req2, map[string]string{
		"POST_ID": "2",
	})
	w2 = httptest.NewRecorder()
	dBase.Db.(*mocks.PostRepo).On("GetByID", "2").Return(nil, fmt.Errorf("no user"))
	service.Downvote(w2, b)

	resp = w2.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	title = "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestUnVote(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	ansP := p
	ansP.Views++
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	fd := &forms.VoteForm{
		ID:   "1",
		Vote: 1,
	}

	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	dBase.Db.(*mocks.PostRepo).On("DecreaseVote", fd, p).Return(p)
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(ansP, nil)

	req2 := httptest.NewRequest("GET", `/api/post/1/upvote`, nil)
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	w2 := httptest.NewRecorder()
	service.Unvote(w2, b)

	resp := w2.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestBadUnVote(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	ansP := p
	ansP.Views++
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)
	fd := &forms.VoteForm{
		ID:   "1",
		Vote: 1,
	}

	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	dBase.Db.(*mocks.PostRepo).On("DecreaseVote", fd, p).Return(p)
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(nil, fmt.Errorf("no user"))
	req2 := httptest.NewRequest("GET", `/api/post/1/upvote`, nil)
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	w2 := httptest.NewRecorder()
	service.Unvote(w2, b)

	resp := w2.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

	b = mux.SetURLVars(req2, map[string]string{
		"POST_ID": "2",
	})
	w2 = httptest.NewRecorder()
	dBase.Db.(*mocks.PostRepo).On("GetByID", "2").Return(nil, fmt.Errorf("no user"))
	service.Unvote(w2, b)

	resp = w2.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	title = "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestAddComment(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	w := httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(1), "/api/user/1").Return(session.NewSession(uint32(2)), nil)

	ses.Create(w, uint32(1), "/api/user/1")
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	userForm := forms.UserForm{
		ID:    "2",
		Login: "ayta",
	}
	p.Score = 1
	p.Votes = []*forms.VoteForm{
		{
			ID:   "2",
			Vote: 1,
		},
	}
	p.CreatedBy = userForm
	p.CurrentTime = "2"
	p.Text = ""
	p.Type = "url"
	p.URL = "privet"
	req2 := httptest.NewRequest("POST", `/api/post/1`, strings.NewReader(`{"comment": "zxc"}`))
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	p.Comments = append(p.Comments, comments.Comment{
		ID:          strconv.Itoa(p.ComCount + 1),
		CreatedBy:   userForm,
		Description: "zxc",
		CurrentTime: "2",
	})
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(p, nil)
	w1 := httptest.NewRecorder()
	service.AddComment(w1, b)

	resp := w1.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestGetAllCategoryId(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	expectedPosts, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	ansP := p
	ansP.Views++
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	dBase.Db.(*mocks.PostRepo).On("GetAll").Return(expectedPosts, nil)
	dBase.Db.(*mocks.PostRepo).On("GetPostsCategory", "music").Return(expectedPosts, nil)

	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	dBase.Db.(*mocks.PostRepo).On("IncreaseViews", p).Return(nil)
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(ansP, nil)

	req := httptest.NewRequest("GET", "/api/posts/", nil)
	w := httptest.NewRecorder()
	service.GetAllPosts(w, req)

	req1 := httptest.NewRequest("GET", `/api/posts/music`, nil)
	a := mux.SetURLVars(req1, map[string]string{
		"CATEGORY_NAME": "music",
	})
	w1 := httptest.NewRecorder()
	service.GetCategory(w1, a)

	req2 := httptest.NewRequest("GET", `/api/post/1`, nil)
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	w2 := httptest.NewRecorder()
	service.GetPost(w2, b)
}

func TestAddEmptyComment(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	w := httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(1), "/api/user/1").Return(session.NewSession(uint32(2)), nil)

	ses.Create(w, uint32(1), "/api/user/1")
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	userForm := forms.UserForm{
		ID:    "2",
		Login: "ayta",
	}
	p.Score = 1
	p.Votes = []*forms.VoteForm{
		{
			ID:   "2",
			Vote: 1,
		},
	}
	p.CreatedBy = userForm
	p.CurrentTime = "2"
	p.Text = ""
	p.Type = "url"
	p.URL = "privet"
	req2 := httptest.NewRequest("POST", `/api/post/1`, strings.NewReader(`{"comment": ""}`))
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	p.Comments = append(p.Comments, comments.Comment{
		ID:          strconv.Itoa(p.ComCount + 1),
		CreatedBy:   userForm,
		Description: "zxc",
		CurrentTime: "2",
	})
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(p, nil)
	w1 := httptest.NewRecorder()
	service.AddComment(w1, b)

	resp := w1.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestBadAddComment(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	w := httptest.NewRecorder()
	ses.EXPECT().Create(w, uint32(1), "/api/user/1").Return(session.NewSession(uint32(2)), nil)

	ses.Create(w, uint32(1), "/api/user/1")
	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	userForm := forms.UserForm{
		ID:    "2",
		Login: "ayta",
	}
	p.Score = 1
	p.Votes = []*forms.VoteForm{
		{
			ID:   "2",
			Vote: 1,
		},
	}
	p.CreatedBy = userForm
	p.CurrentTime = "2"
	p.Text = ""
	p.Type = "url"
	p.URL = "privet"
	req2 := httptest.NewRequest("POST", `/api/post/1`, strings.NewReader(`{"commet: ""}`))
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b := mux.SetURLVars(req2, map[string]string{
		"POST_ID": "2",
	})
	dBase.Db.(*mocks.PostRepo).On("GetByID", "2").Return(nil, fmt.Errorf("no user"))
	p.Comments = append(p.Comments, comments.Comment{
		ID:          strconv.Itoa(p.ComCount + 1),
		CreatedBy:   userForm,
		Description: "zxc",
		CurrentTime: "2",
	})
	//dBase.Db.(*mocks.PostRepo).On("Update", p).Return(p, nil)
	w1 := httptest.NewRecorder()
	service.AddComment(w1, b)
	b = mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	service.AddComment(w1, b)

	req2 = httptest.NewRequest("POST", `/api/post/1`, strings.NewReader(`{"comment": "q"}`))
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	b = mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(nil, fmt.Errorf("no user"))
	service.AddComment(w1, b)

	resp := w1.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestDeleteComment(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}

	userForm := forms.UserForm{
		ID:    "2",
		Login: "ayta",
	}
	p.Comments = append(p.Comments, comments.Comment{
		ID:          strconv.Itoa(p.ComCount + 1),
		CreatedBy:   userForm,
		Description: "zxc",
		CurrentTime: "2",
		PostID:      "1",
	})

	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	req := httptest.NewRequest("DELETE", `/api/post/1/1`, nil)
	//req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	a := mux.SetURLVars(req, map[string]string{
		"POST_ID":    "1",
		"COMMENT_ID": "1",
	})
	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	//p.Comments = nil
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(p, nil)
	w1 := httptest.NewRecorder()
	service.DeleteComment(w1, a)

	resp := w1.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}

}
func TestBadDeleteComment(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}
	p.Comments = nil

	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	req := httptest.NewRequest("DELETE", `/api/post/1/1`, nil)
	//req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	a := mux.SetURLVars(req, map[string]string{
		"POST_ID":    "1",
		"COMMENT_ID": "1",
	})
	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	w1 := httptest.NewRecorder()
	service.DeleteComment(w1, a)

	dBase.Db.(*mocks.PostRepo).On("GetByID", "2").Return(nil, fmt.Errorf("no user"))
	a = mux.SetURLVars(req, map[string]string{
		"POST_ID":    "2",
		"COMMENT_ID": "1",
	})
	w1 = httptest.NewRecorder()
	service.DeleteComment(w1, a)

	resp := w1.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	title := "s"
	if !bytes.Contains(body, []byte(title)) {
		t.Errorf("no text found")
		return
	}
}

func TestBadUpdateDeleteComment(t *testing.T) {
	dBase := repo.InitMyRepoTest()

	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()
	_, p := GetPost()
	ses := session.NewMockSessionRepo(ctrl)
	service := &PostsHandler{
		PostRepo:       dBase,
		Logger:         zap.NewNop().Sugar(), // не пишет логи
		SessionManager: ses,
	}

	userForm := forms.UserForm{
		ID:    "2",
		Login: "ayta",
	}
	p.Comments = append(p.Comments, comments.Comment{
		ID:          strconv.Itoa(p.ComCount + 1),
		CreatedBy:   userForm,
		Description: "zxc",
		CurrentTime: "2",
		PostID:      "1",
	})

	dBase.Db.(*mocks.PostRepo).On("GetLastId").Return(1)
	dBase.Db.(*mocks.PostRepo).On("Add", p).Return(nil)

	req := httptest.NewRequest("DELETE", `/api/post/1/1`, nil)
	//req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTIyNzYzMTEsImlhdCI6MTY1MjE4OTkxMSwidXNlciI6eyJpZCI6IjIiLCJ1c2VybmFtZSI6ImF5dGEifX0.8thii0EYhBR_HuR11JcbFfVZpwok0fAeNK9iItUmO7o")
	a := mux.SetURLVars(req, map[string]string{
		"POST_ID":    "1",
		"COMMENT_ID": "1",
	})
	dBase.Db.(*mocks.PostRepo).On("GetByID", "1").Return(p, nil)
	//p.Comments = nil
	dBase.Db.(*mocks.PostRepo).On("Update", p).Return(nil, fmt.Errorf("no user"))
	w1 := httptest.NewRecorder()
	service.DeleteComment(w1, a)

}
