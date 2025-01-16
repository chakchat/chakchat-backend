// CAUTION: this file is only sketch
package model

import "github.com/google/uuid"

type OridinaryChat struct {
	ChatID        uuid.UUID `gorm:"primaryKey"`
	LeastMemberID uuid.UUID
	MostMemberID  uuid.UUID
	BlockedBy     uuid.UUID
}

type SecretChat struct {
	ChatID        uuid.UUID `gorm:"primaryKey"`
	LeastMemberID uuid.UUID
	MostMemberID  uuid.UUID
}

type Group struct {
	GroupID uuid.UUID
	Name    string
	AdminID uuid.UUID
}

type SecretGroup struct {
	GroupID uuid.UUID
	Name    string
	AdminID uuid.UUID
}

// type TextMessage struct {
// 	chatUpdate
// 	MessageID uint64

// 	Text string

// 	SentAtTimestamp     int64
// 	ModifiedAtTimestamp int64
// 	DeletedAtTimestamp  int64
// }

// type Reaction struct {
// 	chatUpdate
// 	ReactionID uint64

// 	EmojiCode          int64
// 	SentAtTimestamp    int64
// 	DeletedAtTimestamp int64
// }

// type chatUpdate struct {
// 	SenderID uuid.UUID
// 	ChatID   uuid.UUID
// }
