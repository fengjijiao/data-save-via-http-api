package conf

import (
	"gopkg.in/yaml.v2"
)

func Load(context []byte) error {
	err := yaml.Unmarshal([]byte(context), &Config)
    if err != nil {
        return err
    }
	Config.setDefaults()
	return nil
}

func Parser(context []byte) (*ConfInfo, error) {
	var newConfInfo ConfInfo
	err := yaml.Unmarshal(context, &newConfInfo)
    if err != nil {
        return nil, err
    }
	newConfInfo.setDefaults()
	return &newConfInfo, nil
}