package entity

type UserExtraInfo struct {
	UserId int `db:"user_id"`
	//关注用户数
	FollowUserCount int `db:"follow_user_count"`
	//收藏帖子数量
	Collections int `db:"collections"`
	//粉丝数
	FollowedUserCount int `db:"followed_user_count"`
	//未读消息数--普通消息
	UnreadMessageCount int `db:"unread_message_count"`
	//未读消息数--站点消息
	UnreadAdminMessageCount int `db:"unread_admin_message_count"`
}
