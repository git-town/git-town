package forgedomain

// DetectedForgeType is the forge type that actually exists.
// Its the one the user has configured,
// or if that is "auto", the automatically detected forge type.
type DetectedForgeType ForgeType

// ForgeType converts this value into a ForgeType.
func (self DetectedForgeType) ForgeType() ForgeType { return ForgeType(self) }
