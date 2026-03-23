/**
 * Template Store for Gantt Chart
 * Manages project templates with Pinia
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

export const useTemplateStore = defineStore('template', () => {
  // State
  const templates = ref([])
  const categories = ref([
    { id: 'software', name: 'Software Development', icon: 'code' },
    { id: 'construction', name: 'Construction', icon: 'building' },
    { id: 'marketing', name: 'Marketing Campaign', icon: 'megaphone' },
    { id: 'event', name: 'Event Planning', icon: 'calendar' },
    { id: 'custom', name: 'Custom', icon: 'document' }
  ])
  const selectedTemplate = ref(null)
  const recentTemplates = ref([])
  const loading = ref(false)
  const searchQuery = ref('')
  const selectedCategory = ref('all')

  // Storage key
  const STORAGE_KEY = 'gantt-templates'
  const RECENT_KEY = 'gantt-recent-templates'

  // Computed
  const filteredTemplates = computed(() => {
    let filtered = templates.value

    // Filter by category
    if (selectedCategory.value !== 'all') {
      filtered = filtered.filter(t => t.category === selectedCategory.value)
    }

    // Filter by search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      filtered = filtered.filter(t =>
        t.name.toLowerCase().includes(query) ||
        t.description?.toLowerCase().includes(query)
      )
    }

    return filtered
  })

  const templatesByCategory = computed(() => {
    const grouped = {}
    categories.value.forEach(cat => {
      grouped[cat.id] = templates.value.filter(t => t.category === cat.id)
    })
    return grouped
  })

  /**
   * Load templates from localStorage
   */
  function loadTemplates() {
    try {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) {
        templates.value = JSON.parse(stored)
      }

      const recent = localStorage.getItem(RECENT_KEY)
      if (recent) {
        recentTemplates.value = JSON.parse(recent)
      }
    } catch (error) {
      console.error('Error loading templates:', error)
      ElMessage.error('Failed to load templates')
    }
  }

  /**
   * Save templates to localStorage
   */
  function saveTemplates() {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(templates.value))
    } catch (error) {
      console.error('Error saving templates:', error)
      ElMessage.error('Failed to save templates')
    }
  }

  /**
   * Add a new template
   */
  function addTemplate(template) {
    const newTemplate = {
      id: `template-${Date.now()}`,
      name: template.name,
      description: template.description || '',
      category: template.category || 'custom',
      tasks: template.tasks || [],
      dependencies: template.dependencies || [],
      resources: template.resources || [],
      settings: template.settings || {},
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      ...template
    }

    templates.value.push(newTemplate)
    saveTemplates()
    ElMessage.success('Template created successfully')

    return newTemplate
  }

  /**
   * Update an existing template
   */
  function updateTemplate(id, updates) {
    const index = templates.value.findIndex(t => t.id === id)
    if (index === -1) {
      ElMessage.error('Template not found')
      return null
    }

    templates.value[index] = {
      ...templates.value[index],
      ...updates,
      updatedAt: new Date().toISOString()
    }

    saveTemplates()
    ElMessage.success('Template updated successfully')

    return templates.value[index]
  }

  /**
   * Delete a template
   */
  function deleteTemplate(id) {
    const index = templates.value.findIndex(t => t.id === id)
    if (index === -1) {
      ElMessage.error('Template not found')
      return false
    }

    templates.value.splice(index, 1)
    saveTemplates()
    ElMessage.success('Template deleted successfully')

    return true
  }

  /**
   * Duplicate a template
   */
  function duplicateTemplate(id) {
    const template = templates.value.find(t => t.id === id)
    if (!template) {
      ElMessage.error('Template not found')
      return null
    }

    const duplicate = {
      ...template,
      id: `template-${Date.now()}`,
      name: `${template.name} (Copy)`,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }

    templates.value.push(duplicate)
    saveTemplates()
    ElMessage.success('Template duplicated successfully')

    return duplicate
  }

  /**
   * Select a template
   */
  function selectTemplate(id) {
    const template = templates.value.find(t => t.id === id)
    if (!template) {
      ElMessage.error('Template not found')
      return
    }

    selectedTemplate.value = template
    addToRecent(template)
  }

  /**
   * Add template to recent
   */
  function addToRecent(template) {
    // Remove if already exists
    const index = recentTemplates.value.findIndex(t => t.id === template.id)
    if (index !== -1) {
      recentTemplates.value.splice(index, 1)
    }

    // Add to beginning
    recentTemplates.value.unshift(template)

    // Keep only 5 recent
    if (recentTemplates.value.length > 5) {
      recentTemplates.value = recentTemplates.value.slice(0, 5)
    }

    // Save to localStorage
    try {
      localStorage.setItem(RECENT_KEY, JSON.stringify(recentTemplates.value))
    } catch (error) {
      console.error('Error saving recent templates:', error)
    }
  }

  /**
   * Get template by ID
   */
  function getTemplate(id) {
    return templates.value.find(t => t.id === id)
  }

  /**
   * Create template from current project
   */
  function createFromProject(projectData, options = {}) {
    const template = {
      name: options.name || 'Untitled Template',
      description: options.description || '',
      category: options.category || 'custom',
      tasks: projectData.tasks || [],
      dependencies: projectData.dependencies || [],
      resources: projectData.resources || [],
      settings: {
        viewMode: projectData.viewMode,
        dayWidth: projectData.dayWidth,
        ...options.settings
      }
    }

    return addTemplate(template)
  }

  /**
   * Apply template to project
   */
  function applyTemplateToProject(templateId, projectData) {
    const template = getTemplate(templateId)
    if (!template) {
      ElMessage.error('Template not found')
      return null
    }

    // Create a copy of template data
    const appliedData = {
      tasks: template.tasks.map(task => ({
        ...task,
        id: `task-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
      })),
      dependencies: template.dependencies.map(dep => ({
        ...dep,
        id: `dep-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
      })),
      resources: template.resources || [],
      settings: template.settings || {}
    }

    // Update task dependencies to use new task IDs
    const taskIdMap = new Map()
    template.tasks.forEach((oldTask, index) => {
      taskIdMap.set(oldTask.id, appliedData.tasks[index].id)
    })

    appliedData.dependencies = appliedData.dependencies.map(dep => ({
      ...dep,
      predecessorId: taskIdMap.get(dep.predecessorId) || dep.predecessorId,
      successorId: taskIdMap.get(dep.successorId) || dep.successorId
    }))

    addToRecent(template)
    ElMessage.success(`Template "${template.name}" applied successfully`)

    return appliedData
  }

  /**
   * Validate template data
   */
  function validateTemplate(template) {
    const errors = []

    if (!template.name || template.name.trim() === '') {
      errors.push('Template name is required')
    }

    if (!template.category || !categories.value.find(c => c.id === template.category)) {
      errors.push('Valid category is required')
    }

    if (!template.tasks || !Array.isArray(template.tasks) || template.tasks.length === 0) {
      errors.push('Template must have at least one task')
    }

    if (template.dependencies) {
      template.dependencies.forEach((dep, index) => {
        if (!dep.predecessorId || !dep.successorId) {
          errors.push(`Dependency ${index + 1} is missing predecessor or successor`)
        }
      })
    }

    return {
      valid: errors.length === 0,
      errors
    }
  }

  /**
   * Export template to JSON
   */
  function exportTemplate(id) {
    const template = getTemplate(id)
    if (!template) {
      ElMessage.error('Template not found')
      return null
    }

    const dataStr = JSON.stringify(template, null, 2)
    const dataBlob = new Blob([dataStr], { type: 'application/json' })
    const url = URL.createObjectURL(dataBlob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${template.name.replace(/\s+/g, '-').toLowerCase()}.json`
    link.click()
    URL.revokeObjectURL(url)

    ElMessage.success('Template exported successfully')
  }

  /**
   * Import template from JSON
   */
  function importTemplate(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()

      reader.onload = (e) => {
        try {
          const template = JSON.parse(e.target.result)
          const validation = validateTemplate(template)

          if (!validation.valid) {
            ElMessage.error(`Invalid template: ${validation.errors.join(', ')}`)
            reject(new Error(validation.errors.join(', ')))
            return
          }

          // Generate new ID to avoid conflicts
          template.id = `template-${Date.now()}`
          const added = addTemplate(template)
          resolve(added)
        } catch (error) {
          ElMessage.error('Failed to parse template file')
          reject(error)
        }
      }

      reader.onerror = () => {
        ElMessage.error('Failed to read template file')
        reject(new Error('Failed to read file'))
      }

      reader.readAsText(file)
    })
  }

  /**
   * Clear all templates
   */
  function clearAllTemplates() {
    templates.value = []
    recentTemplates.value = []
    localStorage.removeItem(STORAGE_KEY)
    localStorage.removeItem(RECENT_KEY)
    ElMessage.success('All templates cleared')
  }

  /**
   * Initialize store
   */
  function initialize() {
    loadTemplates()

    // Add default templates if none exist
    if (templates.value.length === 0) {
      addDefaultTemplates()
    }
  }

  /**
   * Add default templates
   */
  function addDefaultTemplates() {
    const defaultTemplates = [
      {
        name: 'Software Development',
        description: 'Standard software development lifecycle',
        category: 'software',
        tasks: [
          { id: 'task-1', name: 'Requirements Gathering', duration: 5, start: 0 },
          { id: 'task-2', name: 'System Design', duration: 7, start: 5 },
          { id: 'task-3', name: 'Development', duration: 14, start: 12 },
          { id: 'task-4', name: 'Testing', duration: 7, start: 26 },
          { id: 'task-5', name: 'Deployment', duration: 3, start: 33 }
        ],
        dependencies: [
          { id: 'dep-1', predecessorId: 'task-1', successorId: 'task-2', type: 'finish-to-start' },
          { id: 'dep-2', predecessorId: 'task-2', successorId: 'task-3', type: 'finish-to-start' },
          { id: 'dep-3', predecessorId: 'task-3', successorId: 'task-4', type: 'finish-to-start' },
          { id: 'dep-4', predecessorId: 'task-4', successorId: 'task-5', type: 'finish-to-start' }
        ]
      },
      {
        name: 'Marketing Campaign',
        description: 'Complete marketing campaign workflow',
        category: 'marketing',
        tasks: [
          { id: 'task-1', name: 'Market Research', duration: 5, start: 0 },
          { id: 'task-2', name: 'Strategy Development', duration: 4, start: 5 },
          { id: 'task-3', name: 'Content Creation', duration: 10, start: 9 },
          { id: 'task-4', name: 'Campaign Launch', duration: 1, start: 19 },
          { id: 'task-5', name: 'Monitoring & Optimization', duration: 14, start: 20 }
        ],
        dependencies: [
          { id: 'dep-1', predecessorId: 'task-1', successorId: 'task-2', type: 'finish-to-start' },
          { id: 'dep-2', predecessorId: 'task-2', successorId: 'task-3', type: 'finish-to-start' },
          { id: 'dep-3', predecessorId: 'task-3', successorId: 'task-4', type: 'finish-to-start' },
          { id: 'dep-4', predecessorId: 'task-4', successorId: 'task-5', type: 'finish-to-start' }
        ]
      }
    ]

    defaultTemplates.forEach(t => {
      addTemplate(t)
    })
  }

  return {
    // State
    templates,
    categories,
    selectedTemplate,
    recentTemplates,
    loading,
    searchQuery,
    selectedCategory,

    // Computed
    filteredTemplates,
    templatesByCategory,

    // Actions
    loadTemplates,
    saveTemplates,
    addTemplate,
    updateTemplate,
    deleteTemplate,
    duplicateTemplate,
    selectTemplate,
    getTemplate,
    createFromProject,
    applyTemplateToProject,
    validateTemplate,
    exportTemplate,
    importTemplate,
    clearAllTemplates,
    initialize
  }
})
