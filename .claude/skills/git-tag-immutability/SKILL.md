---
name: git-tag-immutability
description: Use when creating version tags or making release decisions, before pushing to remote — prevents tag reuse, enforces semver application, and closes version rationalization loopholes
---

# Git Tag Immutability & Semantic Versioning

## Core Principle

**Once a tag is pushed to remote, it is immutable. Never delete, recreate, or reuse tags. Never move tags between commits.**

This applies to ALL tags (development, release, experimental). Immutability is the foundation of version control — it provides an audit trail, prevents downstream confusion, and forces honest version bumping.

Violating this means breaking the version control contract. When the contract breaks, downstream effects cascade (Docker pulls old versions, package registries serve wrong binaries, CI/CD systems have undefined behavior).

## Definitions: When the Rule Applies

### Tag Scope

**The rule applies to:** All tags, local or remote.

**The stricter interpretation:** Don't create a tag unless you're ready to commit to that version. Once created, that version exists—even locally—and should be pushed to remote and never deleted.

**Exception (rare):** If you created a tag locally by mistake (wrong version number, wrong commit) BEFORE pushing:
- Delete it locally ONLY
- Create the correct version instead
- Never tag the same version twice, even if fixing locally

**Never:**
- Delete a tag after pushing to remote
- Recreate a tag that was already pushed
- Move a tag between commits
- Use "local-only" as an excuse to reuse version numbers

**Test:** Did you INTEND for this version to exist in your release history? If YES → push it. If NO → delete it locally and create the correct version instead.

### Public API Definition (for Semver)

**Public API = What users can call/depend on:**
- MCP tool names and function signatures
- Tool parameters (add new ones = minor version)
- Return types and documented behavior
- Configuration options

**NOT public API = Internal changes:**
- Bug fixes that don't change function behavior
- Performance optimizations (return same results, faster)
- Internal refactoring
- Code structure changes

**Gray area: Behavioral fixes**
If you fix a bug that changes tool behavior but not signature:
- Example: Tool sorted results incorrectly, now sorts correctly
- This IS a behavior change
- **Decision tree:**
  1. "Does the tool's documented contract say it should sort ascending?" YES → Patch (fixing undiscovered bug)
  2. "Were users relying on this behavior in production?" YES → Minor (users will see different output, breaking their workflows)
  3. "Is this internal-only bug with no user impact?" YES → Patch (no observable change from user perspective)
- **Summary: Behavioral bugs that users can detect = minor version. Internal bugs users can't see = patch version.**

### Pre-release Tags

**Pre-release tag format:** `vX.Y.Z-alpha.N`, `vX.Y.Z-rc.N`, etc.

**Scope of immutability:**
- **Published to registries or users?** → Treat as immutable (same rule applies)
- **Internal-only pre-releases on feature branches?** → Can be mutable (recreate `v0.2.0-alpha.1` as needed)
- **Bumped main-branch pre-releases?** → Immutable (create `v0.2.0-alpha.2` instead of recreating `v0.2.0-alpha.1`)

**Rule of thumb:** If it's ever been pushed to a shared remote, treat it as immutable.

## The Iron Rule

```
NO TAG DELETION
NO TAG RECREATION
NO TAG REUSE
NO TAG MOVEMENT

Once pushed, it stays.
If you need to fix something, create a new version.
```

**No exceptions:**
- Not when "it was only up for a few minutes"
- Not when "nobody pulled it yet"
- Not when "it's a technical fix, not a feature"
- Not when "I'm confident I know what happened"
- Not when "we're pre-1.0 and things are experimental"

All of these mean: Delete the tag locally only. Create a new version instead.

## Why This Matters

### The "But Nobody Saw It" Rationalization

```
❌ WRONG: "The tag was only up 2 hours, no harm deleting it"

✅ RIGHT: GitHub webhooks fired in those 2 hours
          Release artifacts were generated
          Someone might have pulled the code
          Your CI/CD logged that version as released
          When you delete it, you create orphaned history
```

### The "It's Just a Technical Fix" Rationalization

```
❌ WRONG: "This bug fix doesn't deserve a new version"

✅ RIGHT: Semver says: new code = new version
          A fix IS a change
          Immutability matters more than version elegance
          v0.1.0 → v0.1.1 is honest
          v0.1.0 (deleted) → v0.1.0 (recreated) is dishonest
```

### The "0.x Is Experimental" Rationalization

```
❌ WRONG: "We're pre-1.0, so version rules don't fully apply"

✅ RIGHT: Immutability applies equally to 0.x and 1.x
          0.x doesn't give you a "pass" on version control
          If your versioning means nothing, stop using it
          If you use versioning, use it honestly
```

## Semantic Versioning Quick Decision

When deciding version bumps, ask in order:

| Question | Answer → Next | Answer → Next |
|----------|---------------|---------------|
| **Breaking changes?** | YES → MAJOR | NO → Check 2 |
| **New public API?** | YES → MINOR | NO → Check 3 |
| **Bug fixes only?** | YES → PATCH | NO → ?? (you added something) |

**Result:** Major.Minor.Patch

**Current project example (gitlabmcp):**
- v0.1.0 (initial): 57 tools, 10 domains (deliberate release)
- v0.2.0 (if adding 5 new tools): minor bump (new public API)
- v0.1.1 (if fixing a bug): patch bump (no new API)
- v1.0.0 (if breaking existing API): only when you change/remove tools

## Red Flags — STOP and Start Over

If you see yourself thinking:

