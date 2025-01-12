package dbrepo

import (
	"context"
	"errors"
	"github.com/indexcoder/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (m *postgresRepo) AllUsers() bool {
	return true
}

func (m *postgresRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newID int
	stmt := `insert into reservations (room_id, first_name, last_name, email, phone, start_date, end_date, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt, res.RoomID, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, time.Now(), time.Now()).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (m *postgresRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `Insert into room_restrictions (room_id, reservation_id, restriction_id, start_date, end_date, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, stmt, res.RoomID, res.ReservationID, res.RestrictionID, res.StartDate, res.EndDate, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var numRows int
	query := `select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date;`
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

func (m *postgresRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select r.id, r.room_name from rooms r where r.id not in (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (m *postgresRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := `select r.id, r.room_name from rooms r where r.id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName)
	if err != nil {
		return room, err
	}
	return room, nil
}

func (m *postgresRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select u.id, u.first_name, u.last_name, u.email, u.phone, u.access_level from users u where u.id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.AccessLevel)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (m *postgresRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1, last_name=$2, email=$3, phone=$4, access_level=$5`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Phone, u.AccessLevel)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return id, "", errors.New("wrong password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil
}

func (m *postgresRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date,
                     r.created_at, rm.id, rm.room_name
from reservations r left join rooms rm on r.room_id = rm.id order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(&reservation.ID, &reservation.FirstName, &reservation.LastName, &reservation.Email, &reservation.Phone,
			&reservation.StartDate, &reservation.EndDate, &reservation.CreatedAt, &reservation.RoomID, &reservation.Room.RoomName)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresRepo) AllNewReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.created_at, r.processed, rm.id, rm.room_name from reservations r left join rooms rm on r.room_id = rm.id where processed = 0 order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(&reservation.ID, &reservation.FirstName, &reservation.LastName, &reservation.Email, &reservation.Phone, &reservation.StartDate, &reservation.EndDate, &reservation.CreatedAt, &reservation.Processed, &reservation.RoomID, &reservation.Room.RoomName)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.created_at, r.processed, rm.id, rm.room_name from reservations r left join rooms rm on r.room_id = rm.id where r.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&reservation.ID, &reservation.FirstName, &reservation.LastName, &reservation.Email, &reservation.Phone, &reservation.StartDate, &reservation.EndDate, &reservation.CreatedAt, &reservation.Processed, &reservation.Room.ID, &reservation.Room.RoomName)
	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

func (m *postgresRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5 where id=$6`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Phone, u.UpdatedAt, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from reservations where id = $1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update reservations set processed=$1 where id = $2`
	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select id, room_name, created_at from rooms order by room_name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}

	defer rows.Close()
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName, &room.CreatedAt)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil

}

func (m *postgresRepo) GetRestrictionsForRoomsByDate(roomID int, startDate, endDate time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var roomRestrictions []models.RoomRestriction

	query := `select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date from room_restrictions where $1 < end_date and $2 >= start_date and room_id = $3`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction

		err := rows.Scan(&r.ID, &r.ReservationID, &r.RestrictionID, &r.RoomID, &r.StartDate, &r.EndDate)
		if err != nil {
			return nil, err
		}
		roomRestrictions = append(roomRestrictions, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roomRestrictions, nil
}

func (m *postgresRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from room_restrictions where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
