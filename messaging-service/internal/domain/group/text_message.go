package group

import "github.com/chakchat/chakchat-backend/messaging-service/internal/domain"

func (g *GroupChat) NewTexMessage(sender domain.UserID, text string, replyTo *Message)
