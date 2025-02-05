package group

import (
	"errors"
	"slices"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

const (
	maxGroupNameLen   = 50
	maxDescriptionLen = 300
)

var (
	ErrAdminNotMember = errors.New("group members doesn't include admin")

	ErrGroupNameEmpty   = errors.New("group name is empty")
	ErrGroupNameTooLong = errors.New("group name is too long")
	ErrGroupDescTooLong = errors.New("group description is too long")

	ErrUserAlreadyMember = errors.New("user is already a member of a chat")
	ErrMemberIsAdmin     = errors.New("group member is admin")

	ErrGroupPhotoEmpty = errors.New("group photo is empty")
)

type URL string

type GroupChat struct {
	domain.Chat
	Admin   domain.UserID
	Members []domain.UserID

	Name        string
	Description string
	GroupPhoto  URL
}

func NewGroupChat(admin domain.UserID, members []domain.UserID, name string) (*GroupChat, error) {
	if err := validateGroupInfo(name, ""); err != nil {
		return nil, err
	}

	if !slices.Contains(members, admin) {
		return nil, ErrAdminNotMember
	}

	normMembers := normilizeMembers(members)

	return &GroupChat{
		Chat: domain.Chat{
			ID: domain.NewChatID(),
		},
		Admin:       admin,
		Members:     normMembers,
		Name:        name,
		Description: "",
		GroupPhoto:  "",
	}, nil
}

func (g *GroupChat) UpdateInfo(name, description string) error {
	if err := validateGroupInfo(name, description); err != nil {
		return err
	}

	g.Name = name
	g.Description = description
	return nil
}

func (g *GroupChat) UpdatePhoto(photo URL) error {
	g.GroupPhoto = photo
	return nil
}

func (g *GroupChat) DeletePhoto() error {
	if g.GroupPhoto == "" {
		return ErrGroupPhotoEmpty
	}

	g.GroupPhoto = ""
	return nil
}

func (g *GroupChat) AddMember(newMember domain.UserID) error {
	if g.IsMember(newMember) {
		return ErrUserAlreadyMember
	}

	g.Members = append(g.Members, newMember)
	return nil
}

func (g *GroupChat) DeleteMember(member domain.UserID) error {
	if g.Admin == member {
		return ErrMemberIsAdmin
	}

	i := slices.Index(g.Members, member)
	if i == -1 {
		return domain.ErrUserNotMember
	}

	g.Members = slices.Delete(g.Members, i, i+1)
	return nil
}

func (g *GroupChat) IsMember(user domain.UserID) bool {
	return slices.Contains(g.Members, user)
}

func (g *GroupChat) ChatID() domain.ChatID {
	return g.Chat.ID
}

func (g *GroupChat) ValidateCanSend(sender domain.UserID) error {
	if !g.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	return nil
}

func normilizeMembers(members []domain.UserID) []domain.UserID {
	met := make(map[domain.UserID]struct{}, len(members))
	normMembers := make([]domain.UserID, 0, len(members))

	for _, member := range members {
		if _, ok := met[member]; !ok {
			normMembers = append(normMembers, member)
			met[member] = struct{}{}
		}
	}

	return normMembers
}

func validateGroupInfo(name, description string) error {
	var errs []error
	if name == "" {
		errs = append(errs, ErrGroupNameEmpty)
	}
	if len(name) > maxGroupNameLen {
		errs = append(errs, ErrGroupNameTooLong)
	}
	if len(description) > maxDescriptionLen {
		errs = append(errs, ErrGroupDescTooLong)
	}
	return errors.Join(errs...)
}