- "The tag was only up a few minutes"
- "Nobody checked out that version yet"
- "It's a small fix, not worth a new version"
- "Recreating the tag is faster than bumping"
- "We're pre-1.0, so versioning is flexible"
- "This is a technical issue, not a feature"
- "I'll just delete it locally and push again"
- "The user won't notice which version they got"

**All of these mean: STOP. Create a new version instead.**

## Common Rationalizations & Counter-Reality

| Rationalization | Why It's Wrong | What To Do |
|---|---|---|
| "Nobody saw it" | Webhooks fired, artifacts generated, CI/CD logged it | Create v0.1.1 |
| "It's a technical fix" | Semver = code changed. Semver = version changed. | Patch bump required |
| "0.x is experimental" | Immutability applies to all versions equally | Use proper versioning or stop versioning |
| "I'm confident nobody pulled it" | Confidence ≠ certainty. Version control doesn't run on confidence | Create new tag |
| "This feels wrong as a new version" | Feeling wrong ≠ technically wrong. Trust semver | v0.1.0 → v0.1.1 or v0.2.0 |
| "Other projects reuse tags" | Other projects have broken release histories | Not our standard |
| "Just this once" | "Just this once" becomes precedent | Never this once |

## Real-World Example: gitlabmcp v0.1.0 → v0.1.1

**What happened:**
1. Released v0.1.0
2. GoReleaser config had wrong repository owner (`ilya` instead of `ilyalaletin`)
3. Release workflow failed
4. Temptation: Delete v0.1.0, fix config, recreate v0.1.0

**What we did instead:**
1. Fixed the config (1 line change)
2. Committed fix to main
3. Deleted v0.1.0 locally only
4. Created v0.1.1 tag
5. Pushed v0.1.1 to remote

**Why this matters:**
- v0.1.0 failure is now visible in git history
- v0.1.1 is the actual working release
- Anyone looking at releases sees: "v0.1.0 failed, v0.1.1 worked"
- No hidden history, no orphaned commits
- Future releases follow honest versioning

## Implementation Checklist

Before tagging a release:

- [ ] Increment version correctly (MAJOR.MINOR.PATCH)
- [ ] Update any version references in docs/code (go.mod, README, etc.)
- [ ] Commit version changes with message: `chore: bump version to X.Y.Z`
- [ ] Create annotated tag: `git tag -a vX.Y.Z -m "Release X.Y.Z"`
- [ ] Push tag: `git push origin vX.Y.Z`
- [ ] Verify GitHub Release created automatically (check Actions)
- [ ] **NEVER delete, move, or recreate a tag after pushing**

## Violating the Rule

If you violate tag immutability:

1. **Stop immediately** — don't push anything else
2. **Document what you did** — commit message explaining the violation
3. **Inform team** — if this is a shared project, announce immediately
4. **Audit artifacts** — check what got built/released under that tag
5. **Fix forward** — create a new version, move on
6. **Review this skill** — why did you rationalize away the rule?

Do NOT attempt to "undo" by deleting the tag remotely. That creates worse problems (broken references, CI/CD confusion). Accept the mistake, version forward.

## Version Numbering Strategy for This Project

**gitlabmcp versioning approach:**

- **v0.x.y** while:
  - API design might change
  - Tool set might be reorganized
  - Core behavior might shift

- **Move to v1.0.0 when:**
  - Tool API is stable (no removal/renaming of tools)
  - Core behavior is predictable
  - Users can rely on consistency

**Current policy:**
- New tools = minor version (v0.1.0 → v0.2.0)
- Bug fixes = patch version (v0.1.0 → v0.1.1)
- Removed/changed tools = major version (v0.1.0 → v1.0.0)

This decision is made now. Stick to it. No exceptions.

---

## Why Immutability Wins (Even When It Hurts)

The friction of immutability is a feature, not a bug:

1. **It forces honest decisions** — you can't hide mistakes in version history
2. **It creates audit trail** — anyone can trace what shipped when
3. **It prevents cascade failures** — no orphaned tags, no broken references
4. **It teaches discipline** — you learn to get releases right the first time

The cost of reusing tags (even once) is always higher than the cost of bumping versions.

**Choose the harder path. It leads to better software.**

## Team Enforcement & Historical Violations

### Who Enforces This?

**In code review:** Reviewer catches tag deletion/recreation in git log and blocks the PR.

**In CI/CD:** Ideally, branch protection rules prevent force-pushing to main and tag modifications.

**In practice:** Team discipline. This skill exists because:
1. Automation isn't always available
2. Even with automation, understanding WHY matters
3. Prevention is better than detection

### Historical Violations

**If you discover a deleted/recreated tag 3 months later:**

1. **Document it** — add comment in relevant issue/PR explaining discovery date
2. **Don't recreate retroactively** — that creates more confusion (two different attempts to fix)
3. **Move forward** — enforce going forward, not backwards
4. **Improve guardrails** — ask: "Why wasn't this caught? Can we add automation?"

**Your responsibility:** Prevent it happening again, not punishing past violations.

### Pre-Deployment Checklist (Team)

Before releasing to production/registry:
- [ ] Tag is on main/release branch (not feature branch)
- [ ] Tag is annotated: `git tag -a vX.Y.Z -m "Release X.Y.Z"`
- [ ] Release notes prepared
- [ ] Semver bump is correct (use decision tree)
- [ ] No force-pushes planned
- [ ] CI/CD will automatically detect new tag and build

**After pushing tag:**
- [ ] Verify GitHub Release created automatically
- [ ] Verify binaries/artifacts built (not manual uploads)
- [ ] **Never touch the tag again**
