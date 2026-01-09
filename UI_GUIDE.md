# AI Factory - Professional UI Guide

## Overview

The new AI Factory UI is a modern, professional interface designed for ease of use by you and your partners. It features real-time updates, plan approval workflows, image uploads, and comprehensive project management.

---

## Key Features

### 🎨 **Modern Design**
- Dark theme optimized for long sessions
- Professional color scheme with gradient accents
- Smooth animations and transitions
- Responsive layout for all screen sizes

### 🔴 **Real-Time Updates**
- Live WebSocket connection shows connection status
- Instant notifications when phases complete
- Activity feed updates automatically
- No page refresh needed

### 📋 **Plan Approval Workflow**
- Review AI-generated implementation plans
- See files to create/modify, tech stack, testing strategy
- Approve or reject with feedback
- Plan regenerates based on your input

### 📸 **Image Upload**
- Drag-and-drop design mockups
- Multiple image support
- Preview before submission
- Images automatically referenced in project description

---

## UI Sections

### 1. Header
**Location:** Top of page, always visible

**Components:**
- **AI Factory Logo** - Branding with "Professional MVP Generator with Triple Guarantee™"
- **System Status** - Green badge shows "Online" when backend is connected
- **WebSocket Status** - Purple badge shows real-time connection status

**Status Indicators:**
- 🟢 **Online** - System ready
- 🟣 **WebSocket: Connected** - Real-time updates active
- 🔴 **Offline** - Check if server is running

---

### 2. Tab Navigation

**Three Main Tabs:**

#### **Create Project** (Default)
Where you create new MVP projects

#### **Projects**
View all your projects and their status

#### **Activity Log**
Full history of all system events

---

## Creating Your First Project

### Step 1: Project Details

1. Enter **Project Name**
   - Example: "E-commerce Platform"
   - Keep it descriptive and professional

2. Write **Project Description**
   - Be specific about features needed
   - Include tech stack preferences
   - Mention any special requirements
   - Example:
   ```
   Build a modern e-commerce platform with React, Node.js, and PostgreSQL.
   Include:
   - User authentication (JWT)
   - Product catalog with search
   - Shopping cart functionality
   - Stripe payment integration
   - Admin dashboard
   - Responsive design
   ```

### Step 2: Attach Design Mockups (Optional)

**Two Ways to Upload:**

1. **Drag & Drop**
   - Drag image files directly onto the upload area
   - Multiple files supported

2. **Click to Upload**
   - Click the upload area
   - Select one or more image files
   - Supports PNG, JPG up to 10MB each

**Preview:**
- Thumbnails appear below upload area
- Click ×  to remove unwanted images
- Images are automatically referenced in project description

### Step 3: Create Project

Click **🚀 Create Project** button

**What Happens:**
1. Project is created
2. Discovery phase starts automatically
3. Real-time updates show in Activity Feed
4. Phases progress: Discovery → Validation → Planning
5. System pauses at "Waiting Approval" for your review

---

## Plan Approval Process

### When Plan is Ready

**Modal Appears Automatically** showing:

#### 📋 Implementation Approach
High-level strategy for building your project

#### 📁 Files to Create
List of all new files the AI will generate:
- `src/components/ProductCard.jsx`
- `server/routes/api.js`
- `database/schema.sql`
- etc.

#### ✏️ Files to Modify
Existing files that will be updated (if any)

#### 🛠️ Technology Stack
Technologies that will be used:
- React 18
- Node.js/Express
- PostgreSQL
- Tailwind CSS
- etc.

#### 🧪 Testing Strategy
How tests will be implemented:
- Unit tests with Jest
- Integration tests
- E2E tests with Cypress
- etc.

#### Complexity Badge
- 🟢 **Low** - Simple project
- 🟡 **Medium** - Moderate complexity
- 🔴 **High** - Complex architecture

#### Estimated Time
Expected generation time: "45 mins", "2 hours", etc.

### Your Options

