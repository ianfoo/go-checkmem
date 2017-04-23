# go-checkmem

This is a simple program to repeatedly inspect the **expvars** of a Go program,
and report a subset of the runtime.MemStats that is exposed by default. The
full MemStats includes some large arrays that are of limited value when just
trying to determine how much memory a Go program is using. This program will
report a selected subset in JSON format.

By default, it will check every 30s. This can be controlled with the `INTERVAL`
environment variable. The default address that is checked is **localhost:6060**,
and this can be set with the `ADDR` environment variable.

## Notes

This program may also produce log messages of its own, but these are sent to
stderr by default instead of stdout.
