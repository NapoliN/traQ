package router

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
)

func TestPostFile(t *testing.T) {
	e, cookie, mw, assert, _ := beforeTest(t)

	body, boundary := createFormFile(t)

	req := httptest.NewRequest("POST", "http://test", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)

	if cookie != nil {
		req.Header.Add("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := mw(PostFile)(c)
	assert.NoError(err)

	assert.Equal(http.StatusCreated, rec.Code, rec.Body.String())
}

func TestGetFileByID(t *testing.T) {
	e, cookie, mw, assert, _ := beforeTest(t)

	file := mustMakeFile(t)

	c, rec := getContext(e, t, cookie, nil)
	c.SetPath("/:fileID")
	c.SetParamNames("fileID")
	c.SetParamValues(file.ID)

	requestWithContext(t, mw(GetFileByID), c)
	if assert.EqualValues(http.StatusOK, rec.Code) {
		assert.Equal("test message", rec.Body.String())
	}

	c, rec = getContext(e, t, cookie, nil)
	c.SetPath("/:fileID")
	c.SetParamNames("fileID")
	c.SetParamValues(file.ID)
	c.Request().URL.RawQuery = "dl=1"

	requestWithContext(t, mw(GetFileByID), c)
	if assert.EqualValues(http.StatusOK, rec.Code) {
		assert.EqualValues(fmt.Sprintf("attachment; filename=%s", file.Name), rec.Header().Get(echo.HeaderContentDisposition))

	}

}
func TestDeleteFileByID(t *testing.T) {
	e, cookie, mw, assert, _ := beforeTest(t)

	file := mustMakeFile(t)

	c, rec := getContext(e, t, cookie, nil)
	c.SetPath("/:fileID")
	c.SetParamNames("fileID")
	c.SetParamValues(file.ID)

	requestWithContext(t, mw(DeleteFileByID), c)
	assert.EqualValues(http.StatusNoContent, rec.Code, rec.Body.String())
}
func TestGetMetaDataByFileID(t *testing.T) {
	e, cookie, mw, assert, _ := beforeTest(t)

	file := mustMakeFile(t)

	c, rec := getContext(e, t, cookie, nil)
	c.SetPath("/:fileID")
	c.SetParamNames("fileID")
	c.SetParamValues(file.ID)

	requestWithContext(t, mw(GetMetaDataByFileID), c)
	assert.EqualValues(http.StatusOK, rec.Code, rec.Body.String())
}

func createFormFile(t *testing.T) (*bytes.Buffer, string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", "test.txt")
	require.NoError(t, err)

	fh, err := os.Open("../LICENSE")
	require.NoError(t, err)
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	require.NoError(t, err)

	bodyWriter.Close()
	return bodyBuf, bodyWriter.Boundary()
}