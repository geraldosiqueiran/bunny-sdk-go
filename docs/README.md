# Documentation Index

Welcome to the Bunny.net Go SDK documentation. Choose your path based on your needs.

## Quick Navigation

### I'm a new user - where do I start?
1. Read **[../README.md](../README.md)** - Quick start guide with examples
2. Find your use case (Stream, Storage, Shield, Scripting, or Containers)
3. Copy the example and modify for your needs

### I'm contributing code - what do I need to know?
1. Read **[codebase-summary.md](./codebase-summary.md)** - Understand the structure
2. Follow **[code-standards.md](./code-standards.md)** - Naming conventions and patterns
3. Use the checklists in code-standards.md before submitting PRs

### I want to understand how it all works
1. Read **[system-architecture.md](./system-architecture.md)** - Deep dive into design
2. Learn the request flow, error handling, and service patterns
3. Understand design decisions and their rationale

### I'm a project manager or stakeholder
1. Review **[project-overview-pdr.md](./project-overview-pdr.md)** - Project scope and requirements
2. Check **[project-roadmap.md](./project-roadmap.md)** - Timeline and milestones
3. Track progress against success criteria

---

## Documentation Files

### [codebase-summary.md](./codebase-summary.md) (295 lines)
**WHAT:** How the codebase is organized

- Directory structure with descriptions
- File naming conventions
- Package responsibilities
- Architecture patterns
- Authentication model
- Dependencies

**When to read:** Before writing code or understanding structure

---

### [code-standards.md](./code-standards.md) (512 lines)
**HOW:** Standards for writing code in this project

- Naming conventions (types, functions, variables, files)
- Package structure patterns
- Error handling patterns
- Type design patterns
- Testing patterns (table-driven, mocking)
- JSON tag conventions
- Performance considerations
- Code review checklist

**When to read:** Before submitting code changes

---

### [project-overview-pdr.md](./project-overview-pdr.md) (325 lines)
**SCOPE:** Project requirements and goals

- Project purpose and goals
- Feature set overview
- Non-functional requirements
- API coverage status
- Success criteria
- Development phases
- Security & compliance
- Known limitations

**When to read:** Understand project scope and requirements

---

### [system-architecture.md](./system-architecture.md) (668 lines)
**WHY:** Deep understanding of design and implementation

- High-level architecture overview
- Client hierarchy
- Request flow (9-step diagram)
- Authentication flow
- Error handling flow
- Service pattern explanation
- Type system design
- Data flow examples
- Storage regions
- Concurrency & context handling
- Design decisions with rationale
- Extension points
- Performance characteristics

**When to read:** Want to understand how everything works together

---

### [project-roadmap.md](./project-roadmap.md) (453 lines)
**FUTURE:** Timeline, milestones, and planned work

- Release timeline (v0.1 Alpha, v1.0 Stable)
- Feature phases (all 7 phases with completion status)
- Planned features (Pull Zones, DNS, Billing, Abuse APIs)
- Performance improvements (retry logic, rate limiting, connection pooling)
- Milestone targets
- Known issues with workarounds
- Success metrics
- Maintenance plan

**When to read:** Track progress or plan contributions

---

## Document Relationships

```
README.md (ENTRY POINT)
  │
  ├─ codebase-summary.md ──────┐
  │    (How it's organized)     │
  │                             │
  ├─ code-standards.md ─────────┼─ (Read when contributing)
  │    (How to write code)      │
  │                             │
  ├─ system-architecture.md ────┼─ (Read for deep understanding)
  │    (How it works)           │
  │                             │
  ├─ project-overview-pdr.md ───┼─ (Read to understand scope)
  │    (What we're building)    │
  │                             │
  └─ project-roadmap.md ────────┘ (Read to track progress)
       (Where we're going)
```

---

## By Audience

### New Users
1. README.md - Get started
2. Find your API (Stream, Storage, Shield, Scripting, Containers)
3. Copy examples and adapt

### New Developers
1. codebase-summary.md - Understand structure
2. code-standards.md - Learn conventions
3. system-architecture.md - Deep dive
4. Contribute following checklist in code-standards.md

### Project Managers
1. project-overview-pdr.md - Understand scope
2. project-roadmap.md - Track milestones
3. Monitor success criteria

### Code Reviewers
1. code-standards.md - Verify against checklist
2. system-architecture.md - Validate design decisions
3. codebase-summary.md - Check organization

### Maintainers
1. All documents - Weekly review
2. Update roadmap on releases
3. Track maintenance metrics

---

## Key Information Quick Reference

### Authentication

| Service | Key Type | Where to Find |
|---------|----------|---------------|
| Stream API | Library API Key | Library settings → API |
| Storage Zones | Global API Key | Account Settings → API |
| File Operations | Zone Password | Zone → FTP & API Access |
| Shield/WAF | Global API Key | Account Settings → API |
| Scripting | Global API Key | Account Settings → API |
| Containers | Global API Key | Account Settings → API |

### API Packages

| Package | Methods | File |
|---------|---------|------|
| Stream | 23+ | `stream/` |
| Storage (Zones) | 8 | `storage/` |
| Storage (Files) | 5 | `storage/` |
| Shield/WAF | 50+ | `shield/` |
| Scripting | 23 | `scripting/` |
| Containers | 40+ | `containers/` |

### Design Patterns

- **Functional Options** - Client configuration via options
- **Service Interface** - Interfaces for mockability
- **Adapter Pattern** - Bridge public API to internal HTTP
- **Generic Pagination** - Type-safe PaginatedResponse[T]
- **Error Hierarchy** - APIError with typed subtypes
- **Resource Scoping** - Services parameterized with IDs

---

## File Size & Quality

All documentation files are:
- ✅ Under 800 lines of code (LOC limit)
- ✅ Self-contained and readable independently
- ✅ Cross-referenced with related documents
- ✅ Verified against actual codebase
- ✅ Professional formatting with clear structure

---

## Maintenance & Updates

### When to Update Documentation

- **New API feature** → Update codebase-summary.md, project-roadmap.md
- **Code standard change** → Update code-standards.md
- **New design pattern** → Update system-architecture.md
- **Release milestone** → Update project-roadmap.md
- **Requirement change** → Update project-overview-pdr.md

### Review Schedule
- Monthly: Check for outdated information
- Per release: Update version and roadmap

---

## Questions & Feedback

If you can't find what you're looking for:

1. Check the **Quick Navigation** section above
2. Search within each document (Ctrl+F / Cmd+F)
3. Review **system-architecture.md** for design decisions
4. Ask in GitHub issues with reference to relevant section

---

**Last Updated:** Feb 5, 2026
**Status:** Complete documentation suite for v0.1 Alpha
