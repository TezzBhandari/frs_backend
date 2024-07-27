package utils

import "fmt"

func InvalidJsonMsg() string {
	return "invalid json body"
}

func FailedResponseMsg() string {
	return "failed to write response"
}

func InvalidUserIdMsg() string {
	return "invalid user id"
}

func InvalidFundRaiserIdMsg() string {
	return "invalid fund raiser id"
}

func DoesNotExistMsg(v string) string {
	return fmt.Sprintf("%s does not exist", v)
}
