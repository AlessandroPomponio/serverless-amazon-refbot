// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package urlwork contains the functions to handle urls
// and to generate referral links.
package urlwork

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/retgits/bitly/client"
	"github.com/retgits/bitly/client/bitlinks"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/repository"
)

//GetRefURL tries to generate an Amazon referral link.
func GetRefURL(link string, referral string, b *client.Client) (string, error) {

	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", errors.Errorf("%s is not a valid URL: %s", link, err)
	}

	// Default to http scheme in case the field is missing.
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}

	parsedURL, err = unshortenURL(parsedURL)
	if err != nil {
		return "", errors.Errorf("Unable to unshorten URL %s: %s", link, err)
	}

	//It has to be an AmazonDomain URL
	if !strings.HasSuffix(parsedURL.Host, repository.AmazonDomain) {
		return "", errors.Errorf("Amazon domain not supported for URL %s", parsedURL.String())
	}

	// Build the referral URL.
	parsedURL.Path = cutPathAtASIN(parsedURL.Path)
	parsedURL.RawQuery = "&tag=" + referral
	parsedURL.Fragment = ""

	bLinks := bitlinks.New(b)
	return shortenURL(parsedURL.String(), bLinks)

}

// cutPathAtASIN returns a copy of the provided path up to the ASIN.
func cutPathAtASIN(path string) string {

	var builder strings.Builder

	// We want to rebuild the path segment by segment.
	// We will stop once we've reached the ASIN code.
	pathSegments := strings.Split(path, "/")
	for index, part := range pathSegments {

		builder.WriteString(part)
		builder.WriteString("/")

		if index > 0 && (pathSegments[index-1] == "product" || pathSegments[index-1] == "dp" || pathSegments[index-1] == "d") {
			break
		}

	}

	return builder.String()

}

// shortenURL shortens a URL using Bitly.
func shortenURL(u string, b *bitlinks.Bitlinks) (string, error) {

	toShort := bitlinks.ShortenRequest{
		GroupGUID: "",
		LongURL:   u,
		Domain:    "bit.ly",
	}

	shortURL, err := b.ShortenLink(&toShort)
	if err != nil {
		err = errors.Errorf("Error while shortening the URL: %s", err)
		return "", err
	}

	return shortURL.Link, nil

}

// unshortenURL performs a HEAD request to unshorten an URL.
func unshortenURL(u *url.URL) (*url.URL, error) {

	result, err := http.Head(u.String())
	if err != nil {
		return nil, err
	}

	return result.Request.URL, nil

}
