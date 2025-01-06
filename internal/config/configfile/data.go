package configfile

// Data defines the Go equivalent of the TOML file content.
type Data struct {
	Branches                 *Branches     `toml:"branches"`
	Create                   *Create       `toml:"create"`
	CreatePrototypeBranches  *bool         `toml:"create-prototype-branches"`
	Hosting                  *Hosting      `toml:"hosting"`
	PushHook                 *bool         `toml:"push-hook"`
	PushNewbranches          *bool         `toml:"push-new-branches"`
	Ship                     *Ship         `toml:"ship"`
	ShipDeleteTrackingBranch *bool         `toml:"ship-delete-tracking-branch"`
	ShipStrategy             *string       `toml:"ship-strategy"`
	Sync                     *Sync         `toml:"sync"`
	SyncStrategy             *SyncStrategy `toml:"sync-strategy"`
	SyncTags                 *bool         `toml:"sync-tags"`
	SyncUpstream             *bool         `toml:"sync-upstream"`
}

type Branches struct {
	ContributionRegex *string  `toml:"contribution-regex"`
	DefaultType       *string  `toml:"default-type"`
	FeatureRegex      *string  `toml:"feature-regex"`
	Main              *string  `toml:"main"`
	ObservedRegex     *string  `toml:"observed-regex"`
	PerennialRegex    *string  `toml:"perennial-regex"`
	Perennials        []string `toml:"perennials"`
}

func (self Branches) IsEmpty() bool {
	return self.Main == nil && len(self.Perennials) == 0
}

type Create struct {
	NewBranchType   *string `toml:"new-branch-type"`
	PushNewbranches *bool   `toml:"push-new-branches"`
}

type Hosting struct {
	DevRemote      *string `toml:"dev-remote"`
	OriginHostname *string `toml:"origin-hostname"`
	Platform       *string `toml:"platform"`
}

func (self Hosting) IsEmpty() bool {
	return self.Platform == nil && self.OriginHostname == nil
}

type Ship struct {
	DeleteTrackingBranch *bool   `toml:"delete-tracking-branch"`
	Strategy             *string `toml:"strategy"`
}

type Sync struct {
	FeatureStrategy   *string `toml:"feature-strategy"`
	PerennialStrategy *string `toml:"perennial-strategy"`
	PrototypeStrategy *string `toml:"prototype-strategy"`
	PushHook          *bool   `toml:"push-hook"`
	Tags              *bool   `toml:"tags"`
	Upstream          *bool   `toml:"upstream"`
}

type SyncStrategy struct {
	FeatureBranches   *string `toml:"feature-branches"`
	PerennialBranches *string `toml:"perennial-branches"`
	PrototypeBranches *string `toml:"prototype-branches"`
}

func (self SyncStrategy) IsEmpty() bool {
	return self.FeatureBranches == nil && self.PerennialBranches == nil
}
