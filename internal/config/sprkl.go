package config

import (
	"encoding/json"
	"os"
	"regexp"
)

const (
	env_SPRKL_NODE_IMAGE     = "SPRKL_NODE_IMAGE"
	env_SPRKL_INJECT_TO_PODS = "SPRKL_INJECT_TO_PODS"
	env_SPRKL_POD_EXCLUDE    = "SPRKL_POD_EXCLUDE"
	env_SPRKL_POD_INCLUDE    = "SPRKL_POD_INCLUDE"
	env_SPRKL_NS_EXCLUDE     = "SPRKL_NS_EXCLUDE"
	env_SPRKL_NS_INCLUDE     = "SPRKL_NS_INCLUDE"
)

var _instance SprklConfig = SprklConfig{
	SprklNodeImage:    getSprklNodeImage(),
	SprklInjectToPods: getSprklInjectToPods(),
	SprklPodExclude:   getSprklPodExclude(),
	SprklPodInclude:   getSprklPodInclude(),
	SprklNsExclude:    getSprklNsExclude(),
	SprklNsInclude:    getSprklNsInclude(),
	BlacklistedNamespaces: map[string]bool{
		"kube-system":                   true,
		"kube-public":                   true,
		"kube-node-lease":               true,
		"cert-manager":                  true,
		"opentelemetry-operator-system": true,
	},
}

// TODO: immutable
type SprklConfig struct {
	SprklNodeImage        string            `json:"sprkl-node-image"`
	SprklInjectToPods     map[string]string `json:"sprkl-inject-to-pods"`
	SprklPodExclude       *regexp.Regexp    `json:"sprkl-pod-exclude"`
	SprklPodInclude       *regexp.Regexp    `json:"sprkl-pod-include"`
	SprklNsExclude        *regexp.Regexp    `json:"sprkl-ns-exclude"`
	SprklNsInclude        *regexp.Regexp    `json:"sprkl-ns-include"`
	BlacklistedNamespaces map[string]bool
}

func GetSprklConfig() SprklConfig {
	return _instance
}

func getSprklNodeImage() string {
	// TODO: enforce building image with this env
	return os.Getenv(env_SPRKL_NODE_IMAGE)
}

func getSprklPodExclude() *regexp.Regexp {
	val := os.Getenv(env_SPRKL_POD_EXCLUDE)
	regex, err := regexp.Compile(val)
	if err != nil {
		return regexp.MustCompile("")
	}

	return regex
}

func getSprklPodInclude() *regexp.Regexp {
	val := os.Getenv(env_SPRKL_POD_INCLUDE)
	if len(val) == 0 {
		val = "(.+)"
	}

	regex, err := regexp.Compile(val)
	if err != nil {
		return regexp.MustCompile("(.+)")
	}

	return regex
}

func getSprklNsExclude() *regexp.Regexp {
	val := os.Getenv(env_SPRKL_NS_EXCLUDE)
	regex, err := regexp.Compile(val)
	if err != nil {
		return regexp.MustCompile("")
	}

	return regex
}

func getSprklNsInclude() *regexp.Regexp {
	val := os.Getenv(env_SPRKL_NS_INCLUDE)
	if len(val) == 0 {
		val = "(.+)"
	}

	regex, err := regexp.Compile(val)
	if err != nil {
		return regexp.MustCompile("(.+)")
	}

	return regex
}

func getSprklInjectToPods() map[string]string {
	varsToInject := make(map[string]string)
	// INFO: in case of error => return empty config map
	_ = json.Unmarshal([]byte(os.Getenv(env_SPRKL_INJECT_TO_PODS)), &varsToInject)
	return varsToInject
}
