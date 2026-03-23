/**
 * Gantt Chart Guided Tour Steps
 * Defines tour steps for different tour types and languages
 */

import i18n from '@/i18n'

const { t } = i18n.global

/**
 * Basic tour steps for first-time users
 * Covers essential features and navigation
 */
export const basicTourSteps = [
  {
    target: '#gantt-toolbar',
    content: t('gantt.tour.steps.toolbar.content'),
    header: t('gantt.tour.steps.toolbar.header'),
    params: { placement: 'bottom' },
    before: () => {
      // Ensure toolbar is visible
      return new Promise((resolve) => {
        const checkInterval = setInterval(() => {
          const toolbar = document.querySelector('#gantt-toolbar')
          if (toolbar) {
            clearInterval(checkInterval)
            resolve(true)
          }
        }, 100)
      })
    }
  },
  {
    target: '#task-list',
    content: t('gantt.tour.steps.taskList.content'),
    header: t('gantt.tour.steps.taskList.header'),
    params: { placement: 'right' }
  },
  {
    target: '#gantt-timeline',
    content: t('gantt.tour.steps.timeline.content'),
    header: t('gantt.tour.steps.timeline.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '#gantt-zoom-controls',
    content: t('gantt.tour.steps.zoom.content'),
    header: t('gantt.tour.steps.zoom.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="add-task"]',
    content: t('gantt.tour.steps.addTask.content'),
    header: t('gantt.tour.steps.addTask.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="edit-task"]',
    content: t('gantt.tour.steps.editTask.content'),
    header: t('gantt.tour.steps.editTask.header'),
    params: { placement: 'bottom' },
    before: () => {
      // Select a task first
      const firstTask = document.querySelector('.task-row')
      if (firstTask) {
        firstTask.click()
      }
      return Promise.resolve(true)
    }
  },
  {
    target: '[data-tour="dependencies"]',
    content: t('gantt.tour.steps.dependencies.content'),
    header: t('gantt.tour.steps.dependencies.header'),
    params: { placement: 'left' }
  },
  {
    target: '#gantt-status-bar',
    content: t('gantt.tour.steps.statusBar.content'),
    header: t('gantt.tour.steps.statusBar.header'),
    params: { placement: 'top' }
  },
  {
    target: '[data-tour="view-switcher"]',
    content: t('gantt.tour.steps.viewSwitcher.content'),
    header: t('gantt.tour.steps.viewSwitcher.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="export"]',
    content: t('gantt.tour.steps.export.content'),
    header: t('gantt.tour.steps.export.header'),
    params: { placement: 'bottom' }
  }
]

/**
 * Advanced tour steps for experienced users
 * Covers power features and shortcuts
 */
export const advancedTourSteps = [
  {
    target: '#gantt-toolbar',
    content: t('gantt.tour.advanced.toolbar.content'),
    header: t('gantt.tour.advanced.toolbar.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="bulk-edit"]',
    content: t('gantt.tour.advanced.bulkEdit.content'),
    header: t('gantt.tour.advanced.bulkEdit.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="undo-redo"]',
    content: t('gantt.tour.advanced.undoRedo.content'),
    header: t('gantt.tour.advanced.undoRedo.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="resource-leveling"]',
    content: t('gantt.tour.advanced.resourceLeveling.content'),
    header: t('gantt.tour.advanced.resourceLeveling.header'),
    params: { placement: 'left' }
  },
  {
    target: '[data-tour="critical-path"]',
    content: t('gantt.tour.advanced.criticalPath.content'),
    header: t('gantt.tour.advanced.criticalPath.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="baselines"]',
    content: t('gantt.tour.advanced.baselines.content'),
    header: t('gantt.tour.advanced.baselines.header'),
    params: { placement: 'left' }
  },
  {
    target: '[data-tour="milestones"]',
    content: t('gantt.tour.advanced.milestones.content'),
    header: t('gantt.tour.advanced.milestones.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="constraints"]',
    content: t('gantt.tour.advanced.constraints.content'),
    header: t('gantt.tour.advanced.constraints.header'),
    params: { placement: 'left' }
  },
  {
    target: '[data-tour="calendar"]',
    content: t('gantt.tour.advanced.calendar.content'),
    header: t('gantt.tour.advanced.calendar.header'),
    params: { placement: 'left' }
  },
  {
    target: '[data-tour="collaboration"]',
    content: t('gantt.tour.advanced.collaboration.content'),
    header: t('gantt.tour.advanced.collaboration.header'),
    params: { placement: 'left' }
  },
  {
    target: '[data-tour="keyboard-shortcuts"]',
    content: t('gantt.tour.advanced.keyboardShortcuts.content'),
    header: t('gantt.tour.advanced.keyboardShortcuts.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="smart-suggestions"]',
    content: t('gantt.tour.advanced.smartSuggestions.content'),
    header: t('gantt.tour.advanced.smartSuggestions.header'),
    params: { placement: 'left' }
  }
]

/**
 * Feature-specific tour steps
 * Covers new features or updates
 */
export const featuresTourSteps = [
  {
    target: '#gantt-toolbar',
    content: t('gantt.tour.features.newFeatures.content'),
    header: t('gantt.tour.features.newFeatures.header'),
    params: { placement: 'bottom' }
  },
  {
    target: '[data-tour="templates"]',
    content: t('gantt.tour.features.templates.content'),
    header: t('gantt.tour.features.templates.header'),
    params: { placement: 'left' }
  },
  {
    target: '[data-tour="ai-suggestions"]',
    content: t('gantt.tour.features.aiSuggestions.content'),
    header: t('gantt.tour.features.aiSuggestions.header'),
    params: { placement: 'left' }
  },
  {
    target: '[data-tour="workflow-automation"]',
    content: t('gantt.tour.features.workflowAutomation.content'),
    header: t('gantt.tour.features.workflowAutomation.header'),
    params: { placement: 'left' }
  },
  {
    target: '[data-tour="minimap"]',
    content: t('gantt.tour.features.minimap.content'),
    header: t('gantt.tour.features.minimap.header'),
    params: { placement: 'right' }
  },
  {
    target: '[data-tour="context-menu"]',
    content: t('gantt.tour.features.contextMenu.content'),
    header: t('gantt.tour.features.contextMenu.header'),
    params: { placement: 'bottom' },
    before: () => {
      // Show context menu instructions
      return Promise.resolve(true)
    }
  }
]

/**
 * Get tour steps by type
 * @param {string} tourType - Type of tour (basic, advanced, features)
 * @returns {Array} Array of tour steps
 */
export function getTourSteps(tourType = 'basic') {
  switch (tourType) {
    case 'basic':
      return basicTourSteps
    case 'advanced':
      return advancedTourSteps
    case 'features':
      return featuresTourSteps
    default:
      return basicTourSteps
  }
}

/**
 * Tour configuration options
 */
export const tourConfig = {
  // Show tour on first visit
  autoStart: true,
  // Show tour steps count
  showSteps: true,
  // Allow skipping tour
  allowSkip: true,
  // Show progress bar
  showProgress: true,
  // Tour highlight color
  highlightColor: '#4CAF50',
  // Tour animation duration
  animationDuration: 300,
  // Remember tour progress
  rememberProgress: true,
  // Storage key for tour progress
  storageKey: 'gantt-tour-progress'
}

/**
 * Tour types with metadata
 */
export const tourTypes = {
  basic: {
    id: 'basic',
    name: t('gantt.tour.types.basic.name'),
    description: t('gantt.tour.types.basic.description'),
    duration: '3-5 min',
    steps: basicTourSteps.length
  },
  advanced: {
    id: 'advanced',
    name: t('gantt.tour.types.advanced.name'),
    description: t('gantt.tour.types.advanced.description'),
    duration: '5-8 min',
    steps: advancedTourSteps.length
  },
  features: {
    id: 'features',
    name: t('gantt.tour.types.features.name'),
    description: t('gantt.tour.types.features.description'),
    duration: '2-3 min',
    steps: featuresTourSteps.length
  }
}

/**
 * Get completed tours from localStorage
 * @returns {Object} Object with completed tour types
 */
export function getCompletedTours() {
  try {
    const stored = localStorage.getItem(tourConfig.storageKey)
    return stored ? JSON.parse(stored) : {}
  } catch (error) {
    console.error('Error reading tour progress:', error)
    return {}
  }
}

/**
 * Mark tour as completed
 * @param {string} tourType - Type of tour to mark as completed
 */
export function markTourCompleted(tourType) {
  try {
    const completed = getCompletedTours()
    completed[tourType] = {
      completed: true,
      completedAt: new Date().toISOString()
    }
    localStorage.setItem(tourConfig.storageKey, JSON.stringify(completed))
  } catch (error) {
    console.error('Error saving tour progress:', error)
  }
}

/**
 * Check if tour is completed
 * @param {string} tourType - Type of tour to check
 * @returns {boolean} True if tour is completed
 */
export function isTourCompleted(tourType) {
  const completed = getCompletedTours()
  return !!completed[tourType]?.completed
}

/**
 * Reset tour progress
 * @param {string} tourType - Type of tour to reset (or 'all' for all tours)
 */
export function resetTourProgress(tourType = 'all') {
  try {
    if (tourType === 'all') {
      localStorage.removeItem(tourConfig.storageKey)
    } else {
      const completed = getCompletedTours()
      delete completed[tourType]
      localStorage.setItem(tourConfig.storageKey, JSON.stringify(completed))
    }
  } catch (error) {
    console.error('Error resetting tour progress:', error)
  }
}
