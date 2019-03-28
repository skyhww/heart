package entity

type UserPost struct {
	Id int `db:"id"`
	UserId int `db:"user_id"`
	Content []byte `db:"content"`
	Attach string `db:"attach"`
	CreateTime int `db:"create_time"`
	Enable  int `db:"enable"`
}
