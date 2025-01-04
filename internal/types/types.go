package types

type Student struct {
	ID 		int
	Name 	string `validate:"required"`
	Age 	int `validate:"required"`
	Email 	string `validate:"required"`
}