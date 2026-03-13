package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

type Post struct {
	Title       string
	Slug        string
	Date        time.Time
	Description string
	Tags        []string
	Content     template.HTML
	RawContent  string
}

type Page struct {
	Title      string
	Slug       string
	Content    template.HTML
	RawContent string
}

type PageData struct {
	SiteName  string
	PageTitle string
	MetaDesc  string
	Year      int
	Body      template.HTML
}

type IndexData struct {
	Posts       []Post
	CurrentPage int
	TotalPages  int
	PrevPage    string
	NextPage    string
}

const SiteName = "F4N"

var (
	reCodeSpan      = regexp.MustCompile("`([^`]+)`")
	reBoldItalic    = regexp.MustCompile(`\*\*\*(.+?)\*\*\*`)
	reBold          = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reItalic        = regexp.MustCompile(`\*(.+?)\*`)
	reStrikethrough = regexp.MustCompile(`~~(.+?)~~`)
	reLink          = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reImage         = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	reHR            = regexp.MustCompile(`^(\-{3,}|\*{3,}|_{3,})$`)
	reUL            = regexp.MustCompile(`^[\-\*\+] `)
	reOL            = regexp.MustCompile(`^\d+\. `)
	reSlug          = regexp.MustCompile(`-+`)

	// Templates
	tmplLayout = template.Must(template.New("layout").Parse(layoutTmpl))
	tmplIndex  = template.Must(template.New("index").Parse(indexTmpl))
	tmplPost   = template.Must(template.New("post").Parse(postTmpl))
)

// ─────────────────────────────────────────────
//  Markdown parser (zero dependencies)
// ─────────────────────────────────────────────

