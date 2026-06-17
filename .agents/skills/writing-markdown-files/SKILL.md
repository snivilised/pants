---
name: writing-markdown-files
description: Generate clean Markdown files that avoid common lint errors, especially markdownlint violations.
---

# Markdown Cleaner Skill

## Purpose

Produce Markdown that is structurally clean, consistent, and likely to pass markdown linters such as markdownlint.

## Core behavior

When generating or editing Markdown:

- Preserve valid Markdown semantics.
- Prefer simple, standard syntax.
- Avoid formatting that commonly triggers lint errors.
- Keep structure predictable and consistent.
- Rewrite content when needed to satisfy lint rules instead of forcing exceptions.

## Common lint-error prevention

### Headings

- Use exactly one top-level heading when the document expects one.
- Make the first line a top-level heading unless the document type explicitly says otherwise.
- Put a blank line before and after headings unless the style guide says not to.
- Write headings with a space after the hash marks, for example `# Title`, not `#Title`.
- Keep heading levels in order and avoid skipping levels unless the document structure requires it.

### Blank lines

- Use a single blank line between block elements.
- Avoid multiple consecutive blank lines.
- Put a blank line before lists, code blocks, blockquotes, and tables when needed for readability and parser safety.
- Remove trailing spaces at the end of lines.

### Line length

- Wrap long prose lines to a reasonable width.
- Prefer short sentences and paragraphs.
- Leave URLs intact only when line wrapping would break readability or the target style guide.
- Do not add arbitrary line breaks inside words, links, or inline code.

### Lists

- Use consistent bullet markers within the same list.
- Put a space after the list marker, for example `- Item`, not `-Item`.
- Keep list indentation consistent.
- Use proper nesting for sublists with matching indentation.
- Do not mix numbered and bulleted formats inside one list unless structure requires it.

### Code and inline code

- Wrap commands, filenames, flags, identifiers, and literal values in backticks.
- Do not add spaces directly inside inline code delimiters.
- Use fenced code blocks for multi-line code.
- Add a language tag to fenced code blocks when the language is known.
- Do not indent fenced code blocks unless the surrounding context requires it.

### Punctuation and spacing

- Use one space after sentences.
- Avoid double spaces.
- Keep punctuation outside code spans unless the literal text requires otherwise.
- Do not add extra spaces before commas, periods, colons, or closing brackets.

### Links and images

- Use valid Markdown link syntax with descriptive link text.
- Avoid bare URLs in prose when a link can be used.
- Make image alt text descriptive and concise.
- Ensure reference-style links, if used, are defined and consistent.

### Tables

- Use tables only when they improve clarity.
- Keep table rows and columns aligned.
- Ensure each table row has the same number of columns.
- Avoid line breaks inside table cells unless the renderer supports them reliably.

### Common markdownlint traps

- Do not start headings with no space after `#`.
- Do not place multiple top-level headings unless the document structure requires them.
- Do not use multiple consecutive blank lines.
- Do not leave trailing whitespace.
- Do not create very long unwrapped paragraphs when the target style enforces line length.
- Do not use inconsistent list indentation.
- Do not forget a blank line before and after lists or code blocks when the renderer or style guide expects one.

## Editing workflow

When given Markdown to produce or fix:

1. Read the full document structure first.
2. Identify likely lint problems.
3. Rewrite for correctness and consistency.
4. Preserve meaning while changing formatting.
5. Re-check headings, spacing, lists, links, code blocks, and tables.
6. Prefer lint-safe output over stylistic cleverness.

## Output requirements

- Return valid Markdown only.
- Do not add explanatory prose outside the requested content.
- Do not mention lint rules unless the user asks for them.
- If a conflict exists between style preference and lint safety, choose lint safety.

## Quality checklist

Before finalizing, verify:

- The first heading is valid.
- Headings have proper spacing.
- There is only one blank line between blocks.
- No trailing spaces remain.
- Long paragraphs are wrapped.
- Lists are consistently indented.
- Inline code is used for literal tokens.
- Code fences are properly opened and closed.
- Tables are structurally consistent.
- Links render correctly.
---
