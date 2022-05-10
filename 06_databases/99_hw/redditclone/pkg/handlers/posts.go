package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"redditclone/pkg/comments"
	"redditclone/pkg/errorsForProject"
	"redditclone/pkg/forms"
	"redditclone/pkg/posts"
	"redditclone/pkg/posts/repo"
	"redditclone/pkg/session"
	"strconv"
	"strings"
	"time"
)

type PostsHandler struct {
	Logger         *zap.SugaredLogger
	PostRepo       repo.MyRepo
	SessionManager session.SessionRepo
}

// ВСЕ ГЕТТЕРЫ

func (h *PostsHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	res, err := h.PostRepo.GetAll()
	if err != nil {
		JsonError(w, http.StatusBadRequest, "GetALLPost: "+err.Error(), h.Logger)
		return
	}
	SendSliceRequest(w, "Get AllPosts: ", res, http.StatusOK, h.Logger)

	h.Logger.Infof("Get AllPosts: %v", http.StatusOK)
	return
}

func (h *PostsHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := h.PostRepo.GetPostsCategory(vars["CATEGORY_NAME"])
	if err != nil {
		JsonError(w, http.StatusBadRequest, "GetCategoryPost: "+err.Error(), h.Logger)
		return
	}
	SendSliceRequest(w, "Get Category: ", res, http.StatusOK, h.Logger)

	h.Logger.Infof("Get Category: %v", http.StatusOK)
	return
}

func (h *PostsHandler) GetUserPost(w http.ResponseWriter, r *http.Request) {
	ses, err := h.SessionManager.Check(r)
	vars := mux.Vars(r)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "GetUserPost: "+err.Error(), h.Logger)
		return
	}
	userForm := forms.UserForm{
		ID:    strconv.Itoa(int(ses.UserID)),
		Login: vars["USER_LOGIN"],
	}
	res, err := h.PostRepo.GetPostsByUser(userForm)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "GetUserPost: "+err.Error(), h.Logger)
		return
	}
	SendSliceRequest(w, "Get UserPost: ", res, http.StatusOK, h.Logger)

	h.Logger.Infof("Get UserPost: %v", http.StatusOK)
	return
}

func (h *PostsHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, err := h.PostRepo.GetByID(vars["POST_ID"])
	if err != nil {
		JsonError(w, http.StatusBadRequest, "GetPost: "+posts.ErrNoPost.Error(), h.Logger)
		return
	}
	h.PostRepo.IncreaseViews(post)

	_, err = h.PostRepo.Update(post)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "AddComment update: "+err.Error(), h.Logger)
		return
	}

	SendRequest(w, "GetPost: ", post, http.StatusOK, h.Logger)

	h.Logger.Infof("GetPosts: %v", post.ID)
	return
}

// ДОБАВЛЕНИЕ ИЛИ УДАЛЕНИЕ ПОСТА

func (h *PostsHandler) Add(w http.ResponseWriter, r *http.Request) {
	var newPost *posts.Post

	fd := &forms.PostFirstData{}
	err := json.NewDecoder(r.Body).Decode(&fd)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "Add: Cant Decode", h.Logger)
		return
	}

	userForm, _, errForm := GetUserForm(w, r, h.Logger) // получение юзера и времени, ошибка отправляется прям там
	if errForm != nil {
		return
	}
	newPost = &posts.Post{
		Category:    fd.Category,
		Title:       fd.Title,
		CreatedBy:   userForm,
		Type:        fd.TypeOfPost,
		CurrentTime: "2",
	}

	if fd.TypeOfPost == "text" {
		newPost.Text = fd.Text
	} else {
		newPost.URL = fd.Url
	}

	_, err = strconv.Atoi(userForm.ID)
	if err != nil {
		h.Logger.Infof("cant atoi")
		return
	}

	err = h.PostRepo.Add(newPost)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "AddPost: "+err.Error(), h.Logger)
		return
	}
	SendRequest(w, "Add post: ", newPost, http.StatusCreated, h.Logger)

	h.Logger.Infof("Added post: %v", newPost.ID)
	return
}

