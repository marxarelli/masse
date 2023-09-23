package v1

// Config represents the top-most user-facing Phyton configuration
type Config struct {
	Parameters *Parameters `yaml:"parameters"`
	// Options    *build_v1.Options `yaml:"options"`
	// Targets    *build_v1.Targets `yaml:"targets"`
}

// Parameters represents build-time arguments
type Parameters map[string]string
