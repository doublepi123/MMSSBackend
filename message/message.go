package message

type Message struct {
	Msg string
}

const successMsg = "Success"

func Success() Message {
	return Message{
		Msg: successMsg,
	}
}
