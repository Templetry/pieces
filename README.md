# Templetry common pieces

Pieces that are not specific to one form ([ADR-0016](https://github.com/Templetry/wiki/blob/main/adr/0016-common-pieces.md)): they live here once and any compatible project can adopt them.

| Piece | Applies to | What it adds |
|---|---|---|
| [`renovate/`](renovate/) | any template | Renovate config: grouped scheduled PRs, immediate security fixes, monthly lockfile refresh |
| [`audit-trail-go-sqlite/`](audit-trail-go-sqlite/) | `go-rest-sqlite` | Append-only audit table, repository and read-only endpoint |
| [`agent-pointers/`](agent-pointers/) | any template | One instruction file per AI tool, each pointing at the project's `AGENTS.md`, plus an MCP config so an agent can adopt pieces and run updates itself |

## How it works

A piece directory carries a `piece.yml` with an optional `applies_to` list of **template names** (the `name` field of a form's `template.yml`). Empty means universal. Several directories may declare the same piece `name` with disjoint `applies_to` — one implementation per ecosystem — and a project simply asks for the name:

```sh
templetry pieces ./my-project      # form pieces + the common ones that fit
templetry add renovate ./my-project
```

Form-local pieces win over common ones on a name clash: a form shipping its own implementation is making a deliberate statement.

Adopting a common piece records **this repository** as its source, so `templetry update` follows it here — fix a piece once and every project that adopted it sees the fix.

## Verification

Every piece here is applied in CI the way a user applies it: a target form is
rendered from the catalog, the piece is fetched **from the commit under test**
through the registry, and the result is built and tested. A common piece is the
only Templetry code that lands inside somebody else’s existing project, so a
local shortcut that cannot fail the way an install fails would not be worth
much.

## Contributing

Keep the rule that makes pieces safe: a piece may only add files that do not exist in the project, and may only touch shared files through declared patches. Wiring that lives in code needs the form to expose a socket.

The full guide: [authoring pieces](https://github.com/Templetry/wiki/blob/main/guide/authoring-pieces.md).
