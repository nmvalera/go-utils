package httptestutils

import (
	http "net/http"

	kkrtgomock "github.com/kkrt-labs/go-utils/gomock"
	gomock "go.uber.org/mock/gomock"
	"gopkg.in/h2non/gock.v1"
)

// NewGockRequest creates a new gock request
func NewGockRequest() *gock.Request {
	req := gock.NewRequest()
	req.Response = gock.NewResponse()
	return req
}

// DoGock declares a call to Do with a mocked gock request and responder
func (mr *MockSenderMockRecorder) DoGock(req *gock.Request) *gomock.Call {
	return mr.Do(kkrtgomock.GockMatcher(req)).DoAndReturn(func(r *http.Request) (*http.Response, error) { return gock.Responder(r, req.Response, nil) })
}
