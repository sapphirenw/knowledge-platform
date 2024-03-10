package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"math"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Chunks string `s` into a list if strings into equal lengths with a max size of `n`
// If the string is not divisible by n, then the items will be as close in length as possible
func ChunkStringEqualUntilN(s string, n int) []string {
	if len(s) == 0 || n <= 0 {
		return []string{}
	}

	totalLength := len(s)
	numParts := int(math.Ceil(float64(totalLength) / float64(n)))
	evenLength := totalLength / numParts
	extraChars := totalLength % numParts

	var parts []string
	start := 0
	for i := 0; i < numParts; i++ {
		partLength := evenLength
		if i < extraChars {
			partLength++ // Distribute extra characters among the first few parts
		}
		end := start + partLength
		if end > totalLength {
			end = totalLength
		}
		parts = append(parts, s[start:end])
		start = end
	}

	return parts
}

// Converts a slice of one numeric type to another numeric type using generics.
func ConvertSlice[T any, U any](list []T, convert func(T) U) []U {
	result := make([]U, len(list))
	for i, v := range list {
		result[i] = convert(v)
	}
	return result
}

func GoogleUUIDToPGXUUID(googleUUID uuid.UUID) pgtype.UUID {
	var pid pgtype.UUID
	pid.Scan(googleUUID.String())
	return pid
}

func PGXUUIDToGoogleUUID(pgxUUID pgtype.UUID) (uuid.UUID, error) {
	tmp, err := pgxUUID.Value()
	if err != nil {
		return uuid.New(), err
	}
	return uuid.Parse(tmp.(string))
}

func CleanInput(input string) string {
	// Replace all whitespace characters with a space
	input = strings.Join(strings.Fields(input), " ")

	// Replace multiple consecutive spaces with a single space
	re := regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")

	return input
}

func GenerateFingerprint(input []byte) string {
	hash := sha256.Sum256(input)
	return base64.StdEncoding.EncodeToString(hash[:])
}
