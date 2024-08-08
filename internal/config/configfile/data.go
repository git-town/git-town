package configfile

// Data defines the Go equivalent of the TOML file content.
type Data struct {
	Branches                 *Branches     `toml:"branches"`
	Hosting                  *Hosting      `toml:"hosting"`
	PushHook                 *bool         `toml:"push-hook"`
	PushNewbranches          *bool         `toml:"push-new-branches"`
	ShipDeleteTrackingBranch *bool         `toml:"ship-delete-tracking-branch"`
	SyncStrategy             *SyncStrategy `toml:"sync-strategy"`
	SyncTags                 *bool         `toml:"sync-tags"`
	SyncUpstream             *bool         `toml:"sync-upstream"`
}

type Branches struct {
	Main           *string  `toml:"main"`
	PerennialRegex *string  `toml:"perennial-regex"`
	Perennials     []string `toml:"perennials"`
}

func (self Branches) IsEmpty() bool {
	return self.Main == nil && len(self.Perennials) == 0
}

type Hosting struct {
	OriginHostname *string `toml:"origin-hostname"`
	Platform       *string `toml:"platform"`
}

func (self Hosting) IsEmpty() bool {
	return self.Platform == nil && self.OriginHostname == nil
}

type SyncStrategy struct {
	FeatureBranches   *string `toml:"feature-branches"`
	PerennialBranches *string `toml:"perennial-branches"`
}

func (self SyncStrategy) IsEmpty() bool {
	return self.FeatureBranches == nil && self.PerennialBranches == nil
}
