package errmap

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
)

type Response struct {
	Code int
	Body restapi.ErrorResponse
}

func Map(err error) Response {
	var domErr domain.Error
	if errors.As(err, &domErr) {
		if resp, ok := domainErrMap[domErr]; ok {
			return resp
		}
		panic(fmt.Errorf("domain error is not mapped: %s", domErr))
	}

	var servicesErr services.Error
	if errors.As(err, &servicesErr) {
		if resp, ok := servicesErrMap[servicesErr]; ok {
			return resp
		}
		panic(fmt.Errorf("services error is not mapped: %s", servicesErr))
	}

	// Other errors are counted as internal server errors.
	return Response{
		Code: http.StatusInternalServerError,
		Body: restapi.ErrorResponse{
			ErrorType:    restapi.ErrTypeInternal,
			ErrorMessage: "Internal Server Error",
		},
	}
}

var servicesErrMap = map[services.Error]Response{
	services.ErrChatAlreadyExists: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "chat_already_exists",
			ErrorMessage: "Chat already exists",
		},
	},
	services.ErrChatNotFound: {
		Code: http.StatusNotFound,
		Body: restapi.ErrorResponse{
			ErrorType:    "chat_not_found",
			ErrorMessage: "Chat is not found",
		},
	},
	services.ErrFileNotFound: {
		Code: http.StatusNotFound,
		Body: restapi.ErrorResponse{
			ErrorType:    "file_not_found",
			ErrorMessage: "File is not found",
		},
	},
	services.ErrInvalidPhoto: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "invalid_photo",
			ErrorMessage: "Invalid photo",
		},
	},
	services.ErrMessageNotFound: {
		Code: http.StatusNotFound,
		Body: restapi.ErrorResponse{
			ErrorType:    "message_not_found",
			ErrorMessage: "Message is not found",
		},
	},
	services.ErrReactionNotFound: {
		Code: http.StatusNotFound,
		Body: restapi.ErrorResponse{
			ErrorType:    "reaction_not_found",
			ErrorMessage: "Reaction is not found",
		},
	},
	services.ErrSecretUpdateNotFound: {
		Code: http.StatusNotFound,
		Body: restapi.ErrorResponse{
			ErrorType:    "secret_update_not_found",
			ErrorMessage: "Secret update is not found",
		},
	},
}

var domainErrMap = map[domain.Error]Response{
	domain.ErrAdminNotMember: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "admin_not_member",
			ErrorMessage: "Admin is not a member of a group",
		},
	},
	domain.ErrAlreadyBlocked: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "already_blocked",
			ErrorMessage: "Already blocked",
		},
	},
	domain.ErrAlreadyUnblocked: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "already_unblocked",
			ErrorMessage: "Already unblocked",
		},
	},
	domain.ErrChatBlocked: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "chat_blocked",
			ErrorMessage: "Chat is blocked",
		},
	},
	domain.ErrChatWithMyself: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "chat_with_myself",
			ErrorMessage: "Cannot create a chat with myself",
		},
	},
	domain.ErrFileTooBig: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "file_too_big",
			ErrorMessage: "File is too big",
		},
	},
	domain.ErrGroupDescTooLong: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "group_description_too_long",
			ErrorMessage: "Description is too long",
		},
	},
	domain.ErrGroupNameEmpty: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "group_name_empty",
			ErrorMessage: "Group name is empty",
		},
	},
	domain.ErrGroupNameTooLong: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "group_name_too_long",
			ErrorMessage: "Group name is too long",
		},
	},
	domain.ErrGroupPhotoEmpty: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "group_photo_empty",
			ErrorMessage: "Group photo is empty",
		},
	},
	domain.ErrMemberIsAdmin: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "member_is_admin",
			ErrorMessage: "Member is admin",
		},
	},
	domain.ErrReactionNotFromUser: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "reaction_not_from_user",
			ErrorMessage: "Reaction is not from this user",
		},
	},
	domain.ErrSenderNotAdmin: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "sender_not_admin",
			ErrorMessage: "Sender is not admin",
		},
	},
	domain.ErrTextEmpty: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "text_empty",
			ErrorMessage: "Text is empty",
		},
	},
	domain.ErrTooManyTextRunes: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "too_many_characters",
			ErrorMessage: "Too many characters",
		},
	},
	domain.ErrUpdateDeleted: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "update_deleted",
			ErrorMessage: "Update is deleted",
		},
	},
	domain.ErrUpdateNotFromChat: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "update_not_from_chat",
			ErrorMessage: "Update is not from chat",
		},
	},
	domain.ErrUserAlreadyMember: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "user_already_member",
			ErrorMessage: "User ia already a member",
		},
	},
	domain.ErrUserNotMember: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "user_not_member",
			ErrorMessage: "User is not a member",
		},
	},
	domain.ErrUserNotSender: {
		Code: http.StatusBadRequest,
		Body: restapi.ErrorResponse{
			ErrorType:    "user_not_sender",
			ErrorMessage: "User is not a sender",
		},
	},
}
