package handlers

import (
	"github.com/indexcoder/bookings/pkg/config"
	"github.com/indexcoder/bookings/pkg/models"
	"github.com/indexcoder/bookings/pkg/render"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{App: a}
}

func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, req *http.Request) {
	remoteIp := req.RemoteAddr
	m.App.Session.Put(req.Context(), "remote_ip", remoteIp)

	render.Template(w, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, req *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World again."

	remoteIp := m.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.Template(w, "about.page.html", &models.TemplateData{StringMap: stringMap})
}
