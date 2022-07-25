package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ControllerTest struct {
	name       string
	method     string
	url        string
	statuscode int
	reqBody    *bytes.Buffer
	pre        func()
	post       func()
	ExpectErr  bool
}

func TestGetEmployee(t *testing.T) {

	temp := jsonMarshalFunc
	//jsonStr := []byte(`{"id":3,"name":"test1","DOJ": "2006-03-02T15:04:05Z","skillset":["go","testing"]}`)
	tests := []ControllerTest{
		{
			name:       "Get All Employees Successfull",
			method:     http.MethodGet,
			statuscode: http.StatusOK,
			url:        "/portal/api/v1/employee",
		},
		{
			name:       "Get Single Employee Successfull",
			method:     http.MethodGet,
			statuscode: http.StatusOK,
			url:        "/portal/api/v1/employee/test1",
		},
		{
			name:       "Get Single employee success but not found",
			method:     http.MethodGet,
			statuscode: http.StatusOK,
			url:        "/portal/api/v1/employee/test123",
		},
		{
			name:       "Get Single employee internal server error",
			method:     http.MethodGet,
			statuscode: http.StatusInternalServerError,
			url:        "/portal/api/v1/employee",
			pre: func() {
				jsonMarshalFunc = func(v interface{}) ([]byte, error) {
					return nil, errors.New("error while marshalling")
				}
			},
			post: func() {
				jsonMarshalFunc = temp
			},
		},
		{
			name:       "Delete Single Employee Successfull",
			method:     http.MethodDelete,
			statuscode: http.StatusOK,
			url:        "/portal/api/v1/employee/test1",
		},
		/* PUT TEST NEEDS TO BE MOCK
		{
			name:       "Update Existing Employee Successfull",
			method:     http.MethodPut,
			statuscode: http.StatusNoContent,
			url:        "/portal/api/v1/employee/test1",
			reqBody:    bytes.NewBuffer(jsonStr),
		},*/

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.pre != nil {
				tt.pre()
			}
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal("")
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := Handlers()
			handler.ServeHTTP(rr, req)
			if st := rr.Code; st != tt.statuscode {
				t.Errorf("test name :: %v , expected status %d got %d", tt.name, tt.statuscode, st)
			}
			if tt.post != nil {
				tt.post()
			}
		})
	}
}

//MonkeyPatching
//needs to be completed

func TestPostEmployee(t *testing.T) {
	jsonStr := []byte(`{"id":1,"name":"testUser","skillset":["go","testing"]}`)
	emptyJson := []byte(``)
	testcases := []ControllerTest{
		{
			name:       "Post Employee Success",
			method:     http.MethodPost,
			url:        "/portal/api/v1/employee",
			statuscode: http.StatusOK,
			reqBody:    bytes.NewBuffer(jsonStr),
		},
		{
			name:       "Post Employee Bad Request",
			method:     http.MethodPost,
			url:        "/portal/api/v1/employee",
			statuscode: http.StatusBadRequest,
			reqBody:    bytes.NewBuffer(emptyJson),
		},
		{
			name:       "Post Employee method not allowed",
			method:     http.MethodPost,
			url:        "/portal/api/v1/employ",
			statuscode: http.StatusNotFound,
			reqBody:    bytes.NewBuffer(jsonStr),
		},
		{
			name:       "Post Employee method not allowed",
			method:     http.MethodPatch,
			url:        "/portal/api/v1/employee",
			statuscode: http.StatusMethodNotAllowed,
			reqBody:    bytes.NewBuffer(jsonStr),
		},
	}

	for _, tt := range testcases {
		req, err := http.NewRequest(tt.method, tt.url, tt.reqBody)
		if err != nil {
			t.Fatal("")
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := Handlers()
		handler.ServeHTTP(rr, req)
		if st := rr.Code; st != tt.statuscode {
			t.Errorf("test name :: %v , expected status %d got %d", tt.name, tt.statuscode, st)
		}
	}
}
