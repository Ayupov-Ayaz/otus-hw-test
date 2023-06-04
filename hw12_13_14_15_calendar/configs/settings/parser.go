package settings

import (
	"errors"
	"fmt"
	"os"

	env8 "github.com/caarlos0/env/v8"
	yaml3 "gopkg.in/yaml.v3"
)

func UnmarshalYaml(data []byte, cfg interface{}) error {
	if err := yaml3.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return nil
}

func UnmarshalYamlFile(yamlFile string, cfg interface{}) error {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Если файла нет, то просто возвращаем nil
			return nil
		}

		// Если файл есть, но не удалось его прочитать, возвращаем ошибку.
		return err
	}

	return UnmarshalYaml(data, cfg)
}

func UnmarshalEnv(prefix string, cfg interface{}) error {
	opts := env8.Options{
		Prefix: prefix,
	}

	return env8.ParseWithOptions(cfg, opts)
}
