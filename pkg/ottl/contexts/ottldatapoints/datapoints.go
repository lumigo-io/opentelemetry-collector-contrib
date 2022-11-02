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

package ottldatapoints // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoints"

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/internal/ottlcommon"
)

var _ ottlcommon.ResourceContext = TransformContext{}
var _ ottlcommon.InstrumentationScopeContext = TransformContext{}

type TransformContext struct {
	dataPoint            interface{}
	metric               pmetric.Metric
	metrics              pmetric.MetricSlice
	instrumentationScope pcommon.InstrumentationScope
	resource             pcommon.Resource
}

func NewTransformContext(dataPoint interface{}, metric pmetric.Metric, metrics pmetric.MetricSlice, instrumentationScope pcommon.InstrumentationScope, resource pcommon.Resource) TransformContext {
	return TransformContext{
		dataPoint:            dataPoint,
		metric:               metric,
		metrics:              metrics,
		instrumentationScope: instrumentationScope,
		resource:             resource,
	}
}

func (ctx TransformContext) GetDataPoint() interface{} {
	return ctx.dataPoint
}

func (ctx TransformContext) GetInstrumentationScope() pcommon.InstrumentationScope {
	return ctx.instrumentationScope
}

func (ctx TransformContext) GetResource() pcommon.Resource {
	return ctx.resource
}

func (ctx TransformContext) GetMetric() pmetric.Metric {
	return ctx.metric
}

func (ctx TransformContext) GetMetrics() pmetric.MetricSlice {
	return ctx.metrics
}

func NewParser(functions map[string]interface{}, telemetrySettings component.TelemetrySettings) ottl.Parser[TransformContext] {
	return ottl.NewParser[TransformContext](functions, parsePath, parseEnum, telemetrySettings)
}

var symbolTable = map[ottl.EnumSymbol]ottl.Enum{
	"FLAG_NONE":              0,
	"FLAG_NO_RECORDED_VALUE": 1,
}

func init() {
	for k, v := range ottlcommon.MetricSymbolTable {
		symbolTable[k] = v
	}
}

func parseEnum(val *ottl.EnumSymbol) (*ottl.Enum, error) {
	if val != nil {
		if enum, ok := symbolTable[*val]; ok {
			return &enum, nil
		}
		return nil, fmt.Errorf("enum symbol, %s, not found", *val)
	}
	return nil, fmt.Errorf("enum symbol not provided")
}

func parsePath(val *ottl.Path) (ottl.GetSetter[TransformContext], error) {
	if val != nil && len(val.Fields) > 0 {
		return newPathGetSetter(val.Fields)
	}
	return nil, fmt.Errorf("bad path %v", val)
}

func newPathGetSetter(path []ottl.Field) (ottl.GetSetter[TransformContext], error) {
	switch path[0].Name {
	case "resource":
		return ottlcommon.ResourcePathGetSetter[TransformContext](path[1:])
	case "instrumentation_scope":
		return ottlcommon.ScopePathGetSetter[TransformContext](path[1:])
	case "metric":
		return ottlcommon.MetricPathGetSetter[TransformContext](path[1:])
	case "attributes":
		mapKey := path[0].MapKey
		if mapKey == nil {
			return accessAttributes(), nil
		}
		return accessAttributesKey(mapKey), nil
	case "start_time_unix_nano":
		return accessStartTimeUnixNano(), nil
	case "time_unix_nano":
		return accessTimeUnixNano(), nil
	case "value_double":
		return accessDoubleValue(), nil
	case "value_int":
		return accessIntValue(), nil
	case "exemplars":
		return accessExemplars(), nil
	case "flags":
		return accessFlags(), nil
	case "count":
		return accessCount(), nil
	case "sum":
		return accessSum(), nil
	case "bucket_counts":
		return accessBucketCounts(), nil
	case "explicit_bounds":
		return accessExplicitBounds(), nil
	case "scale":
		return accessScale(), nil
	case "zero_count":
		return accessZeroCount(), nil
	case "positive":
		if len(path) == 1 {
			return accessPositive(), nil
		}
		switch path[1].Name {
		case "offset":
			return accessPositiveOffset(), nil
		case "bucket_counts":
			return accessPositiveBucketCounts(), nil
		}
	case "negative":
		if len(path) == 1 {
			return accessNegative(), nil
		}
		switch path[1].Name {
		case "offset":
			return accessNegativeOffset(), nil
		case "bucket_counts":
			return accessNegativeBucketCounts(), nil
		}
	case "quantile_values":
		return accessQuantileValues(), nil
	}
	return nil, fmt.Errorf("invalid path expression %v", path)
}

