package message

type Message struct {
	Code int
	Msg  string
}

const successCode = 0
const successMsg = "Success"
const failCode = 1
const failMsg = "Fail"

func Success() Message {
	return Message{
		Code: successCode,
		Msg:  successMsg,
	}
}

func Fail() Message {
	return Message{
		Code: failCode,
		Msg:  failMsg,
	}
}
