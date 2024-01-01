package configfile

// Data is the unvalidated data as read by the TOML parser.
type Data struct {
	Branches                 Branches      `toml:"branches"`
	CodeHosting              *CodeHosting  `toml:"code-hosting"`
	SyncStrategy             *SyncStrategy `toml:"sync-strategy"`
	PushNewbranches          *bool         `toml:"push-new-branches"`
	ShipDeleteTrackingBranch *bool         `toml:"ship-delete-remote-branch"`
	SyncUpstream             *bool         `toml:"sync-upstream"`
}

type Branches struct {
	Main       *string  `toml:"main"`
	Perennials []string `toml:"perennials"`
}

type CodeHosting struct {
	Platform       *string `toml:"platform"`
	OriginHostname *string `toml:"origin-hostname"`
}

type SyncStrategy struct {
	FeatureBranches   *string `toml:"feature-branches"`
	PerennialBranches *string `toml:"perennial-branches"`
}
