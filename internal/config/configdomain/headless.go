package configdomain

// Headless indicates whether the propose command should skip opening a browser.
type Headless bool

func (self Headless) Enabled() bool {
	return bool(self)
}
