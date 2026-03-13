# Design: Tag Pages with Pagination

## 1. Overview
Implement dedicated static pages for each tag found in blog posts, including full pagination (10 posts per page).

## 2. Tag Collection
- Iterate through all loaded posts.
- Maintain a `map[string][]Post` where keys are tags.
- Ensure tags are case-insensitive or standardized (slugified).

## 3. URL Structure
- Main tag page: `static/tag-[slug].html`.
- Paginated tag pages: `static/tag-[slug]-page-[n].html`.
- Links to tags in templates will point to `/static/tag-[slug].html`.

## 4. Architecture Refactoring
- Introduce `renderPaginatedIndex` function in `ssg/main.go`.
- This function will:
  - Take a list of posts, a file prefix, and metadata.
  - Handle chunking (10 posts/page).
  - Generate proper `PrevPage` and `NextPage` URLs based on the prefix.
  - Call `renderToFile` for each page chunk.

## 5. UI Updates
- Update `indexTmpl` and `postTmpl`:
  - Change `<span class="tag">` to `<a href="/static/tag-{{.}}.html" class="tag">`.
- Update `IndexData` struct:
  - Add `PageHeading` string.
- Tag pages will display `<h1>Articles avec le tag: [Tag]</h1>` above the list.

## 6. Success Criteria
- Clicking any tag on a post or the index takes the user to a filtered list.
- Tag pages with > 10 posts have working pagination links.
- "Articles avec le tag: [Tag]" is clearly visible on filtered pages.
- Navigation between tag pages and back to the main index works correctly.
