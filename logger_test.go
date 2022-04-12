package log

import (
	"bytes"
	"context"
	"github.com/c2fo/testify/assert"
	"log"
	"net/http"
	"testing"
	"time"
)

type obj struct {
	Name    string `json:"name"`
	Count   int    `json:"count"`
	Enabled bool   `json:"enabled"`
}

var (
	sampleObjects = []*obj{
		{"a", 1, true},
		{"b", 2, false},
		{"c", 3, true},
		{"d", 4, false},
		{"e", 5, true},
		{"f", 6, false},
		{"g", 7, true},
		{"h", 8, false},
		{"i", 9, true},
		{"j", 0, false},
	}
	sampleArray   = make([]int, 10000)
	sampleString  = "some string with a somewhat realistic length"
	sampleDataMap = map[string]interface{}{
		"a": "a",
		"b": 1000,
		"c": time.Now(),
	}
)

var (
	sampleContext      context.Context
	requestWithContext *http.Request
	logger             Logger
)

func init() {
	logger = NewLogger(sampleString)
	sampleContext = logger.BuildContextDataAndSetValue("11")

	request, _ := http.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(`{}`)))

	data := make(map[string]string, 0)
	data["http_method"] = request.Method
	data["language"] = "language_code"

	requestWithContext = logger.SetContextDataAndSetValue(request, data, "12")
}

func BenchmarkLogNative_InfoSimple(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Println(sampleContext, sampleString)
	}
}

func BenchmarkLog_InfoSimple(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info(sampleContext, sampleString)
	}
}

func BenchmarkLog_InfoSimpleWithFormat(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Infof(sampleContext, "%s", sampleString)
	}
}

func BenchmarkLogNative_LogLargeArray(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Println(sampleContext, sampleArray)
	}
}

func BenchmarkLog_LogLargeArray(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info(sampleContext, sampleArray)
	}
}

func BenchmarkLogNative_InfoWithComplexArgs(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Println(sampleContext, sampleObjects)
	}
}

func BenchmarkLog_InfoWithComplexArgs(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info(sampleContext, sampleObjects)
	}
}

func BenchmarkLog_InfoMapWithComplexArgs(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.InfoMap(sampleContext, sampleDataMap, sampleObjects)
	}
}

func BenchmarkLog_LogRequest(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.LogRequest(requestWithContext.Context(), requestWithContext)
	}
}

func TestNoCollisionWhenBuildContextData(t *testing.T) {
	type thisKeyType string
	var (
		thisKey      thisKeyType = thisKeyType(ContextDataMapKey) // to ensure the value equal
		thisKeyValue             = "muhammad_fakhri"

		randomID = "randomID"
	)

	ctx := logger.BuildContextDataAndSetValue(randomID)
	newCtx := context.WithValue(ctx, thisKey, thisKeyValue)
	contextDataFromLogger := newCtx.Value(ContextDataMapKey).(map[string]string)

	// Ensure there is no collision
	assert.Equal(t, thisKeyValue, newCtx.Value(thisKey))
	assert.Equal(t, randomID, contextDataFromLogger[ContextIdKey])
}
