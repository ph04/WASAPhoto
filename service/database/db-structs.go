package database

type DatabaseLogin struct {
	Username string `json:"username"`
}

func DatabaseLoginDefault() DatabaseLogin {
	return DatabaseLogin{
		Username: "",
	}
}

type DatabaseUser struct {
	Id       uint32 `json:"id"`
	Username string `json:"username"`
}

func DatabaseUserDefault() DatabaseUser {
	return DatabaseUser{
		Id:       0,
		Username: "",
	}
}

type DatabasePhoto struct {
	Id           uint32       `json:"id"`
	User         DatabaseUser `json:"user"`
	Url          string       `json:"url"`
	Date         string       `json:"date"`
	LikeCount    int          `json:"like_count"`
	CommentCount int          `json:"comment_count"`
	LikeStatus   bool         `json:"like_status"`
}

func DatabasePhotoDefault() DatabasePhoto {
	return DatabasePhoto{
		Id:           0,
		User:         DatabaseUserDefault(),
		Url:          "",
		Date:         "",
		LikeCount:    0,
		CommentCount: 0,
		LikeStatus:   false,
	}
}

type DatabaseComment struct {
	Id          uint32        `json:"id"`
	User        DatabaseUser  `json:"user"`
	Photo       DatabasePhoto `json:"photo"`
	Date        string        `json:"date"`
	CommentBody string        `json:"comment_body"`
}

func DatabaseCommentDefault() DatabaseComment {
	return DatabaseComment{
		Id:          0,
		User:        DatabaseUserDefault(),
		Photo:       DatabasePhotoDefault(),
		Date:        "",
		CommentBody: "",
	}
}

type DatabaseProfile struct {
	User           DatabaseUser    `json:"user"`
	Photos         []DatabasePhoto `json:"photos"`
	PhotoCount     int             `json:"photo_count"`
	FollowersCount int             `json:"followers_count"`
	FollowingCount int             `json:"following_count"`
	FollowStatus   bool            `json:"follow_status"`
	BanStatus      bool            `json:"ban_status"`
}

func DatabaseProfileDefault() DatabaseProfile {
	emptyArray := make([]DatabasePhoto, 0)

	return DatabaseProfile{
		User:           DatabaseUserDefault(),
		Photos:         emptyArray,
		PhotoCount:     0,
		FollowersCount: 0,
		FollowingCount: 0,
		FollowStatus:   false,
		BanStatus:      false,
	}
}

type DatabaseStream struct {
	User   DatabaseUser    `json:"user"`
	Photos []DatabasePhoto `json:"photos"`
}

func DatabaseStreamDefault() DatabaseStream {
	emptyArray := make([]DatabasePhoto, 0)

	return DatabaseStream{
		User:   DatabaseUserDefault(),
		Photos: emptyArray,
	}
}

type DatabaseUserList struct {
	Users []DatabaseUser `json:"users"`
}

func DatabaseUserListDefault() DatabaseUserList {
	emptyArray := make([]DatabaseUser, 0)

	return DatabaseUserList{
		Users: emptyArray,
	}
}

type DatabaseCommentList struct {
	Comments []DatabaseComment `json:"comments"`
}

func DatabaseCommentListDefault() DatabaseCommentList {
	emptyArray := make([]DatabaseComment, 0)

	return DatabaseCommentList{
		Comments: emptyArray,
	}
}
