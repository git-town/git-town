package configfile

// Data defines the Go equivalent of the TOML file content.
type Data struct {
	Branches                 *Branches     `toml:"branches"`
	Create                   *Create       `toml:"create"`
	CreatePrototypeBranches  *bool         `toml:"create-prototype-branches"`
	Hosting                  *Hosting      `toml:"hosting"`
	Propose                  *Propose      `toml:"propose"`
	PushHook                 *bool         `toml:"push-hook"`
	PushNewBranches          *bool         `toml:"push-new-branches"`
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
	DisplayTypes      *string  `toml:"display-types"`
	FeatureRegex      *string  `toml:"feature-regex"`
	Main              *string  `toml:"main"`
	ObservedRegex     *string  `toml:"observed-regex"`
	Order             *string  `toml:"order"`
	PerennialRegex    *string  `toml:"perennial-regex"`
	Perennials        []string `toml:"perennials"`
	UnknownType       *string  `toml:"unknown-type"`
}

func (self Branches) IsEmpty() bool {
	return self.Main == nil && len(self.Perennials) == 0
}

type Create struct {
	BranchPrefix     *string `toml:"branch-prefix"`
	NewBranchType    *string `toml:"new-branch-type"`
	PushNewBranches  *bool   `toml:"push-new-branches"`
	ShareNewBranches *string `toml:"share-new-branches"`
	Stash            *bool   `toml:"stash"`
}

type Hosting struct {
	Browser         *string `toml:"browser"`
	DevRemote       *string `toml:"dev-remote"`
	ForgeType       *string `toml:"forge-type"`
	GithubConnector *string `toml:"github-connector"`
	GitlabConnector *string `toml:"gitlab-connector"`
	OriginHostname  *string `toml:"origin-hostname"`
	Platform        *string `toml:"platform"`
}

func (self Hosting) IsEmpty() bool {
	return self.ForgeType == nil && self.OriginHostname == nil && self.Platform == nil
}

type Propose struct {
	Breadcrumb *string `toml:"breadcrumb"`
	Direction  *string `toml:"direction"`
	Lineage    *string `toml:"lineage"`
}

type Ship struct {
	DeleteTrackingBranch *bool   `toml:"delete-tracking-branch"`
	IgnoreUncommitted    *bool   `toml:"ignore-uncommitted"`
	Strategy             *string `toml:"strategy"`
}

type Sync struct {
	AutoResolve       *bool   `toml:"auto-resolve"`
	AutoSync          *bool   `toml:"auto-sync"`
	Detached          *bool   `toml:"detached"`
	FeatureStrategy   *string `toml:"feature-strategy"`
	PerennialStrategy *string `toml:"perennial-strategy"`
	PrototypeStrategy *string `toml:"prototype-strategy"`
	PushBranches      *bool   `toml:"push-branches"`
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
