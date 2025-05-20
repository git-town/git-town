package configdomain

type UpdateOutdatedSettings bool

func (self UpdateOutdatedSettings) IsTrue() bool {
	return bool(self)
}

const (
	UpdateOutdatedYes UpdateOutdatedSettings = true
	UpdateOutdatedNo  UpdateOutdatedSettings = false
)
