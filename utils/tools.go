package utils

import (
	"encoding/json"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ConvertStructData(object interface{}, targetStruct interface{}) interface{} {
	// Use reflection to get the type of the targetted struct
	targetStructType := reflect.TypeOf(targetStruct)
	// Create a new variable of the same type as the targetted struct
	targetStructData := reflect.New(targetStructType).Interface()
	// Convert object to fit into the targetted struct
	temporaryVariable, _ := json.Marshal(object)
	json.Unmarshal(temporaryVariable, &targetStructData)

	return targetStructData
}

func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomStr := make([]byte, length)
	for i := range randomStr {
		randomStr[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomStr)
}

// Generates a random integer with a specified number of digits
func GetRandomInt(size int) uint32 {
	if size <= 0 {
		return 0
	}

	// Calculate the min and max possible values for the specified size
	min := intPow(10, size-1)
	max := intPow(10, size) - 1

	// Initialize the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random integer within the range [min, max]
	return uint32(rand.Intn(max-min+1) + min)
}

// intPow calculates the power of base^exponent for integers
func intPow(base, exponent int) int {
	result := 1
	for i := 0; i < exponent; i++ {
		result *= base
	}
	return result
}

// PASSWORD HASHING
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	log.Println(err)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// UUID PARSER
func ParseUUID(input string) (*uuid.UUID, *ErrorResponse) {
	uuidVal, err := uuid.Parse(input)
	if err != nil {
		errData := RequestErr(ERR_INVALID_VALUE, "Invalid UUID")
		return nil, &errData
	}
	return &uuidVal, nil
}

// Check if keys exist in map
func KeysExistInMap(keys []string, myMap map[string]interface{}) bool {
	for _, key := range keys {
		if _, ok := myMap[key]; !ok {
			return false
		}
	}
	return true
}
