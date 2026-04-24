---
name: "🖥️ Frontend User Story"
about: "Create or update a feature for the frontend application"
title: "🖥️ [Brief, descriptive title of the Frontend Technical Story]"
labels: ["frontend", "enhancement"]
assignees: []
---

# 🖥️ Frontend User Story

## 🎯 User Story
As a **[type of user]**,\n 
I want to **[goal/desire]**,/n 
so that I can **[benefit/value]**.

## 🧭 Acceptance Criteria
- [ ] The feature should be responsive and match the design in Figma.
- [ ] The user can perform [specific action] easily.
- [ ] The UI/UX follows the design system guidelines.
- [ ] The feature works on all major browsers (Chrome, Firefox, Safari, Edge).
- [ ] Accessibility (a11y) standards are met (WCAG 2.1).

## 🧩 UI Design / References
- Figma Link: [Paste here]
- Screenshots / Mockups: [Attach if needed]

## 🧪 Testing Notes
- [ ] Unit tests cover main logic and components.
- [ ] End-to-end tests (Cypress/Playwright) verify user flow.
- [ ] Manual testing steps: [describe]

## ⚙️ Technical Notes
- Framework: [React / Vue / Angular]
- Components affected: [list or link to components]
- State management: [Redux / Zustand / Context / Pinia]
- API endpoints to consume: [list endpoints]

## 🚀 Definition of Done
- [ ] Code reviewed and approved.
- [ ] All tests passing (unit + E2E).
- [ ] Deployed to staging and validated by QA.
- [ ] Merged to `main`.

## 📝 API Contract

| Description | Method | HTTP Status Codes | Route               | Response     |
|-------------|--------|-------------------|---------------------|--------------|
| Get By Id   | Get    | 200, 404, 500     | /api/v1/example/:id | `ExampleDTO` |

ExampleDTO:

```json
{
  guid: "string"
  title: "string"
}
```