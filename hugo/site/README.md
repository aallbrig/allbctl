# allbctl Documentation Site

This is the Hugo-based documentation site for allbctl.

## Requirements

- Hugo Extended v0.120.0 or later
- Git (for theme submodule)

**Note**: Currently using hugo-book theme v8 for compatibility with Hugo 0.123.7.

## Local Development

```bash
# Navigate to site directory
cd hugo/site

# Start Hugo server
hugo server --buildDrafts

# View at http://localhost:1313
```

## Building

```bash
# Build static site
hugo

# Output will be in public/ directory
```

## Theme

Uses [Hugo Book](https://github.com/alex-shpak/hugo-book) theme (v8).

**Version Note**: Using v8 for compatibility with Hugo 0.123.7. The theme is checked out as a git submodule at the v8 tag instead of the latest version which requires Hugo 0.146.0+.

## Structure

```
content/
├── _index.md                  # Homepage
└── docs/
    ├── getting-started/
    │   ├── _index.md
    │   └── installation.md
    ├── commands/
    │   └── _index.md
    ├── status/
    │   ├── _index.md
    │   ├── runtimes.md
    │   ├── projects.md
    │   ├── packages.md
    │   └── databases.md
    └── bootstrap/
        └── _index.md
```

## Adding Content

```bash
# Create new page
hugo new docs/section/page.md

# Edit frontmatter and content
```

## Deployment

The site can be deployed to:
- GitHub Pages
- Netlify
- Vercel
- Any static hosting service

## Configuration

See `hugo.toml` for site configuration including:
- Base URL
- Theme settings
- Menu items
- Markup options
