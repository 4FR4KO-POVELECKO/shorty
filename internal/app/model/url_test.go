package model_test

import (
	"ShortURL/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL_ValidateURL(t *testing.T) {
	testCases := []struct {
		name    string
		payload *model.URL
		err     bool
	}{
		{
			name: "valid",
			payload: &model.URL{
				OriginURL: "google.com",
				ShortURL:  "JL_Cfwf951",
			},
			err: false,
		},
		{
			name: "empty",
			payload: &model.URL{
				OriginURL: "",
				ShortURL:  "",
			},
			err: true,
		},
		{
			name: "invalid origin url",
			payload: &model.URL{
				OriginURL: "googlecom",
				ShortURL:  "XYZabc123_",
			},
			err: true,
		},
		{
			name: "short url < 10",
			payload: &model.URL{
				OriginURL: "google.com",
				ShortURL:  "JLCfwf_95",
			},
			err: true,
		},
		{
			name: "short url > 10",
			payload: &model.URL{
				OriginURL: "google.com",
				ShortURL:  "XYZabc1234_",
			},
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.payload.ValidateURL()
			eq := false

			if err != nil {
				eq = assert.Error(t, err)
			}

			assert.Equal(t, tc.err, eq)
		})
	}
}

func TestURL_ValidateOriginURL(t *testing.T) {
	testCases := []struct {
		name    string
		payload *model.URL
		err     bool
	}{
		{
			name: "valid",
			payload: &model.URL{
				OriginURL: "google.com",
			},
			err: false,
		},
		{
			name: "empty",
			payload: &model.URL{
				OriginURL: "",
			},
			err: true,
		},
		{
			name: "invalid origin url",
			payload: &model.URL{
				OriginURL: "googlecom",
			},
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.payload.ValidateOriginURL()
			eq := false

			if err != nil {
				eq = assert.Error(t, err)
			}

			assert.Equal(t, tc.err, eq)
		})
	}
}

func TestURL_ValidateShortURL(t *testing.T) {
	testCases := []struct {
		name    string
		payload *model.URL
		err     bool
	}{
		{
			name: "valid",
			payload: &model.URL{
				ShortURL: "JL_Cfwf951",
			},
			err: false,
		},
		{
			name: "empty",
			payload: &model.URL{
				ShortURL: "",
			},
			err: true,
		},
		{
			name: "short url < 10",
			payload: &model.URL{
				OriginURL: "google.com",
				ShortURL:  "JLCfwf_95",
			},
			err: true,
		},
		{
			name: "short url > 10",
			payload: &model.URL{
				OriginURL: "google.com",
				ShortURL:  "XYZabc1234_",
			},
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.payload.ValidateShortURL()
			eq := false

			if err != nil {
				eq = assert.Error(t, err)
			}

			assert.Equal(t, tc.err, eq)
		})
	}
}
