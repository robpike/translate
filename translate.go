// Copyright 2012 The rspace Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Translate uses the Google translate API from the command line to translate
// its arguments. By default it auto-detects the input language and translates
// to English.

package main // import "robpike.io/cmd/translate"

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	// https://developers.google.com/translate/v2/using_rest#supported-query-params
	key    = flag.String("key", "", "Google API key (defaults to $GOOGLEAPIKEY)")
	target = flag.String("to", "en", "destination language (two-letter code)")
	source = flag.String("from", "", "source language (two-letter code); auto-detected by default")
)

type Response struct {
	Data struct {
		Translations []Translation
	}
}

type Translation struct {
	TranslatedText         string
	DetectedSourceLanguage string
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	v := make(url.Values)
	if *key == "" {
		*key = os.Getenv("GOOGLEAPIKEY")
	}
	if *key == "" {
		log.Fatal("$GOOGLEAPIKEY not set")
	}
	v.Set("key", *key)
	v.Set("target", *target)
	if *source != "" {
		v.Set("source", *source)
	}
	v.Set("q", strings.Join(flag.Args(), " "))
	url := "https://www.googleapis.com/language/translate/v2?" + v.Encode()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		log.Fatal(err)
	}
	for _, t := range r.Data.Translations {
		fmt.Printf("%s (%s)\n", html.UnescapeString(t.TranslatedText), t.DetectedSourceLanguage)
	}
}
