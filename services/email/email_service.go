package services

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(to, subject, body string) error {
	// Crear el objeto "From" para el remitente
	from := mail.NewEmail("Example User", "jeipige@gmail.com")
	// Crear el objeto "To" para el destinatario
	toEmail := mail.NewEmail("Recipient", to)
	// El contenido del correo en texto plano y HTML
	plainTextContent := body
	htmlContent := "<strong>" + body + "</strong>"

	// Crear el mensaje utilizando la API de SendGrid
	message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)

	// Crear el cliente de SendGrid con la clave API
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	// Enviar el mensaje y obtener la respuesta
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Error al enviar correo: %v", err)
		return err
	}

	// Loguear la respuesta de SendGrid
	log.Printf("Correo enviado con éxito: %v", response.StatusCode)
	log.Printf("Cuerpo de la respuesta: %v", response.Body)
	log.Printf("Cabeceras de la respuesta: %v", response.Headers)

	return nil
}
