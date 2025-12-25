package configdomain

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// UseMessage indicates how to set commit message.
// The following variants are supported:
// - EditDefaultMessage
// - UseCustomMessage(CommitMessage)
// - UseDefaultMessage
//
// See `--message`, `--edit` and `--no-edit` flags in `git commit` and `git merge`.
type UseMessage struct {
	customMessage      Option[gitdomain.CommitMessage]
	editDefaultMessage bool
}

// EditDefaultMessage indicates that the default message should be edited.
func EditDefaultMessage() UseMessage {
	return UseMessage{
		customMessage:      None[gitdomain.CommitMessage](),
		editDefaultMessage: true,
	}
}

// UseCustomMessage indicates that the custom message should be used as is without opening an editor.
func UseCustomMessage(message gitdomain.CommitMessage) UseMessage {
	return UseMessage{
		customMessage:      Some(message),
		editDefaultMessage: false, // This value is unused.
	}
}

// UseCustomMessageOr returns UseCustomMessage(message) if message is Some, or other.
func UseCustomMessageOr(message Option[gitdomain.CommitMessage], other UseMessage) UseMessage {
	if m, has := message.Get(); has {
		return UseCustomMessage(m)
	}
	return other
}

// UseDefaultMessage indicates that the default message should be used as is without opening an editor.
func UseDefaultMessage() UseMessage {
	return UseMessage{
		customMessage:      None[gitdomain.CommitMessage](),
		editDefaultMessage: false,
	}
}

// UseMessageWithFallbackToDefault returns UseMessage based on an optional message and fallbackToDefault.
// - If message is Some(...) then this is equivalent to UseCustomMessage(...).
// - If message is None and fallbackToDefault then this is equivalent to UseDefaultMessage().
// - If message is None and !fallbackToDefault then this is equivalent to EditDefaultMessage().
func UseMessageWithFallbackToDefault(message Option[gitdomain.CommitMessage], fallbackToDefault bool) UseMessage {
	return UseMessage{
		customMessage:      message,
		editDefaultMessage: !fallbackToDefault,
	}
}

// GetCustomMessage returns a copy of the custom message and IsCustomMessage().
func (self *UseMessage) GetCustomMessage() (gitdomain.CommitMessage, bool) {
	// editDefaultMessage is irrelevant when a custom message is set.
	return self.customMessage.Get()
}

// GetCustomMessageOrPanic returns a copy of the custom message if IsCustomMessage() or panics.
func (self *UseMessage) GetCustomMessageOrPanic() gitdomain.CommitMessage {
	// editDefaultMessage is irrelevant when a custom message is set.
	if message, is := self.customMessage.Get(); is {
		return message
	}
	panic("UseMessage is not UseCustomMessage")
}

// IsCustomMessage indicates that the custom message should be used as is without opening an editor.
func (self *UseMessage) IsCustomMessage() bool {
	// editDefaultMessage is irrelevant when a custom message is set.
	return self.customMessage.IsSome()
}

// IsEditDefault indicates that the default message should be edited.
func (self *UseMessage) IsEditDefault() bool {
	if self.customMessage.IsSome() {
		return false
	}
	return self.editDefaultMessage
}

// IsUseDefault indicates that the default message should be used as is without opening an editor.
func (self *UseMessage) IsUseDefault() bool {
	if self.customMessage.IsSome() {
		return false
	}
	return !self.editDefaultMessage
}
