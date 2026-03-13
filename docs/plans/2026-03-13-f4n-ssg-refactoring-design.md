# Design: F4N SSG Refactoring

## 1. Overview
Rename the blog to "F4N", add an "A propos" page from `pages/about.md`, and implement 10-post pagination.

## 2. Renaming & Navigation
- Change `SiteName` constant from `"Mon Blog & CV"` to `"F4N"`.
- Update `layoutTmpl` navigation menu:
  - Add link to `/static/about.html` labeled "A propos".

## 3. Pages & About Rendering
- Define `Page` struct for generic static pages (Title, Slug, Date, Content).
- Add `loadPages(dir string)` function to process Markdown files in `pages/`.
- Render `pages/about.md` to `static/about.html`.
- Use `postTmpl` for rendering static pages to maintain consistency.

## 4. Pagination
- Constant `PostsPerPage = 10`.
- The main `index.html` will show the most recent 10 posts.
- Subsequent pages will be rendered as `static/page-2.html`, `static/page-3.html`, etc.
- Update `indexTmpl` with navigation:
  - "Précédent" link (where applicable).
  - Current page indicator ("Page X sur Y").
  - "Suivant" link (where applicable).

## 5. Success Criteria
- Site title is "F4N" everywhere.
- "A propos" menu item exists and displays content from `pages/about.md`.
- Blog posts are paginated 10 per page.
- Pagination links work correctly between root and static folder.
