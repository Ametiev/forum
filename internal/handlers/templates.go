package handlers

import (
	"path/filepath"
	"text/template"
	"time"

	"dyelesho/forum/internal/models"
)

type TemplateData struct {
	CurrentYear     int
	Post            *models.Post
	Posts           []*models.Post
	Form            any
	IsAuthenticated bool
	Category        string
	Categories      []string
	Comments        []models.Comment
	ErrorStruct     *ErrorStruct
	CommentError    bool
}

func HumanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": HumanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
