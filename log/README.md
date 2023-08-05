# logs.gigfinder.util

Provides a thread safe and fast logging system.

## Exposed API

### Values

#### `const int` V
logging level 0

#### `const int` VV
logging level 1

#### `const int` VVV
logging level 2

#### `var int` Verbosity
The package's logging level


### Functions

#### Start `() void`
Opens the background logging system.

#### Close `() void`
Safely closes the background logging system.

#### Msg `(level int, text string) void`
Logs a message at the specified level

#### Msgf `(level int, text string, args ...interface{}) void`
Logs a message at the specified level, and formats it using the Printf syntax.