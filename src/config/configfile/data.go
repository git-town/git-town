package configfile

// Data defines the Go equivalent of the TOML file content.
type Data struct {
	Branches                 *Branches     `toml:"branches"`
	CodeHosting              *CodeHosting  `toml:"code-hosting"`
	SyncStrategy             *SyncStrategy `toml:"sync-strategy"`
	PushHook                 *bool         `toml:"push-hook"`
	PushNewbranches          *bool         `toml:"push-new-branches"`
	ShipDeleteTrackingBranch *bool         `toml:"ship-delete-remote-branch"`
	SyncBeforeShip           *bool         `toml:"sync-before-ship"`
	SyncUpstream             *bool         `toml:"sync-upstream"`
}

type Branches struct {
	Main       *string  `toml:"main"`
	Perennials []string `toml:"perennials"`
}

func (self Branches) IsEmpty() bool {
	return self.Main == nil && len(self.Perennials) == 0
}

type CodeHosting struct {
	Platform       *string `toml:"platform"`
	OriginHostname *string `toml:"origin-hostname"`
}

func (self CodeHosting) IsEmpty() bool {
	return self.Platform == nil && self.OriginHostname == nil
}

type SyncStrategy struct {
	FeatureBranches   *string `toml:"feature-branches"`
	PerennialBranches *string `toml:"perennial-branches"`
}

func (self SyncStrategy) IsEmpty() bool {
	return self.FeatureBranches == nil && self.PerennialBranches == nil
}