func accessAttributes() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			switch ctx.GetDataPoint().(type) {
			case pmetric.NumberDataPoint:
				return ctx.GetDataPoint().(pmetric.NumberDataPoint).Attributes(), nil
			case pmetric.HistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.HistogramDataPoint).Attributes(), nil
			case pmetric.ExponentialHistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Attributes(), nil
			case pmetric.SummaryDataPoint:
				return ctx.GetDataPoint().(pmetric.SummaryDataPoint).Attributes(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			switch ctx.GetDataPoint().(type) {
			case pmetric.NumberDataPoint:
				if attrs, ok := val.(pcommon.Map); ok {
					attrs.CopyTo(ctx.GetDataPoint().(pmetric.NumberDataPoint).Attributes())
				}
			case pmetric.HistogramDataPoint:
				if attrs, ok := val.(pcommon.Map); ok {
					attrs.CopyTo(ctx.GetDataPoint().(pmetric.HistogramDataPoint).Attributes())
				}
			case pmetric.ExponentialHistogramDataPoint:
				if attrs, ok := val.(pcommon.Map); ok {
					attrs.CopyTo(ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Attributes())
				}
			case pmetric.SummaryDataPoint:
				if attrs, ok := val.(pcommon.Map); ok {
					attrs.CopyTo(ctx.GetDataPoint().(pmetric.SummaryDataPoint).Attributes())
				}
			}
			return nil
		},
	}
}

func accessAttributesKey(mapKey *string) ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			switch ctx.GetDataPoint().(type) {
			case pmetric.NumberDataPoint:
				return ottlcommon.GetMapValue(ctx.GetDataPoint().(pmetric.NumberDataPoint).Attributes(), *mapKey), nil
			case pmetric.HistogramDataPoint:
				return ottlcommon.GetMapValue(ctx.GetDataPoint().(pmetric.HistogramDataPoint).Attributes(), *mapKey), nil
			case pmetric.ExponentialHistogramDataPoint:
				return ottlcommon.GetMapValue(ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Attributes(), *mapKey), nil
			case pmetric.SummaryDataPoint:
				return ottlcommon.GetMapValue(ctx.GetDataPoint().(pmetric.SummaryDataPoint).Attributes(), *mapKey), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			switch ctx.GetDataPoint().(type) {
			case pmetric.NumberDataPoint:
				ottlcommon.SetMapValue(ctx.GetDataPoint().(pmetric.NumberDataPoint).Attributes(), *mapKey, val)
			case pmetric.HistogramDataPoint:
				ottlcommon.SetMapValue(ctx.GetDataPoint().(pmetric.HistogramDataPoint).Attributes(), *mapKey, val)
			case pmetric.ExponentialHistogramDataPoint:
				ottlcommon.SetMapValue(ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Attributes(), *mapKey, val)
			case pmetric.SummaryDataPoint:
				ottlcommon.SetMapValue(ctx.GetDataPoint().(pmetric.SummaryDataPoint).Attributes(), *mapKey, val)
			}
			return nil
		},
	}
}

