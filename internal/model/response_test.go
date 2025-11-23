package model

import (
	"testing"
)

func TestSuccessResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	response := SuccessResponse(data)

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Data == nil {
		t.Error("expected data to be set")
	}

	if response.Error != nil {
		t.Error("expected error to be nil")
	}
}

func TestSuccessResponseWithMeta(t *testing.T) {
	data := []string{"item1", "item2"}
	meta := &Meta{
		Page:       1,
		PerPage:    10,
		Total:      100,
		TotalPages: 10,
	}

	response := SuccessResponseWithMeta(data, meta)

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Data == nil {
		t.Error("expected data to be set")
	}

	if response.Meta == nil {
		t.Error("expected meta to be set")
	}

	if response.Meta.Page != 1 {
		t.Errorf("expected page 1, got %d", response.Meta.Page)
	}

	if response.Meta.Total != 100 {
		t.Errorf("expected total 100, got %d", response.Meta.Total)
	}
}

func TestErrorResponse(t *testing.T) {
	response := ErrorResponse("ERROR_CODE", "error message")

	if response.Success {
		t.Error("expected success to be false")
	}

	if response.Data != nil {
		t.Error("expected data to be nil")
	}

	if response.Error == nil {
		t.Fatal("expected error to be set")
	}

	if response.Error.Code != "ERROR_CODE" {
		t.Errorf("expected code ERROR_CODE, got %s", response.Error.Code)
	}

	if response.Error.Message != "error message" {
		t.Errorf("expected message 'error message', got %s", response.Error.Message)
	}
}
