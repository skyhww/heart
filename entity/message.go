package entity

import (
	"time"
	"github.com/jmoiron/sqlx"
)

type Message struct {
	Id         int64      `db:"id"`
	Message    string     `db:"message"`
	FromUser   int64      `db:"from_user"`
	ToUser     int64      `db:"to_user"`
	Read       int        `db:"read"`
	CreateTime *time.Time `db:"create_time"`
	Url        string     `db:"url"`
}

func (message *Message) IsRead() bool {
	return message.Read == 1
}

func NewMessage(message, url string, from, to int64) *Message {
	m := &Message{}
	m.Message = message
	m.FromUser = from
	m.ToUser = to
	m.Url = url
	return m
}

type MessagePersist interface {
	//保存文本信息
	Save(message *Message) error
	GetUnreadMessage(userId int64) (*[]Message, error)
	//保存流信息
	//SaveStreamMessage(message *Message) error
	GetUnreadCount(userId int64) (int, error)
}

type MessageDao struct {
	DB *sqlx.DB
}

func (messageDao *MessageDao) Save(message *Message) error {
	tx := messageDao.DB.MustBegin()
	r, err := tx.Exec("insert into message(message,from_user,to_user,read,create_time,url) values(:message,:from_user,:to_user,0,:create_time,:url)", message)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer tx.Commit()
	message.Id, _ = r.LastInsertId()
	return nil
}

func (messageDao *MessageDao) GetUnreadMessage(userId int64) (*[]Message, error) {
	m := &[]Message{}
	err := messageDao.DB.Select(m, "select * from message where to_user=? order by create_time desc")
	if err != nil {
		return nil, err
	}
	return m, nil
}
func (messageDao *MessageDao) GetUnreadCount(userId int64) (int, error) {
	m := 0
	err := messageDao.DB.Select(m, "select count(id) from message where to_user=? order by create_time desc")
	if err != nil {
		return 0, err
	}
	return m, nil
}

/*func (messageDao *MessageDao) SaveStreamMessage(message *Message) error {
	if message.StreamMessage == nil {
		return nil
	}
	tx := messageDao.DB.MustBegin()
	r, err := tx.Exec("insert into message(from_user,to_user,read,create_time) values(:from_user,:to_user,0,:create_time)", message)
	if err != nil {
		tx.Rollback()
		return err
	}
	message.StreamMessage.MessageId, err = r.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	r2, err := tx.Exec("insert into message_stream(message_id,url) values(:message_id,:url)", message.StreamMessage)
	if err != nil {
		tx.Rollback()
		return err
	}
	message.StreamMessage.Id, err = r2.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
*/
