package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"dyelesho/forum/internal/models"
	"dyelesho/forum/internal/validator"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Posts         *models.Model
	TemplateCache map[string]*template.Template
	Users         *models.UserModel
	categories    []string
	Reactions     *models.ReactionModel
}

type CategoriesForm struct {
	Technology    string
	Entertainment string
	Travel        string
	Health        string
	validator.Validator
}

type CommentCreateForm struct {
	CContent string
	validator.Validator
}

type PostCreateForm struct {
	Title      string
	Content    string
	Category   string
	Categories []string
	validator.Validator
}

type ErrorStruct struct {
	Status int
	Text   string
}

type UserSignupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		category := r.URL.Query().Get("category")

		session, err := app.CheckSession(w, r)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}

		var posts []*models.Post

		switch category {
		case "created":
			if session == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			posts, err = app.Posts.GetPostsByUser(session.UserName)

		case "liked":
			if session == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			posts, err = app.Posts.GetPostsByUserReaction(session.UserID)

		case "latest":
			posts, err = app.Posts.Latest()

		default:
			posts, err = app.Posts.Latest()
		}

		if err != nil {
			app.ServerError(w, err, r)
			return
		}

		data := app.NewTemplateData(r)
		data.Posts = posts
		data.IsAuthenticated = session != nil

		app.Render(w, http.StatusOK, "home.html", data, r)
	case http.MethodPost:
		session, err := app.CheckSession(w, r)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}

		var posts []*models.Post
		r.ParseForm()

		form := &CategoriesForm{
			Technology:    r.FormValue("Technology"),
			Entertainment: r.FormValue("Entertainment"),
			Travel:        r.FormValue("Travel"),
			Health:        r.FormValue("Health"),
		}

		posts, err = app.Posts.Latest()
		if err != nil {
			app.ServerError(w, err, r)
			return
		}

		var filteredPosts []*models.Post

		for i := range posts {
			switch {
			case form.Technology != "" && posts[i].Category != "" && strings.Contains(posts[i].Category, form.Technology):
				filteredPosts = append(filteredPosts, posts[i])
			case form.Travel != "" && posts[i].Category != "" && strings.Contains(posts[i].Category, form.Travel):
				filteredPosts = append(filteredPosts, posts[i])
			case form.Health != "" && posts[i].Category != "" && strings.Contains(posts[i].Category, form.Health):
				filteredPosts = append(filteredPosts, posts[i])
			case form.Entertainment != "" && posts[i].Category != "" && strings.Contains(posts[i].Category, form.Entertainment):
				filteredPosts = append(filteredPosts, posts[i])
			}
		}

		if filteredPosts == nil {
			data := app.NewTemplateData(r)
			data.Posts = []*models.Post{} 
			data.IsAuthenticated = session != nil

			app.Render(w, http.StatusOK, "home.html", data, r)
			return
		}

		data := app.NewTemplateData(r)
		data.Posts = filteredPosts
		data.IsAuthenticated = session != nil

		app.Render(w, http.StatusOK, "home.html", data, r)
	}
}

func (app *Application) PostView(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/post/view/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		app.NotFound(w, r)
		return
	}

	post, err := app.Posts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w, r)
		} else {
			app.ServerError(w, err, r)
		}
		return
	}

	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	comments, err := app.Posts.GetComments(post.ID)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	data := app.NewTemplateData(r)
	data.Post = post
	data.Comments = comments

	if session != nil {
		data.Post.IsAuthenticated = true
		for i := range data.Comments {
			data.Comments[i].IsAuthenticated = true
		}
	}

	app.Render(w, http.StatusOK, "view.html", data, r)
}

