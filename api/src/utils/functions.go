package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Chunks string `s` into a list if strings into equal lengths with a max size of `n`
// If the string is not divisible by n, then the items will be as close in length as possible
func ChunkStringEqualUntilN(s string, n int) []string {
	if len(s) == 0 || n <= 0 {
		return []string{}
	}

	totalRuneCount := utf8.RuneCountInString(s) // Count of runes instead of bytes
	numParts := int(math.Ceil(float64(totalRuneCount) / float64(n)))
	evenLength := totalRuneCount / numParts
	extraChars := totalRuneCount % numParts

	var parts []string
	start := 0
	for i := 0; i < numParts; i++ {
		partLength := evenLength
		if i < extraChars {
			partLength++ // Distribute extra characters among the first few parts
		}

		end := start
		count := 0
		for count < partLength && end < len(s) {
			_, size := utf8.DecodeRuneInString(s[end:])
			end += size
			count++
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

func GoogleUUIDPtrToPGXUUID(googleUUID *uuid.UUID) pgtype.UUID {
	var pid pgtype.UUID
	// allow for nil passed
	if googleUUID == nil {
		return pid
	}
	pid.Scan(googleUUID.String())
	return pid
}

func GoogleUUIDToPGXUUID(googleUUID uuid.UUID) pgtype.UUID {
	var pid pgtype.UUID
	pid.Scan(googleUUID.String())
	return pid
}

func PGXUUIDToGoogleUUID(pgxUUID pgtype.UUID) *uuid.UUID {
	if !pgxUUID.Valid {
		return nil
	}
	uid, _ := uuid.FromBytes(pgxUUID.Bytes[:])
	return &uid
}

func PGXUUIDFromString(src string) pgtype.UUID {
	var pid pgtype.UUID
	pid.Scan(src)
	return pid
}

func StringFromPGXUUID(pgxUUID pgtype.UUID) string {
	uid, err := uuid.FromBytes(pgxUUID.Bytes[:])
	if err != nil {
		return ""
	}
	return uid.String()
}

func CleanInput(input string) string {
	// Replace all whitespace characters with a space
	input = strings.Join(strings.Fields(input), " ")

	// Replace multiple consecutive spaces with a single space
	re := regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")

	// remove all invalid utf-8 characters
	valid := make([]rune, 0, len(input))
	for len(input) > 0 {
		r, size := utf8.DecodeRuneInString(input)
		if r != utf8.RuneError {
			valid = append(valid, r)
		}
		input = input[size:]
	}
	return string(valid)
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

// Convert between 2 structs with equal values using reflection
func ReflectStructs[A any, B any](a A) B {
	var b B
	aVal := reflect.ValueOf(a).Elem()
	bPtr := reflect.New(reflect.TypeOf(b).Elem())
	bVal := bPtr.Elem()
	for i := 0; i < aVal.NumField(); i++ {
		bVal.Field(i).Set(aVal.Field(i))
	}
	return bPtr.Interface().(B)
}
