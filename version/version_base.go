package version

const AppName = "autocost"

var (
	// The compilation date. This will be filled in by the compiler.
	BuildDate string

	// The git commit that was compiled. This will be filled in by the compiler.
	GitCommit   string
	GitBranch   string
	GitDescribe string

	Version           = "0.1.2"
	VersionPrerelease = ""
	VersionMetadata   = ""
)
