package types

type Student struct {
	ID 		int64
	Name 	string `validate:"required"`
	Age 	int `validate:"required"`
	Email 	string `validate:"required"`
}