package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"net/url"
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
	return hex.EncodeToString(hash[:])
}

func GenerateRandomFingerprint() string {
	return GenerateFingerprint([]byte(GenerateRandomString(64)))
}

// Parses the protocol and the domain name from the website
func ParseWebsiteInformation(inputURL string) (protocol string, domain string, err error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", "", err
	}
	domain = u.Hostname()
	protocol = u.Scheme
	return protocol, domain, nil
}

func GenerateRandomString(n int) string {
	const letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}
