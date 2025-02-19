package cmd

// flags
var (
	autoInstall    bool
	defaultFlag    bool
	forceFlag      bool
	listFlag       bool
	jsonFlag       bool
	productionFlag bool
	verifiedFlag   bool
)

// options
var (
	accountName     string
	assistantId     string
	bearerToken     string
	cfgFile         string
	copyDirectory   string
	directory       string
	endpointUrl     string
	environmentYaml string
	environmentFile string
	granularity     string
	ignores         []string
	robotFile       string
	robotId         string
	robotName       string
	port            string
	runTask         string
	shellDirectory  string
	templateName    string
	gracePeriod     int
	validityTime    int
	workspaceId     string
	wskey           string
	zipfile         string
)