func (app *Application) CreateComment(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/post/view/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		app.NotFound(w, r)
		return
	}

	comment := r.FormValue("comment")

	if comment == "" || strings.TrimSpace(comment) == "" || utf8.RuneCountInString(comment) > 300 || countLines(comment) > 15 {
		post, err := app.Posts.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(w, r)
			} else {
				app.ServerError(w, err, r)
			}
			return
		}

		comments, err := app.Posts.GetComments(post.ID)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}

		data := &TemplateData{
			CurrentYear:     time.Now().Year(),
			Post:            post,
			IsAuthenticated: true,
			Comments:        comments,
			CommentError:    true,
		}

		app.Render(w, http.StatusOK, "view.html", data, r)
		return
	}

	comment = strings.Replace(comment, "\n", "<br>", -1)
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentInput := models.Comment{
		Author:   session.UserName,
		CContent: comment,
		PostID:   id,
	}

	err = app.Posts.PostComment(commentInput)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func countLines(comment string) int {
	lines := strings.Split(comment, "\n")
	return len(lines)
}

func (app *Application) PostCreate(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Categories = app.categories
	data.Form = PostCreateForm{}
	app.Render(w, http.StatusOK, "create.html", data, r)
}

func (app *Application) PostCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, r)
		return
	}

	catForm := &CategoriesForm{
		Technology:    r.FormValue("Technology"),
		Entertainment: r.FormValue("Entertainment"),
		Travel:        r.FormValue("Travel"),
		Health:        r.FormValue("Health"),
	}

	form := &PostCreateForm{
		Title:    r.PostForm.Get("title"),
		Content:  r.PostForm.Get("content"),
		Category: getCats(catForm),
	}
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.CheckFormValue(form.Category), "cats", "At least one category should be checked")
	if !form.Valid() || !catForm.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		data.Categories = app.categories
		app.Render(w, http.StatusUnprocessableEntity, "create.html", data, r)
		return
	}

	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userName := session.UserName
	id, err := app.Posts.Insert(form.Title, form.Content, form.Category, userName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func getCats(form *CategoriesForm) string {
	var cats string
	if form.Entertainment != "" {
		cats += form.Entertainment + " "
	}
	if form.Health != "" {
		cats += form.Health + " "
	}
	if form.Technology != "" {
		cats += form.Technology + " "
	}
	if form.Travel != "" {
		cats += form.Travel + " "
	}

	return cats
}

func (app *Application) UserSignup(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = UserSignupForm{}
	app.Render(w, http.StatusOK, "signup.html", data, r)
}

func (app *Application) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, r)
		return
	}
	form := UserSignupForm{
		Name:     strings.ToLower(r.PostForm.Get("name")),
		Email:    strings.ToLower(r.PostForm.Get("email")),
		Password: r.PostForm.Get("password"),
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.ValidUsername(form.Name), "name", "Invalid username format")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, "signup.html", data, r)
		return
	}
	err = app.Users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEntry) {
			form.AddFieldError("name", "Username or email addres is already in use")
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "signup.html", data, r)
		} else {
			app.ServerError(w, err, r)
		}
		return
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) UserLogin(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = UserLoginForm{}
	app.Render(w, http.StatusOK, "login.html", data, r)
}

func (app *Application) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, r)
		return
	}
	form := UserLoginForm{
		Email:    strings.ToLower(r.PostForm.Get("email")),
		Password: r.PostForm.Get("password"),
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, "login.html", data, r)
		return
	}
	id, err := app.Users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "login.html", data, r)
		} else {
			app.ServerError(w, err, r)
		}
		return
	}

	userName, err := app.Users.GetUserNameByEmail(form.Email)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	token, expiration, err := app.Posts.CreateSession(id, userName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) UserLogout(w http.ResponseWriter, r *http.Request) {
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	if session != nil {
		err = app.Posts.DeleteSessionByUserId(session.UserID)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
	}
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().AddDate(-1, 0, 0),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}
	token := cookie.Value

	session, err := app.Posts.GetSessionFromToken(token)
	if err != nil {
		return false
	}

	if session.ExpirationDate.Before(time.Now()) {
		return false
	}

	return true
}

func (app *Application) ErrorHandler(w http.ResponseWriter, errorNum int, r *http.Request) {
	data := app.NewTemplateData(r)
	Res := &ErrorStruct{
		Status: errorNum,
		Text:   http.StatusText(errorNum),
	}
	data.ErrorStruct = Res
	err := app.renderErr(w, http.StatusUnprocessableEntity, "error.html", data, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}