**✅ Approve Plan**
- Click green "Approve Plan" button
- Code generation starts immediately
- Progress shown in real-time

**❌ Reject & Revise**
- Click red "Reject & Revise" button
- Enter feedback explaining what to change
- AI regenerates plan with your feedback
- New plan presented for approval

---

## Activity Feed

### Live Activity Panel
**Location:** Right side of Create Project tab

**Shows:**
- ✅ Success events (green) - Phase completed
- ℹ️ Info events (blue) - Phase started, system updates
- ⚠️ Warning events (orange) - Errors or issues

**Features:**
- Auto-updates via WebSocket
- Shows last 20 events
- Timestamps on each event
- Event counter at top

### Full Activity Log Tab
**Location:** Activity Log tab

**Shows:**
- Complete history of all events
- Same color coding as live feed
- Scrollable list
- Clear button to reset

---

## Projects Tab

### Your Projects List

**Each Project Card Shows:**
- **Project Name** - Bold title
- **Status Badge** - Active/Complete/Blocked
- **Description** - First 150 characters
- **Created Date** - When project was created
- **Current Phase** - Discovery, Planning, CodeGen, etc.
- **Thinking Mode Badge** - Fast/Normal/Extended (if shown)
- **Phase Indicator** - Pulsing dot shows active phase

**Hover Effects:**
- Card lifts slightly
- Border turns blue
- Smooth animation

**Click to View:**
- Opens project details
- Shows plan if awaiting approval
- Displays full status

---

## Real-Time Features

### WebSocket Updates

**Automatic Notifications For:**
- Phase transitions
- Status changes
- Errors or warnings
- Completion events

**No Refresh Needed:**
- Activity feed updates instantly
- Project status changes live
- Phase indicators update automatically

### Status Indicators

**Connection Status:**
- **Green dot pulsing** - Connected
- **Red dot** - Disconnected
- **Yellow dot** - Reconnecting

**System automatically reconnects** if connection drops

---

## Professional Features

### 1. Thinking Modes (Automatic)

**Shown as badge on project cards:**

- **Fast** 🔵 - Simple projects
  - Quick responses
  - Minimal reasoning
  - Low complexity (score 1-3)

- **Normal** 🟢 - Standard projects
  - Balanced reasoning
  - Good for most projects
  - Medium complexity (score 4-6)

- **Extended** 🟣 - Complex projects
  - Deep step-by-step thinking
  - Considers edge cases
  - High complexity (score 7-10)

**Auto-selected based on your project description**

### 2. Git Worktree Isolation

**When code generation starts:**
- Creates isolated git branch: `ai-factory/{project-id}`
- All changes go to separate worktree
- Safe from conflicts with your main branch
- Auto-commits with proper attribution
- Cleanup on project completion

**You don't see this** - it happens automatically in the background

### 3. Triple Guarantee System™

**Every project verified:**
- ✅ Build Guarantee - Code compiles
- ✅ Runtime Guarantee - App starts
- ✅ Test Guarantee - Tests pass

**Quality Certificate** generated automatically

---

## Tips for Best Results

### Writing Project Descriptions

**Be Specific:**
✅ "Build a task manager with React, user auth, and real-time updates"
❌ "Make me a to-do app"

**Include Tech Stack:**
✅ "Use TypeScript, Next.js 14, Prisma, and PostgreSQL"
❌ "Use modern technologies"

**List Features:**
✅ "Include: drag-and-drop, filters, tags, due dates, notifications"
❌ "With some useful features"

**Mention Testing:**
✅ "Include unit tests with Jest, >80% coverage"
❌ "Should have tests"

### Using Images Effectively

**Best Practices:**
- Upload UI mockups or wireframes
- Include design references
- Show desired layout/style
- Clear, high-resolution images
- Reference specific UI elements in description

**Example:**
"Build the landing page shown in the attached mockup. Match the color scheme and layout exactly."

### Reviewing Plans

