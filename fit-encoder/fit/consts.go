package fit

import (
	"errors"

	"github.com/renbou/jogmock/fit-encoder/fit/types"
)

const (
	FIT_MESG_NUM_FILE_ID     types.FitUint16 = 0
	FIT_MESG_NUM_SESSION     types.FitUint16 = 18
	FIT_MESG_NUM_LAP         types.FitUint16 = 19
	FIT_MESG_NUM_RECORD      types.FitUint16 = 20
	FIT_MESG_NUM_EVENT       types.FitUint16 = 21
	FIT_MESG_NUM_DEVICE_INFO types.FitUint16 = 23
	FIT_MESG_NUM_ACTIVITY    types.FitUint16 = 34
	FIT_MESG_NUM_FIELD_DESC  types.FitUint16 = 206
	FIT_MESG_NUM_DEV_DATA_ID types.FitUint16 = 207
)

var (
	ErrInvalidLocalMsgType = errors.New("invalid local message type (> 15)")
	ErrInvalidMsgType      = errors.New("invalid message type (> 1)")
	ErrInvalidMsgSpecific  = errors.New("only definition message can have msg specific set")
	ErrInvalidMessage      = errors.New("message is not a valid definition or data message type")
	ErrFieldNumMismatch    = errors.New("unexpected number of fields")
)
