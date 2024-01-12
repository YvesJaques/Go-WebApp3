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

func (m *postgresDBRepo) GetAnArticle() (int, int, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `select id, user_id, title, content from posts LIMIT 1`

	row := m.DB.QueryRowContext(ctx, query)

	var id int
	var userID int
	var title string
	var content string

	err := row.Scan(&id, &userID, &title, &content)
	if err != nil {
		return id, userID, "", "", err
	}
	return id, userID, title, content, nil
}

func (m *postgresDBRepo) GetArticles(optional_amount ...int) (models.ArticleList, error) {
	amount := 3
	if len(optional_amount) > 0 {
		amount = optional_amount[0]
	}

	query := `select id, user_id, title, content from posts ORDER BY id DESC LIMIT $1`

	var artList models.ArticleList
	rows, err := m.DB.Query(query, amount)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, uID int
		var title, content string
		err := rows.Scan(&id, &uID, &title, &content)
		if err != nil {
			panic(err)
		}
		artList.ID = append(artList.ID, id)
		artList.UserID = append(artList.UserID, uID)
		artList.Title = append(artList.Title, title)
		artList.Content = append(artList.Content, content)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return artList, nil
}
