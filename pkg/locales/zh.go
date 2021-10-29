package locales

func init() {
	Locales["zh"] = map[string]string{
		"reader.wrong_captcha":           "验证码错误！",
		"reader.email_subject_send_code": "[{title}] 读者登录验证",
		"reader.email_body_send_code":    "您好，您的验证码为{code}，验证码的有效期为10分钟。",
		"reader.email_send_failed":       "邮件发送失败，请稍后再试。",
		"reader.wrong_validate_code":     "邮件验证码错误或已失效。",
		"reader.inactive":                "您的账号被禁用了，请联系网站管理员。",
		"reader.exists_email":            "该邮件已存在。",
		"reader.exists_nickname":         "该昵称已存在。",
	}
}
