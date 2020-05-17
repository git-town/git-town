package steps

// FeatureState contains state that is shared by all scenarios in a feature.
type FeatureState struct {

	// state contains the state for the currently executing scenario in this feature.
	// Godog executes each feature concurrently, but the scenarios in a feature sequentially.
	// This means there is always only one active scenario for each feature.
	state scenarioState
}
