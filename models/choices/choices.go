package choices

type ReactionChoice string

const (
	RLIKE  ReactionChoice = "LIKE"
	RLOVE  ReactionChoice = "LOVE"
	RHAHA  ReactionChoice = "HAHA"
	RWOW   ReactionChoice = "WOW"
	RSAD   ReactionChoice = "SAD"
	RANGRY ReactionChoice = "ANGRY"
)

type NotificationChoice string

const (
	NREACTION NotificationChoice = "REACTION"
	NCOMMENT  NotificationChoice = "COMMENT"
	NREPLY    NotificationChoice = "REPLY"
	NADMIN    NotificationChoice = "ADMIN"
)

type FriendStatusChoice string

const (
	FPENDING  FriendStatusChoice = "PENDING"
	FACCEPTED FriendStatusChoice = "ACCEPTED"
)