func parseMarkdown(src string) template.HTML {
	lines := strings.Split(src, "\n")
	var buf bytes.Buffer
	inCode := false
	inUL := false
	inOL := false
	inBlockquote := false
	inPara := false
	var paraBuf strings.Builder

	inlineFormat := func(s string) string {
		// Code spans
		s = reCodeSpan.ReplaceAllString(s, "<code>$1</code>")
		// Bold + italic
		s = reBoldItalic.ReplaceAllString(s, "<strong><em>$1</em></strong>")
		// Bold
		s = reBold.ReplaceAllString(s, "<strong>$1</strong>")
		// Italic
		s = reItalic.ReplaceAllString(s, "<em>$1</em>")
		// Strikethrough
		s = reStrikethrough.ReplaceAllString(s, "<del>$1</del>")
		// Links
		s = reLink.ReplaceAllString(s, `<a href="$2">$1</a>`)
		// Images
		s = reImage.ReplaceAllString(s, `<img src="$2" alt="$1">`)
		return s
	}

	flushPara := func() {
		if inPara {
			content := strings.TrimSpace(paraBuf.String())
			if content != "" {
				buf.WriteString("<p>" + inlineFormat(template.HTMLEscapeString(content)) + "</p>\n")
			}
			paraBuf.Reset()
			inPara = false
		}
	}

	closeList := func() {
		if inUL {
			buf.WriteString("</ul>\n")
			inUL = false
		}
		if inOL {
			buf.WriteString("</ol>\n")
			inOL = false
		}
	}

	closeBlockquote := func() {
		if inBlockquote {
			buf.WriteString("</blockquote>\n")
			inBlockquote = false
		}
	}

	for _, line := range lines {
		// Fenced code blocks
		if strings.HasPrefix(line, "```") {
			flushPara()
			if inCode {
				buf.WriteString("</code></pre>\n")
				inCode = false
			} else {
				closeList()
				closeBlockquote()
				lang := strings.TrimPrefix(line, "```")
				if lang != "" {
					buf.WriteString(fmt.Sprintf(`<pre><code class="language-%s">`, template.HTMLEscapeString(lang)))
				} else {
					buf.WriteString("<pre><code>")
				}
				inCode = true
			}
			continue
		}
		if inCode {
			buf.WriteString(template.HTMLEscapeString(line) + "\n")
			continue
		}

		// Horizontal rule
		if reHR.MatchString(strings.TrimSpace(line)) {
			flushPara()
			closeList()
			closeBlockquote()
			buf.WriteString("<hr>\n")
			continue
		}

		// Blockquote
		if strings.HasPrefix(line, "> ") {
			flushPara()
			closeList()
			if !inBlockquote {
				buf.WriteString("<blockquote>\n")
				inBlockquote = true
			}
			buf.WriteString("<p>" + inlineFormat(template.HTMLEscapeString(line[2:])) + "</p>\n")
			continue
		}
		closeBlockquote()

		// Headings
		if strings.HasPrefix(line, "#") {
			flushPara()
			if strings.HasPrefix(line, "###### ") {
				closeList()
				buf.WriteString("<h6>" + inlineFormat(template.HTMLEscapeString(line[7:])) + "</h6>\n")
				continue
			}
			if strings.HasPrefix(line, "##### ") {
				closeList()
				buf.WriteString("<h5>" + inlineFormat(template.HTMLEscapeString(line[6:])) + "</h5>\n")
				continue
			}
			if strings.HasPrefix(line, "#### ") {
				closeList()
				buf.WriteString("<h4>" + inlineFormat(template.HTMLEscapeString(line[5:])) + "</h4>\n")
				continue
			}
			if strings.HasPrefix(line, "### ") {
				closeList()
				buf.WriteString("<h3>" + inlineFormat(template.HTMLEscapeString(line[4:])) + "</h3>\n")
				continue
			}
			if strings.HasPrefix(line, "## ") {
				closeList()
				buf.WriteString("<h2>" + inlineFormat(template.HTMLEscapeString(line[3:])) + "</h2>\n")
				continue
			}
			if strings.HasPrefix(line, "# ") {
				closeList()
				buf.WriteString("<h1>" + inlineFormat(template.HTMLEscapeString(line[2:])) + "</h1>\n")
				continue
			}
		}

		// Unordered list
		if reUL.MatchString(line) {
			flushPara()
			closeBlockquote()
			if !inUL {
				if inOL {
					buf.WriteString("</ol>\n")
					inOL = false
				}
				buf.WriteString("<ul>\n")
				inUL = true
			}
			buf.WriteString("<li>" + inlineFormat(template.HTMLEscapeString(line[2:])) + "</li>\n")
			continue
		}

		// Ordered list
		if reOL.MatchString(line) {
			flushPara()
			closeBlockquote()
			if !inOL {
				if inUL {
					buf.WriteString("</ul>\n")
					inUL = false
				}
				buf.WriteString("<ol>\n")
				inOL = true
			}
			idx := strings.Index(line, ". ")
			buf.WriteString("<li>" + inlineFormat(template.HTMLEscapeString(line[idx+2:])) + "</li>\n")
			continue
		}

		closeList()

		// Empty line → paragraph break
		if strings.TrimSpace(line) == "" {
			flushPara()
			buf.WriteString("\n")
			continue
		}

		if inPara {
			paraBuf.WriteString(" ")
		}
		paraBuf.WriteString(line)
		inPara = true
	}

	flushPara()
	closeList()
	closeBlockquote()
	if inCode {
		buf.WriteString("</code></pre>\n")
	}

	return template.HTML(buf.String())
}

// ─────────────────────────────────────────────
//  Front-matter parser
// ─────────────────────────────────────────────

