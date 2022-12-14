package mackerel

type MackerelConfig struct {
	ApiKey             string   `env:"EXT_MACKEREL_API_KEY"`
	ApiKeySSMParamName string   `env:"EXT_MACKEREL_API_KEY_SSM"`
	RoleFullnames      []string `env:"EXT_MACKEREL_ROLE_FULL_NAMES" envSeparator:","`
}
