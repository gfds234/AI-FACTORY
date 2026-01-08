# Partner Testing Guide - AI FACTORY First User Test

## Overview

This guide is for conducting the **first usability test** with a partner. The goal is to validate the workflow, identify UX issues, and gather feedback before cloud deployment.

---

## Test Setup

### Current Access
- **Public URL**: https://condensative-irreplaceable-makhi.ngrok-free.dev
- **Local URL**: http://localhost:8080 (if running locally)
- **Status**: Orchestrator running, ngrok tunnel active

### Partner Prerequisites
- Web browser (Chrome, Firefox, Edge)
- No technical knowledge required
- 30-45 minutes for testing session

---

## Partner Onboarding Script

### Step 1: Share Access (5 min)

**Message template to send partner:**

```
Hi [Partner Name],

I'd love your help testing AI FACTORY - an autonomous MVP generation tool. You'll be the first user to test it!

🔗 Access here: https://condensative-irreplaceable-makhi.ngrok-free.dev

What to expect:
- You'll create a simple project idea (like "todo app")
- Watch AI FACTORY autonomously build it through 6 phases
- See real-time progress in List or Kanban view
- Download the completed project

Time needed: 30-45 minutes
No technical skills required - just your honest feedback!

Let me know when you're ready and I'll walk you through it.
```

### Step 2: Initial Walkthrough (10 min)

**Screen share with partner and guide them through:**

1. **Landing page** - First impressions?
   - Is it clear what AI FACTORY does?
   - Is the UI intuitive?

2. **Toggle views** - Click "📊 Kanban Board" and "📋 List View"
   - Which view do you prefer?
   - Is the difference clear?

3. **Create Project button** - Find and click it
   - Is it obvious where to start?

### Step 3: Create First Project (15 min)

**Test Scenario 1: Simple Todo App**

Guide partner to create project:
- **Name**: "My First Todo App"
- **Description**: "A simple todo list web app with React. Users can add, complete, and delete tasks."

**What to observe:**
- Do they understand the form fields?
- Are they confident clicking "Create Project"?
- Do they know what will happen next?

**Expected behavior:**
- Project appears in list/kanban view
- Phase starts at "Discovery"
- Real-time updates show progress
- They can see the project moving through phases

**Questions to ask during execution:**
1. "What do you think is happening right now?"
2. "Is the progress clear?"
3. "What would you like to see that's missing?"
4. "Is anything confusing?"

### Step 4: Monitor Progress (10 min)

**Test Scenario 2: Phase Transitions**

Watch together as project moves through phases:
- Discovery → Validation → Planning → Code Gen → Review → QA → Docs → Complete

**What to observe:**
- Do they understand what each phase does?
- Is the progress bar meaningful?
- Do they know how long it will take?
- Are error messages (if any) helpful?

**Questions to ask:**
1. "Can you explain what's happening in [current phase]?"
2. "Do you feel informed about progress?"
3. "What additional information would help?"

### Step 5: Download Project (5 min)

**Test Scenario 3: Getting Results**

Once project completes:
- Guide them to download/export feature
- Have them download the ZIP file
- Ask them to open and explore files

**What to observe:**
- Is download button easy to find?
- Is ZIP extraction clear?
- Do they understand what they received?

**Questions to ask:**
1. "Is this what you expected?"
2. "What would you do with this code?"
3. "Do you feel confident the code works?"

---

## Feedback Collection

### During Testing - Note Taking Template

```markdown
## Partner Testing Session - [Date]

### Partner Info
- Name: [Partner Name]
- Background: [Technical/Non-technical]
- Use case: [Why they're interested in AI FACTORY]

### First Impressions (0-5 min)
- Landing page reaction:
- Initial confusion points:
- What stood out positively:
- What stood out negatively:

### Creating Project (5-15 min)
- Understood form fields: YES / NO / PARTIALLY
- Confident to proceed: YES / NO
- Confusion points:
- Suggestions:

### Monitoring Progress (15-25 min)
- Understood phases: YES / NO / PARTIALLY
- Found progress clear: YES / NO
- Wanted more info about:
- Suggestions:

### Downloading Results (25-30 min)
- Found download easily: YES / NO
- Understood deliverable: YES / NO
- Satisfied with result: YES / NO
- Suggestions:

### Overall Feedback (30-45 min)
- What worked well:
- What was confusing:
- What was missing:
- Would they use it again: YES / NO / MAYBE
- Would they pay for it: YES / NO / MAYBE (if yes, how much?)
- Net Promoter Score (0-10):
```

