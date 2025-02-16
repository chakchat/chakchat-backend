package secgroup

import (
	"slices"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type SecretGroupChat struct {
	domain.SecretChat
	Admin   domain.UserID
	Members []domain.UserID

	Name        string
	Description string
	GroupPhoto  domain.URL
}

func NewSecretGroupChat(admin domain.UserID, members []domain.UserID, name string) (*SecretGroupChat, error) {
	if err := domain.ValidateGroupInfo(name, ""); err != nil {
		return nil, err
	}

	if !slices.Contains(members, admin) {
		return nil, domain.ErrAdminNotMember
	}

	normMembers := domain.NormilizeMembers(members)

	return &SecretGroupChat{
		SecretChat: domain.SecretChat{
			Chat: domain.Chat{
				ID: domain.NewChatID(),
			},
		},
		Admin:       admin,
		Members:     normMembers,
		Name:        name,
		Description: "",
		GroupPhoto:  "",
	}, nil
}

func (g *SecretGroupChat) Delete(sender domain.UserID) error {
	if sender != g.Admin {
		return domain.ErrSenderNotAdmin
	}
	return nil
}

func (g *SecretGroupChat) UpdateInfo(sender domain.UserID, name, description string) error {
	if sender != g.Admin {
		return domain.ErrSenderNotAdmin
	}

	if err := domain.ValidateGroupInfo(name, description); err != nil {
		return err
	}

	g.Name = name
	g.Description = description
	return nil
}

func (g *SecretGroupChat) UpdatePhoto(sender domain.UserID, photo domain.URL) error {
	if sender != g.Admin {
		return domain.ErrSenderNotAdmin
	}

	g.GroupPhoto = photo
	return nil
}

func (g *SecretGroupChat) DeletePhoto(sender domain.UserID) error {
	if sender != g.Admin {
		return domain.ErrSenderNotAdmin
	}

	if g.GroupPhoto == "" {
		return domain.ErrGroupPhotoEmpty
	}

	g.GroupPhoto = ""
	return nil
}

func (g *SecretGroupChat) AddMember(sender domain.UserID, newMember domain.UserID) error {
	if sender != g.Admin {
		return domain.ErrSenderNotAdmin
	}

	if g.IsMember(newMember) {
		return domain.ErrUserAlreadyMember
	}

	g.Members = append(g.Members, newMember)
	return nil
}

func (g *SecretGroupChat) DeleteMember(sender domain.UserID, member domain.UserID) error {
	if sender != g.Admin {
		return domain.ErrSenderNotAdmin
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

func (c *SecretGroupChat) SetExpiration(sender domain.UserID, exp *time.Duration) error {
	if c.Admin != sender {
		return domain.ErrSenderNotAdmin
	}
	c.Exp = exp
	return nil
}

func (g *SecretGroupChat) IsMember(user domain.UserID) bool {
	return slices.Contains(g.Members, user)
}

func (g *SecretGroupChat) ValidateCanSend(sender domain.UserID) error {
	if !g.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	return nil
}
