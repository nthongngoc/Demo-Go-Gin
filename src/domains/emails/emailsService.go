package emails

// import (
// 	"github.com/khoa5773/go-server/src/configs"
// 	"github.com/sendgrid/sendgrid-go"
// 	"github.com/sendgrid/sendgrid-go/helpers/mail"
// )

// func SendEmails(to mail.Email, subject, content string) (bool, error) {
// 	from := mail.NewEmail("Ngoc", "ngoc.nthongngoc@gmail.com")
// 	message := mail.NewSingleEmail(from, subject, &to, "", content)
// 	client := sendgrid.NewSendClient(configs.ConfigsService)
// 	response, err := client.Send(message)
// 	if err != nil || response.StatusCode != 200 {
// 		return false, err
// 	}
// 	return true, nil
// }
