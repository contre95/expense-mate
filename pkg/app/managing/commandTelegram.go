package managing

import (
	"expenses-app/pkg/app"
)

type TelegramCommander struct {
	logger app.Logger
	cmds   chan string
}

type CommandResp struct {
	Msg string
}

func NewTelegramCommander(l app.Logger, c chan string) *TelegramCommander {
	return &TelegramCommander{l, c}
}

func (tc *TelegramCommander) Command(cmd string) (*CommandResp, error) {
	tc.cmds <- cmd
	tc.logger.Debug("Command sent: %s", cmd)
	received := <-tc.cmds
	tc.logger.Debug("Received: %s", received)
	return &CommandResp{Msg: received}, nil
}
