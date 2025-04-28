package coverage

import (
	"bufio"
	"bytes"
	"fmt"
)

func mergeCoverages(coverages [][]byte) ([]byte, error) {
	var merged bytes.Buffer
	headerWritten := false

	for _, coverage := range coverages {
		scanner := bufio.NewScanner(bytes.NewReader(coverage))
		for scanner.Scan() {
			line := scanner.Bytes()

			if len(line) == 0 {
				continue
			}

			if bytes.HasPrefix(line, []byte("mode:")){
				if headerWritten {
					continue
				}
				merged.Write(line)
				merged.WriteByte('\n')
				headerWritten = true
				continue
			}

			merged.Write(line)
			merged.WriteByte('\n')
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to scan coverage content: %w", err)
		}
	}

	return merged.Bytes(), nil
}
