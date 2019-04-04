package entity

type UserPost struct {
	Id         int    `db:"id"`
	UserId     int    `db:"user_id"`
	Content    string `db:"content"`
	CreateTime int    `db:"create_time"`
	Enable     int    `db:"enable"`
}

type PostAttach struct {
	Id         int64 `db:"id"`
	PostId     int64 `db:"post_id"`
	No         int   `db:"no"`
	CreateTime int   `db:"create_time"`
	Enable     int   `db:"enable"`
}
