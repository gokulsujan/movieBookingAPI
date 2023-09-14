package controller

import (
	"crypto/rand"
	"math/big"
	"net/smtp"
	"os"
	"strconv"
)

// Getting the otp and sending the otp to the user
func GetOTP(name, email string) string {
	otp, err := getRandNum()
	if err != nil {
		panic(err)
	}
	sendOtp(name, otp, email)
	return otp
}

// Getting a random number for otp. This function helps get otp to generate the a random otp
func getRandNum() (string, error) {
	otp, err := rand.Int(rand.Reader, big.NewInt(8999))
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(otp.Int64()+1000, 10), nil
}

// sending generated otp to the user mail using smtp
func sendOtp(name, otp, email string) {
	SMTPemail := os.Getenv("Email")
	pass := os.Getenv("pass")
	auth := smtp.PlainAuth("", SMTPemail, pass, "smtp.gmail.com")
	msg := "Subject: WebPortal OTP\nHey " + name + "Your OTP is " + otp
	err := smtp.SendMail("smtp.gmail.com:587", auth, SMTPemail, []string{email}, []byte(msg))
	if err != nil {
		panic(err)
	}
}