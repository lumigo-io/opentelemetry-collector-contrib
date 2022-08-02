// Copyright 2021 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package udplogreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver"

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"gopkg.in/yaml.v2"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/input/udp"
)

const (
	typeStr   = "udplog"
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a factory for udp receiver
func NewFactory() component.ReceiverFactory {
	return adapter.NewFactory(ReceiverType{}, stability)
}

// ReceiverType implements adapter.LogReceiverType
// to create a udp receiver
type ReceiverType struct{}

// Type is the receiver type
func (f ReceiverType) Type() config.Type {
	return typeStr
}

// CreateDefaultConfig creates a config with type and version
func (f ReceiverType) CreateDefaultConfig() config.Receiver {
	return &UDPLogConfig{
		BaseConfig: adapter.BaseConfig{
			ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
			Operators:        adapter.OperatorConfigs{},
		},
		Input: adapter.InputConfig{},
	}
}

// BaseConfig gets the base config from config, for now
func (f ReceiverType) BaseConfig(cfg config.Receiver) adapter.BaseConfig {
	return cfg.(*UDPLogConfig).BaseConfig
}

// UDPLogConfig defines configuration for the udp receiver
type UDPLogConfig struct {
	adapter.BaseConfig `mapstructure:",squash"`
	Input              adapter.InputConfig `mapstructure:",remain"`
}

// DecodeInputConfig unmarshals the input operator
func (f ReceiverType) DecodeInputConfig(cfg config.Receiver) (*operator.Config, error) {
	logConfig := cfg.(*UDPLogConfig)
	yamlBytes, _ := yaml.Marshal(logConfig.Input)
	inputCfg := udp.NewConfig("udp_input")

	if err := yaml.Unmarshal(yamlBytes, &inputCfg); err != nil {
		return nil, err
	}

	return &operator.Config{Builder: inputCfg}, nil
}
