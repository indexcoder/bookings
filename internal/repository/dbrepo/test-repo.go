package dbrepo

import (
	"errors"
	"github.com/indexcoder/bookings/internal/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("cannot insert a reservation twice")
	}
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	if res.RoomID == 1000 {
		return errors.New("cannot insert a restriction twice")
	}
	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("id out of range")
	}
	return room, nil
}

func (m *testDBRepo) GetUserById(id int) (models.User, error) {
	var u models.User
	return u, nil
}
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservation models.Reservation
	return reservation, nil
}

func (m *testDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRestrictionsForRoomsByDate(roomID int, startDate, endDate time.Time) ([]models.RoomRestriction, error) {
	var roomRestrictions []models.RoomRestriction
	return roomRestrictions, nil
}
