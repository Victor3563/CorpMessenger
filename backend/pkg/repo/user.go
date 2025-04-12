//REalesation User config + funk like Create, Auth, Update and delete

package repo

import (
	"database/sql"
	"errors"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Repository инкапсулирует работу с БД
type Repository struct {
	DB *sql.DB
}

// NewRepository возвращает новый объект Repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// Create new user
func (r *Repository) CreateUser(username, password, email string) (User, error) {
	// Check if username already exists
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	if err != nil {
		return User{}, err
	}
	if exists {
		return User{}, errors.New("username already exists")
	}

	var user User
	query := `INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3)
	RETURNING id, username, password, email;`

	err = r.DB.QueryRow(query, username, password, email).Scan(&user.ID, &user.Name, &user.Password, &user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// AuthenticateUser checks if a username and password are correct
func (r *Repository) AuthenticateUser(username, password string) (User, error) {
	var user User
	query := `SELECT id, username, password, email FROM users WHERE username=$1`
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Name, &user.Password, &user.Email)
	if err != nil {
		return User{}, errors.New("invalid username or password")
	}
	if user.Password != password {
		return User{}, errors.New("password")
	}
	return user, nil
}

// UpdateUser updates user attributes in the store
func (r *Repository) UpdateUser(id int, newUsername, newPassword, newEmail string) error {
	query := `
			UPDATE users
			SET username=$1, password=$2, email=$3, updated_at=NOW()
			WHERE id = $4
	`
	res, err := r.DB.Exec(query, newUsername, newPassword, newEmail, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// DeleteUser removes a user from the Data
func (r *Repository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id=$1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// создали отдельный класс чтобы точно не возращать лишнюю инфу, например пароль
type UserInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// FindUser find a user from the Data
func (r *Repository) FindUser(name string) ([]User, error) {
	// SQL-запрос для выбора пользователей, чье имя начинается с name
	query := `SELECT id, username, email FROM users WHERE username LIKE $1 || '%'`
	rows, err := r.DB.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// FindUser find a user from the Data
func (r *Repository) FindUserbyID(user_id int) (User, error) {
	var user User
	// SQL-запрос для выбора пользователей, чье имя начинается с name
	query := `SELECT id, username, email, password FROM users WHERE id = $1`
	err := r.DB.QueryRow(query, user_id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