func parseFrontMatter(content string) (map[string]string, string) {
	meta := map[string]string{}
	if !strings.HasPrefix(content, "---") {
		return meta, content
	}
	rest := content[3:]
	end := strings.Index(rest, "\n---")
	if end == -1 {
		return meta, content
	}
	front := rest[:end]
	body := rest[end+4:]
	for _, line := range strings.Split(front, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			meta[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return meta, strings.TrimSpace(body)
}

// ─────────────────────────────────────────────
//  Slug helper
// ─────────────────────────────────────────────

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	slug := reSlug.ReplaceAllString(b.String(), "-")
	return strings.Trim(slug, "-")
}

// ─────────────────────────────────────────────
//  Post loader
// ─────────────────────────────────────────────

func loadPosts(dir string) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading posts dir: %w", err)
	}
	var posts []Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s: %v\n", path, err)
			continue
		}
		meta, body := parseFrontMatter(string(data))

		title := meta["title"]
		if title == "" {
			title = strings.TrimSuffix(e.Name(), ".md")
		}
		slug := meta["slug"]
		if slug == "" {
			slug = slugify(strings.TrimSuffix(e.Name(), ".md"))
		}
		desc := meta["description"]

		var tags []string
		if t := meta["tags"]; t != "" {
			for _, tag := range strings.Split(t, ",") {
				tags = append(tags, strings.TrimSpace(tag))
			}
		}

		var date time.Time
		if d := meta["date"]; d != "" {
			var err error
			date, err = time.Parse("2006-01-02", d)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: invalid date format in %s: %v\n", path, err)
				if info, err := e.Info(); err == nil {
					date = info.ModTime()
				}
			}
		} else {
			if info, err := e.Info(); err == nil {
				date = info.ModTime()
			}
		}

		posts = append(posts, Post{
			Title:       title,
			Slug:        slug,
			Date:        date,
			Description: desc,
			Tags:        tags,
			Content:     parseMarkdown(body),
			RawContent:  body,
		})
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})
	return posts, nil
}

func loadPages(dir string) ([]Page, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading pages dir: %w", err)
	}
	var pages []Page
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		meta, body := parseFrontMatter(string(data))
		title := meta["title"]
		if title == "" {
			title = strings.TrimSuffix(e.Name(), ".md")
		}
		slug := meta["slug"]
		if slug == "" {
			slug = slugify(strings.TrimSuffix(e.Name(), ".md"))
		}
		pages = append(pages, Page{
			Title:      title,
			Slug:       slug,
			Content:    parseMarkdown(body),
			RawContent: body,
		})
	}
	return pages, nil
}

