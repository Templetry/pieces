// A piece that exists on disk but is not in the catalog is invisible; one in
// the catalog but not on disk is a broken fetch for whoever tries it. Both
// are silent, so both are checked here.
import fs from "node:fs";

const errors = [];
const err = (m) => errors.push(m);

const dirs = fs.readdirSync(".", { withFileTypes: true })
  .filter((e) => e.isDirectory() && !e.name.startsWith(".") && fs.existsSync(`${e.name}/piece.yml`))
  .map((e) => e.name);

if (!dirs.length) err("no piece directories found");

const field = (txt, name) => {
  for (const raw of txt.split("\n")) {
    const line = raw.trim();
    if (line.startsWith(name + ":")) return line.slice(name.length + 1).trim();
  }
  return null;
};

const reg = await (await fetch("https://raw.githubusercontent.com/Templetry/catalog/main/registry.json")).json();
const declared = new Map((reg.pieces ?? []).map((p) => [p.path, p]));

// Template names the catalog knows, so applies_to cannot point at nothing.
const templates = new Set();
for (const p of reg.parents ?? []) for (const f of p.forms ?? []) templates.add(f.name);

const readme = fs.readFileSync("README.md", "utf8");

for (const dir of dirs) {
  const txt = fs.readFileSync(`${dir}/piece.yml`, "utf8");
  const name = field(txt, "name");
  if (!name) err(`${dir}/piece.yml: no name`);
  if (!field(txt, "description")) err(`${dir}/piece.yml: no description`);

  const entry = declared.get(dir);
  if (!entry) {
    err(`${dir} is not listed in the catalog registry — nobody can adopt it`);
  } else if (entry.name !== name) {
    err(`${dir}: piece.yml says "${name}", registry says "${entry.name}"`);
  }

  if (!readme.includes(`\`${dir}/\``)) err(`README does not list ${dir}/`);

  // applies_to entries are template names, and a typo there silently makes
  // the piece unavailable everywhere instead of failing loudly.
  const lines = txt.split("\n");
  const start = lines.findIndex((l) => l.trim().startsWith("applies_to:"));
  if (start >= 0) {
    for (let i = start + 1; i < lines.length; i++) {
      const l = lines[i];
      if (!l.startsWith("  ") && l.trim()) break;
      const m = l.trim();
      if (!m.startsWith("- ")) continue;
      const target = m.slice(2).trim();
      if (!templates.has(target)) err(`${dir}: applies_to "${target}" is not a template in the catalog`);
    }
  }
}

for (const [path] of declared) {
  if (!dirs.includes(path)) err(`registry lists "${path}" but this repo has no such piece`);
}

for (const e of errors) console.log(`::error::${e}`);
console.log(`\n${dirs.length} pieces checked · ${errors.length} errors`);
process.exitCode = errors.length ? 1 : 0;
