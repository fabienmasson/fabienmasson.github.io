# Tag Pages with Pagination Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement static tag pages with 10-post pagination.

**Architecture:** Refactor the pagination logic into a reusable function `renderPaginatedIndex` and update `main()` to generate pages for each unique tag.

**Tech Stack:** Go (Standard Library), HTML/CSS.

---

### Task 1: Update IndexData and Templates

**Files:**
- Modify: `ssg/main.go:41-47` (Update IndexData)
- Modify: `ssg/main.go:740-806` (Update Templates)

**Step 1: Add PageHeading to IndexData**

Modify `IndexData` struct:
```go
type IndexData struct {
	Posts       []Post
	CurrentPage int
	TotalPages  int
	PrevPage    string
	NextPage    string
	PageHeading string // Added
}
```

**Step 2: Update indexTmpl to include Heading and linkable tags**

- Add `{{if .PageHeading}}<h1 class="page-heading">{{.PageHeading}}</h1>{{end}}` at the top.
- Change `<span class="tag">{{.}}</span>` to `<a href="/static/tag-{{slugify .}}.html" class="tag">{{.}}</a>`.

**Step 3: Update postTmpl to include linkable tags**

- Change `<span class="tag">{{.}}</span>` to `<a href="/static/tag-{{slugify .}}.html" class="tag">{{.}}</a>`.

*Note: I need to add slugify to the template FuncMap.*

**Step 4: Commit**

```bash
git add ssg/main.go
git commit -m "feat: make tags clickable and add PageHeading to IndexData"
```

---

### Task 2: Implement renderPaginatedIndex Helper

**Files:**
- Modify: `ssg/main.go` (Add function before main)

**Step 1: Create renderPaginatedIndex function**

```go
func renderPaginatedIndex(posts []Post, filePrefix, title, desc, heading string) error {
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
		if filePrefix == "" { // Main index
			if pageNum == 1 {
				fileName = "index.html"
			} else {
				fileName = filepath.Join("static", fmt.Sprintf("page-%d.html", pageNum))
			}
		} else { // Tag pages
			if pageNum == 1 {
				fileName = filepath.Join("static", fmt.Sprintf("%s.html", filePrefix))
			} else {
				fileName = filepath.Join("static", fmt.Sprintf("%s-page-%d.html", filePrefix, pageNum))
			}
		}

		data := IndexData{
			Posts:       pagePosts,
			CurrentPage: pageNum,
			TotalPages:  totalPages,
			PageHeading: heading,
		}

		// Prev Link
		if pageNum > 1 {
			if filePrefix == "" {
				if pageNum == 2 {
					data.PrevPage = "/index.html"
				} else {
					data.PrevPage = fmt.Sprintf("/static/page-%d.html", pageNum-1)
				}
			} else {
				if pageNum == 2 {
					data.PrevPage = fmt.Sprintf("/static/%s.html", filePrefix)
				} else {
					data.PrevPage = fmt.Sprintf("/static/%s-page-%d.html", filePrefix, pageNum-1)
				}
			}
		}

		// Next Link
		if pageNum < totalPages {
			if filePrefix == "" {
				data.NextPage = fmt.Sprintf("/static/page-%d.html", pageNum+1)
			} else {
				data.NextPage = fmt.Sprintf("/static/%s-page-%d.html", filePrefix, pageNum+1)
			}
		}

		err := renderToFile(fileName, tmplIndex, data, title, desc)
		if err != nil {
			return err
		}
	}
	return nil
}
```

**Step 2: Commit**

```bash
git add ssg/main.go
git commit -m "feat: add renderPaginatedIndex helper function"
```

---

### Task 3: Refactor main() and Generate Tag Pages

**Files:**
- Modify: `ssg/main.go` (Update main and template initialization)

**Step 1: Add slugify to Template FuncMap**

In `var (...)` section, update template parsing to include `Funcs`:
```go
	funcMap = template.FuncMap{
		"slugify": slugify,
	}
	tmplLayout = template.Must(template.New("layout").Funcs(funcMap).Parse(layoutTmpl))
	tmplIndex  = template.Must(template.New("index").Funcs(funcMap).Parse(indexTmpl))
	tmplPost   = template.Must(template.New("post").Funcs(funcMap).Parse(postTmpl))
```

**Step 2: Update main() to use helper and collect tags**

```go
	// 1. Render Index
	err = renderPaginatedIndex(posts, "", "Accueil", "Mon portfolio et blog personnel", "")
	if err != nil {
		fmt.Printf("Error rendering index: %v\n", err)
		os.Exit(1)
	}

	// 1.5 Render Tag Pages
	tagMap := make(map[string][]Post)
	for _, p := range posts {
		for _, t := range p.Tags {
			tagMap[t] = append(tagMap[t], p)
		}
	}
	for tag, tagPosts := range tagMap {
		slug := "tag-" + slugify(tag)
		heading := fmt.Sprintf("Articles avec le tag: <em>%s</em>", tag)
		err = renderPaginatedIndex(tagPosts, slug, "Tag: "+tag, "Articles taggués avec "+tag, heading)
		if err != nil {
			fmt.Printf("Error rendering tag page %s: %v\n", tag, err)
		}
	}
```

**Step 3: Commit**

```bash
git add ssg/main.go
git commit -m "feat: generate paginated tag pages"
```

---

### Task 4: Final Verification

**Step 1: Run build**

Run: `go run ssg/*.go`

**Step 2: Verify links**

- Check `index.html`: Are tags links? Do they point to `static/tag-X.html`?
- Check a tag page: Is there a heading? Does pagination work (if you add temporary posts)?

**Step 3: Commit final build**

```bash
git add index.html static/*.html
git commit -m "chore: final build with tag pages"
```
