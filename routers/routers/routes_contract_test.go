package routers

import (
	"aofs/internal/proto"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitRoute_RegistersExpectedAPIs(t *testing.T) {
	r := InitRoute()
	registered := map[string]struct{}{}
	for _, route := range r.Routes() {
		registered[fmt.Sprintf("%s %s", route.Method, route.Path)] = struct{}{}
	}

	expected := []string{
		"GET /space/v1/api/inner/file/info",
		"POST /space/v1/api/inner/file/infos",
		"GET /space/v1/api/file/info",
		"GET /space/v1/api/file/list",
		"POST /space/v1/api/file/rename",
		"POST /space/v1/api/file/copy",
		"POST /space/v1/api/file/move",
		"POST /space/v1/api/file/delete",
		"GET /space/v1/api/file/download",
		"GET /space/v1/api/file/search",
		"GET /space/v1/api/file/thumb",
		"GET /space/v1/api/file/compressed",
		"POST /space/v1/api/file/vod/symlink",
		"POST /space/v1/api/folder/create",
		"GET /space/v1/api/folder/info",
		"POST /space/v1/api/user/init",
		"POST /space/v1/api/user/delete",
		"GET /space/v1/api/user/storage",
		"GET /space/v1/api/sync/synced",
		"GET /space/v1/api/recycled/clear",
		"POST /space/v1/api/recycled/clear",
		"POST /space/v1/api/recycled/restore",
		"GET /space/v1/api/recycled/list",
		"POST /space/v1/api/multipart/create",
		"POST /space/v1/api/multipart/delete",
		"GET /space/v1/api/multipart/list",
		"POST /space/v1/api/multipart/upload",
		"POST /space/v1/api/multipart/complete",
		"GET /space/v1/api/status",
		"GET /space/v1/api/async/task",
	}

	for _, api := range expected {
		if _, ok := registered[api]; !ok {
			t.Fatalf("missing route: %s", api)
		}
	}
}

func TestStatus_WithoutUserID_ReturnsOK(t *testing.T) {
	r := InitRoute()
	req := httptest.NewRequest(http.MethodGet, "/space/v1/api/status", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected http code: got=%d want=%d", rec.Code, http.StatusOK)
	}

	var rsp proto.Rsp
	if err := json.Unmarshal(rec.Body.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if rsp.Code != proto.CodeOk {
		t.Fatalf("unexpected business code: got=%d want=%d", rsp.Code, proto.CodeOk)
	}
}

func TestFileList_WithoutUserID_ReturnsParamErr(t *testing.T) {
	r := InitRoute()
	req := httptest.NewRequest(http.MethodGet, "/space/v1/api/file/list", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected http code: got=%d want=%d", rec.Code, http.StatusOK)
	}

	var rsp proto.Rsp
	if err := json.Unmarshal(rec.Body.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if rsp.Code != proto.CodeParamErr {
		t.Fatalf("unexpected business code: got=%d want=%d", rsp.Code, proto.CodeParamErr)
	}
}
