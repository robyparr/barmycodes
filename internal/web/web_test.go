package web_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/robyparr/barmycodes/internal"
	"github.com/robyparr/barmycodes/internal/web"
)

func TestRouter(t *testing.T) {
	router := web.NewRouter()

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

func escapeBase64String(str string) string {
	return strings.ReplaceAll(str, "+", "&#43;")
}
