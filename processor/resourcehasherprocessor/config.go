// Copyright The OpenTelemetry Authors
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

package resourcehasherprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcehasherprocessor"

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/config"
) // Config defines configuration for Resource Hasher processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	MaximumCacheSize         int                      `mapstructure:"max_cache_size"`
	MaximumCacheEntryAge     time.Duration            `mapstructure:"max_cache_entry_age"`
}

// Validate config
func (cfg *Config) Validate() error {
	if cfg.MaximumCacheSize < 1 {
		return fmt.Errorf("the minimum cache size is 1")
	}

	return nil
}
