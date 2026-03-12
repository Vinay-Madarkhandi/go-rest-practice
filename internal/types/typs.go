package types

type Student struct {
	Id    int64  `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required"`
}

type StudentCreatedDTO struct {
	Success bool    `json:"success"`
	Data    Student `json:"data"`
}
