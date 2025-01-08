package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/indexcoder/bookings/internal/config"
	"github.com/indexcoder/bookings/internal/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/features", handlers.Repo.Features)

	mux.Get("/search-availability", handlers.Repo.Search)
	mux.Post("/search-availability", handlers.Repo.SearchPost)
	mux.Post("/availability-json", handlers.Repo.AvailabilityJson)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login", handlers.Repo.PostLogin)
	mux.Get("/logout", handlers.Repo.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-calendar", handlers.Repo.AdminCalendarReservations)
		mux.Post("/reservations-calendar", handlers.Repo.AdminPostCalendarReservations)

		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservations)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservations)
		mux.Get("/process-reservation/{src}/{id}", handlers.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}", handlers.Repo.AdminDeleteReservation)

	})

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
