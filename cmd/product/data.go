package product

type ReleaseInfo struct {
	Name        string
	ReleaseDate string
	LatestName  string
	LatestDate  string
	EoasFrom    string
	EolFrom     string
	LTS         bool
}

type productReleases struct {
	Name     string
	Releases []ReleaseInfo
}

type productResult struct {
	Name        string
	EolLabel    string
	ReleaseName string
	ReleaseDate string
	EolFrom     string
}
