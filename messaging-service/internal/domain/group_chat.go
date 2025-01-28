package domain

import (
	"errors"
	"slices"

	"github.com/google/uuid"
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

type GroupChat struct {
	ID      ChatID
	Admin   UserID
	Members []UserID

	Secret      bool
	Name        string
	Description string
	GroupPhoto  URL
	CreatedAt   Timestamp
}

func NewGroupChat(admin UserID, members []UserID, name string) (*GroupChat, error) {
	return newGroup(admin, members, name, false)
}

func NewSecretGroupChat(admin UserID, members []UserID, name string) (*GroupChat, error) {
	return newGroup(admin, members, name, true)
}

func newGroup(admin UserID, members []UserID, name string, secret bool) (*GroupChat, error) {
	if err := validateGroupInfo(name, ""); err != nil {
		return nil, err
	}

	if !slices.Contains(members, admin) {
		return nil, ErrAdminNotMember
	}

	normMembers := normilizeMembers(members)

	return &GroupChat{
		ID:          ChatID(uuid.New()),
		Admin:       admin,
		Members:     normMembers,
		Secret:      secret,
		Name:        name,
		Description: "",
		GroupPhoto:  "",
		// Maybe this should not be set here
		CreatedAt: Timestamp(TimeFunc().Unix()),
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

func (g *GroupChat) AddMember(newMember UserID) error {
	if g.IsMember(newMember) {
		return ErrUserAlreadyMember
	}

	g.Members = append(g.Members, newMember)
	return nil
}

func (g *GroupChat) DeleteMember(member UserID) error {
	if g.Admin == member {
		return ErrMemberIsAdmin
	}

	i := slices.Index(g.Members, member)
	if i == -1 {
		return ErrUserNotMember
	}

	g.Members = slices.Delete(g.Members, i, i+1)
	return nil
}

func (g *GroupChat) IsMember(user UserID) bool {
	return slices.Contains(g.Members, user)
}

func normilizeMembers(members []UserID) []UserID {
	met := make(map[UserID]struct{}, len(members))
	normMembers := make([]UserID, 0, len(members))

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
