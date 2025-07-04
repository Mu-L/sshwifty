// Sshwifty - A Web SSH client
//
// Copyright (C) 2019-2023 Ni Rui <ranqus@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

//go:generate go run ./static_page_generater ../../.tmp/dist ./static_pages.go
//go:generate go fmt ./static_pages.go

package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/nirui/sshwifty/application/log"
)

type staticData struct {
	data        []byte
	dataHash    string
	compressed  []byte
	created     time.Time
	contentType string
}

func (s staticData) hasCompressed() bool {
	return len(s.compressed) > 0
}

func staticFileExt(fileName string) string {
	extIdx := strings.LastIndex(fileName, ".")
	if extIdx < 0 {
		return ""
	}
	return strings.ToLower(fileName[extIdx:])
}

func serveStaticCacheData(
	dataName string,
	fileExt string,
	w http.ResponseWriter,
	r *http.Request,
	l log.Logger,
) error {
	if fileExt == ".html" || fileExt == ".htm" {
		return ErrNotFound
	}
	return serveStaticCachePage(dataName, w, r, l)
}

func serveStaticCachePage(
	dataName string,
	w http.ResponseWriter,
	r *http.Request,
	l log.Logger,
) error {
	d, dFound := staticPages[dataName]
	if !dFound {
		return ErrNotFound
	}
	selectedData := d.data
	if clientSupportGZIP(r) && d.hasCompressed() {
		selectedData = d.compressed
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("Content-Encoding", "gzip")
	}
	w.Header().Add("Cache-Control", "public, max-age=5184000")
	w.Header().Add("Content-Type", d.contentType)
	_, wErr := w.Write(selectedData)
	return wErr
}

func serveStaticPage(
	dataName string,
	code int,
	w http.ResponseWriter,
	r *http.Request,
	l log.Logger,
) error {
	d, dFound := staticPages[dataName]
	if !dFound {
		return ErrNotFound
	}
	selectedData := d.data
	if clientSupportGZIP(r) && d.hasCompressed() {
		selectedData = d.compressed
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("Content-Encoding", "gzip")
	}
	w.Header().Add("Content-Type", d.contentType)
	w.WriteHeader(code)
	_, wErr := w.Write(selectedData)
	return wErr
}
