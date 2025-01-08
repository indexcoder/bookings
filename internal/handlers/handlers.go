package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/indexcoder/bookings/internal/config"
	"github.com/indexcoder/bookings/internal/driver"
	"github.com/indexcoder/bookings/internal/forms"
	"github.com/indexcoder/bookings/internal/helpers"
	"github.com/indexcoder/bookings/internal/models"
	"github.com/indexcoder/bookings/internal/render"
	"github.com/indexcoder/bookings/internal/repository"
	"github.com/indexcoder/bookings/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{App: app, DB: dbrepo.NewPostgresRepo(app, db.SQL)}
}

func NewTestRepo(app *config.AppConfig) *Repository {
	return &Repository{App: app, DB: dbrepo.NewTestingDBRepo(app)}
}

func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, req *http.Request) {
	m.DB.AllUsers()
	render.Template(w, req, "home.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "about.html", &models.TemplateData{})
}

func (m *Repository) Features(w http.ResponseWriter, req *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World again Features."
	remoteIp := m.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	render.Template(w, req, "features.html", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) Search(w http.ResponseWriter, req *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World again Contact."
	remoteIp := m.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	render.Template(w, req, "search.html", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) SearchPost(w http.ResponseWriter, req *http.Request) {

	sd := req.Form.Get("start")
	ed := req.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(req.Context(), "error", "There are no rooms available")
		http.Redirect(w, req, "/search-availability", http.StatusSeeOther)
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(req.Context(), "reservation", res)

	render.Template(w, req, "choose-room.html", &models.TemplateData{Data: data})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Msg       string `json:"msg"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) AvailabilityJson(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		resp := jsonResponse{OK: false, Msg: "Internal server error"}
		out, _ := json.MarshalIndent(resp, "", "   ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := req.Form.Get("start")
	ed := req.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(req.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{OK: false, Msg: "Error connecting to database"}
		out, _ := json.MarshalIndent(resp, "", "   ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Msg:       "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, req *http.Request) {

	roomId, err := strconv.Atoi(chi.URLParam(req, "id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomId
	m.App.Session.Put(req.Context(), "reservation", res)

	http.Redirect(w, req, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) Reservation(w http.ResponseWriter, req *http.Request) {

	res, ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(req.Context(), "error", "Невозможно получить бронирование из сеанса")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(req.Context(), "error", "Комната не найдено!")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(req.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, req, "make-reservation.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		m.App.Session.Put(req.Context(), "error", "can't parse form")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := req.Form.Get("start_date")
	ed := req.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(req.Context(), "error", "can't parse start date")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(req.Context(), "error", "can't parse end date")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(req.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(req.Context(), "error", "invalid date!")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		RoomID:    roomID,
		StartDate: startDate,
		EndDate:   endDate,
		FirstName: req.Form.Get("first_name"),
		LastName:  req.Form.Get("last_name"),
		Email:     req.Form.Get("email"),
		Phone:     req.Form.Get("phone"),
	}

	form := forms.New(req.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, req, "make-reservation.html", &models.TemplateData{Form: form, Data: data})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(req.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(req.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	// send notification

	htmMessage := fmt.Sprintf(`<strong>Reservation Confirmation</strong><br>Dear %s:, <br> This is confirm your reservation from %s to %s.`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		From:     reservation.Email,
		To:       "recipient@example.com",
		Subject:  "Reservation Confirmation",
		Content:  htmMessage,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	m.App.Session.Put(req.Context(), "reservation", reservation)
	http.Redirect(w, req, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, req *http.Request) {
	res, ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(req.Context(), "error", "can't get from session")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(req.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = res

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, req, "reservation-summary.html", &models.TemplateData{Data: data, StringMap: stringMap})

}

func (m *Repository) BookRoom(w http.ResponseWriter, req *http.Request) {
	roomID, _ := strconv.Atoi(req.URL.Query().Get("id"))
	sd := req.URL.Query().Get("s")
	ed := req.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(req.Context(), "reservation", res)

	http.Redirect(w, req, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) Login(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "login.html", &models.TemplateData{Form: forms.New(nil)})
}

func (m *Repository) PostLogin(w http.ResponseWriter, req *http.Request) {
	_ = m.App.Session.RenewToken(req.Context())

	err := req.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(req.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, req, "login.html", &models.TemplateData{Form: form})
		return
	}

	id, _, err := m.DB.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		log.Println(err)
		m.App.Session.Put(req.Context(), "error", "authentication failed")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(req.Context(), "user_id", id)
	m.App.Session.Put(req.Context(), "flash", "Logged in successfully!")
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, req *http.Request) {
	_ = m.App.Session.Destroy(req.Context())
	_ = m.App.Session.RenewToken(req.Context())
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, req *http.Request) {
	render.Template(w, req, "admin-dashboard.html", &models.TemplateData{Form: forms.New(nil)})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, req *http.Request) {

	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, req, "admin-all-reservations.html", &models.TemplateData{Data: data, Form: forms.New(nil)})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, req *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, req, "admin-new-reservations.html", &models.TemplateData{Data: data, Form: forms.New(nil)})
}

func (m *Repository) AdminShowReservations(w http.ResponseWriter, req *http.Request) {

	exploded := strings.Split(req.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	srs := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = srs

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, req, "admin-show-reservations.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

func (m *Repository) AdminPostShowReservations(w http.ResponseWriter, req *http.Request) {

	exploded := strings.Split(req.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	srs := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = srs

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
	}

	res.FirstName = req.FormValue("first_name")
	res.LastName = req.FormValue("last_name")
	res.Email = req.FormValue("email")
	res.Phone = req.FormValue("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(req.Context(), "flash", "Изменено!")

	http.Redirect(w, req, fmt.Sprintf("/admin/reservations-%s", srs), http.StatusSeeOther)

}

func (m *Repository) AdminProcessReservation(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(req, "id"))
	src := chi.URLParam(req, "src")

	_ = m.DB.UpdateProcessedForReservation(id, 1)
	m.App.Session.Put(req.Context(), "flash", "Reservation marked as processed!")
	http.Redirect(w, req, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(req, "id"))
	src := chi.URLParam(req, "src")

	_ = m.DB.DeleteReservation(id)
	m.App.Session.Put(req.Context(), "flash", "Reservation deleted!")
	http.Redirect(w, req, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminCalendarReservations(w http.ResponseWriter, req *http.Request) {

	now := time.Now()

	if req.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(req.URL.Query().Get("y"))
		month, _ := strconv.Atoi(req.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMont := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMont.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMont; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		restrictions, err := m.DB.GetRestrictionsForRoomsByDate(x.ID, firstOfMont, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
				blockMap[y.StartDate.Format("2006-01-2")] = y.RestrictionID
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(req.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)

	}

	render.Template(w, req, "admin-calendar-reservations.html", &models.TemplateData{StringMap: stringMap, Data: data, IntMap: intMap, Form: forms.New(nil)})
}

func (m *Repository) AdminPostCalendarReservations(w http.ResponseWriter, req *http.Request) {
	log.Println("Works")
}
