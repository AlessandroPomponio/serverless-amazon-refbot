// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package urlwork

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/repository"
)

func getPathFromString(rawURL string, t *testing.T) string {

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		t.Errorf("getPathFromString: unable to parse raw url %s: %s", rawURL, err)
	}

	return parsedURL.Path

}

func getURLFromString(rawURL string, t *testing.T) *url.URL {

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		t.Errorf("getURLFromString: unable to parse raw url %s: %s", rawURL, err)
	}

	return parsedURL

}

func Test_cutPathAtASIN(t *testing.T) {

	repository.AmazonDomain = "amazon.it"

	firstURL := "https://www.amazon.it/Buono-Regalo-Amazon-it-Da-stampare/dp/B005VEAJK6/a-lot-of-things-that-should-not-be-in-the-output-of-the-function"
	secondURL := "https://www.amazon.it/dp/B078WST5RK/a-lot-of-things-that-should-not-be-in-the-output-of-the-function"
	thirdURL := "https://www.amazon.it/gp/aw/d/B0794VJ18B/a-lot-of-things-that-should-not-be-in-the-output-of-the-function"
	fourthURL := "https://www.amazon.it/gp/product/B0794VJ18B/a-lot-of-things-that-should-not-be-in-the-output-of-the-function"

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "domain/name/dp/ASIN",
			path: getPathFromString(firstURL, t),
			want: "/Buono-Regalo-Amazon-it-Da-stampare/dp/B005VEAJK6/",
		},
		{
			name: "domain/dp/ASIN",
			path: getPathFromString(secondURL, t),
			want: "/dp/B078WST5RK/",
		},
		{
			name: "domain/gp/aw/d/ASIN",
			path: getPathFromString(thirdURL, t),
			want: "/gp/aw/d/B0794VJ18B/",
		},
		{
			name: "domain/gp/product/ASIN",
			path: getPathFromString(fourthURL, t),
			want: "/gp/product/B0794VJ18B/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cutPathAtASIN(tt.path); got != tt.want {
				t.Errorf("cutPathAtASIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unshortenURL(t *testing.T) {

	tests := []struct {
		name    string
		args    *url.URL
		want    *url.URL
		wantErr bool
	}{
		{
			name:    "Amazon short URL https",
			args:    getURLFromString("https://amzn.to/2lVEfGs", t),
			want:    getURLFromString("https://www.amazon.it/gp/product/B0794VJ18B/a-lot-of-things-that-should-not-be-in-the-output-of-the-function", t),
			wantErr: false,
		},
		{
			// We want an error here but unshortenURL will be called
			// by GetRefURL, which will perform the appropriate checks.
			name:    "Amazon short URL no scheme",
			args:    getURLFromString("amzn.to/2lVEfGs", t),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Amazon short URL http",
			args:    getURLFromString("http://amzn.to/2lVEfGs", t),
			want:    getURLFromString("https://www.amazon.it/gp/product/B0794VJ18B/a-lot-of-things-that-should-not-be-in-the-output-of-the-function", t),
			wantErr: false,
		},
		{
			name:    "Bitly self link",
			args:    getURLFromString("https://bitly.is/1g3AhR6", t),
			want:    getURLFromString("https://bitly.com/", t),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unshortenURL(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("unshortenURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unshortenURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
