package inout

type Modules struct {
	Name      string   `yaml:"name"`
	Resources []string `yaml:"resources"`
}

type YamlMapping struct {
	Modules []Modules `yaml:"modules"`
	Confg   []string  `yaml:"config"`
}

type ParsedTf struct {
	Providers []string
	Resources []string
}

type Resource struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
}

type ModuleResource struct {
	Module       string
	ResourceType string
	Quantity     string
}

type CsvResources struct {
	Resource string `json:"Resource"`
	Module   string `json:"Module"`
	Quantity int    `json:"Quantity"`
}

type BlockInnerKey struct {
	MainKey        string `json:"MainKey"`
	InnerKey       string `json:"InnerKey"`
	SecondInnerKey string `json:"SecondInnerKey"`
	Line           string `json:"Line"`
}

type Imports struct {
	Resource_key int
	Resource_id  string
}

type UnmappedOutputs struct {
	ResourceName     string
	ResourceVariable string
}

type Outputs struct {
	OuputModule     string
	OputputResource string
	OuptutModuleRef string
}

type Template struct {
	//Subscription_ID   string //no
	UnmappedResources []string
	NotFoundResources []UnmappedOutputs
	FoundResources    []UnmappedOutputs
	//ResourceGroups    []string //no
}
