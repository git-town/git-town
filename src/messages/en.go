package messages

const (
	AddOrRemove                         = `invalid argument %q. Please provide either "add" or "remove"`
	BranchAlreadyExists                 = "there is already a branch %q"
	BranchNotFound                      = "there is no branch %q"
	BranchNotInSync                     = "%q is not in sync with its tracking branch, please sync the branches before renaming"
	CacheUsedBeforeInit                 = "using a cached value before initialization"
	CanOnlyDiffFeatureBranches          = "you can only diff-parent feature branches"
	CanOnlyKillFeatureBranches          = "you can only kill feature branches"
	CannotContinue                      = "nothing to continue"
	CannotDeleteRemoteBranchWhenOffline = "cannot delete remote branch %q in offline mode"
	CannotLoadRunstate                  = "cannot load previous run state: %w"
	CannotRenameMainBranch              = "the main branch cannot be renamed"
	NothingToAbort                      = "nothing to abort"
	OnlyFeatureBranchesCanHaveParents   = "the branch %q is not a feature branch. Only feature branches can have parent branches"
	OpenInBrowser                       = "Please open in a browser: %s\n"
	RenameToSameName                    = "cannot rename branch to current name"
	ResolveConflictsBeforeContinuing    = "you must resolve the conflicts before continuing"
	UnknownArgument                     = "unknown argument: %q"
	UnknownCompletionType               = "unknown completion type: %q"
	WarnRenamePerennialBranch           = "%q is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'"
	YesOrNo                             = `invalid argument: %q. Please provide either "yes" or "no".\n`
)