func (h *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isSuccess := h.PostRepo.Delete(vars["POST_ID"])
	if !isSuccess {
		JsonError(w, http.StatusBadRequest, "DELETE: "+errorsForProject.ErrCantDelete.Error(), h.Logger)
		return
	}

	resp, err := json.Marshal(map[string]string{
		"message": "success",
	})
	if err != nil {
		JsonError(w, http.StatusBadRequest, "DELETE: "+errorsForProject.ErrCantMarshal.Error(), h.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	h.Logger.Infof("Deleted Post: %v", vars["POST_ID"])
	return

}

//  ДОБАВЛЕНИЕ ИЛИ УДАЛЕНИЕ КОММЕНТАРИЯ

func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, err := h.PostRepo.GetByID(vars["POST_ID"])
	if err != nil {
		JsonError(w, http.StatusBadRequest, "AddComment: "+posts.ErrNoPost.Error(), h.Logger)
		return
	}
	fd := &forms.CommentForm{}
	err = json.NewDecoder(r.Body).Decode(&fd)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "AddComment: Cant Decode", h.Logger)
		return
	}

	if fd.Description == "" {
		resp, errMarshal := json.Marshal(map[string]interface{}{
			"errors": []errorsForProject.RegisterError{{
				Msg:      "is required",
				Location: "body",
				Param:    "comment",
			},
			}})
		if errMarshal != nil {
			JsonError(w, http.StatusBadRequest, "Register: "+errorsForProject.ErrCantMarshal.Error(), h.Logger)
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(resp)
		return
	}

	userForm, _, errForm := GetUserForm(w, r, h.Logger) // получение юзера и времени, ошибка отправляется прям там
	if errForm != nil {
		return
	}

	newComment := &comments.Comment{
		ID:          strconv.Itoa(post.ComCount + 1),
		CreatedBy:   userForm,
		Description: fd.Description,
		CurrentTime: "2",
	}
	post.ComCount++
	post.Comments = append(post.Comments, *newComment)
	post1, err := h.PostRepo.Update(post)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "AddComment update: "+err.Error(), h.Logger)
		return
	}

	SendRequest(w, "Add comment: ", post1, http.StatusCreated, h.Logger)
	h.Logger.Infof("added comment: %v with PostId: %v", newComment.ID, vars["POST_ID"])
	return
}

func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, err := h.PostRepo.GetByID(vars["POST_ID"])
	if err != nil {
		JsonError(w, http.StatusBadRequest, "AddComment: "+posts.ErrNoPost.Error(), h.Logger)
		return
	}
	pos, IsSuccess := comments.Delete(post.Comments, vars["COMMENT_ID"])
	if !IsSuccess {
		JsonError(w, http.StatusBadRequest, "Delete Comment: "+errorsForProject.ErrCantDelete.Error(), h.Logger)
		return
	}
	post.Comments = pos
	post1, err := h.PostRepo.Update(post)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "Delete comment update: "+err.Error(), h.Logger)
		return
	}
	SendRequest(w, "Delete Comment: ", post1, http.StatusOK, h.Logger)

	h.Logger.Infof("Deleted comment: %v with PostID: %v", vars["COMMENT_ID"], vars["POST_ID"])
	return
}

// ФУНКЦИИ С VOTE

func (h *PostsHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	fd, post, userId, postId, err1 := GetParamForVote(w, r, h, 1)
	if err1 != nil {
		JsonError(w, http.StatusBadRequest, "Delete comment update: "+err1.Error(), h.Logger)
		return
	}
	resp := h.PostRepo.IncreaseVote(fd, post)
	post1, err := h.PostRepo.Update(resp)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "Delete comment update: "+err.Error(), h.Logger)
		return
	}

	SendRequest(w, "Upvote: ", post1, http.StatusOK, h.Logger)

	h.Logger.Infof("Upvote with UserId: %v with PostID: %v", userId, postId)
	return
}

