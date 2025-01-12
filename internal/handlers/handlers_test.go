package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/indexcoder/bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type postData struct {
	key string
	val string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/features", "GET", http.StatusOK},
	{"non-existent", "/green-noo-ser/this/is", "GET", http.StatusNotFound},

	//{"features", "/search-availability", "GET", http.StatusOK},
	//{"search", "/search", "POST", []postData{
	//	{key: "start", val: "2021-01-01"},
	//	{key: "end", val: "2021-01-02"},
	//}, http.StatusOK},
	//{"search-json", "/search-json", "POST", []postData{
	//	{key: "start", val: "2021-01-01"},
	//	{key: "end", val: "2021-01-02"},
	//}, http.StatusOK},
	//{"about-post", "/about", "POST", []postData{
	//	{key: "start", val: "2021-01-01"},
	//	{key: "end", val: "2021-01-02"},
	//}, http.StatusOK},
}

func TestNewHandler(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, tc := range theTests {
		resp, err := ts.Client().Get(ts.URL + tc.url)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}
		if resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, tc.expectedStatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	handler := http.HandlerFunc(Repo.Reservation)

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusTemporaryRedirect)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {

	//reqBody := "start_date=2026-01-02"
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2026-01-02")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Asan")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Maratov")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Asan@mail.com")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	postedDate := url.Values{}
	postedDate.Add("start_date", "2026-01-02")
	postedDate.Add("end_date", "2026-01-02")
	postedDate.Add("first_name", "Asan")
	postedDate.Add("last_name", "Maratov")
	postedDate.Add("email", "Asan@mail.com")
	postedDate.Add("phone", "1234567890")
	postedDate.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong status code: got %v want %v", rr.Code, http.StatusSeeOther)
	}

	// test for missing post Body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response coder for missing post body status code: got %v want %v", rr.Code, http.StatusSeeOther)
	}

	// test for invalid start date
	postedDate.Set("start_date", "invalid")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response coder for invalid start date code: got %v want %v", rr.Code, http.StatusSeeOther)
	}

	// test for invalid end date
	postedDate.Set("end_date", "invalid")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response coder for invalid end date code: got %v want %v", rr.Code, http.StatusSeeOther)
	}

	// test for invalid room id
	postedDate.Set("room_id", "invalid")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response coder for invalid room id code: got %v want %v", rr.Code, http.StatusSeeOther)
	}

	// test for invalid data
	postedDate.Set("start_date", "2026-01-02")
	postedDate.Set("end_date", "2026-01-02")
	postedDate.Set("first_name", "Asan")
	postedDate.Set("last_name", "N")
	postedDate.Set("email", "Asan@mail.com")
	postedDate.Set("phone", "1234567890")
	postedDate.Set("room_id", "1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response coder for invalid Data code: got %v want %v", rr.Code, http.StatusSeeOther)
	}

	// test for failed to insert reservation into database
	postedDate.Set("room_id", "2")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler when failed when trying to fail inserting reservation code: got %v want %v", rr.Code, http.StatusSeeOther)
	}

	// test for failed to insert restriction into database
	postedDate.Set("room_id", "1000")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler when failed when trying to fail inserting reservation code: got %v want %v", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_AvailabilityJson(t *testing.T) {
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create request
	req, _ := http.NewRequest("POST", "/availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	handler := http.HandlerFunc(Repo.AvailabilityJson)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var j jsonResponse

	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Errorf("Error unmarshalling json response: %v", err)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
