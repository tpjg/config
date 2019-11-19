# config

Configuration of Golang apps from environment flags.

Reads the configuration into a struct, allows easy application configuration and clean code (single argument with complete config).

# example

```go
import (
    "github.com/tpjg/config"
    "fmt"
)

type TestStruct struct {
	S          string        `flag:"TEST_STRING"`
	I          int           `flag:"TEST_INT"`
	I64        int64         `flag:"TEST_INT64"`
	F64        float64       `flag:"TEST_FLOAT64"`
	D          time.Duration `flag:"TEST_DURATION"`
	UI         uint          `flag:"TEST_UINT"`
	UI64       uint64        `flag:"TEST_UINT64"`
	unexported string        `flag:"CANNOT_BE_SET"`
}

func main() {
    ts := TestStruct{}
	config.ReadStructFromEnv(&ts)

    fmt.Printf("Config = %v\n", ts)
}
```

To  allow overridding the flags with arguments on the command line use:

```go
config.ReadStructFromEnvOverrideWithArgs(&ts)
```
