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
