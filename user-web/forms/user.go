package forms

type LoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password  string `form:"password" json:"password" binding:"required,min=5,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=4,max=4"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}
