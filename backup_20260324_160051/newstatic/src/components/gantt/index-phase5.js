/**
 * Gantt Chart Components - Phase 5 Exports
 * UX Optimization and Automation Components
 */

// Overlays
export { default as GuidedTour } from './overlays/GuidedTour.vue'
export { default as ContextMenu } from './overlays/ContextMenu.vue'
export { default as Minimap } from './overlays/Minimap.vue'

// Panels
export { default as SmartSuggestionsPanel } from './panels/SmartSuggestionsPanel.vue'
export { default as SuggestionCard } from './panels/SuggestionCard.vue'

// Dialogs
export { default as TemplateManagerDialog } from './dialogs/TemplateManagerDialog.vue'

// Utilities
export {
  getTourSteps,
  tourConfig,
  tourTypes,
  getCompletedTours,
  markTourCompleted,
  isTourCompleted,
  resetTourProgress
} from '@/utils/tourSteps'

export {
  WorkflowAutomationEngine,
  getWorkflowEngine,
  executeAutomation,
  RuleTypes,
  TriggerTypes,
  defaultRules
} from '@/utils/workflowAutomation'

export {
  analyzeSchedule,
  predictDelays,
  optimizeSchedule,
  createSuggestion
} from '@/utils/aiOptimizer'

// Stores
export { default as useTemplateStore } from '@/stores/templateStore'
