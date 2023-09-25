package inout

type Modules struct {
	Name       string   `yaml:"name"`
	EntryPoint string   `yaml:"entry_point"`
	Resources  []string `yaml:"resources"`
}

type CommonVar struct {
	Name  string
	Value []string
}

type BackendConf struct {
	Resource_group_name  string `yaml:"resource_group_name"`
	Storage_account_name string `yaml:"storage_account_name"`
	Container_name       string `yaml:"container_name"`
	Key_prefix           string `yaml:"key_prefix"`
}

type CommonVar struct {
	Name  string
	Value []string
}

type YamlMapping struct {
	Modules    []Modules   `yaml:"modules"`
	CommonVars []CommonVar `yaml:"common"`
	Confg      []string    `yaml:"config"`
	Backend    BackendConf `yaml:"backend"`
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
