package locales

func init() {
	Locales["en"] = map[string]string{
		"reader.wrong_captcha":           "Wrong captcha!",
		"reader.email_subject_send_code": "[{title}] Reader Login Validation",
		"reader.email_body_send_code":    "Hi, your validate code is {code}, and the expiration time is 10 minutes.",
		"reader.email_send_failed":       "Email sent failed, please try again latter.",
		"reader.wrong_validate_code":     "The email validate code is incorrect or expired.",
		"reader.inactive":                "Your account is inactive, please contact the website administrator.",
		"reader.exists_email":            "This email is existed.",
		"reader.exists_nickname":         "This nickname is existed.",
	}
}
