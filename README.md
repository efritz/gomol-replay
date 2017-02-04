# Gomol-Replay

[![GoDoc](https://godoc.org/github.com/efritz/gomol-replay?status.svg)](https://godoc.org/github.com/efritz/gomol-replay)
[![Build Status](https://secure.travis-ci.org/efritz/gomol-replay.png)](http://travis-ci.org/efritz/gomol-replay)
[![codecov.io](http://codecov.io/github/efritz/gomol-replay/coverage.svg?branch=master)](http://codecov.io/github/efritz/gomol-replay?branch=master)

A logging adapter for [gomol](https://github.com/aphistic/gomol) to support replaying
a sequence of log messages but at a higher log level. 

*Intended use case:* Each request in an HTTP server has a unique log adapter which
traces the request. This adapter generally logs at the DEBUG level. When a request
encounters an error or is being served slowly, the entire trace can be replayed at
a higher level so the entire context is available for analysis.

## Example

```go
adapter := NewReplayAdapter(
    logger,           // gomol logger or adapter
    gomol.LevelDebug, // track debug messages for replay
    gomol.LevelInfo,  // also track info messages
)

// ...

if requestIsTakingLong() {
    // Re-log journaled messages at warning level
    adapter.Replay(gomol.LevelWarning)
}
```

Messages which are replayed at a higher level will keep the original message timestamp
(if supplied), or use the time the `Log` message was invoked (if not supplied). Each 
message will also be sent with an additional attribute called `replayed-from-level` with
a value equal to the original level of the message.

## License

Copyright (c) 2017 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
