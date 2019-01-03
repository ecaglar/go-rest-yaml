//package model defines objects (models)
package model

type Metadata struct {
	Title       string           `yaml:"title"`
	Version     string           `yaml:"version"`
	Company     string           `yaml:"company"`
	Website     string           `yaml:"website"`
	Source      string           `yaml:"source"`
	License     string           `yaml:"license"`
	Maintainers []MaintainPerson `yaml:"maintainers"`
	Description string           `yaml:"description"`
}

type MaintainPerson struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}
