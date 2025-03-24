package load

import (
	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
)

func WithEnv(key, value string) Option {
	return WithEnvMap(map[string]string{key: value})
}

func WithEnvMap(env map[string]string) Option {
	return func(_ string, cfg *load.Config, modcfg *modconfig.Config) error {
		for key, value := range env {
			ev := key + "=" + value
			cfg.Env = append(cfg.Env, ev)
			modcfg.Env = append(modcfg.Env, ev)
		}
		return nil
	}
}
