# const

Generate Go constants with helper methods from a plain text list. The tool can also write a matching Swift enum.

You write this:

```
package status

//name=Status type=int

new      // title{New} desc{Waiting for review} iota=1
active   // title{Active} #visible code{200}
banned   // title{Banned} desc{Banned by moderator} code{403}
```

And get a gofmt-ed `status.go` with the type, the constants and helpers: `String()`, `Title()`, `Desc()`, `Valid()`, `StatusIndexes()`, `StatusIndex("active")`, `IsVisible()`, `Code()`, and so on. Each generated file starts with a comment that lists every function it contains.

## Install

```sh
go install github.com/monopolly/const@latest
```

## Usage

```sh
const status.const        # writes status.go (and status.swift if #swift is set)
```

The output files sit next to the input and share its name with the extension replaced, so `status.const` becomes `status.go`. Don't give the input a `.go` extension: the tool refuses to overwrite its own input.

You can also run it with `go generate`:

```go
//go:generate const status.const
```

## Input format

```
package <go package>

//<list options>

<value>   // <field options>
<value>   // <field options>
```

- The first line is the Go package of the generated file.
- The first `//` line holds the list options.
- Each following line is one value. Anything after `//` holds the field options and is kept as the constant's comment.
- Empty lines are ignored.

### List options

| option | description |
|---|---|
| `name=Status` | Creates a named type `type Status int`. Helpers become methods `(a Status) Title()`, and package level functions get the type prefix: `StatusIndexes()`, `StatusValid()`, `StatusIndex()`. This lets you keep several lists in one package. Without `name=`, constants are plain `int` or `string` values and helpers are plain functions `Title(a int)`. |
| `type=string` | Constant type: `int` (default), `int8`…`int64`, `uint`…`uint64`, `string` |
| `#simple` | Constants only, no functions |
| `#up` | Upper case Go names: `NEW`, `SENSOR_ALARM_ANY` |
| `#index` | Extra int constants `StatusIndexNew`, `StatusIndexActive`… starting from 1 |
| `#swift` | Also write a `.swift` enum |

### Field options

| option | example | result |
|---|---|---|
| `title{}` | `title{Real Estate}` | `Title()`. Falls back to the Go name when not set |
| `desc{}` | `desc{Houses and land}` | `Desc()`. Generated only if at least one value has `desc{}` |
| `list{}` | `list{brand, model}` | `List() []string`. Generated only if used |
| `name{}` | `art // name{artwork}` | Replaces the string value |
| `go{}` | `art // go{Artwork}` | Replaces the Go constant name |
| `iota=N` | `transferred // iota=101` | Starts a new const block from N (int types) |
| `#tag` or `@tag` | `#virtual` | `IsVirtual()`, `Virtual()`, `StatusVirtualIndexes()`, `StatusVirtualList()` |
| `key{number}` | `code{401}` | `Code() int` |
| `key{text}` | `label{Premium}` | `Label() string` |
| `key{true/false}` | `paid{true}` | `Paid() bool` |
| `@key=value` | `@code=401` | Same as `key{value}`. The value can't contain spaces |
| `_` | `_` | Skips a value but keeps its iota position |

#### Custom functions

Any `key{value}` that isn't a built-in option becomes a function. The first letter of the key is capitalized: `somefeature{134}` becomes `Somefeature() int`. The return type is picked from the values:

- `int` if all values are numbers
- `bool` if all values are `true` or `false`
- `string` otherwise, including mixed values

Values that don't set the key return `0`, `""` or `false`.

```
free     // limit{10}  paid{false} label{Starter}
pro      // limit{100} paid{true}
```

```go
func (a Plan) Limit() int     // Free: 10, Pro: 100
func (a Plan) Paid() bool     // Pro: true
func (a Plan) Label() string  // Free: "Starter"
```

## Any strings as values

With `type=string`, a value can be any string without spaces, for example a NATS subject. To use spaces or `//`, put the value in quotes: `"wine & spirits" // title{Wine}`. The tool builds a readable Go name like this:

- Words are split on any character that is not a letter or digit (`.` `-` `_` `$` `{}` `/` …).
- `*` becomes `Any` and `>` becomes `All` (NATS wildcards).
- Each word is capitalized. Upper case words stay upper case (`BEMT`), unless the whole value is upper case.
- A leading digit gets an `N` prefix.

| value | Go | `#up` | Swift |
|---|---|---|---|
| `sensor.alarm.smoke` | `SensorAlarmSmoke` | `SENSOR_ALARM_SMOKE` | `sensorAlarmSmoke` |
| `sensor.alarm.*` | `SensorAlarmAny` | `SENSOR_ALARM_ANY` | `sensorAlarmAny` |
| `sensor.>` | `SensorAll` | `SENSOR_ALL` | `sensorAll` |
| `myapp.us-east.users.profile.update` | `MyappUsEastUsersProfileUpdate` | `MYAPP_US_EAST_USERS_PROFILE_UPDATE` | `myappUsEastUsersProfileUpdate` |
| `obs.gnss.BEMT` | `ObsGnssBEMT` | `OBS_GNSS_BEMT` | `obsGnssBEMT` |
| `wasmbus.ctl.v1.{lattice}.scale` | `WasmbusCtlV1LatticeScale` | `WASMBUS_CTL_V1_LATTICE_SCALE` | `wasmbusCtlV1LatticeScale` |
| `$SYS.REQ.SERVER.PING` | `SysReqServerPing` | `SYS_REQ_SERVER_PING` | `sysReqServerPing` |
| `user_id` | `UserId` | `USER_ID` | `userId` |

The constant value is always the original string: `SensorAlarmAny = Subject("sensor.alarm.*")`. If two values produce the same name (`orders.*` and `orders.any`), generation stops with an error. Rename one of them with `go{}`.

## Generated functions

For `name=Status type=int`:

| function | description |
|---|---|
| `(Status) String() string` | String value |
| `(Status) Title() string` | `title{}`, or the Go name |
| `(Status) Desc() string` | `desc{}`, only if used |
| `(Status) Int() int` | Value as int (int types) |
| `(Status) Valid() bool` | Is the value known |
| `StatusValid(v int) bool` | Is a raw value known |
| `StatusValidIndexType(v Status) bool` | Is the value known |
| `StatusIndexes() []Status` | All values |
| `StatusIndex(v string) int` | String to value, `0` if unknown (int types) |
| `(Status) List() []string` | `list{}`, only if used |
| `Status<Tag>Indexes()`, `Status<Tag>List()` | Values with `#tag` |
| `(Status) Is<Tag>() bool`, `(Status) <Tag>() bool` | Has `#tag` |
| `(Status) <Key>() int / string / bool` | Custom `key{value}` |

The comment at the top of each generated file lists exactly what was generated, with signatures.

## Swift

`#swift` writes an enum with `rawValue`, `id`, `title`, `desc` (only if used) and a `Bool` property for every tag. Swift keywords are escaped (`` `public` ``).

```swift
enum Subject: String {
	case sensorAlarmSmoke = "sensor.alarm.smoke", sensorAlarmAny = "sensor.alarm.*", ...
	var id: String { ... }
	var title: String { ... }
}
```

## Examples

The [examples](examples) directory has inputs and their generated output:

- [status](examples/status): int list with `iota`, tags, `desc{}` and custom functions
- [subject](examples/subject): NATS subjects as string values
- [category](examples/category): string list with titles, `list{}`, bool, int and string custom functions, and a quoted value

## Development

```sh
go test ./...            # regenerates examples/* and compares with committed output
go test -update          # rewrite examples after changing the generator
go build ./...           # makes sure the generated examples compile
```
