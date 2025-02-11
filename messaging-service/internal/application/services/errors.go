package services

import "errors"

var (
	ErrAdminNotMember = errors.New("service: admin is not group member")

	ErrGroupNameEmpty   = errors.New("service: group name is empty")
	ErrGroupNameTooLong = errors.New("service: group name is too long")
	ErrGroupDescTooLong = errors.New("service: group description is too long")

	ErrUserAlreadyMember = errors.New("service: user is already a member of a chat")
	ErrMemberIsAdmin     = errors.New("service: group member is admin")

	ErrGroupPhotoEmpty = errors.New("service: group photo is empty")
)

var (
	ErrFileNotFound = errors.New("service: file not found")
	ErrInvalidPhoto = errors.New("service: invalid photo")
)

var (
	ErrChatNotFound         = errors.New("service: chat not found")
	ErrChatAlreadyBlocked   = errors.New("service: chat already blocked")
	ErrChatAlreadyUnblocked = errors.New("service: chat already unblocked")
	ErrUserNotMember        = errors.New("service: user not chat member")

	ErrChatAlreadyExists = errors.New("service: chat already exists")
	ErrChatWithMyself    = errors.New("service: chat with myself")

	ErrInternal = errors.New("service: unknown error")
)