const baseCSS = `
@import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:ital,wght@0,300;0,400;0,500;0,600;1,300;1,400&display=swap');

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

:root {
  --bg:       #f7f6f3;
  --surface:  #eeecea;
  --ink:      #181714;
  --ink-2:    #6a6762;
  --accent:   #c8502a;
  --border:   #e0ddd7;
  --radius:   3px;
  --max:      700px;
  --mono:     'IBM Plex Mono', monospace;
}

html { scroll-behavior: smooth; }

body {
  background: var(--bg);
  color: var(--ink);
  font-family: var(--mono);
  font-size: 14px;
  line-height: 1.8;
  font-weight: 300;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* ── Layout ── */
.site-header {
  padding: 2.5rem 2rem 2rem;
  border-bottom: 1px solid var(--border);
}
.header-inner {
  max-width: var(--max);
  margin: 0 auto;
  display: flex;
  align-items: baseline;
  gap: 2.5rem;
}
.site-title {
  font-family: var(--mono);
  font-size: 1rem;
  font-weight: 600;
  letter-spacing: -0.01em;
  color: var(--ink);
  text-decoration: none;
}
.site-title:hover { color: var(--accent); }

nav { display: flex; gap: 1.5rem; }
nav a {
  font-size: 0.82rem;
  font-weight: 500;
  letter-spacing: 0.07em;
  text-transform: uppercase;
  color: var(--ink-2);
  text-decoration: none;
  transition: color .15s;
}
nav a:hover { color: var(--accent); }

main {
  flex: 1;
  max-width: var(--max);
  width: 100%;
  margin: 0 auto;
  padding: 4rem 2rem 6rem;
}

footer {
  border-top: 1px solid var(--border);
  padding: 2rem;
  text-align: center;
  font-size: 0.8rem;
  color: var(--ink-2);
}

/* ── Index ── */
.page-heading {
  font-family: var(--mono);
  font-size: clamp(1.6rem, 4vw, 2.2rem);
  font-weight: 300;
  line-height: 1.2;
  letter-spacing: -0.03em;
  margin-bottom: 3rem;
  color: var(--ink);
}
.page-heading em {
  font-style: italic;
  color: var(--accent);
}

.post-list { list-style: none; }

.post-item {
  display: grid;
  grid-template-columns: 5rem 1fr;
  gap: 0 1.5rem;
  align-items: start;
  padding: 1.75rem 0;
  border-bottom: 1px solid var(--border);
}
.post-item:first-child { border-top: 1px solid var(--border); }

.post-date-col {
  padding-top: .18rem;
  font-size: 0.72rem;
  font-weight: 500;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--ink-2);
  line-height: 1.4;
}
.post-date-col .day { font-size: 1.6rem; font-weight: 300; letter-spacing: -0.04em; color: var(--ink); line-height: 1; display: block; margin-bottom: .1rem; }

.post-title-link {
  font-family: var(--mono);
  font-size: 1rem;
  font-weight: 500;
  letter-spacing: -0.02em;
  line-height: 1.4;
  color: var(--ink);
  text-decoration: none;
  display: block;
  margin-bottom: .4rem;
  transition: color .15s;
}
.post-title-link:hover { color: var(--accent); }

.post-description {
  font-size: 0.9rem;
  color: var(--ink-2);
  line-height: 1.6;
  margin-bottom: .6rem;
}

.tags { display: flex; flex-wrap: wrap; gap: .4rem; }
.tag {
  font-size: 0.68rem;
  font-weight: 500;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  background: var(--surface);
  color: var(--ink-2);
  padding: .15rem .55rem;
  border-radius: 2px;
  border: 1px solid var(--border);
  text-decoration: none;
  transition: background .15s, color .15s;
}
.tag:hover { background: var(--accent); color: #fff; border-color: var(--accent); }

/* ── Single post ── */
.post-header { margin-bottom: 3rem; }
.post-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.78rem;
  font-weight: 500;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--ink-2);
  margin-bottom: 1rem;
}
.post-meta .sep { opacity: .35; }

.post-title {
  font-family: var(--mono);
  font-size: clamp(1.4rem, 3.5vw, 1.9rem);
  font-weight: 400;
  line-height: 1.2;
  letter-spacing: -0.03em;
  margin-bottom: 1.2rem;
}

.post-description-lead {
  font-size: 1.05rem;
  color: var(--ink-2);
  line-height: 1.65;
  border-left: 3px solid var(--accent);
  padding-left: 1rem;
}

/* ── Prose ── */
.prose p { margin-bottom: 1.4rem; }
.prose h1, .prose h2, .prose h3, .prose h4, .prose h5, .prose h6 {
  font-family: var(--mono);
  font-weight: 500;
  letter-spacing: -0.02em;
  line-height: 1.3;
  margin: 2.5rem 0 .9rem;
  color: var(--ink);
}
.prose h2 { font-size: 1.65rem; }
.prose h3 { font-size: 1.3rem; }
.prose h4 { font-size: 1.1rem; }

.prose a { color: var(--accent); text-underline-offset: 3px; }
.prose a:hover { opacity: .75; }

.prose strong { font-weight: 500; }
.prose em { font-style: italic; }
.prose del { opacity: .5; }

.prose ul, .prose ol {
  padding-left: 1.4rem;
  margin-bottom: 1.4rem;
}
.prose li { margin-bottom: .3rem; }

.prose blockquote {
  border-left: 3px solid var(--accent);
  margin: 1.8rem 0;
  padding: .6rem 1.2rem;
  background: var(--surface);
  border-radius: 0 var(--radius) var(--radius) 0;
  color: var(--ink-2);
  font-style: italic;
}

.prose code {
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: .82em;
  background: var(--surface);
  border: 1px solid var(--border);
  padding: .1em .4em;
  border-radius: var(--radius);
}

.prose pre {
  background: var(--ink);
  color: #e8e4df;
  border-radius: 6px;
  padding: 1.4rem 1.5rem;
  overflow-x: auto;
  margin: 1.8rem 0;
  font-size: .83em;
  line-height: 1.7;
}
.prose pre code {
  background: none;
  border: none;
  padding: 0;
  color: inherit;
  font-size: 1em;
}

.prose hr {
  border: none;
  border-top: 1px solid var(--border);
  margin: 2.5rem 0;
}

.prose img { max-width: 100%; border-radius: 6px; margin: 1.5rem 0; }

/* ── Back link ── */
.back-link {
  display: inline-flex;
  align-items: center;
  gap: .4rem;
  font-size: 0.8rem;
  font-weight: 500;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--ink-2);
  text-decoration: none;
  margin-bottom: 3rem;
  transition: color .15s;
}
.back-link:hover { color: var(--accent); }
.back-link svg { transition: transform .15s; }
.back-link:hover svg { transform: translateX(-3px); }

/* ── Empty state ── */
.empty {
  text-align: center;
  padding: 4rem 0;
  color: var(--ink-2);
}
.empty h2 { font-family: var(--mono); font-size: 1.3rem; font-weight: 400; margin-bottom: .6rem; }

@media (max-width: 520px) {
  .post-item { grid-template-columns: 1fr; }
  .post-date-col { display: flex; align-items: baseline; gap: .5rem; margin-bottom: .4rem; }
  .post-date-col .day { font-size: 1.1rem; }
}
`

