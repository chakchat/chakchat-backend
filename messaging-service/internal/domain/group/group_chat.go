package group

import (
	"slices"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type GroupChat struct {
	domain.Chat
	Admin   domain.UserID
	Members []domain.UserID

	Name        string
	Description string
	GroupPhoto  domain.URL
}

func NewGroupChat(admin domain.UserID, members []domain.UserID, name string) (*GroupChat, error) {
	if err := domain.ValidateGroupInfo(name, ""); err != nil {
		return nil, err
	}

	if !slices.Contains(members, admin) {
		return nil, domain.ErrAdminNotMember
	}

	normMembers := domain.NormilizeMembers(members)

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

func (g *GroupChat) Delete(sender domain.UserID) error {
	if sender != g.Admin {
		return domain.ErrNotAdmin
	}
	return nil
}

func (g *GroupChat) UpdateInfo(sender domain.UserID, name, description string) error {
	if sender != g.Admin {
		return domain.ErrNotAdmin
	}

	if err := domain.ValidateGroupInfo(name, description); err != nil {
		return err
	}

	g.Name = name
	g.Description = description
	return nil
}

func (g *GroupChat) UpdatePhoto(sender domain.UserID, photo domain.URL) error {
	if sender != g.Admin {
		return domain.ErrNotAdmin
	}

	g.GroupPhoto = photo
	return nil
}

func (g *GroupChat) DeletePhoto(sender domain.UserID) error {
	if sender != g.Admin {
		return domain.ErrNotAdmin
	}

	if g.GroupPhoto == "" {
		return domain.ErrGroupPhotoEmpty
	}

	g.GroupPhoto = ""
	return nil
}

func (g *GroupChat) AddMember(sender domain.UserID, newMember domain.UserID) error {
	if sender != g.Admin {
		return domain.ErrNotAdmin
	}

	if g.IsMember(newMember) {
		return domain.ErrUserAlreadyMember
	}

	g.Members = append(g.Members, newMember)
	return nil
}

func (g *GroupChat) DeleteMember(sender domain.UserID, member domain.UserID) error {
	if sender != g.Admin {
		return domain.ErrNotAdmin
	}

	if g.Admin == member {
		return domain.ErrMemberIsAdmin
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

func (g *GroupChat) ValidateCanSend(sender domain.UserID) error {
	if !g.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	return nil
}
