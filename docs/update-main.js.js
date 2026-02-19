/**
 * Gantt Chart Integration - Main.js Update Script
 *
 * This script contains all the required registrations and imports
 * for the Gantt Chart Editor integration.
 *
 * Run this script to update your main.js file:
 * node docs/update-main.js.js
 */

const fs = require('fs')
const path = require('path')

const mainJsPath = path.join(__dirname, '../newstatic/src/main.js')

// Read current main.js
let mainJsContent = ''
try {
  mainJsContent = fs.readFileSync(mainJsPath, 'utf-8')
  console.log('✓ Read current main.js')
} catch (error) {
  console.error('✗ Failed to read main.js:', error.message)
  process.exit(1)
}

// Check if already updated
if (mainJsContent.includes('vue-virtual-scroller') &&
    mainJsContent.includes('vue-tour')) {
  console.log('\n⚠ main.js already contains Gantt integration code')
  console.log('Skipping update...\n')
  process.exit(0)
}

// Required imports to add
const importsToAdd = `
// ==================== Gantt Chart Editor Integration ====================
// Virtual Scroller (REQUIRED for performance with large datasets)
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import { RecycleScroller } from 'vue-virtual-scroller'

// Guided Tour (for first-time user onboarding)
import VueTour from 'vue-tour'
import 'vue-tour/dist/vue-tour.css'

// Gantt stores (if needed globally)
// import { useGanttStore } from '@/stores/ganttStore'
// import { useUndoRedoStore } from '@/stores/undoRedoStore'
// import { useCalendarStore } from '@/stores/calendarStore'
// import { useCollaborationStore } from '@/stores/collaborationStore'
// import { useTemplateStore } from '@/stores/templateStore'

`

// Component registrations to add
const registrationsToAdd = `
  // ==================== Gantt Chart Components ====================
  // Virtual Scroller - CRITICAL for performance
  app.component('RecycleScroller', RecycleScroller)

  // Vue Tour - Guided tour for new users
  app.use(VueTour)

  // Initialize Gantt stores (optional, can be done in components)
  // const ganttStore = useGanttStore()
  // const undoRedoStore = useUndoRedoStore()
  // const calendarStore = useCalendarStore()
  // const collaborationStore = useCollaborationStore()
  // const templateStore = useTemplateStore()

  // Provide stores globally (optional)
  // app.provide('ganttStore', ganttStore)
  // app.provide('undoRedoStore', undoRedoStore)
  // app.provide('calendarStore', calendarStore)
  // app.provide('collaborationStore', collaborationStore)
  // app.provide('templateStore', templateStore)
`

// Find where to insert (after ElementPlus imports, before app creation)
const importInsertIndex = mainJsContent.indexOf("import App from")
if (importInsertIndex === -1) {
  console.error('✗ Could not find insertion point for imports')
  process.exit(1)
}

// Insert imports
mainJsContent = mainJsContent.slice(0, importInsertIndex) +
               importsToAdd +
               mainJsContent.slice(importInsertIndex)

// Find where to insert registrations (before app.mount)
const mountIndex = mainJsContent.indexOf('app.mount(')
if (mountIndex === -1) {
  console.error('✗ Could not find app.mount()')
  process.exit(1)
}

// Insert registrations
mainJsContent = mainJsContent.slice(0, mountIndex) +
               registrationsToAdd +
               '\n' +
               mainJsContent.slice(mountIndex)

// Backup original file
const backupPath = mainJsPath + '.backup'
try {
  fs.writeFileSync(backupPath, fs.readFileSync(mainJsPath))
  console.log('✓ Created backup at:', backupPath)
} catch (error) {
  console.error('✗ Failed to create backup:', error.message)
}

// Write updated main.js
try {
  fs.writeFileSync(mainJsPath, mainJsContent, 'utf-8')
  console.log('✓ Updated main.js with Gantt integration code')
} catch (error) {
  console.error('✗ Failed to write main.js:', error.message)
  process.exit(1)
}

console.log('\n✅ Integration complete!')
console.log('\nNext steps:')
console.log('1. Review the changes in main.js')
console.log('2. Uncomment store initialization if needed')
console.log('3. Restart dev server: npm run dev')
console.log('4. Test Gantt chart functionality\n')
console.log('If something goes wrong, restore from backup:')
console.log(`  cp ${backupPath} ${mainJsPath}\n`)

// Verify updated content
console.log('Updated main.js preview:')
console.log('─'.repeat(60))
const lines = mainJsContent.split('\n')
const ganttLines = lines.filter(line =>
  line.includes('vue-virtual-scroller') ||
  line.includes('vue-tour') ||
  line.includes('RecycleScroller')
)
ganttLines.slice(0, 10).forEach(line => console.log(line))
if (ganttLines.length > 10) {
  console.log('... and more')
}
console.log('─'.repeat(60))
