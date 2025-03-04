module github.com/Clement-Jean/protein

go 1.24

require (
	github.com/google/go-cmp v0.6.0
	golang.org/x/exp v0.0.0-20250215185904-eff6e970281f
)

require (
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/telemetry v0.0.0-20240521205824-bda55230c457 // indirect
	golang.org/x/tools v0.30.0 // indirect
)

tool (
	golang.org/x/tools/cmd/deadcode
	golang.org/x/tools/cmd/stringer
	golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment
)
