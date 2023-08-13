package main

// Dict
type Dict struct {
	Word        string         `yaml:"word"`
	Bopomofo    string         `yaml:"bopomofo"`
	Description string         `yaml:"description"`
	Examples    []*DictExample `yaml:"examples"`
	Distinguish string         `yaml:"distinguish"`

	Code         string              `yaml:"-"`
	ExampleStr   string              `yaml:"-"`
	WordTeardown []*DictWordTeardown `yaml:"-"`
}

// DictExample
type DictExample struct {
	Words       []string `yaml:"words"`
	Description string   `yaml:"description"`
	Correct     string   `yaml:"correct"`
	Incorrect   string   `yaml:"incorrect"`
}

// DictWordTeardown
type DictWordTeardown struct {
	Character string
	Bopomofo  []string
	Accent    string
}

// SimpleDict
type SimpleDict struct {
	Code       string `json:"c"`
	Word       string `json:"n"`
	ExampleStr string `json:"e"`
}
