package mailer

import (
	"fmt"

	"github.com/praveennagaraj97/online-consultation/pkg/env"
)

type EmailMetaInfoTemplateData struct {
	MetaTitle string
}

type ConfirmEmailTemplateData struct {
	MetaInfo               *EmailMetaInfoTemplateData
	Title                  string
	Message                string
	VerifyLink             string
	FooterMessage          string
	NotificationPreviewMsg string
	ButtonName             string
}

func GetRegisterEmailTemplateData(name string) *ConfirmEmailTemplateData {

	return &ConfirmEmailTemplateData{
		MetaInfo: &EmailMetaInfoTemplateData{
			MetaTitle: "Online Consultation | Welcome",
		},
		NotificationPreviewMsg: "Welcome to Online Consultation, Confirm your email address",
		Title:                  fmt.Sprintf("Hi %s, Welcome to Online Consultation", name),
		Message: `Tap the button below to confirm your email address. If you
                  didn't create an account with
                  Online Consultation, you can safely delete
                  this email.`,
		VerifyLink:    env.GetEnvVariable("CLIENT_VERIFY_EMAIL_LINK"),
		FooterMessage: `You received this email because you registered to Online consultation.`,
		ButtonName:    "Confirm Email Address",
	}

}

func GetVerifyEmailTemplateData(name, verifyLink string) *ConfirmEmailTemplateData {
	return &ConfirmEmailTemplateData{
		MetaInfo: &EmailMetaInfoTemplateData{
			MetaTitle: "Online Consultation | Verify Email Address",
		},
		NotificationPreviewMsg: "Verify your email address.",
		Title:                  fmt.Sprintf("Hi %s,\n", name),
		Message: `Tap the button below to verify your email address. If you
                  didn't request to verify your email address for
                  Online Consultation, you can safely delete
                  this email.`,
		VerifyLink: verifyLink,
		FooterMessage: `You received this email because we received a request for
                  email verify from your account.`,
		ButtonName: "Verify Now",
	}
}
