package minter

type ChanMinter struct {
	Message string
}

func NewMessage(text string) *ChanMinter {
	return &ChanMinter{
		Message: text,
	}
}
