package dbrepo

import (
	"context"
	"errors"
	"time"
	"web3/models"

	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) InsertPost(newPost models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `insert into posts(title, content, user_id) values ($1, $2, $3)`

	_, err := m.DB.ExecContext(ctx, query, newPost.Title, newPost.Content, newPost.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `select name, email, password, acct_created, last_login, user_type, id from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var user models.User

	err := row.Scan(
		&user.Name,
		&user.Email,
		&user.Password,
		&user.AcctCreated,
		&user.LastLogin,
		&user.UserType,
		&user.ID,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `update users set name = $1, email = $2, last_login = $3, user_type = $4`

	_, err := m.DB.ExecContext(ctx, query, u.Name, u.Email, time.Now(), u.UserType)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) AuthenticateUser(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	query := `select id, password from users where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("password is incorrect")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil
}
