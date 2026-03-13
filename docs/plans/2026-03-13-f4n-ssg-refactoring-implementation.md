# F4N SSG Refactoring Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rename the blog to "F4N", add an "A propos" page, and implement 10-post pagination.

**Architecture:** Update the Go SSG in `ssg/main.go` to handle a second content type (static pages) and split the post list into paginated index files.

**Tech Stack:** Go (Standard Library), HTML/CSS (Vanilla).

---

### Task 1: Rename Site & Update Navigation

**Files:**
- Modify: `ssg/main.go:34-34` (Constant SiteName)
- Modify: `ssg/main.go:678-680` (Layout navigation)

**Step 1: Update SiteName constant**

Modify line 34:
```go
const SiteName = "F4N"
```

**Step 2: Update layout template navigation**

Modify lines 678-680:
```go
      <nav>
        <a href="/index.html">Articles</a>
        <a href="/static/about.html">A propos</a>
      </nav>
```

**Step 3: Commit**

```bash
git add ssg/main.go
git commit -m "feat: rename site to F4N and add 'A propos' to navigation"
```

---

### Task 2: Define Page Struct & Page Loader

**Files:**
- Modify: `ssg/main.go:24-25` (Add Page struct)
- Modify: `ssg/main.go:370-371` (Add loadPages function)

**Step 1: Define Page struct**

```go
type Page struct {
	Title      string
	Slug       string
	Content    template.HTML
	RawContent string
}
```

**Step 2: Add loadPages function**

```go
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
```

**Step 3: Commit**

```bash
git add ssg/main.go
git commit -m "feat: add Page struct and loadPages function"
```

---

### Task 3: Render Static Pages

**Files:**
- Modify: `ssg/main.go:769-802` (Update main function)

**Step 1: Load and render pages in main()**

```go
	// 2. Load and render Static Pages
	pages, err := loadPages("pages")
	if err != nil {
		fmt.Printf("Error loading pages: %v\n", err)
	} else {
		for _, p := range pages {
			outputPath := filepath.Join("static", p.Slug+".html")
			// Reuse post layout for pages
			err = renderToFile(outputPath, tmplPost, struct{ Post Post }{Post: Post{Title: p.Title, Content: p.Content}}, p.Title, "")
			if err != nil {
				fmt.Printf("Error rendering page %s: %v\n", p.Slug, err)
			}
		}
	}
```

**Step 2: Run and verify build**

Run: `go run ssg/*.go`
Check: `static/about.html` exists and contains "TBD".

**Step 3: Commit**

```bash
git add ssg/main.go
git commit -m "feat: render static pages from pages/ directory"
```

---

### Task 4: Pagination Logic - Data Structure

**Files:**
- Modify: `ssg/main.go:26-32` (Update PageData)

**Step 1: Add Pagination info to template data**

```go
type IndexData struct {
	Posts       []Post
	CurrentPage int
	TotalPages  int
	PrevPage    string
	NextPage    string
}
```

**Step 2: Commit**

```bash
git commit -am "refactor: add IndexData for pagination support"
```

---

### Task 5: Pagination Implementation - Chunking

**Files:**
- Modify: `ssg/main.go:769-802` (Update main function)

**Step 1: Implement chunking logic in main()**

```go
	const postsPerPage = 10
	totalPosts := len(posts)
	totalPages := (totalPosts + postsPerPage - 1) / postsPerPage

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
```

**Step 2: Commit**

```bash
git commit -am "feat: implement pagination logic in SSG main"
```

---

### Task 6: Pagination Implementation - Template

**Files:**
- Modify: `ssg/main.go:690-716` (Update indexTmpl)

**Step 1: Add pagination links to index template**

```go
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
```

**Step 2: Verify build and pagination layout**

Run: `go run ssg/*.go`

**Step 3: Commit**

```bash
git commit -am "feat: add pagination controls to index template"
```

---

### Task 7: Cleanup & Test

**Files:**
- Modify: `ssg/main_test.go` (Optional: verify basic logic if existing)

**Step 1: Run full build and check all links**

Check:
1. `index.html` (Title, Articles menu)
2. `static/about.html` (Content, menu link)
3. Pagination works if > 10 posts exist.

**Step 2: Commit final build**

```bash
git add index.html static/*.html
git commit -m "chore: final build with F4N branding and pagination"
```
