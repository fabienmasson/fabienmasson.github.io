# Agent Guide: Static Blog Project

This repository contains a personal project to build a static site generator (SSG) in Go, primarily for a blog and CV.

## Project Objectives

- **SSG Implementation**: Create a static site generator in Go.
- **Output Structure**:
    - `index.html` generated at the project root.
    - Other pages, JS, and CSS generated in the `static/` directory.
- **Generator Location**: All SSG source code must reside in the `ssg/` directory.
- **Technical Stack**:
    - **SSG**: Written in Go (Standard Library preferred).
    - **Frontend**: HTML, CSS, and Vanilla JavaScript.
    - **Dependencies**: NO external JavaScript libraries allowed. Zero-dependency approach.
- **Theme**: Minimalist and modern.
- **Typography**: Must use **IBM Plex Mono**.

## Repository Structure

- `ssg/`: Go source code for the site generator.
- `posts/`: Markdown files (`.md`) serving as the blog content.
- `static/`: Generated secondary pages and assets (CSS, JS).
- `cv/`: Original CV files (XML/XSLT/CSS).
- `index.html`: The main landing page (generated).

## Development Commands

### Go (SSG)
- **Run Generator**: `go run ssg/*.go`
- **Build**: `go build -o bin/ssg ./ssg`
- **Test All**: `go test ./ssg/...`
- **Test Single Package**: `go test ./ssg/parser`
- **Lint**: `go fmt ./ssg/... && go vet ./ssg/...`

### PDF Generation (CV)
- **Generate PDF**: `weasyprint index.html cv/cv2026.pdf`
- **Requirements**: WeasyPrint must be installed on the system.

## Code Style Guidelines

### Go Development
- **Standard Library First**: Avoid external dependencies. Use `html/template`, `os`, `path/filepath`, `regexp`, etc.
- **Formatting**: Always run `go fmt` before committing.
- **Naming**: 
    - Use `camelCase` for internal variables and functions.
    - Use `PascalCase` for exported types and functions if needed (though the SSG might be internal).
- **Error Handling**: 
    - Explicitly handle all errors: `if err != nil { ... }`.
    - Wrap errors with context: `fmt.Errorf("parsing markdown: %w", err)`.
- **Project Layout**: Keep a simple structure. If the project grows, split into modules:
  - `ssg/main.go`: Entry point and CLI logic.
  - `ssg/parser.go`: Markdown to HTML logic.
  - `ssg/template.go`: Template engine wrapper.
- **Comments**: Focus on "why" for complex logic (like custom markdown parsing).

### Frontend Development
- **CSS**: 
    - Use CSS Variables for colors and themes.
    - Minimalist design: focus on whitespace and typography.
    - Responsive design: ensure mobile-friendliness.
- **JavaScript**: 
    - Vanilla JS only.
    - No frameworks (React, Vue, etc.) or libraries (jQuery).
    - Use modern ES6+ syntax.
- **Typography Integration**:
    ```css
    @import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@300;400;600&display=swap');
    :root {
      --font-mono: 'IBM Plex Mono', monospace;
    }
    body { font-family: var(--font-mono); }
    ```
- **Accessibility (A11y)**:
  - Use semantic HTML tags (`<article>`, `<nav>`, `<header>`, `<footer>`).
  - Ensure proper contrast for text.
  - Add `alt` attributes to images.

## Markdown Specifications
The SSG should support:
- **Front-matter**: YAML-style metadata at the top of `.md` files.
  ```markdown
  ---
  title: Post Title
  date: 2026-03-12
  description: Short summary
  tags: tag1, tag2
  slug: my-post-slug
  ---
  ```
- **Standard Features**: Headings, Bold/Italic, Code blocks (with syntax highlighting hints), Lists, Blockquotes, Links, and Images.

## Implementation Roadmap (Internal Reference)
1. **Parser Migration**: Move the prototype parser from `posts/` content into `ssg/`.
2. **Template Integration**: Implement a robust template engine using `html/template` for base layout and post pages.
3. **Asset Handling**: Ensure CSS and JS from `static/` are correctly linked in generated HTML.
4. **Build System**: Set up a watcher or a simple command to rebuild the site.
5. **CV Integration**: Optionally integrate the CV layout into the SSG output if applicable.

## Testing Strategy
- **Markdown Parser**: Write unit tests for each markdown element (bold, italic, code blocks, etc.).
- **Template Rendering**: Verify that variables are correctly interpolated in the layout.
- **File System**: Ensure files are generated in the correct locations (`/` and `static/`).

## Git Protocol
- Create descriptive commit messages (e.g., `feat: implement markdown parser`, `fix: correct page break in pdf`).
- Always run `go fmt` and `go vet` before committing.
- Do not commit generated binaries (`bin/`).
- Use branch-based development for major features.

## Philosophies
- **Simplicity Over Complexity**: If a feature can be implemented with a simple regex in Go, prefer that over importing a heavy markdown library.
- **Performance**: Static site generation should be fast. Avoid unnecessary I/O.
- **Portability**: The project should be easy to run on any machine with Go installed.

---
*Note: This file is intended for agentic workflows. Follow these instructions strictly to maintain project consistency.*
