---
name: "⬆️ Dependency Update"
about: "Review and validate dependency updates created by Dependabot"
title: "⬆️ chore(deps): bump <package-name> from <old-version> to <new-version>"
labels: ["chore", "dependencies"]
assignees: []
---

# ⬆️ Dependency Update

## 🧾 Summary

- [Pull Request](URL)

1. Replace URL with the Pull Request URL
2. Edit the PR (Pull Request)
3. Copy the MD (Markdown) from the PR
4. Replace this three lines with the copied content


---

## 🧪 Validation (Required before merge)
- [ ] Project builds successfully
- [ ] All tests pass
- [ ] No regressions observed
- [ ] Critical flows tested manually
- [ ] Edge cases verified
- [ ] Do not update if code breaks:
  - [ ] Talk to PM
  - [ ] Close it with note.

---

## 🚨 Pre-Merge Testing (Mandatory)
- [ ] Fully tested locally or in staging environment
- [ ] Verified no impact on core functionality
- [ ] Checked logs / console for hidden errors
- [ ] Confirmed compatibility with related dependencies

---

## ⚠️ Impact
- [ ] No breaking changes
- [ ] Potential breaking changes (details below)

**Details (if applicable):**  
<!-- Add migration steps, deprecations, or risks -->

---

## 🔍 Notes
<!-- Optional: security fixes, performance improvements, release notes -->

---

## ✅ Checklist
- [ ] Code reviewed
- [ ] Approved
- [ ] Thoroughly tested (required)
- [ ] Merged
