# SONDA

[![CI](https://github.com/KarasuJager/sonda/actions/workflows/ci.yml/badge.svg)](https://github.com/KarasuJager/sonda/actions/workflows/ci.yml)

SONDA is a focused command-line tool for reproducible HTTP differential testing.
It sends a baseline `GET` request, changes one existing query parameter, sends
the mutated request, and reports the observable differences.

> observe | mutate | compare | reproduce

SONDA records evidence. It does not decide whether a difference is a
vulnerability. A changed status, response body, or timing can have many causes;
security meaning depends on application context and human investigation.

## What differential HTTP testing means

Differential testing compares two controlled observations of the same system.
For SONDA v1.0, the first observation is the original URL and the second is the
same URL with one query parameter replaced. SONDA compares:

- HTTP status code
- response body length
- SHA-256 body fingerprint
- Sorensen-Dice body similarity over byte bigrams
- elapsed request time

The output describes what changed, not why it changed. This separation keeps
the evidence reproducible and avoids presenting heuristics as security facts.

## Requirements

- Go 1.27 or newer
- Python 3 only for the optional local lab

SONDA uses only the Go standard library.

## Installation and build

Clone the repository and build the CLI:

```sh
git clone https://github.com/KarasuJager/sonda.git
cd sonda
go build -o bin/sonda ./cmd/sonda
```

On Windows:

```powershell
go build -o bin\sonda.exe ./cmd/sonda
```

After a v1.0.0 tag is published, Go can install that exact release directly:

```sh
go install github.com/KarasuJager/sonda/cmd/sonda@v1.0.0
```

## CLI usage

Collect only a baseline fingerprint:

```sh
sonda --url "https://example.test/search?q=admin"
```

Compare a baseline with a query-parameter mutation:

```sh
sonda \
  --url "https://example.test/search?q=admin&page=1" \
  --param q \
  --value SONDA_TEST
```

The named parameter must already exist in the URL. SONDA replaces all values
for that parameter with the supplied value and preserves unrelated parameters.

Options:

| Option | Meaning |
| --- | --- |
| `--url` | Target URL. Required. |
| `--param` | Existing query parameter to mutate. |
| `--value` | Replacement value. Required with `--param`. |
| `--header` | Custom `Name: Value` request header. Repeatable. |
| `--timeout` | Request timeout in seconds. Default: `10`. |
| `--proxy` | HTTP proxy URL, such as `http://127.0.0.1:8080`. |
| `--json` | Emit the differential report as JSON. Requires a mutation. |

Every request uses `User-Agent: SONDA/1.0`.

## JSON output

Use `--json` to write a machine-readable report to standard output:

```sh
sonda \
  --url "http://127.0.0.1:9090/?q=admin" \
  --param q \
  --value SONDA_TEST \
  --json > result.json
```

The report contains the original and mutated URLs, mutation metadata, both
fingerprints, and the computed differences. Body similarity is a number from
`0` (no shared byte bigrams) to `1` (identical bodies). Timing fields are
milliseconds and can be negative in the diff when the mutated request is
faster. See [`examples/result.json`](examples/result.json) for a complete
example; exact timing values vary by run.

## Custom headers

Repeat `--header` to reproduce authentication, content negotiation, virtual
hosting, or other request context:

```sh
sonda \
  --url "https://example.test/search?q=admin" \
  --param q \
  --value SONDA_TEST \
  --header "Authorization: Bearer test-token" \
  --header "Accept-Language: en" \
  --header "Host: virtual.example.test"
```

Headers are applied identically to the baseline and mutated requests. Avoid
placing secrets in shell history or shared reports.

## Burp Suite integration

With Burp's proxy listener on `127.0.0.1:8080`, route both requests through it:

```sh
sonda \
  --url "http://127.0.0.1:9090/?q=admin" \
  --param q \
  --value SONDA_TEST \
  --proxy "http://127.0.0.1:8080"
```

The target remains the application URL; `--proxy` identifies Burp. Keep
interception off for unattended runs, or forward both requests manually. For
HTTPS targets, Go must trust Burp's CA certificate through the operating-system
trust configuration.

## Local deterministic lab

The included lab returns stable bodies for known `q` values and listens only on
`127.0.0.1:9090`.

Start it in one terminal:

```sh
python lab/lab.py
```

Then run SONDA in another:

```sh
go run ./cmd/sonda \
  --url "http://127.0.0.1:9090/?q=admin" \
  --param q \
  --value SONDA_TEST
```

Use `py lab\lab.py` on Windows if the Python launcher provides that command.
The response content is deterministic; measured timings are not.

## Architecture

```text
cmd/sonda/              CLI parsing and orchestration
internal/httpclient/    GET requests, headers, timeout, and proxy handling
internal/mutation/      Query-parameter replacement
internal/fingerprint/   Status, length, SHA-256, and timing snapshot
internal/similarity/    Byte-bigram Sorensen-Dice coefficient
internal/diff/          Snapshot and body comparison
internal/report/        JSON report model and encoding
lab/                    Deterministic local HTTP lab
examples/               Example output
```

The packages follow the runtime flow: observe the baseline, mutate the URL,
observe again, compare the two observations, and render the result.

## Testing and validation

Run the release checks from the repository root:

```sh
gofmt -w .
go vet ./...
go test ./...
go build ./cmd/sonda
```

The GitHub Actions workflow checks formatting, runs `go vet`, and executes the
test suite on each push and pull request.

## Authorized use

Use SONDA only on systems you own or have explicit permission to test. You are
responsible for complying with applicable laws, contracts, rules of engagement,
and rate limits.

## v1.0 limitations

SONDA v1.0 is intentionally narrow:

- HTTP `GET` only
- one existing query parameter mutated per run
- no request bodies, crawling, payload database, or multi-step workflows
- no vulnerability scanning or automatic vulnerability classification
- redirects follow Go's standard client behavior
- response bodies over 10 MiB are rejected
- timing measurements include network and proxy variability
- TLS verification uses the operating-system trust configuration
- reports contain differential metadata, not response bodies or request headers

These constraints are part of the frozen v1.0 scope, not an invitation to infer
findings that SONDA does not produce.

## License

SONDA is available under the [MIT License](LICENSE).
