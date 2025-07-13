package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"

	"github.com/gin-gonic/gin"
)

func PrintJSON(v any) {
	jsonBytes, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(jsonBytes))
}

func CapitalizeFirstLetter(s string) string {
	if s == "" {
		return ""
	}

	// Convert string to a slice of runes to handle Unicode characters correctly
	runes := []rune(s)

	// If the first character is a letter, convert it to uppercase
	if unicode.IsLetter(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
	}

	return string(runes)
}

func BaseResponseError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"data":    nil,
		"message": message,
		"status":  "error",
	})
}

func BaseResponsePaginateSuccess(c *gin.Context, message string, data interface{}, count int, next, prev *string) {
	if data == nil {
		c.JSON(200, gin.H{
			"data": gin.H{
				"items": []interface{}{}, // default empty array
				"count": count,
				"next":  next,
				"prev":  prev,
			},
			"message": message,
			"status":  "ok",
		})
		return
	}

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// If it's a slice and nil, return empty array
	if t.Kind() == reflect.Slice && v.IsNil() {
		c.JSON(200, gin.H{
			"data": gin.H{
				"items": []interface{}{}, // empty JSON array
				"count": count,
				"next":  next,
				"prev":  prev,
			},
			"message": message,
			"status":  "ok",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": gin.H{
			"items": data,
			"count": count,
			"next":  next,
			"prev":  prev,
		},
		"message": message,
		"status":  "ok",
		"error":   nil,
	})
}

func BaseResponseDetailSuccess(c *gin.Context, message string, data interface{}, prevData, nextData *BaseResourceNavigation) {
	c.JSON(200, gin.H{
		"data": gin.H{
			"item": data,
			"prev": prevData,
			"next": nextData,
		},
		"message": message,
		"status":  "success",
	})
}
