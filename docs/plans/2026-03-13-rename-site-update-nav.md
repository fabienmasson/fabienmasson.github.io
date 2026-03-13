# Rename Site & Update Navigation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rename the site to "F4N" and update the navigation menu to include an "A propos" link.

**Architecture:** Update the global `SiteName` constant and the HTML layout template within the SSG's main entry point.

**Tech Stack:** Go (Standard Library)

---

### Task 1: Update SiteName Constant

**Files:**
- Modify: `ssg/main.go:34-34`

**Step 1: Read the file to verify content**

Run: `grep -n "const SiteName" /home/fabien/Projets/Divers/cv/ssg/main.go`

**Step 2: Update SiteName constant**

Modify line 34 in `/home/fabien/Projets/Divers/cv/ssg/main.go` to:
```go
const SiteName = "F4N"
```

**Step 3: Verify the change**

Run: `grep 'const SiteName = "F4N"' /home/fabien/Projets/Divers/cv/ssg/main.go`
Expected: Match found.

---

### Task 2: Update Layout Template Navigation

**Files:**
- Modify: `ssg/main.go:678-680`

**Step 1: Read the template section**

Run: `cat /home/fabien/Projets/Divers/cv/ssg/main.go | sed -n '670,690p'`

**Step 2: Update layout template navigation**

Modify lines 678-680 in `/home/fabien/Projets/Divers/cv/ssg/main.go` to:
```html
      <nav>
        <a href="/index.html">Articles</a>
        <a href="/static/about.html">A propos</a>
      </nav>
```

**Step 3: Verify the change**

Run: `grep -A 3 '<nav>' /home/fabien/Projets/Divers/cv/ssg/main.go`
Expected: Output shows the updated navigation links.

---

### Task 3: Build and Verify

**Step 1: Run Go vet and fmt**

Run: `go fmt ./ssg/... && go vet ./ssg/...`
Expected: No errors.

**Step 2: Run the generator**

Run: `go run ssg/*.go`
Expected: Site generates successfully (no errors in terminal).

**Step 3: Verify generated index.html**

Run: `grep 'F4N' index.html && grep 'A propos' index.html`
Expected: Matches found in the generated output.

---

### Task 4: Commit Changes

**Step 1: Add and commit**

Run: `git add ssg/main.go && git commit -m "feat: rename site to F4N and add 'A propos' to navigation"`
Expected: Commit successful.
