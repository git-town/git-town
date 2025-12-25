package configdomain

type UpdateOutdatedSettings bool

func (self UpdateOutdatedSettings) ShouldUpdateOutdatedSettings() bool {
	return bool(self)
}

const (
	UpdateOutdatedYes UpdateOutdatedSettings = true
	UpdateOutdatedNo  UpdateOutdatedSettings = false
)