const layoutTmpl = `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.PageTitle}}</title>
  <meta name="description" content="{{.MetaDesc}}">
  <style>` + baseCSS + `</style>
</head>
<body>
  <header class="site-header">
    <div class="header-inner">
      <a class="site-title" href="/index.html">{{.SiteName}}</a>
      <nav>
        <a href="/index.html">Articles</a>
        <a href="/static/about.html">A propos</a>
      </nav>
    </div>
  </header>
  <main>{{.Body}}</main>
  <footer>
    &copy; {{.Year}} {{.SiteName}}
  </footer>
</body>
</html>`

const indexTmpl = `
{{if .Posts}}
<ul class="post-list">
  {{range $i, $p := .Posts}}
  <li class="post-item">
    <div class="post-date-col">
      <span class="day">{{$p.Date.Format "02"}}</span>
      {{$p.Date.Format "Jan"}}<br>{{$p.Date.Format "2006"}}
    </div>
    <div>
      <a class="post-title-link" href="/static/{{$p.Slug}}.html">{{$p.Title}}</a>
      {{if $p.Description}}<p class="post-description">{{$p.Description}}</p>{{end}}
      {{if $p.Tags}}
      <div class="tags">
        {{range $p.Tags}}<span class="tag">{{.}}</span>{{end}}
      </div>
      {{end}}
    </div>
  </li>
  {{end}}
</ul>

{{if gt .TotalPages 1}}
<div class="pagination" style="display: flex; justify-content: space-between; align-items: center; margin-top: 3rem; padding-top: 2rem; border-top: 1px solid var(--border);">
  <div>
    {{if .PrevPage}}
    <a href="{{.PrevPage}}" class="tag" style="text-transform: none;">&lt;&lt; Précédent</a>
    {{end}}
  </div>
  <div style="font-size: 0.8rem; color: var(--ink-2); font-weight: 500;">
    Page {{.CurrentPage}} sur {{.TotalPages}}
  </div>
  <div>
    {{if .NextPage}}
    <a href="{{.NextPage}}" class="tag" style="text-transform: none;">Suivant &gt;&gt;</a>
    {{end}}
  </div>
</div>
{{end}}

{{else}}
<div class="empty">
  <h2>Aucun article pour l'instant</h2>
  <p>Créez un fichier <code>.md</code> dans le dossier <code>posts/</code> pour commencer.</p>
</div>
{{end}}`

