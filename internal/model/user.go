package model

import (
	"time"

	"github.com/uptrace/bun"
)

// User Структура данных с информацией о пользователе
type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            uint64     `bun:"id,pk,autoincrement"`
	Email         string     `bun:"email,unique"`
	Phone         string     `bun:"phone,unique"`
	Password      string     `bun:"password,notnull"`
	CreatedAt     *time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `bun:"updated_at,nullzero"`
	DeletedAt     *time.Time `bun:"deleted_at,soft_delete,nullzero"`
	BlockedAt     *time.Time `bun:"blocked_at,nullzero"`
	Profile       *Profile   `bun:"rel:has-one,join:id=user_id"`
	//Roles         []Role     `bun:"m2m:auth_user_roles,join:User=Role"`
}

type Profile struct {
	bun.BaseModel `bun:"table:profiles"`
	ID            uint64     `bun:"id,pk,autoincrement"`
	UserID        uint64     `bun:"user_id"`
	Name          string     `bun:"name"`
	Surname       string     `bun:"surname"`
	Patronymic    string     `bun:"patronymic"`
	Sex           bool       `bun:"sex,nullzero,notnull,default:true"`
	Birthday      *time.Time `bun:"birthday,nullzero"`
	Country       string     `bun:"country"`
	City          string     `bun:"city"`
	Address       string     `bun:"address"`
}
