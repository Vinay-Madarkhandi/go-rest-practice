package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/config"
	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/types"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Db *sql.DB
}

func New(config *config.Config) (*MySQL, error) {
	// Open connection with MySQL DB
	db, err := sql.Open("mysql", config.DataBaseDSN)
	if err != nil {
		return nil, err
	}

	// Ensure that the connection works
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	slog.Info("connected to database")
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS students (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		age INT NOT NULL,
		email VARCHAR(255) NOT NULL
	)	
	`)
	if err != nil {
		return nil, err
	}

	return &MySQL{db}, nil
}

func (m *MySQL) CreateStudent(name string, email string, age int) (int64, error) {
	query := "INSERT INTO students (name, email, age) VALUES (?, ?, ?)"
	db, err := m.Db.Exec(query, name, email, age)
	if err != nil {
		return 0, err
	}
	return db.LastInsertId()
}

func (m *MySQL) GetById(id int64) (types.Student, error) {
	stmt, err := m.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil

}