const postTmpl = `
<a class="back-link" href="/index.html">
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
  Tous les articles
</a>
<article>
  <header class="post-header">
    <div class="post-meta">
      <time>{{.Post.Date.Format "2 January 2006"}}</time>
      {{if .Post.Tags}}<span class="sep">·</span>
      <div class="tags">{{range .Post.Tags}}<span class="tag">{{.}}</span>{{end}}</div>
      {{end}}
    </div>
    <h1 class="post-title">{{.Post.Title}}</h1>
    {{if .Post.Description}}
    <p class="post-description-lead">{{.Post.Description}}</p>
    {{end}}
  </header>
  <div class="prose">{{.Post.Content}}</div>
</article>`

func renderToFile(filePath string, t *template.Template, data any, title, desc string) error {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("executing content template: %w", err)
	}

	pd := PageData{
		SiteName:  SiteName,
		PageTitle: title,
		MetaDesc:  desc,
		Year:      time.Now().Year(),
		Body:      template.HTML(buf.String()),
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	if err := tmplLayout.Execute(f, pd); err != nil {
		f.Close()
		return fmt.Errorf("executing layout: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("closing file: %w", err)
	}
	return nil
}

func main() {
	fmt.Println("SSG Build Started")

	posts, err := loadPosts("posts")
	if err != nil {
		fmt.Printf("Error loading posts: %v\n", err)
		os.Exit(1)
	}

	// Create static directory
	if err := os.MkdirAll("static", 0755); err != nil {
		fmt.Printf("Error creating static dir: %v\n", err)
		os.Exit(1)
	}

	// 1. Render Index
	const postsPerPage = 10
	totalPosts := len(posts)
	totalPages := (totalPosts + postsPerPage - 1) / postsPerPage

	if totalPosts == 0 {
		totalPages = 1
	}

	for i := 0; i < totalPages; i++ {
		start := i * postsPerPage
		end := start + postsPerPage
		if end > totalPosts {
			end = totalPosts
		}

		pagePosts := posts[start:end]
		pageNum := i + 1

		var fileName string
		if pageNum == 1 {
			fileName = "index.html"
		} else {
			fileName = filepath.Join("static", fmt.Sprintf("page-%d.html", pageNum))
		}

		data := IndexData{
			Posts:       pagePosts,
			CurrentPage: pageNum,
			TotalPages:  totalPages,
		}
		if pageNum > 1 {
			if pageNum == 2 {
				data.PrevPage = "/index.html"
			} else {
				data.PrevPage = fmt.Sprintf("/static/page-%d.html", pageNum-1)
			}
		}
		if pageNum < totalPages {
			data.NextPage = fmt.Sprintf("/static/page-%d.html", pageNum+1)
		}

		err = renderToFile(fileName, tmplIndex, data, "Articles", "Mon blog")
		if err != nil {
			fmt.Printf("Error rendering page %d: %v\n", pageNum, err)
		}
	}

	// 2. Render Posts
	for _, p := range posts {
		outputPath := filepath.Join("static", p.Slug+".html")
		err = renderToFile(outputPath, tmplPost, struct{ Post Post }{p}, p.Title, p.Description)
		if err != nil {
			fmt.Printf("Error rendering post %s: %v\n", p.Slug, err)
			continue
		}
	}

	// 3. Render Static Pages
	pages, _ := loadPages("pages")
	for _, p := range pages {
		outputPath := filepath.Join("static", p.Slug+".html")
		// Reuse tmplPost by casting Page to Post
		err = renderToFile(outputPath, tmplPost, struct{ Post Post }{Post: Post{Title: p.Title, Content: p.Content, Date: time.Now()}}, p.Title, "")
		if err != nil {
			fmt.Printf("Error rendering page %s: %v\n", p.Slug, err)
			continue
		}
	}

	fmt.Printf("Build Complete (%d posts, %d pages)\n", len(posts), len(pages))
}
