package messages

const (
	InputAddOrRemove = `invalid argument %q. Please provide either "add" or "remove"`
	InputYesOrNo     = `invalid argument: %q. Please provide either "yes" or "no".\n`

	AbortNothingToDo                    = "nothing to abort"
	ArgumentUnknown                     = "unknown argument: %q"
	BranchAlreadyExists                 = "there is already a branch %q"
	BranchDoesntExist                   = "there is no branch %q"
	BrowserOpen                         = "Please open in a browser: %s\n"
	CacheUsedUnitialized                = "using a cached value before initialization"
	CompletionTypeUnknown               = "unknown completion type: %q"
	ContinueNothingToDo                 = "nothing to continue"
	ContinueUnresolvedConflicts         = "you must resolve the conflicts before continuing"
	DiffNoFeatureBranch                 = "you can only diff-parent feature branches"
	DeleteRemoteBranchCannotWhenOffline = "cannot delete remote branch %q in offline mode"
	KillOnlyFeatureBranches             = "you can only kill feature branches"
	RenameBranchNotInSync               = "%q is not in sync with its tracking branch, please sync the branches before renaming"
	RenameMainBranch                    = "the main branch cannot be renamed"
	RenameToSameName                    = "cannot rename branch to current name"
	RunstateLoadProblem                 = "cannot load previous run state: %w"
	SetParentNoFeatureBranch            = "the branch %q is not a feature branch. Only feature branches can have parent branches"
	RenamePerennialBranchWarning        = "%q is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'"
)
