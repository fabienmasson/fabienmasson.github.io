# Static Site Generator (SSG) Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a minimalist, zero-dependency Go SSG that generates a static blog from Markdown files.

**Architecture:** A single-file Go program (`ssg/main.go`) that parses Markdown posts with regex, extracts metadata, and renders them into embedded HTML templates. It outputs `index.html` at the root and posts in the `static/` directory.

**Tech Stack:** Go (Standard Library only), HTML, CSS (IBM Plex Mono).

---

### Task 1: Setup and Hello World

**Files:**
- Create: `ssg/main.go`

**Step 1: Write the minimal implementation**

```go
package main

import "fmt"

func main() {
	fmt.Println("SSG Build Started")
}
```

**Step 2: Run to verify**

Run: `go run ssg/main.go`
Expected: Output "SSG Build Started"

**Step 3: Commit**

```bash
git add ssg/main.go
git commit -m "feat: initial ssg setup"
```

---

### Task 2: Models and Templates

**Files:**
- Modify: `ssg/main.go`

**Step 1: Define types and embedded templates**

Add the `Post` struct and the template constants (`layoutTmpl`, `indexTmpl`, `postTmpl`) to `ssg/main.go`. Use the templates from the prototype in `posts/20260311-tech-blog.md` but adapt them for static links (e.g., links should point to `static/slug.html`).

**Step 2: Verify compilation**

Run: `go build -o bin/ssg ./ssg`
Expected: Build success.

**Step 3: Commit**

```bash
git add ssg/main.go
git commit -m "feat: add models and embedded templates"
```

---

### Task 3: Markdown and Front-matter Parsers (TDD)

**Files:**
- Modify: `ssg/main.go`
- Create: `ssg/main_test.go`

**Step 1: Write tests for parsers**

Add tests for `parseMarkdown` and `parseFrontMatter` to `ssg/main_test.go`.

**Step 2: Run tests to verify they fail**

Run: `go test ./ssg/...`
Expected: FAIL (functions not defined)

**Step 3: Implement parsers**

Copy and adapt `parseMarkdown` and `parseFrontMatter` from the prototype into `ssg/main.go`.

**Step 4: Run tests to verify they pass**

Run: `go test ./ssg/...`
Expected: PASS

**Step 5: Commit**

```bash
git add ssg/main.go ssg/main_test.go
git commit -m "feat: implement markdown and front-matter parsing with tests"
```

---

### Task 4: Post Loading Logic

**Files:**
- Modify: `ssg/main.go`

**Step 1: Implement loadPosts**

Implement the `loadPosts` function that reads the `posts/` directory, parses each file, and returns a sorted slice of `Post` objects.

**Step 2: Add simple print in main to verify**

Update `main()` to call `loadPosts("posts")` and print the count.

**Step 3: Verify**

Run: `go run ssg/main.go`
Expected: "Loaded X posts"

**Step 4: Commit**

```bash
git add ssg/main.go
git commit -m "feat: implement post loading and sorting"
```

---

### Task 5: Static Generation Orchestration

**Files:**
- Modify: `ssg/main.go`

**Step 1: Implement rendering and file writing**

Implement `renderToFile` and the build loop in `main()`.
- Ensure `static/` exists.
- Write `index.html` to root.
- Write each post to `static/[slug].html`.

**Step 2: Verify generation**

Run: `go run ssg/main.go`
Expected: `index.html` created at root, HTML files created in `static/`.

**Step 3: Commit**

```bash
git add ssg/main.go
git commit -m "feat: implement static file generation"
```

---

### Task 6: Final Verification and Cleanup

**Files:**
- Modify: `ssg/main.go`
- Modify: `cv.css` (if needed for typography consistency)

**Step 1: Verify all links**

Check that `index.html` links correctly to `static/post-slug.html` and that post pages link back to `/index.html`.

**Step 2: Run full build and lint**

Run: `go fmt ./ssg/... && go vet ./ssg/... && go run ssg/main.go`

**Step 3: Commit**

```bash
git add .
git commit -m "chore: final ssg cleanup and verification"
```
