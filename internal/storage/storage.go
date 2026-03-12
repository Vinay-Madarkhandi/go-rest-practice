package storage

import "github.com/Vinay-Madarkhandi/go-rest-practice/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetById(id int64) (types.Student, error)
}