func accessStartTimeUnixNano() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			switch ctx.GetDataPoint().(type) {
			case pmetric.NumberDataPoint:
				return ctx.GetDataPoint().(pmetric.NumberDataPoint).StartTimestamp().AsTime().UnixNano(), nil
			case pmetric.HistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.HistogramDataPoint).StartTimestamp().AsTime().UnixNano(), nil
			case pmetric.ExponentialHistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).StartTimestamp().AsTime().UnixNano(), nil
			case pmetric.SummaryDataPoint:
				return ctx.GetDataPoint().(pmetric.SummaryDataPoint).StartTimestamp().AsTime().UnixNano(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newTime, ok := val.(int64); ok {
				switch ctx.GetDataPoint().(type) {
				case pmetric.NumberDataPoint:
					ctx.GetDataPoint().(pmetric.NumberDataPoint).SetStartTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, newTime)))
				case pmetric.HistogramDataPoint:
					ctx.GetDataPoint().(pmetric.HistogramDataPoint).SetStartTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, newTime)))
				case pmetric.ExponentialHistogramDataPoint:
					ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).SetStartTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, newTime)))
				case pmetric.SummaryDataPoint:
					ctx.GetDataPoint().(pmetric.SummaryDataPoint).SetStartTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, newTime)))
				}
			}
			return nil
		},
	}
}

func accessTimeUnixNano() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			switch ctx.GetDataPoint().(type) {
			case pmetric.NumberDataPoint:
				return ctx.GetDataPoint().(pmetric.NumberDataPoint).Timestamp().AsTime().UnixNano(), nil
			case pmetric.HistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.HistogramDataPoint).Timestamp().AsTime().UnixNano(), nil
			case pmetric.ExponentialHistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Timestamp().AsTime().UnixNano(), nil
			case pmetric.SummaryDataPoint:
				return ctx.GetDataPoint().(pmetric.SummaryDataPoint).Timestamp().AsTime().UnixNano(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newTime, ok := val.(int64); ok {
				switch ctx.GetDataPoint().(type) {
				case pmetric.NumberDataPoint:
					ctx.GetDataPoint().(pmetric.NumberDataPoint).SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, newTime)))
				case pmetric.HistogramDataPoint:
					ctx.GetDataPoint().(pmetric.HistogramDataPoint).SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, newTime)))
				case pmetric.ExponentialHistogramDataPoint:
					ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, newTime)))
				case pmetric.SummaryDataPoint:
					ctx.GetDataPoint().(pmetric.SummaryDataPoint).SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, newTime)))
				}
			}
			return nil
		},
	}
}

func accessDoubleValue() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if numberDataPoint, ok := ctx.GetDataPoint().(pmetric.NumberDataPoint); ok {
				return numberDataPoint.DoubleValue(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newDouble, ok := val.(float64); ok {
				if numberDataPoint, ok := ctx.GetDataPoint().(pmetric.NumberDataPoint); ok {
					numberDataPoint.SetDoubleValue(newDouble)
				}
			}
			return nil
		},
	}
}

func accessIntValue() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if numberDataPoint, ok := ctx.GetDataPoint().(pmetric.NumberDataPoint); ok {
				return numberDataPoint.IntValue(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newInt, ok := val.(int64); ok {
				if numberDataPoint, ok := ctx.GetDataPoint().(pmetric.NumberDataPoint); ok {
					numberDataPoint.SetIntValue(newInt)
				}
			}
			return nil
		},
	}
}

func accessExemplars() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			switch ctx.GetDataPoint().(type) {
			case pmetric.NumberDataPoint:
				return ctx.GetDataPoint().(pmetric.NumberDataPoint).Exemplars(), nil
			case pmetric.HistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.HistogramDataPoint).Exemplars(), nil
			case pmetric.ExponentialHistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Exemplars(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newExemplars, ok := val.(pmetric.ExemplarSlice); ok {
				switch ctx.GetDataPoint().(type) {
				case pmetric.NumberDataPoint:
					newExemplars.CopyTo(ctx.GetDataPoint().(pmetric.NumberDataPoint).Exemplars())
				case pmetric.HistogramDataPoint:
					newExemplars.CopyTo(ctx.GetDataPoint().(pmetric.HistogramDataPoint).Exemplars())
				case pmetric.ExponentialHistogramDataPoint:
					newExemplars.CopyTo(ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Exemplars())
				}
			}
			return nil
		},
	}
}