**Check These:**
- ✅ All required features mentioned
- ✅ Tech stack matches your preferences
- ✅ File structure makes sense
- ✅ Testing approach is adequate
- ✅ No missing dependencies

**Reject If:**
- ❌ Wrong technologies chosen
- ❌ Missing critical features
- ❌ Files don't match project structure
- ❌ Complexity seems off

---

## Keyboard Shortcuts

- **Tab** - Switch between sections
- **Enter** - Submit form (when focused)
- **Esc** - Close modal
- **Ctrl+R** - Refresh projects list

---

## Troubleshooting

### "Offline" Status

**Problem:** Red status badge shows "Offline"

**Solution:**
1. Check if server is running
2. Run: `ai_factory_final.exe -mode=server`
3. Wait for "AI Studio Orchestrator on port 8080"
4. Refresh browser

### WebSocket Disconnected

**Problem:** Purple badge shows "Disconnected"

**Solution:**
- System auto-reconnects in 5 seconds
- No action needed
- If persists, refresh page

### No Projects Loading

**Problem:** "No projects yet" shown when you have projects

**Solution:**
1. Click 🔄 Refresh button
2. Check server console for errors
3. Verify Ollama is running

### Plan Modal Not Appearing

**Problem:** Phase shows "waiting_approval" but no modal

**Solution:**
1. Click the project card in Projects tab
2. Modal should appear
3. Refresh page if needed

### Image Upload Failed

**Problem:** Image doesn't upload

**Solution:**
- Check image size (<10MB)
- Verify format (PNG, JPG only)
- Check server console for errors
- Try different image

---

## For Partners & Demonstrations

### Showcasing Features

**Start with a Simple Project:**
1. Create "React Counter App"
2. Show real-time phase updates
3. Approve plan when prompted
4. Watch code generate live
5. Show quality certificate

**Then Show Complex Project:**
1. Create e-commerce platform with mockups
2. Demonstrate thinking mode (Extended)
3. Reject first plan to show revision
4. Approve revised plan
5. Show git worktree isolation

### Highlighting Value

**Key Points:**
- **Speed**: From idea to code in minutes
- **Quality**: Triple Guarantee ensures it works
- **Safety**: Git isolation prevents conflicts
- **Control**: Plan approval gives oversight
- **Visibility**: Real-time updates show progress
- **Professional**: Quality certificates for clients

### Pricing Context

**When discussing pricing:**
- "This generates production-ready MVPs"
- "Every project includes Triple Guarantee"
- "Quality certificates prove it works"
- "Compare to weeks of manual development"
- "$800-$2,500 per MVP with automated quality verification"
- "Money-back guarantee because we verify quality"

---

## Advanced Features

### Activity Log Export

**Coming Soon:**
- Export activity log as CSV
- Share with team members
- Audit trail for projects

### Project Templates

**Coming Soon:**
- Save project configurations
- Reuse common setups
- Share templates with team

### Batch Operations

**Coming Soon:**
- Create multiple projects
- Bulk approve/reject
- Parallel generation

---

## System Requirements

### Browser Compatibility
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

### Network
- Stable internet connection
- WebSocket support required
- Port 8080 accessible

### Backend
- AI Factory server running
- Ollama installed and running
- Sufficient disk space for projects

---

## Getting Help

### Check Logs
- Server console shows detailed logs
- Activity feed shows user-facing events
- Browser console shows client errors

### Common Issues
- 90% of issues: Server not running
- 5%: Ollama not running
- 5%: Network/firewall issues

### Verify Setup
```bash
# 1. Check Ollama
ollama list

# 2. Start server
ai_factory_final.exe -mode=server

# 3. Open browser
http://localhost:8080
```

---

## Summary

The new professional UI makes AI Factory:
- **Easy to use** - Intuitive interface, no training needed
- **Professional** - Impressive to partners and clients
- **Powerful** - All features accessible, real-time updates
- **Reliable** - Clear status indicators, error handling

**Perfect for demonstrations, daily use, and partner presentations!**
