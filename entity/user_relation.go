package entity
/**

import (
	"time"
	"heart/service/common"
	"github.com/jmoiron/sqlx"
	"database/sql"
)

type UserRelation struct {
	Id         int64      `db:"id" json:"id"`
	UserId     int64      `db:"user_id" json:"user_id"`
	FriendId   int64      `db:"friend_id" json:"friend_id"`
	Enable     int        `db:"enable" json:"-"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
}

type UserRelationPersist interface {
	//保存关系
	Add(userId, friendId int64) error
	//分页查询关注的用户
	GetFriend(userId int64, page *base.Page) ([]*UserRelation, error)
	//是关注了某个用户
	ExistsRelation(userId, friendId int64) (*UserRelation, error)
	//删除
	Delete(relation *UserRelation) error

	Get(id int64) (*UserRelation, error)

	FollowedCount()
	UserFollowInfoPersist
}

type RelationDao struct {
	DB *sqlx.DB
}

func (relationDao *RelationDao) Add(userId, friendId int64) error {
	tx := relationDao.DB.MustBegin()
	n := time.Now()
	relation := &UserRelation{UserId: userId, FriendId: friendId, CreateTime: &n}
	r, err := tx.NamedExec("update user_relation set enable=1,create_time=:create_time  where user_id=:user_id and friend_id=:friend_id", relation)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows == 1 {
		return tx.Commit()
	} else {
		_, err := tx.NamedExec("insert into user_relation(user_id,friend_id,enable,create_time) values(:user_id,:friend_id,1,:create_time)", relation)
		if err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	}
}
func (relationDao *RelationDao) GetFriend(userId int64, page *base.Page) ([]*UserRelation, error) {
	r := &[]*UserRelation{}
	err := relationDao.DB.Select(r, "select * from user_relation where user_id=? order by create_time desc ", userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return *r, nil
}
func (relationDao *RelationDao) ExistsRelation(userId, friendId int64) (*UserRelation, error) {
	r := &UserRelation{}
	err := relationDao.DB.Get(r, "select * from user_relation where user_id=? and friend_id=? and enable=1 ")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return r, nil
}
func (relationDao *RelationDao) Delete(relation *UserRelation) error {
	tx := relationDao.DB.MustBegin()
	_, err := tx.Exec("update user_relation set enable=0 where id=?", relation.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
func (relationDao *RelationDao) Get(id int64) (*UserRelation, error) {
	r := &UserRelation{}
	err := relationDao.DB.Get(r, "select * from user_relation where id=? ", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return r, nil
}
**/