### Post-Test Survey

**Send this after the session:**

```markdown
# AI FACTORY Feedback Survey

Thank you for testing! Please answer these quick questions:

## Usability (1-5 stars)
- How easy was it to create a project? ⭐⭐⭐⭐⭐
- How clear was the progress tracking? ⭐⭐⭐⭐⭐
- How intuitive was the interface? ⭐⭐⭐⭐⭐
- Overall experience? ⭐⭐⭐⭐⭐

## Features
What did you like most?
[Free text]

What was most confusing?
[Free text]

What's missing?
[Free text]

## Value Proposition
Would you use AI FACTORY for real projects?
[ ] Yes, definitely
[ ] Maybe, if improved
[ ] No

If this saved you 10+ hours, what would you pay per MVP?
[ ] $0 (wouldn't pay)
[ ] $50-100
[ ] $200-500
[ ] $500-1000
[ ] $1000-2500
[ ] $2500+

## Open Feedback
Any other thoughts?
[Free text]
```

---

## Test Scenarios (Copy-Paste Ready)

### Scenario 1: Todo App
```
Name: Simple Todo List
Description: A web-based todo app with React. Users can add new tasks, mark them as complete, and delete tasks. Include basic styling with CSS.
```

### Scenario 2: Landing Page
```
Name: Product Landing Page
Description: A modern landing page for a SaaS product. Include hero section, features, pricing table, and contact form. Use Tailwind CSS for styling.
```

### Scenario 3: REST API
```
Name: User Management API
Description: A REST API for user management with Node.js and Express. Include endpoints for user registration, login (JWT auth), get user profile, and update profile. Use MongoDB for storage.
```

### Scenario 4: Dashboard
```
Name: Analytics Dashboard
Description: A simple analytics dashboard with React and Chart.js. Display sample metrics like user growth, revenue trends, and top features. Include date range selector.
```

---

## Success Metrics

### Qualitative Goals
- ✅ Partner completes project creation without help
- ✅ Partner understands what's happening during execution
- ✅ Partner finds and downloads result successfully
- ✅ Partner expresses interest in using again

### Quantitative Benchmarks
- **Time to first project**: < 2 minutes
- **Confusion incidents**: < 3 during session
- **Questions asked**: < 5 clarifying questions
- **Overall satisfaction**: ≥ 4/5 stars
- **Would recommend**: ≥ 7/10 (NPS)

---

## Common Issues & Solutions

### Issue: "I don't know what to create"
**Solution**: Provide test scenarios above, start with Scenario 1

### Issue: "How long will this take?"
**Solution**: Add estimated time to UI (future improvement)

### Issue: "Is it working? Nothing's happening"
**Solution**: Improve loading indicators (future improvement)

### Issue: "What do I do with the downloaded code?"
**Solution**: Add README to generated projects explaining next steps

### Issue: "Does the code actually work?"
**Solution**: Point to Quality Certificate showing test results

---

## Post-Test Action Items

### Immediate Fixes (Do before next test)
1. [ ] Fix any blocking bugs discovered
2. [ ] Add missing critical features mentioned
3. [ ] Improve confusing UI elements

### Short-term Improvements (Do before cloud launch)
1. [ ] Address top 3 usability issues
2. [ ] Add requested "nice-to-have" features
3. [ ] Improve onboarding flow based on feedback

### Long-term Roadmap (Post-launch)
1. [ ] Advanced features requested
2. [ ] Integration requests
3. [ ] Pricing model validation

---

## Partner Follow-Up Template

**After fixing issues:**

```
Hi [Partner Name],

Thanks again for testing AI FACTORY! Your feedback was incredibly valuable.

Based on your suggestions, I've made these improvements:
✅ [Improvement 1]
✅ [Improvement 2]
✅ [Improvement 3]

Would you be willing to do a quick 10-min follow-up test to see if these changes help?

Also, I'm planning to launch publicly soon. Would you be interested in:
- Being a beta tester for new features?
- Being featured as a case study?
- Getting early access pricing?

Let me know!
```

---

## Next Steps After Successful Test

1. **Iterate**: Fix top issues, improve UX
2. **Expand testing**: 2-3 more partners with different backgrounds
3. **Cloud deployment**: Deploy to Fly.io for permanent access
4. **Public launch**: Open to broader audience
5. **Monetization**: Implement pricing based on feedback

---

**Remember**: This partner is helping YOU build a better product. Their confusion is your opportunity to improve!
