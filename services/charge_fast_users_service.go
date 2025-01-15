package services

import (
	"strconv"
	"time"
)

// ChargeFastUsers realiza la creación masiva de usuarios y asignación de roles con auditoría.
func ChargeFastUsers(users []struct {
	Name     string
	Email    string
	Password string
	Active   bool
	RoleID   uint
}, userID uint) error {
	currentTime := time.Now()

	for _, userInput := range users {
		// Crear usuario (reutilizando `CreateUser` de user_service.go)
		user, err := CreateUser(userInput.Name, userInput.Email, userInput.Password, userInput.Active)
		if err != nil {
			return err
		}

		// Registrar auditoría para creación de usuario
		err = RegisterAudit("INSERT", "Se creó un usuario con email: "+userInput.Email, userID, "SEGURIDAD", currentTime)
		if err != nil {
			return err
		}

		// Asignar rol al usuario (reutilizando `AssignRoleToUser` de user_role.go)
		_, err = AssignRoleToUser(user.ID, userInput.RoleID)
		if err != nil {
			return err
		}

		// Registrar auditoría para asignación de rol
		err = RegisterAudit("INSERT", "Se asignó el rol "+strconv.Itoa(int(userInput.RoleID))+" al usuario con email: "+userInput.Email, userID, "SEGURIDAD", currentTime)
		if err != nil {
			return err
		}
	}

	return nil
}
