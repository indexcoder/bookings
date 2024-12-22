package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/indexcoder/bookings/internal/config"
	"github.com/indexcoder/bookings/internal/helpers"
	"github.com/indexcoder/bookings/internal/models"
	"github.com/indexcoder/bookings/internal/render"
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
	render.Template(w, req, "home.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "about.html", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, req *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World again Contact."
	remoteIp := m.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	render.Template(w, req, "contact.html", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) Features(w http.ResponseWriter, req *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World again Features."
	remoteIp := m.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	render.Template(w, req, "features.html", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) Search(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	start := req.Form.Get("start")
	end := req.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s.", start, end)))
}

type jsonResponse struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

func (m *Repository) SearchJson(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := jsonResponse{
		OK:  true,
		Msg: "ok available",
	}
	out, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
