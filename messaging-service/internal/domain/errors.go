package domain

type Error struct {
	text string
}

func (e Error) Error() string {
	return e.text
}

// NOTE:
// If you add an error here
// You alse should add it to `errmap` package
var (
	ErrAdminNotMember      = Error{"group members doesn't include admin"}
	ErrGroupNameEmpty      = Error{"group name is empty"}
	ErrGroupNameTooLong    = Error{"group name is too long"}
	ErrGroupDescTooLong    = Error{"group description is too long"}
	ErrUserAlreadyMember   = Error{"user is already a member of a chat"}
	ErrMemberIsAdmin       = Error{"group member is admin"}
	ErrGroupPhotoEmpty     = Error{"group photo is empty"}
	ErrChatWithMyself      = Error{"chat with myself"}
	ErrChatBlocked         = Error{"chat is blocked"}
	ErrFileTooBig          = Error{"file is too big"}
	ErrReactionNotFromUser = Error{"the reaction is not from this user"}
	ErrTooManyTextRunes    = Error{"too many runes in text"}
	ErrTextEmpty           = Error{"the text is empty"}
	ErrUserNotSender       = Error{"user is not update's sender"}
	ErrUpdateNotFromChat   = Error{"update is not from this chat"}
	ErrUpdateDeleted       = Error{"update is deleted"}
	ErrAlreadyBlocked      = Error{"chat is already blocked"}
	ErrAlreadyUnblocked    = Error{"chat is already unblocked"}
	ErrSenderNotAdmin      = Error{"sender is not admin"}
	ErrUserNotMember       = Error{"user is not member of a chat"}
	ErrInvalidDeleteMode   = Error{"invalid delete mode"}
	ErrInvalidReactionType = Error{"invalid reaction type"}
)
