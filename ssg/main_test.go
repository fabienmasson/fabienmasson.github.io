package main

import (
	"html/template"
	"reflect"
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantMeta map[string]string
		wantBody string
	}{
		{
			name: "Basic front-matter",
			input: `---
title: Hello World
date: 2026-03-12
---
Body content`,
			wantMeta: map[string]string{
				"title": "Hello World",
				"date":  "2026-03-12",
			},
			wantBody: "Body content",
		},
		{
			name:     "No front-matter",
			input:    "Just body content",
			wantMeta: map[string]string{},
			wantBody: "Just body content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMeta, gotBody := parseFrontMatter(tt.input)
			if !reflect.DeepEqual(gotMeta, tt.wantMeta) {
				t.Errorf("parseFrontMatter() gotMeta = %v, want %v", gotMeta, tt.wantMeta)
			}
			if gotBody != tt.wantBody {
				t.Errorf("parseFrontMatter() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Basic", "Hello World", "hello-world"},
		{"Special chars", "Hello, World!", "hello-world"},
		{"Multiple spaces", "Hello   World", "hello-world"},
		{"Accents (Current behavior)", "L'été arrive", "l-été-arrive"},
		{"Leading/Trailing", "  Hello World  ", "hello-world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slugify(tt.input); got != tt.want {
				t.Errorf("slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMarkdown(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  template.HTML
	}{
		{
			name:  "Heading",
			input: "# Heading 1",
			want:  template.HTML("<h1>Heading 1</h1>\n"),
		},
		{
			name:  "Bold and Italic",
			input: "**Bold** and *Italic*",
			want:  template.HTML("<p><strong>Bold</strong> and <em>Italic</em></p>\n"),
		},
		{
			name:  "Code block",
			input: "```go\nfmt.Println(\"hi\")\n```",
			want:  template.HTML("<pre><code class=\"language-go\">fmt.Println(&#34;hi&#34;)\n</code></pre>\n"),
		},
		{
			name:  "Unordered list",
			input: "- Item 1\n- Item 2",
			want:  template.HTML("<ul>\n<li>Item 1</li>\n<li>Item 2</li>\n</ul>\n"),
		},
		{
			name:  "Link",
			input: "[Google](https://google.com)",
			want:  template.HTML("<p><a href=\"https://google.com\">Google</a></p>\n"),
		},
		{
			name:  "Paragraph wrapping",
			input: "Line 1\nLine 2",
			want:  template.HTML("<p>Line 1 Line 2</p>\n"),
		},
		{
			name:  "Horizontal rule",
			input: "---",
			want:  template.HTML("<hr>\n"),
		},
		{
			name:  "Horizontal rule with stars",
			input: "***",
			want:  template.HTML("<hr>\n"),
		},
		{
			name:  "Horizontal rule with underscores",
			input: "___",
			want:  template.HTML("<hr>\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMarkdown(tt.input)
			if got != tt.want {
				t.Errorf("parseMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
