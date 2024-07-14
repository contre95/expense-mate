package managing

import (
	"expenses-app/pkg/app"
)

type TelegramCommander struct {
	logger   app.Logger
	sends    chan string
	received chan string
}

type CommandResp struct {
	Msg string
}

func NewTelegramCommander(l app.Logger, s, r chan string) *TelegramCommander {
	return &TelegramCommander{l, s, r}
}

func (tc *TelegramCommander) Command(cmd string) (*CommandResp, error) {
	tc.sends <- cmd
	tc.logger.Debug("Command sent: %s", cmd)
	return &CommandResp{Msg: "Command sent"}, nil
}

func (tc *TelegramCommander) CommandWithResponse(cmd string) (*CommandResp, error) {
	tc.sends <- cmd
	tc.logger.Debug("Command sent: %s", cmd)
	answser := <-tc.received // Waits until it gets a messages
	return &CommandResp{Msg: answser}, nil
}