func (h *PostsHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	fd, post, userId, postId, err1 := GetParamForVote(w, r, h, -1)
	if err1 != nil {
		JsonError(w, http.StatusBadRequest, "Delete comment update: "+err1.Error(), h.Logger)
		return
	}
	resp := h.PostRepo.IncreaseVote(fd, post)
	_, err := h.PostRepo.Update(resp)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "Delete comment update: "+err.Error(), h.Logger)
		return
	}

	SendRequest(w, "Downvote: ", resp, http.StatusOK, h.Logger)

	h.Logger.Infof("Downvote with UserId: %v with PostID: %v", userId, postId)
	return
}

func (h *PostsHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	fd, post, userId, postId, err1 := GetParamForVote(w, r, h, 0)
	if err1 != nil {
		JsonError(w, http.StatusBadRequest, "Delete comment update: "+err1.Error(), h.Logger)
		return
	}

	resp := h.PostRepo.DecreaseVote(fd, post)
	_, err := h.PostRepo.Update(post)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "Delete comment update: "+err.Error(), h.Logger)
		return
	}

	SendRequest(w, "Unvote: ", resp, http.StatusOK, h.Logger)

	h.Logger.Infof("Unvote with UserId: %v with PostID: %v", userId, postId)
	return
}

// СТОРОННИЕ ФУНКЦИИ

func GetParamForVote(w http.ResponseWriter, r *http.Request, h *PostsHandler, vote int) (fd *forms.VoteForm, post *posts.Post, UserId string, postId string, errAns error) {
	vars := mux.Vars(r)
	post, err := h.PostRepo.GetByID(vars["POST_ID"])
	if err != nil {
		JsonError(w, http.StatusBadRequest, "AddComment: "+posts.ErrNoPost.Error(), h.Logger)
		return nil, nil, "", "", fmt.Errorf("no user")
	}
	userForm, _, errForm := GetUserForm(w, r, h.Logger) // получение юзера и времени, ошибка отправляется прям там
	if errForm != nil {
		return nil, nil, "", "", fmt.Errorf("no user")
	}
	fd = &forms.VoteForm{
		ID: userForm.ID,
	}
	switch vote {
	case 1:
		fd.Vote = 1
	case -1:
		fd.Vote = -1
	default:
		fd.Vote = 0
	}

	return fd, post, userForm.ID, vars["POST_ID"], nil
}

func SendSliceRequest(w http.ResponseWriter, errStr string, res []*posts.Post, status int, Logger *zap.SugaredLogger) {
	resp, err := json.Marshal(res)
	if err != nil {
		JsonError(w, http.StatusBadRequest, errStr+errorsForProject.ErrCantMarshal.Error(), Logger)
		return
	}

	w.WriteHeader(status)
	w.Write(resp)
}

func SendRequest(w http.ResponseWriter, errStr string, post *posts.Post, status int, Logger *zap.SugaredLogger) {
	resp, err := json.Marshal(post)
	if err != nil {
		JsonError(w, http.StatusBadRequest, errStr+errorsForProject.ErrCantMarshal.Error(), Logger)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

func GetUserForm(w http.ResponseWriter, r *http.Request, Logger *zap.SugaredLogger) (forms.UserForm, []byte, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[strings.Index(tokenString, " ")+1:]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return TokenSecret, nil
	})
	if err != nil {
		JsonError(w, http.StatusBadRequest, "GetUserForm: cant token parse", Logger)
		return forms.UserForm{}, []byte{}, err
	}

	in := claims["user"]
	ourUser := in.(map[string]interface{}) // map с login и Id юзера из токена

	timing, errTime := time.Now().UTC().MarshalText()
	if errTime != nil {
		JsonError(w, http.StatusBadRequest, "Time.MarshalText: year outside of range [0,9999]", Logger)
		return forms.UserForm{}, []byte{}, err
	}
	res := forms.UserForm{
		ID:    fmt.Sprint(ourUser["id"]),
		Login: ourUser["username"].(string),
	}
	return res, timing, nil
}
