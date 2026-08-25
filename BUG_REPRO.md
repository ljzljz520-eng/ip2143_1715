# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	workshopnotice/cmd/workshop	0.002s
ok  	workshopnotice/internal/domain	0.001s
ok  	workshopnotice/internal/flow	0.023s
--- FAIL: Test2143BusinessRegression (0.01s)
    regression_test.go:26: second notice number = 101, want 202
FAIL
FAIL	workshopnotice/internal/flow022	0.008s
?   	workshopnotice/internal/importer	[no test files]
ok  	workshopnotice/internal/httpapi	0.006s
ok  	workshopnotice/internal/report	0.001s
ok  	workshopnotice/internal/store	0.005s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/workshop): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/workshop): exit `0`
- Frontend build (web): exit `0`
