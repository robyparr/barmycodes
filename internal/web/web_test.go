package web_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/robyparr/barmycodes/internal"
	"github.com/robyparr/barmycodes/internal/web"
)

var nowFuncStub = func() time.Time { return time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC) }

func TestMainPage(t *testing.T) {
	router := web.NewRouter(time.Now)

	t.Run("it returns an empty page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		body := response.Body.String()

		assertStatus(t, response.Code, http.StatusOK)
		assertBodyContains(t, body, "Bar My Codes")
		assertBodyContains(t, body, "Your barcodes will showup here.")
	})

	t.Run("it returns Code128 codes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/?type=Code128&b[]=test1&b[]=test2", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		body := response.Body.String()
		wantCodes, _ := internal.GenerateBarcodes([]string{"test1", "test2"}, "Code128")

		assertStatus(t, response.Code, http.StatusOK)
		assertBodyContains(t, body, escapeBase64String(wantCodes[0].Base64()))
		assertBodyContains(t, body, escapeBase64String(wantCodes[1].Base64()))
	})
}

func TestDownloadPNG(t *testing.T) {
	router := web.NewRouter(time.Now)
	request, _ := http.NewRequest(http.MethodGet, "/png?type=Code128&b[]=Test", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	assertStatus(t, response.Code, http.StatusOK)
	assertBodyIsFixture(t, response.Body.Bytes(), "../testdata/barcode_Test.png")
}

func TestDownloadPDF(t *testing.T) {
	router := web.NewRouter(nowFuncStub)

	request, _ := http.NewRequest(http.MethodGet, "/pdf?type=Code128&b[]=Test", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	assertStatus(t, response.Code, http.StatusOK)
	assertBodyIsFixture(t, response.Body.Bytes(), "../testdata/barcode_Test.pdf")
}

func TestUnknownURL(t *testing.T) {
	router := web.NewRouter(time.Now)
	request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	assertStatus(t, response.Code, http.StatusNotFound)
	assertBodyContains(t, response.Body.String(), "404: The page you're looking for could not be found.")
}

func assertStatus(t *testing.T, got int, want int) {
	if got != want {
		t.Errorf("Unexpected response status. Got %d, want %d.", got, want)
	}
}

func assertBodyContains(t *testing.T, body string, str string) {
	if !strings.Contains(body, str) {
		t.Errorf("Expected response body to contain '%s' but did not: \n %s", str, body)
	}
}

func assertBodyIsFixture(t *testing.T, body []byte, fixturePath string) {
	fixture, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("error loading fixture: %s", err)
	}

	if !bytes.Equal(body, fixture) {
		t.Errorf("Body does not match fixture '%s'", fixturePath)
	}
}

func escapeBase64String(str string) string {
	return strings.ReplaceAll(str, "+", "&#43;")
}
