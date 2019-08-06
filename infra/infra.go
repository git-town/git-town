/*
Package infra manages the repository setup for the feature specs.

Each feature spec starts out in a pre-defined GitEnvironment that looks like this:
- The "developer" GitRepository contains the local workspace repo.
	This is where the feature specs run.
- The "origin" GitRepository contains the remote repo for the workspace repo in the "developer" folder.
	This is where pushes from "developer" go to.
- All repos contain a "main" branch that is configured as Git Town's main branch.

GitManager creates new GitEnvironments for tests.
Setting up the standardized environment happens before each scenario, i.e. a lot.
It takes a while and includes numerous disk operations.
To make this process fast, GitManager creates a fresh standard environment in the "memoized" folder using CreateMemoizedEnvironment.
Before each scenario, it copies it copies this folder into the scenario folder using CreateScenarioEnvironment.

The folder structure on disk looks like this:

```
<system tmp folder>
├── GitManager folder
|   └── memoized        # the pre-defined GitEnvironment for scenarios
|       ├── developer
|       └── origin
├── scenario A folder   # GitEnvironment for the currently tested scenario A
|   ├── developer         # the workspace GitRepository for scenario A
|   └── origin            # the origin GitRepository for scenario A
└── scenario B folder   # GitEnvironment for the currently tested scenario B
    ├── developer
		└── origin
```
*/
package infra
