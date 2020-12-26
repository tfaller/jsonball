package operation

import (
	"context"
	"errors"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/propchange"
)

const (
	// AdminCmdRegisterDocType is used to register a new document type
	AdminCmdRegisterDocType = "regDocType"

	// AdminCmdRegisterHandler is used to register a new handler
	AdminCmdRegisterHandler = "regHandler"

	// AdminCmdRequeueHandler is used to requeue all listeners of a handler
	AdminCmdRequeueHandler = "requeueHandler"
)

var (
	// ErrInvalidCommand indicates that the given command is invalid.
	ErrInvalidCommand = errors.New("given command is not a valid admin command")

	// ErrCommandDataMissing indicates that either the complete command is missing
	// or that some data of the command is missing.
	ErrCommandDataMissing = errors.New("command data missing")
)

// AdminCommands executes simple administrative commands
func AdminCommands(ctx context.Context, registry jsonball.Registry, detector propchange.Detector, adminCmd *event.AdminCmd) error {
	if adminCmd == nil {
		return ErrCommandDataMissing
	}

	switch adminCmd.Cmd {

	case AdminCmdRegisterDocType:
		if adminCmd.RegisterDocType == nil {
			return ErrCommandDataMissing
		}
		return registerDocumentType(ctx, registry, adminCmd.RegisterDocType.Type, false)

	case AdminCmdRegisterHandler:
		if adminCmd.RegisterHandler == nil {
			return ErrCommandDataMissing
		}
		return RegisterHandler(ctx, registry, detector, adminCmd.RegisterHandler)

	case AdminCmdRequeueHandler:
		if adminCmd.RequeueHandler == nil {
			return ErrCommandDataMissing
		}
		return RequeueHandler(ctx, detector, adminCmd.RequeueHandler.Handler)
	}

	return ErrInvalidCommand
}
