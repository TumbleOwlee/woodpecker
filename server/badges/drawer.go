// Copyright 2023 The narqo/go-badge Authors. All rights reserved.
// SPDX-License-Identifier: MIT.

package badges

// cspell:words Verdana

import (
	"bytes"
	"html/template"
	"io"
	"sync"

	"go.woodpecker-ci.org/woodpecker/v3/server/badges/fonts"
)

type badge struct {
	Subject string
	Status  string
	Color   Color
	Bounds  bounds
}

type bounds struct {
	// SubjectDx is the width of subject string of the badge.
	SubjectDx float64
	SubjectX  float64
	// StatusDx is the width of status string of the badge.
	StatusDx float64
	StatusX  float64
}

func (b bounds) Dx() float64 {
	return b.SubjectDx + b.StatusDx
}

type badgeDrawer struct {
	tmpl  *template.Template
	mutex *sync.Mutex
}

func (d *badgeDrawer) Render(subject, status string, color Color, w io.Writer) error {
	d.mutex.Lock()
	subjectDx := d.measureString(subject)
	statusDx := d.measureString(status)
	d.mutex.Unlock()

	bdg := badge{
		Subject: subject,
		Status:  status,
		Color:   color,
		Bounds: bounds{
			SubjectDx: subjectDx,
			SubjectX:  subjectDx/2.0 + 1,
			StatusDx:  statusDx,
			StatusX:   subjectDx + statusDx/2.0 - 1,
		},
	}
	return d.tmpl.Execute(w, bdg)
}

func (d *badgeDrawer) RenderBytes(subject, status string, color Color) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := d.Render(subject, status, color, buf)
	return buf.Bytes(), err
}

// shields.io uses Verdana.ttf to measure text width with an extra 10px.
// As we use DejaVuSans.ttf, we have to tune this value a little.
const extraDx = 5

func (d *badgeDrawer) measureString(s string) float64 {
	width := 0
	for _, r := range s {
		if value, exists := fonts.GlyphWidths[r]; exists {
			width += value
		} else {
			width += fonts.GlyphDefaultWidth
		}
	}
	SHIFT := 6
	return float64(width>>SHIFT) + extraDx
}

// RenderBytes renders a badge of the given color, with given subject and status to bytes.
func RenderBytes(subject, status string, color Color) ([]byte, error) {
	drawer, err := initDrawer()
	if err != nil {
		return nil, err
	}
	return drawer.RenderBytes(subject, status, color)
}

var (
	drawer    *badgeDrawer
	initError error
	initOnce  sync.Once
)

func initDrawer() (*badgeDrawer, error) {
	initOnce.Do(func() {
		drawer = &badgeDrawer{
			tmpl:  template.Must(template.New("flat-template").Parse(flatTemplate)),
			mutex: &sync.Mutex{},
		}
		initError = nil
	})
	return drawer, initError
}
