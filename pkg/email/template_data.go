package mailer

import (
	"fmt"
)

type EmailMetaInfoTemplateData struct {
	MetaTitle string
}

type ConfirmEmailTemplateData struct {
	MetaInfo               *EmailMetaInfoTemplateData
	Title                  string
	Message                string
	Href                   string
	FooterMessage          string
	NotificationPreviewMsg string
	ButtonName             string
}

func GetRegisterEmailTemplateData(name, verifyLink string) *ConfirmEmailTemplateData {

	return &ConfirmEmailTemplateData{
		MetaInfo: &EmailMetaInfoTemplateData{
			MetaTitle: "Online Consultation | Welcome",
		},
		NotificationPreviewMsg: "Welcome to Online Consultation, Confirm your email address",
		Title:                  fmt.Sprintf("Hi %s, Welcome to Online Consultation", name),
		Message: `Tap the below button to confirm your email address. If you
                  didn't create an account with
                  Online Consultation, you can safely delete
                  this email.`,
		Href:          verifyLink,
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
		Message: `Tap the below button to verify your email address. If you
                  didn't request to verify your email address for
                  Online Consultation, you can safely delete
                  this email.`,
		Href: verifyLink,
		FooterMessage: `You received this email because we received a request for
                  email verify from your account.`,
		ButtonName: "Verify Now",
	}
}

func GetSignWithEmailLinkTemplateData(name, link string) *ConfirmEmailTemplateData {
	return &ConfirmEmailTemplateData{
		MetaInfo: &EmailMetaInfoTemplateData{
			MetaTitle: "Online Consultation | Sign With Email Link",
		},
		NotificationPreviewMsg: "Click Below to sign in to online consultation",
		Title:                  fmt.Sprintf("Hi %s,\n", name),
		Message: `Tap the below button to sign in to your account. If you
                  didn't request to sign in for
                  Online Consultation, you can safely delete
                  this email.`,
		Href: link,
		FooterMessage: `You received this email because we received a request for
                  sign in from your account.`,
		ButtonName: "Sign In",
	}
}
