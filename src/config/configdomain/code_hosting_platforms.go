package configdomain

type CodeHostingPlatforms []CodeHostingPlatform

func (self CodeHostingPlatforms) Strings() []string {
	result := make([]string, len(self))
	for p, platform := range self {
		result[p] = platform.String()
	}
	return result
}

func AllCodeHostingPlatforms() CodeHostingPlatforms {
	return []CodeHostingPlatform{
		CodeHostingPlatformAutoDetect,
		CodeHostingPlatformBitBucket,
		CodeHostingPlatformGitea,
		CodeHostingPlatformGitHub,
		CodeHostingPlatformGitLab,
	}
}
