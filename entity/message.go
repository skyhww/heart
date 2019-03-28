package entity

import "time"

type Message struct {
	Id         int        `db:"id"`
	Message    string     `db:"message"`
	FromUser   int        `db:"from_user"`
	ToUser     int        `db:"to_user"`
	Read       int        `db:"read"`
	CreateTime *time.Time `db:"create_time"`
}

func NewMessage(message string, from, to int) *Message {
	m := &Message{}
	m.Message = message
	m.FromUser = from
	m.ToUser = to
	return m
}

type MessagePersist interface {
	Save(message Message)
	Get(id int)
}
