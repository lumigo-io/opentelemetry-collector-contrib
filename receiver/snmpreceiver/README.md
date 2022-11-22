# SNMP Receiver

| Status                   |               |
| ------------------------ |---------------|
| Stability                | [development] |
| Supported pipeline types | metrics       |
| Distributions            | [contrib]     |

This receiver fetches stats from a SNMP enabled host using a [golang
snmp client](https://github.com/gosnmp/gosnmp). Metrics are collected
based upon different configurations in the config file.

## Purpose

The purpose of this receiver is to allow users to generically monitor metrics using SNMP.

If one of the specified SNMP data values cannot be loaded on startup, a
warning will be printed, but the application will not fail fast.

## Prerequisites

This receiver supports SNMP versions:

- v1
- v2c
- v3

## Configuration

### Connection Configuration
These configuration options are for connecting to a SNMP host.

- `collection_interval`: (default = `1m`): This receiver collects metrics on an interval. This value must be a string readable by Golang's [time.ParseDuration](https://pkg.go.dev/time#ParseDuration). Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.
- `endpoint` (default: `udp://localhost:161`): SNMP endpoint to connect to in the form of `[udp|tcp][://]{host}[:{port}]`
  - If no scheme is supplied, a default of `udp` is assumed
  - If no port is supplied, a default of `161` is assumed
- `version`: (default = `v2c`): SNMP version options are
  - `v1`: SNMP version 1
  - `v2c`: SNMP version 2c
  - `v3`: SNMP version 3
- `community`: (default = `public`): The community string for the SNMP connection. This is not available for SNMP version `v3`.
- `user`: The user for the SNMP connection. This is only available for SNMP version `v3`.
- `security_level`: (default = `no_auth_no_priv`): The security requirements of the SNMP connection. This is only available for SNMP version `v3`. SNMP `security_level` options are
  - `no_auth_no_priv`: No authentication protocol and no privacy protocol used
  - `auth_no_priv`: Authentication protocol is used but no privacy protocol used
  - `auth_priv`: Both authentication and privacy protocols are used
- `auth_type`: (default = `MD5`): The authentication protocol used for the SNMP connection. This is only available if `security_level` is not set to `no_auth_no_priv`. SNMP `auth_type` options are
  - `MD5`
  - `SHA`
  - `SHA224`
  - `SHA256`
  - `SHA384`
  - `SHA512`
- `auth_password`: The authentication password used for the SNMP connection. This is only available if `security_level` is not set to `no_auth_no_priv`.
- `privacy_type`: (default = `DES`): The privacy protocol used for the SNMP connection. This is only available if `security_level` is set to `auth_priv`. SNMP `privacy_type` options are
  - `DES`
  - `AES`
  - `AES192`
  - `AES256`
  - `AES192c`
  - `AES256c`
- `privacy_password`: The privacy password used for the SNMP connection. This is only available if `security_level` is set to `auth_priv`.

### Metric/Attribute Configuration
These configuration options are for determining what metrics and attributes will be created with what SNMP data

- `resource_attributes`: This may be configured with one or more key value pairs of resource attribute names and resource attribute configurations.
- `attributes` This may be configured with one or more key value pairs of attribute names and attribute configurations
- `metrics`: This is the only required parameter. The must be configured with one or more key value pairs of metric names and metric configuration.

#### Resource Attribute Configuration
Resource attribute configurations are used to define what resource attributes will be used in a collection.

| Field Name           | Description                              | Value        | 
| --                   | --                                       | --           |
| `oid`                  | Required if no `indexed_value_prefix`. This is the column OID in a SNMP table which will use the returned indexed SNMP data to create resource attribute values for unique resources. Metric configurations will reference these resource attribute configurations in order to assign metrics data to resources | string       |
| `indexed_value_prefix` | Required if no `oid`. This is a string prefix which will be added to the indices of returned metric indexed SNMP data to create resource attribute values for unique resources. Metric configurations will reference these resource attribute configurations in order to assign metrics data to resources | string       |
| `description`          | Definition of what the resource attribute represents  | string       |

#### Attribute Configuration
Attribute configurations are used to define what resource attributes will be used in a collection.

| Field Name           | Description                                           | Value                           |
| --                   | --                                                    | --                              |
| `oid`                  | Required if no `indexed_value_prefix` or `enum`. This is the column OID in a SNMP table which will use the returned indexed SNMP data to create attribute values for the attribute. Metric configurations will reference these attribute configurations in order to assign these attributes and indexed data values to metrics and their datapoints | string       |
| `indexed_value_prefix` | Required if no `oid` or `enum`. This is a string prefix which will be added to the indices of returned metric indexed SNMP data to create attribute values the attribute. Metric configurations will reference these attribute configurations in order to assign these attributes and index based value to metrics and their datapoints | string       |
| `enum`                 | Required if no `oid` or `indexed_value_prefix`. This should be a list of values that are possible for this attribute. Metric configurations will reference these attribute configurations in order to assign these attributes and values to metrics and their datapoints | string[]       |
| `description`          | Definition of what the attribute represents           | string       |

#### Metric Configuration

| Field Name  | Description                                                    | Value                       | Default |
| --          | --                                                             | --                          | --      |
| `unit`        | Required. To display what is actually being measured for this metric | string                | 1       |
| `gauge`       | Required if no `sum`. Details that this metric is of the gauge type | GaugeMetric              |         |
| `sum`         | Required if no `gauge`. Details that this metric is of the sum type | SumMetric                |         |
| `column_oids` | Required if no `scalar_oids`. Details that this metric is made from one or more columns in an SNMP table. The returned indexed SNMP data for these OIDs might either be datapoints on a single metrics, or datapoints across multiple metrics attached to different resources depending on the column OID configurations | ColumnOID[] |        |
| `scalar_oids` | Required if no `column_oids`. Details that this metric is made from one or more scalard SNMP values (multiple scalar OIDs would represent multiple datapoints within the same metric) | ScalarOID[]       |       |
| `description` | Definition of what the metric represents                       | string                      |         |

#### GaugeMetric Configuration

| Field Name  | Description                                                    | Value                       | Default |
| --          | --                                                             | --                          | --      |
| `value_type`  | The value type of this metric's data. Can be either `int` or `double` | string             | double   |

#### SumMetric Configuration

| Field Name  | Description                                                    | Value                       | Default |
| --          | --                                                             | --                          | --      |
| `value_type` | The value type of this metric's data. Can be either `int` or `double` | string              | double   |
| `monotonic` | Whether this is a monotonic sum or not                         | bool                        | false   |
| `aggregation` | The aggregation type of this metric's data. Can be either `cumulative` or `delta` | string | cumulative |

#### ScalarOID Configuration

| Field Name  | Description                                                    | Value                       | Default |
| --          | --                                                             | --                          | --      |
| `oid`       | The SNMP scalar OID value to grab data from (must end in .0).  | string                      |         |
| `attributes` | The names of the related attribute enum configurations as well as the values to attach to this returned SNMP scalar data. This can be used to have a metric config with multiple ScalarOIDs as different datapoints with different attributue values within the same metric | Attribute              |    |

#### ColumnOID Configuration

| Field Name  | Description                                                    | Value                       | Default |
| --          | --                                                             | --                          | --      |
| `oid`       | The SNMP scalar OID value to grab data from (must end in .0).  | string                      |         |
| `attributes` | The names of the related attribute configurations as well as the enum values to attach to this returned SNMP indexed data if the attribute configuration has enum data. This can be used to attach a specific metric SNMP column OID to an attribute. In doing so, multiple datapoints for a single metric will be created for each returned SNMP indexed data value for the metric along with different attribute values to differentiate them. This also can be used to have a metric config with multiple ColumnOIDs as different datapoints with different attributue values within the same metric | Attribute[]            |    |
| `resource_attributes` | The names of the related resource attribute configurations. This is used to attach a specific metric SNMP column OID to a resource attribute. In doing so, multiple resources will be created for each returned SNMP indexed data value for the metric | string[]              |    |

#### Attribute

| Field Name  | Description                                                    | Value                       | Default |
| --          | --                                                             | --                          | --      |
| `name`      | The name of the attribute configuration that this data refers to | string                     |         |
| `value`     | If the referred to attribute configuration is of enum type, the specific enum value that should be used for this specific attribute | string        |    |

### Example Configuration

```yaml
receivers:
  snmp:
    collection_interval: 60s
    endpoint: udp://localhost:161
    version: v3
    security_level: auth_priv
    user: otel
    auth_type: "MD5"
    auth_password: $SNMP_AUTH_PASSWORD
    privacy_type: "DES"
    privacy_password: $SNMP_PRIVACY_PASSWORD
    
    resource_attributes:
      resource_attr.name.1:
        indexed_value_prefix: probe
      resource_attr.name.2:
        oid: "1.1.1.1"
    
    attributes:
      attr.name.1:
        value: a2_new_key
        enum:
          - in
          - out
      attr.name.2:
        indexed_value_prefix: device
      attr.name.3:
        oid: "2.2.2.2"
    
    metrics:
      # This metric will have multiple datapoints wil 1 attribute on each.
      # Each datapoint will have a (hopefully) different attribute value
      metric.name.1:
        unit: 1
        sum:
          aggregation: cumulative
          monotonic: true
          value_type: int
        column_oids:
          - oid: "2.2.2.1"
            attributes:
              - name: attr.name.3
      # This metric will have multiple datapoints with 2 attributes on each.
      # Each datapoint will have a guaranteed different attribute indexed value for 1 of the attributes.
      # Half of the datapoints will have the other attribute with a value of "in".
      # The other half will have the other attribute with a value of "out".
      metric.name.2:
        unit: "By"
        gauge:
          value_type: int
        column_oids:
          - oid: "3.3.3.3"
            attributes:
              - name: attr.name.2
              - name: attr.name.1
                value: in
          - oid: "2"
            attributes:
              - name: attr.name.2
              - name: attr.name.1
                value: out
      # This metric will have 2 datapoints with 1 attribute on each
      # One datapoints will have an attribute value of "in".
      # The other will have an attribute value of "out".
      metric.name.3:
        unit: "By"
        sum:
          aggregation: delta
          monotonic: false
          value_type: double
        scalar_oids:
          - oid: "4.4.4.4.0"  
            attributes:
              - name: attr.name.1
                value: in
          - oid: "4.4.4.5.0"  
            attributes:
              - name: aattr.name.1
                value: out
      # This metric will have metrics created with each attached to a different resource.
      # Each resource will have a resource attribute with a guaranteed unique value based on the index.
      metric.name.4:
        unit: "By"
        gauge:
          value_type: int
        column_oids:
          - oid: "5.5.5.5"
            resource_attributes:
              - resource_attr.name.1
      # This metric will have metrics created with each attached to a different resource.
      # Each resource will have a resource attribute with a hopefully unique value.
      metric.name.5:
        unit: "By"
        gauge:
          value_type: int
        column_oids:
          - oid: "1.1.1.2"
            resource_attributes:
              - resource_attr.name.2

```

The full list of settings exposed for this receiver are documented [here](./config.go) with detailed sample configurations [here](./testdata/config.yaml).

[development]: https://github.com/open-telemetry/opentelemetry-collector#development
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib