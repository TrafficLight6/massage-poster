package main

import (
	"fmt"
	"net/smtp"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Username string `yaml:"emailUsername"`
	Password string `yaml:"emailPassword"`
	SmtpHost string `yaml:"smtpHost"`
	SmtpPort string `yaml:"smtpPort"`
}

func sendEmailMassage(from string, to []string, subject string, text string, config Config) bool {
	emailInterface := email.NewEmail()
	emailInterface.From = from
	emailInterface.To = to
	emailInterface.Subject = subject
	emailInterface.Text = []byte(text)
	err := emailInterface.Send(config.SmtpHost+":"+config.SmtpPort, smtp.PlainAuth("", config.Username, config.Password, config.SmtpHost))
	if err != nil {
		fmt.Println("Error sending email", err)
		return false
	}
	return true
}

func getConfig(filePast string) (bool, Config) {
	file, err := os.Open(filePast)
	var config Config

	if err != nil {
		fmt.Println("Error opening file:", err)
		return false, config
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding YAML:", err)
		return false, config
	}

	return true, config
}

func main() {
	success, configData := getConfig("./config.yaml")
	if success == false {
		fmt.Println("Config reading error,shutdowned")
		return
	}
	r := gin.Default()
	r.POST("/send", func(c *gin.Context) {
		from := c.Query("from")
		to := []string{c.Query("to")}
		subject := c.Query("subject")
		text := c.Query("text")
		success = sendEmailMassage(from, to, subject, text, configData)
		if success == false {
			c.JSON(200, gin.H{
				"code":    "4",
				"message": "sending error",
			})
		} else {
			c.JSON(200, gin.H{
				"code":    "2",
				"message": "success",
			})
		}
	})
	r.Run(":7777")
}
