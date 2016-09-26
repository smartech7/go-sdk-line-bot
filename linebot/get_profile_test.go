package linebot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestGetProfile(t *testing.T) {
	type want struct {
		URLPath     string
		RequestBody []byte
		Response    *ProfileResponse
		Error       error
	}
	var testCases = []struct {
		UserID       string
		ResponseCode int
		Response     []byte
		Want         want
	}{
		{
			UserID:       "U0047556f2e40dba2456887320ba7c76d",
			ResponseCode: 200,
			Response:     []byte(`{"userId":"U0047556f2e40dba2456887320ba7c76d","displayName":"BOT API","pictureUrl":"http://dl.profile.line.naver.jp/abcdefghijklmn","statusMessage":"Hello, LINE!"}`),
			Want: want{
				URLPath:     fmt.Sprintf(APIEndpointGetProfile, "U0047556f2e40dba2456887320ba7c76d"),
				RequestBody: []byte(""),
				Response: &ProfileResponse{
					UserID:        "U0047556f2e40dba2456887320ba7c76d",
					DisplayName:   "BOT API",
					PicutureURL:   "http://dl.profile.line.naver.jp/abcdefghijklmn",
					StatusMessage: "Hello, LINE!",
				},
			},
		},
		{
			// Internal server error
			UserID:       "U0047556f2e40dba2456887320ba7c76d",
			ResponseCode: 500,
			Response:     []byte("500 Internal server error"),
			Want: want{
				URLPath:     fmt.Sprintf(APIEndpointGetProfile, "U0047556f2e40dba2456887320ba7c76d"),
				RequestBody: []byte(""),
				Error: &APIError{
					Code: 500,
				},
			},
		},
	}

	var currentTestIdx int
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		tc := testCases[currentTestIdx]
		if r.Method != http.MethodGet {
			t.Errorf("Method %s; want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != tc.Want.URLPath {
			t.Errorf("URLPath %s; want %s", r.URL.Path, tc.Want.URLPath)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(body, tc.Want.RequestBody) {
			t.Errorf("RequestBody %s; want %s", body, tc.Want.RequestBody)
		}
		w.WriteHeader(tc.ResponseCode)
		w.Write(tc.Response)
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		t.Error(err)
	}
	for i, tc := range testCases {
		currentTestIdx = i
		res, err := client.GetProfile(tc.UserID).Do()
		if tc.Want.Error != nil {
			if !reflect.DeepEqual(err, tc.Want.Error) {
				t.Errorf("Error %d %q; want %q", i, err, tc.Want.Error)
			}
		} else {
			if err != nil {
				t.Error(err)
			}
		}
		if !reflect.DeepEqual(res, tc.Want.Response) {
			t.Errorf("Response %d %q; want %q", i, res, tc.Want.Response)
		}
	}
}

func TestGetProfileWithContext(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		time.Sleep(10 * time.Millisecond)
		w.Write([]byte("{}"))
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	_, err = client.GetProfile("U0047556f2e40dba2456887320ba7c76d").WithContext(ctx).Do()
	if err != context.DeadlineExceeded {
		t.Errorf("err %v; want %v", err, context.DeadlineExceeded)
	}
}

func BenchmarkGetProfile(b *testing.B) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Write([]byte(`{"userId":"U","displayName":"A","pictureUrl":"http://","statusMessage":"B"}`))
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.GetProfile("U0047556f2e40dba2456887320ba7c76d").Do()
	}
}
