# AGENTS.md

## Project purpose

SONDA is a small Go tool for reproducible HTTP differential testing. Its
operating philosophy is:

> observe | mutate | compare | reproduce

SONDA gathers differential evidence for human investigation. It does not label
response differences as vulnerabilities.

## Frozen v1.0 scope

Preserve the existing v1.0 feature set:

- HTTP GET baseline
- mutation of one existing query parameter
- status and response-length comparison
- SHA-256 body fingerprinting
- byte-bigram Sorensen-Dice body similarity
- timing delta
- repeatable custom headers
- configurable timeout and HTTP proxy
- Burp Suite interoperability
- JSON reports
- deterministic local Python lab

Do not add AI features, vulnerability classification, a payload database,
automatic scanning, or new offensive capabilities. Avoid new dependencies and
product redesign. Behavior changes require a concrete correctness or security
reason and focused regression coverage.

## Engineering conventions

- Module: `github.com/KarasuJager/sonda`
- User-Agent: `SONDA/1.0`
- Burp proxy convention: `127.0.0.1:8080`
- Local lab convention: `127.0.0.1:9090`
- Keep runtime code in the standard library unless a dependency is essential.
- Keep stdout clean in JSON mode; send diagnostics to stderr.
- Preserve deterministic report fields. Treat timings as observations, not
  conclusions.
- Never commit generated binaries, local reports, caches, or credentials.
- Tests must use local deterministic fixtures such as `httptest`; they must not
  depend on public network services.
- Use only systems that the operator owns or is explicitly authorized to test.

## Validation

Run all checks from the repository root before handing off a change:

```sh
gofmt -w .
go vet ./...
go test ./...
go build ./cmd/sonda
```

For a Linux cross-build from another platform:

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/sonda-linux-amd64 ./cmd/sonda
```

On PowerShell, set the same values through `$env:CGO_ENABLED`, `$env:GOOS`, and
`$env:GOARCH` before running `go build`.

Do not commit automatically unless the user explicitly requests it. Report all
validation failures without hiding or rewriting them.
