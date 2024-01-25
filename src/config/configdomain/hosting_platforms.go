package configdomain

type HostingPlatforms []HostingPlatform

func (self HostingPlatforms) Strings() []string {
	result := make([]string, len(self))
	for p, platform := range self {
		result[p] = platform.String()
	}
	return result
}

func AllCodeHostingPlatforms() HostingPlatforms {
	return []HostingPlatform{
		HostingPlatformAutoDetect,
		HostingPlatformBitBucket,
		HostingPlatformGitea,
		HostingPlatformGitHub,
		HostingPlatformGitLab,
	}
}
