package users

import "github.com/jinzhu/gorm"

type Manager struct {
	*gorm.DB `inject:""`
}

func (m *Manager) FindOrCreateUserForTwitterID(twitterID string) (*User, error) {
	user := &User{}
	m.Where("twitter_id = ?", twitterID).First(user)

	// no user in DB
	if user.TwitterID == "" {
		// create user
		user.TwitterID = twitterID
		m.Create(user)
	}

	return user, nil
}
