package entity

import (
	"time"
	"github.com/jmoiron/sqlx"
)

type Message struct {
	Id         int64      `db:"id" json:"id"`
	Message    *string    `db:"message" json:"message"`
	FromUser   int64      `db:"from_user" json:"from_user"`
	ToUser     int64      `db:"to_user" json:"to_user"`
	Read       int        `db:"read" json:"-"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Url        *string    `db:"url" json:"url"`
}

func (message *Message) IsRead() bool {
	return message.Read == 1
}

func NewMessage(message, url *string, from, to int64) *Message {
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
	GetUnreadMessage(userId, lastId int64) (*Message, error)
	//保存流信息
	//SaveStreamMessage(message *Message) error
	GetUnreadCount(userId int64) (int, error)
	SetRead(userId, messageId int64) (int64, error)
}

type MessageDao struct {
	DB *sqlx.DB
}

func (messageDao *MessageDao) SetRead(userId, messageId int64) (int64, error) {
	tx := messageDao.DB.MustBegin()
	r, err := tx.Exec("update message set read=1 where id=? and to_user=? and read=0", messageId, userId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit()
	return affected, err
}

func (messageDao *MessageDao) Save(message *Message) error {
	tx := messageDao.DB.MustBegin()
	r, err := tx.Exec("insert into message(message,from_user,to_user,read,create_time,url) values(:message,:from_user,:to_user,0,:create_time,:url)", message)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	message.Id, _ = r.LastInsertId()
	return nil
}

func (messageDao *MessageDao) GetUnreadMessage(userId, lastId int64) (*Message, error) {
	m := &Message{}
	err := messageDao.DB.Select(m, "select * from message where to_user=?  and id>?  and read=0 order by create_time desc limit 1", userId, lastId)
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
