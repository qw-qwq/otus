package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jetuuuu/hl_homework/database"
)

type User struct {
	Login     string         `json:"login" db:"login"`
	Password  string         `json:"-" db:"password"`
	FirstName string         `json:"first_name" db:"first_name"`
	LastName  string         `json:"last_name" db:"last_name"`
	Age       uint           `json:"age,string" db:"age"`
	Sex       uint           `json:"sex,string" db:"sex"`
	City      string         `json:"city" db:"city"`
	Hobby     sql.NullString `json:"hobby" db:"hobby"`
	Friends   []string       `json:"friends"`
}

func (db *DB) CreateUser(ctx context.Context, u User) error {
	const query = "insert into users (login, password, first_name, last_name, age, sex, city)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.db.Exec(ctx, query, u.Login, u.Password, u.FirstName, u.LastName, u.Age, u.Sex, u.City, u.Hobby)

	return err
}

func (db *DB) GetUserByID(ctx context.Context, login string) (User, error) {
	const queryUser = "select * from users where login = ?"

	u := User{}
	err := db.db.QueryRow(ctx, queryUser, login).
		Scan(&u.Login, &u.Password, &u.FirstName, &u.LastName, &u.Age, &u.Sex, &u.City, &u.Hobby)
	if err != nil {
		return User{}, err
	}

	rows, err := db.db.Query(ctx, "select user2 from friends where user1 = ?", login)
	if err != nil {
		return User{}, err
	}

	for rows.Next() {
		var f string
		err = rows.Scan(&f)
		if err != nil {
			return User{}, err
		}

		u.Friends = append(u.Friends, f)
	}

	err = rows.Err()
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (db *DB) IsUserExist(ctx context.Context, login, pwd string) bool {
	const query = "select count(*) from users where login = ? and password = ?"

	var cnt int
	_ = db.db.QueryRow(ctx, query, login, pwd).Scan(&cnt)

	return cnt != 0
}

func (db *DB) GetFriendsNames(ctx context.Context, userIDs []string) ([]string, error) {
	query := "select first_name, last_name from users where login in (%s)"

	if len(userIDs) == 1 {
		query = fmt.Sprintf(query, "?")
	} else if len(userIDs) > 1 {
		query = fmt.Sprintf(query, "?"+strings.Repeat(",?", len(userIDs)-1))
	} else {
		query = fmt.Sprintf(query, "''")
	}

	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		args[i] = id
	}

	rows, err := db.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var names []string
	for rows.Next() {
		var firstName, lastName string
		err = rows.Scan(&firstName, &lastName)
		if err != nil {
			return nil, err
		}

		names = append(names, firstName+" "+lastName)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (db *DB) MakeFriends(ctx context.Context, user1, user2 string) error {
	return db.db.Transact(ctx, sql.LevelDefault, func(db *database.DB) error {
		_, err := db.Exec(ctx, "insert into friends (user1, user2) VALUES (?, ?)", user1, user2)
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, "insert into friends (user1, user2) VALUES (?, ?)", user2, user1)
		return err
	})
}
