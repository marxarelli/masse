package load

import "cuelang.org/go/cue/load"

func WithEnv(key, value string) Option {
	return WithEnvMap(map[string]string{key: value})
}

func WithEnvMap(env map[string]string) Option {
	return func(_ string, cfg *load.Config) error {
		for key, value := range env {
			cfg.Env = append(cfg.Env, key+"="+value)
		}
		return nil
	}
}
