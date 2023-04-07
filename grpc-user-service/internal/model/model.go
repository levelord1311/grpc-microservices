package model

type User struct {
	ID       uint64 `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Name     string `db:"name"`
	Surname  string `db:"surname"`
}

type UserEvent struct {
	ID      uint64    `db:"id"`
	UserID  uint64    `db:"user_id"`
	Type    EventType `db:"event_type"`
	Payload []byte    `db:"payload"`
}

type UserEvents struct {
	UserEvent
}

type EventType string

const (
	Created EventType = "created"
	Updated EventType = "updated"
	Removed EventType = "removed"
)
