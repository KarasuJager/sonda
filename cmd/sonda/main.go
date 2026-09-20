package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KarasuJager/sonda/internal/diff"
	"github.com/KarasuJager/sonda/internal/fingerprint"
	"github.com/KarasuJager/sonda/internal/httpclient"
	"github.com/KarasuJager/sonda/internal/mutation"
	"github.com/KarasuJager/sonda/internal/report"
)

type headerFlags []string

func (h *headerFlags) String() string {
	return strings.Join(*h, ", ")
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func main() {
	target := flag.String("url", "", "Target URL")
	param := flag.String("param", "", "Query parameter to mutate")
	value := flag.String("value", "", "Mutation value")

	jsonOutput := flag.Bool(
		"json",
		false,
		"Output differential report as JSON",
	)

	timeoutSeconds := flag.Int(
		"timeout",
		10,
		"HTTP timeout in seconds",
	)

	proxyURL := flag.String(
		"proxy",
		"",
		"HTTP proxy URL, for example http://127.0.0.1:8080",
	)

	var rawHeaders headerFlags

	flag.Var(
		&rawHeaders,
		"header",
		`Custom HTTP header. Repeatable. Example: --header "Authorization: Bearer token"`,
	)

	flag.Parse()

	if *target == "" {
		fmt.Fprintln(os.Stderr, "error: --url is required")
		os.Exit(1)
	}

	if *timeoutSeconds <= 0 {
		fmt.Fprintln(
			os.Stderr,
			"error: --timeout must be greater than zero",
		)
		os.Exit(1)
	}

	if *jsonOutput && *param == "" {
		fmt.Fprintln(
			os.Stderr,
			"error: --json requires --param and --value",
		)
		os.Exit(1)
	}

	if *param != "" && *value == "" {
		fmt.Fprintln(
			os.Stderr,
			"error: --value is required when --param is used",
		)
		os.Exit(1)
	}

	headers, err := parseHeaders(rawHeaders)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error: %v\n",
			err,
		)
		os.Exit(1)
	}

	httpConfig := httpclient.Config{
		Timeout:  time.Duration(*timeoutSeconds) * time.Second,
		ProxyURL: *proxyURL,
		Headers:  headers,
	}

	if !*jsonOutput {
		fmt.Println("SONDA")
		fmt.Println("observe | mutate | compare | reproduce")
		fmt.Println()

		fmt.Printf("[*] Target: %s\n", *target)
		fmt.Printf("[*] Timeout: %ds\n", *timeoutSeconds)

		if *proxyURL != "" {
			fmt.Printf("[*] Proxy: %s\n", *proxyURL)
		}

		if len(headers) > 0 {
			fmt.Printf("[*] Custom headers: %d\n", len(headers))
		}
	}

	baselineResult, err := httpclient.Get(
		*target,
		httpConfig,
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"[!] Baseline request failed: %v\n",
			err,
		)
		os.Exit(1)
	}

	baseline := fingerprint.Generate(
		baselineResult.StatusCode,
		baselineResult.Body,
		baselineResult.Duration,
	)

	if *param == "" {
		printFingerprint("Baseline", baseline)
		return
	}

	mutatedURL, err := mutation.QueryParam(
		*target,
		*param,
		*value,
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"[!] Mutation failed: %v\n",
			err,
		)
		os.Exit(1)
	}

	mutatedResult, err := httpclient.Get(
		mutatedURL,
		httpConfig,
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"[!] Mutated request failed: %v\n",
			err,
		)
		os.Exit(1)
	}

	mutated := fingerprint.Generate(
		mutatedResult.StatusCode,
		mutatedResult.Body,
		mutatedResult.Duration,
	)

	comparison := diff.Compare(
		baseline,
		mutated,
		baselineResult.Body,
		mutatedResult.Body,
	)

	if *jsonOutput {
		r := report.New(
			*target,
			mutatedURL,
			*param,
			*value,
			baseline,
			mutated,
			comparison,
		)

		data, err := report.JSON(r)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"[!] JSON report failed: %v\n",
				err,
			)
			os.Exit(1)
		}

		fmt.Println(string(data))
		return
	}

	printFingerprint("Baseline", baseline)

	fmt.Println()
	fmt.Printf(
		"[*] Mutation: %s=%s\n",
		*param,
		*value,
	)
	fmt.Printf(
		"[*] Mutated URL: %s\n",
		mutatedURL,
	)

	printFingerprint("Mutated", mutated)

	fmt.Println()
	fmt.Println("[Diff]")

	fmt.Printf(
		"Status changed: %t (%d -> %d)\n",
		comparison.StatusChanged,
		comparison.StatusBefore,
		comparison.StatusAfter,
	)

	fmt.Printf(
		"Length changed: %t (%d -> %d, delta %+d)\n",
		comparison.LengthChanged,
		comparison.LengthBefore,
		comparison.LengthAfter,
		comparison.LengthDelta,
	)

	fmt.Printf(
		"SHA-256 changed: %t\n",
		comparison.HashChanged,
	)

	fmt.Printf(
		"Body similarity: %.2f%%\n",
		comparison.BodySimilarity*100,
	)

	fmt.Printf(
		"Time delta: %+d ms\n",
		comparison.TimingDeltaMillis,
	)
}

func parseHeaders(
	rawHeaders []string,
) (http.Header, error) {
	headers := make(http.Header)

	for _, raw := range rawHeaders {
		name, value, ok := strings.Cut(raw, ":")

		if !ok {
			return nil, fmt.Errorf(
				"invalid header %q: expected Name: Value",
				raw,
			)
		}

		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)

		if name == "" {
			return nil, fmt.Errorf(
				"invalid header %q: header name is empty",
				raw,
			)
		}

		headers.Add(name, value)
	}

	return headers, nil
}

func printFingerprint(
	name string,
	fp fingerprint.Fingerprint,
) {
	fmt.Println()
	fmt.Printf("[%s]\n", name)
	fmt.Printf("Status:   %d\n", fp.StatusCode)
	fmt.Printf("Length:   %d bytes\n", fp.BodyLength)
	fmt.Printf("SHA-256:  %s\n", fp.SHA256)
	fmt.Printf("Time:     %d ms\n", fp.DurationMillis)
}
