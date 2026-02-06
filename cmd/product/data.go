package product

type productResult struct {
	Name        string
	EolLabel    string
	ReleaseName string
	ReleaseDate string
	EolFrom     string
}

type ApiRestDescribe struct {
	Result struct {
		Name           string `json:"name"`
		VersionCommand string `json:"versionCommand"`
		Identifiers    []struct {
			Type string `json:"type"`
			Id   string `json:"id"`
		} `json:"identifiers"`
	}
}
