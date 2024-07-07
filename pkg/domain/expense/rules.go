package expense

import "regexp"

type Rule struct {
	ID         string     `validate:"required"`
	Pattern    string     `validate:"required"`
	CategoryID CategoryID `validate:"required"`
	UsersID    []UserID   // This is new
}

// Rules is the rules repository.
type Rules interface {
	All() ([]Rule, error)
	Add(Rule) error
	Delete(id string) error
}

func (r *Rule) Matches(s string) bool {
	matched, err := regexp.MatchString(r.Pattern, s)
	if err != nil {
		// Handle the error appropriately, e.g., log it
		return false
	}
	return matched
}