func accessFlags() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			switch ctx.GetDataPoint().(type) {
			case pmetric.NumberDataPoint:
				return int64(ctx.GetDataPoint().(pmetric.NumberDataPoint).Flags()), nil
			case pmetric.HistogramDataPoint:
				return int64(ctx.GetDataPoint().(pmetric.HistogramDataPoint).Flags()), nil
			case pmetric.ExponentialHistogramDataPoint:
				return int64(ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Flags()), nil
			case pmetric.SummaryDataPoint:
				return int64(ctx.GetDataPoint().(pmetric.SummaryDataPoint).Flags()), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newFlags, ok := val.(int64); ok {
				switch ctx.GetDataPoint().(type) {
				case pmetric.NumberDataPoint:
					ctx.GetDataPoint().(pmetric.NumberDataPoint).SetFlags(pmetric.DataPointFlags(newFlags))
				case pmetric.HistogramDataPoint:
					ctx.GetDataPoint().(pmetric.HistogramDataPoint).SetFlags(pmetric.DataPointFlags(newFlags))
				case pmetric.ExponentialHistogramDataPoint:
					ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).SetFlags(pmetric.DataPointFlags(newFlags))
				case pmetric.SummaryDataPoint:
					ctx.GetDataPoint().(pmetric.SummaryDataPoint).SetFlags(pmetric.DataPointFlags(newFlags))
				}
			}
			return nil
		},
	}
}

func accessCount() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			switch ctx.GetDataPoint().(type) {
			case pmetric.HistogramDataPoint:
				return int64(ctx.GetDataPoint().(pmetric.HistogramDataPoint).Count()), nil
			case pmetric.ExponentialHistogramDataPoint:
				return int64(ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Count()), nil
			case pmetric.SummaryDataPoint:
				return int64(ctx.GetDataPoint().(pmetric.SummaryDataPoint).Count()), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newCount, ok := val.(int64); ok {
				switch ctx.GetDataPoint().(type) {
				case pmetric.HistogramDataPoint:
					ctx.GetDataPoint().(pmetric.HistogramDataPoint).SetCount(uint64(newCount))
				case pmetric.ExponentialHistogramDataPoint:
					ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).SetCount(uint64(newCount))
				case pmetric.SummaryDataPoint:
					ctx.GetDataPoint().(pmetric.SummaryDataPoint).SetCount(uint64(newCount))
				}
			}
			return nil
		},
	}
}

func accessSum() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			switch ctx.GetDataPoint().(type) {
			case pmetric.HistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.HistogramDataPoint).Sum(), nil
			case pmetric.ExponentialHistogramDataPoint:
				return ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).Sum(), nil
			case pmetric.SummaryDataPoint:
				return ctx.GetDataPoint().(pmetric.SummaryDataPoint).Sum(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newSum, ok := val.(float64); ok {
				switch ctx.GetDataPoint().(type) {
				case pmetric.HistogramDataPoint:
					ctx.GetDataPoint().(pmetric.HistogramDataPoint).SetSum(newSum)
				case pmetric.ExponentialHistogramDataPoint:
					ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint).SetSum(newSum)
				case pmetric.SummaryDataPoint:
					ctx.GetDataPoint().(pmetric.SummaryDataPoint).SetSum(newSum)
				}
			}
			return nil
		},
	}
}

func accessExplicitBounds() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if histogramDataPoint, ok := ctx.GetDataPoint().(pmetric.HistogramDataPoint); ok {
				return histogramDataPoint.ExplicitBounds().AsRaw(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newExplicitBounds, ok := val.([]float64); ok {
				if histogramDataPoint, ok := ctx.GetDataPoint().(pmetric.HistogramDataPoint); ok {
					histogramDataPoint.ExplicitBounds().FromRaw(newExplicitBounds)
				}
			}
			return nil
		},
	}
}

