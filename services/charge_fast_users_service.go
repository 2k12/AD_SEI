package services

import (
	"strconv"
	"time"
)

func ChargeFastUsers(users []struct {
	Name     string
	Email    string
	Password string
	Active   bool
	RoleID   uint
}, userID uint) error {
	currentTime := time.Now()

	for _, userInput := range users {
		user, err := CreateUser(userInput.Name, userInput.Email, userInput.Password, userInput.Active)
		if err != nil {
			return err
		}

		err = RegisterAudit("INSERT", "Se creó un usuario con email: "+userInput.Email, userID, "SEGURIDAD", currentTime)
		if err != nil {
			return err
		}

		_, err = AssignRoleToUser(user.ID, userInput.RoleID)
		if err != nil {
			return err
		}

		err = RegisterAudit("INSERT", "Se asignó el rol "+strconv.Itoa(int(userInput.RoleID))+" al usuario con email: "+userInput.Email, userID, "SEGURIDAD", currentTime)
		if err != nil {
			return err
		}
	}

	return nil
}
