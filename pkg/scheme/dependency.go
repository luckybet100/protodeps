package scheme

type Dependency struct {
	Name    string `yaml:"name" json:"name"`
	GitRepo string `yaml:"git" json:"git"`
	Ref     string `yaml:"ref" json:"ref"`
}
