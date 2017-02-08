package users

import "github.com/jinzhu/gorm"

type Manager struct {
	*gorm.DB `inject:""`
}

func (m *Manager) FindOrCreateUser(userAuthID string) (*User, error) {
	user := &User{}
	m.Where("auth_id = ?", userAuthID).First(user)

	// no user in DB
	if user.AuthID == "" {
		// create user
		user.AuthID = userAuthID
		user.IsAdmin = false
		m.Create(user)
	}

	return user, nil
}
