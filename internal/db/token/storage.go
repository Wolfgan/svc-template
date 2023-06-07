package token

// Subject данные для хранения токенов.
type Subject struct {
	ID    uint64
	Email string
	Phone string
	Roles []string
}

// Storage интерфейс хранилища токенов.
type Storage[k string, v *Subject] interface {
	Set(key k, value v) error
	Get(key k) (v, error)
	Del(key k) error
}