func accessBucketCounts() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if histogramDataPoint, ok := ctx.GetDataPoint().(pmetric.HistogramDataPoint); ok {
				return histogramDataPoint.BucketCounts().AsRaw(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newBucketCount, ok := val.([]uint64); ok {
				if histogramDataPoint, ok := ctx.GetDataPoint().(pmetric.HistogramDataPoint); ok {
					histogramDataPoint.BucketCounts().FromRaw(newBucketCount)
				}
			}
			return nil
		},
	}
}

func accessScale() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
				return int64(expoHistogramDataPoint.Scale()), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newScale, ok := val.(int64); ok {
				if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
					expoHistogramDataPoint.SetScale(int32(newScale))
				}
			}
			return nil
		},
	}
}

func accessZeroCount() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
				return int64(expoHistogramDataPoint.ZeroCount()), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newZeroCount, ok := val.(int64); ok {
				if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
					expoHistogramDataPoint.SetZeroCount(uint64(newZeroCount))
				}
			}
			return nil
		},
	}
}

func accessPositive() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
				return expoHistogramDataPoint.Positive(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newPositive, ok := val.(pmetric.ExponentialHistogramDataPointBuckets); ok {
				if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
					newPositive.CopyTo(expoHistogramDataPoint.Positive())
				}
			}
			return nil
		},
	}
}

func accessPositiveOffset() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
				return int64(expoHistogramDataPoint.Positive().Offset()), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newPositiveOffset, ok := val.(int64); ok {
				if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
					expoHistogramDataPoint.Positive().SetOffset(int32(newPositiveOffset))
				}
			}
			return nil
		},
	}
}

func accessPositiveBucketCounts() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
				return expoHistogramDataPoint.Positive().BucketCounts().AsRaw(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newPositiveBucketCounts, ok := val.([]uint64); ok {
				if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
					expoHistogramDataPoint.Positive().BucketCounts().FromRaw(newPositiveBucketCounts)
				}
			}
			return nil
		},
	}
}

func accessNegative() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
				return expoHistogramDataPoint.Negative(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newNegative, ok := val.(pmetric.ExponentialHistogramDataPointBuckets); ok {
				if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
					newNegative.CopyTo(expoHistogramDataPoint.Negative())
				}
			}
			return nil
		},
	}
}

func accessNegativeOffset() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
				return int64(expoHistogramDataPoint.Negative().Offset()), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newNegativeOffset, ok := val.(int64); ok {
				if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
					expoHistogramDataPoint.Negative().SetOffset(int32(newNegativeOffset))
				}
			}
			return nil
		},
	}
}

func accessNegativeBucketCounts() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
				return expoHistogramDataPoint.Negative().BucketCounts().AsRaw(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newNegativeBucketCounts, ok := val.([]uint64); ok {
				if expoHistogramDataPoint, ok := ctx.GetDataPoint().(pmetric.ExponentialHistogramDataPoint); ok {
					expoHistogramDataPoint.Negative().BucketCounts().FromRaw(newNegativeBucketCounts)
				}
			}
			return nil
		},
	}
}

func accessQuantileValues() ottl.StandardGetSetter[TransformContext] {
	return ottl.StandardGetSetter[TransformContext]{
		Getter: func(ctx TransformContext) (interface{}, error) {
			if summaryDataPoint, ok := ctx.GetDataPoint().(pmetric.SummaryDataPoint); ok {
				return summaryDataPoint.QuantileValues(), nil
			}
			return nil, nil
		},
		Setter: func(ctx TransformContext, val interface{}) error {
			if newQuantileValues, ok := val.(pmetric.SummaryDataPointValueAtQuantileSlice); ok {
				if summaryDataPoint, ok := ctx.GetDataPoint().(pmetric.SummaryDataPoint); ok {
					newQuantileValues.CopyTo(summaryDataPoint.QuantileValues())
				}
			}
			return nil
		},
	}
}