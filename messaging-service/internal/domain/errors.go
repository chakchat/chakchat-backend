package domain

import "errors"

var (
	ErrAdminNotMember      = errors.New("group members doesn't include admin")
	ErrGroupNameEmpty      = errors.New("group name is empty")
	ErrGroupNameTooLong    = errors.New("group name is too long")
	ErrGroupDescTooLong    = errors.New("group description is too long")
	ErrUserAlreadyMember   = errors.New("user is already a member of a chat")
	ErrMemberIsAdmin       = errors.New("group member is admin")
	ErrGroupPhotoEmpty     = errors.New("group photo is empty")
	ErrChatWithMyself      = errors.New("chat with myself")
	ErrChatBlocked         = errors.New("chat is blocked")
	ErrFileTooBig          = errors.New("file is too big")
	ErrReactionNotFromUser = errors.New("the reaction is not from this user")
	ErrTooMuchTextRunes    = errors.New("too much runes in text")
	ErrTextEmpty           = errors.New("the text is empty")
	ErrUserNotSender       = errors.New("user is not update's sender")
	ErrUpdateNotFromChat   = errors.New("update is not from this chat")
	ErrUpdateDeleted       = errors.New("update is deleted")
	ErrAlreadyBlocked      = errors.New("chat is already blocked")
	ErrAlreadyUnblocked    = errors.New("chat is already unblocked")
)
