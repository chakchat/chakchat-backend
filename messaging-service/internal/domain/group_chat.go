package domain

import "errors"

const (
	maxGroupNameLen   = 50
	maxDescriptionLen = 300
)

var (
	ErrGroupNameEmpty   = errors.New("group name is empty")
	ErrGroupNameTooLong = errors.New("group name is too long")
	ErrGroupDescTooLong = errors.New("group description is too long")
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

func (g *GroupChat) UpdateInfo(name, description string) error {
	if err := validateGroupInfo(name, description); err != nil {
		return err
	}

	g.Name = name
	g.Description = description
	return nil
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
