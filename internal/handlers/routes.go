package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	app.categories = []string{"Technology", "Travel", "Health", "Entertainment"}
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			remainingPath := r.URL.Path[len("/static/"):]
			r.URL.Path = "/" + remainingPath
			fileServer.ServeHTTP(w, r)
		} else {
			MethodNotAllowedHandler(w, r, []string{http.MethodGet})
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodPost {
			app.Home(w, r)
		} else {
			MethodNotAllowedHandler(w, r, []string{http.MethodGet})
		}
	})

	mux.HandleFunc("/post/view/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			app.PostView(w, r)
		} else if r.Method == http.MethodPost {
			app.RequireAuthentication(http.HandlerFunc(app.CreateComment)).ServeHTTP(w, r)
		} else {
			MethodNotAllowedHandler(w, r, []string{http.MethodGet, http.MethodPost})
		}
	})
	mux.Handle("/post/create", app.RequireAuthentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.PostCreate(w, r)
		case http.MethodPost:
			app.PostCreatePost(w, r)
		default:
			MethodNotAllowedHandler(w, r, []string{http.MethodGet, http.MethodPost})
		}
	})))

	mux.Handle("/likePost", app.RequireAuthentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.CheckSession(w, r)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		userID := session.UserID
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.NotFound(w, r)
			return
		}

		_ = app.Reactions.LikePost(userID, id)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	})))
	mux.Handle("/dislikePost", app.RequireAuthentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.CheckSession(w, r)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		userID := session.UserID
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.NotFound(w, r)
			return
		}

		_ = app.Reactions.DislikePost(userID, id)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	})))
	mux.Handle("/likeComment", app.RequireAuthentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.CheckSession(w, r)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		userID := session.UserID
		commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || commentID < 1 {
			app.NotFound(w, r)
			return
		}
		_ = app.Reactions.LikeComment(userID, commentID)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	})))
	mux.Handle("/dislikeComment", app.RequireAuthentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.CheckSession(w, r)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		userID := session.UserID
		commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || commentID < 1 {
			app.NotFound(w, r)
			return
		}
		_ = app.Reactions.DislikeComment(userID, commentID)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	})))
	mux.HandleFunc("/user/signup", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.UserSignup(w, r)
		case http.MethodPost:
			app.UserSignupPost(w, r)
		default:
			MethodNotAllowedHandler(w, r, []string{http.MethodGet, http.MethodPost})
		}
	})

	mux.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.UserLogin(w, r)
		case http.MethodPost:
			app.UserLoginPost(w, r)
		default:
			MethodNotAllowedHandler(w, r, []string{http.MethodGet, http.MethodPost})
		}
	})

	mux.Handle("/user/logout/", app.RequireAuthentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodGet {
			app.UserLogout(w, r)
		} else {
			MethodNotAllowedHandler(w, r, []string{http.MethodGet})
		}
	})))

	return app.RecoverPanic(app.LogRequest(SecureHeaders(mux)))
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request, allowedMethods []string) {
	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
