package product

type ReleaseInfo struct {
	Name        string `json:"cycle"`
	ReleaseDate string `json:"releaseDate"`
	LatestName  string `json:"latest"`
	LatestDate  string `json:"latestReleaseDate"`
	EoasFrom    string `json:"-"`
	EolFrom     string `json:"eolFrom"`
	LTS         bool   `json:"isLts"`
}

type ProductReleases struct {
	Name     string        `json:"name"`
	Releases []ReleaseInfo `json:"releases"`
}

type productResult struct {
	Name        string
	EolLabel    string
	ReleaseName string
	ReleaseDate string
	EolFrom     string
}

type ApiRespExtended struct {
	Result struct {
		Name     string `json:"name"`
		Releases []struct {
			Name        string `json:"name"`
			ReleaseDate string `json:"releaseDate"`
			Latest      struct {
				Name string `json:"name"`
				Date string `json:"date"`
			} `json:"latest"`
			EoasFrom string `json:"eoasFrom"`
			EolFrom  string `json:"eolFrom"`
			IsLTS    bool   `json:"isLTS"`
		} `json:"releases"`
	}
}

type ApiRestDescribe struct {
	Result struct {
		Name           string               `json:"name"`
		Aliases        []string             `json:"aliases"`
		VersionCommand string               `json:"versionCommand"`
		Identifiers    []DescribeIdentifier `json:"identifiers"`
	}
}

// DescribeIdentifier is a product identifier as returned by the API
type DescribeIdentifier struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}
