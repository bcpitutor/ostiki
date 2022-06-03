package version

var (
	BuildTime     string = ""
	BuildVersion  string = ""
	BuildMachine  string = ""
	CommitDetails string = ""
)

type VDetails struct {
	Version string
	Time    string
	Branch  string
	Details string
	Machine string
}

var VersionDetails = &VDetails{
	Version: BuildVersion,
	Time:    BuildTime,
	Details: CommitDetails,
	Machine: BuildMachine,
}
