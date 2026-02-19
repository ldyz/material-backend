/**
 * Calendar State Management Store
 *
 * Manages project calendars, working days, holidays, and exception dates
 * using Pinia for state management.
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { calendarApi } from '@/api'

export const useCalendarStore = defineStore('calendar', () => {
  // State
  const calendars = ref([])
  const projectCalendar = ref(null)
  const taskCalendars = ref(new Map())
  const loading = ref(false)
  const error = ref(null)

  // Standard calendar presets
  const standardPresets = ref([
    {
      id: 'standard',
      name: 'Standard Work Week',
      description: 'Monday-Friday, 8 hours/day',
      workingDays: [1, 2, 3, 4, 5], // Mon-Fri
      workingHours: { start: '09:00', end: '17:00' },
      holidays: []
    },
    {
      id: '24x7',
      name: '24/7 Operations',
      description: 'All days, 24 hours/day',
      workingDays: [0, 1, 2, 3, 4, 5, 6], // All days
      workingHours: { start: '00:00', end: '23:59' },
      holidays: []
    },
    {
      id: 'custom',
      name: 'Custom Calendar',
      description: 'User-defined working days and hours',
      workingDays: [1, 2, 3, 4, 5],
      workingHours: { start: '09:00', end: '17:00' },
      holidays: []
    }
  ])

  // Computed
  const hasCalendars = computed(() => calendars.value.length > 0)

  const getCalendarById = computed(() => {
    return (id) => calendars.value.find(c => c.id === id)
  })

  const getProjectCalendar = computed(() => {
    return projectCalendar.value || standardPresets.value[0]
  })

  const getTaskCalendar = computed(() => {
    return (taskId) => {
      if (taskCalendars.value.has(taskId)) {
        return taskCalendars.value.get(taskId)
      }
      return projectCalendar.value || standardPresets.value[0]
    }
  })

  // Actions
  const fetchCalendars = async (projectId) => {
    loading.value = true
    error.value = null

    try {
      const response = await calendarApi.list(projectId)
      calendars.value = response.data || []

      // Load project calendar
      if (response.projectCalendar) {
        projectCalendar.value = response.projectCalendar
      }

      // Load task calendars
      if (response.taskCalendars) {
        taskCalendars.value = new Map(
          Object.entries(response.taskCalendars).map(([taskId, calendar]) => [
            parseInt(taskId),
            calendar
          ])
        )
      }

      return calendars.value
    } catch (err) {
      console.error('Failed to fetch calendars:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const createCalendar = async (calendarData) => {
    loading.value = true
    error.value = null

    try {
      const response = await calendarApi.create(calendarData)
      calendars.value.push(response.data)

      return response.data
    } catch (err) {
      console.error('Failed to create calendar:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateCalendar = async (calendarId, updates) => {
    loading.value = true
    error.value = null

    try {
      const response = await calendarApi.update(calendarId, updates)

      const index = calendars.value.findIndex(c => c.id === calendarId)
      if (index !== -1) {
        calendars.value[index] = response.data
      }

      return response.data
    } catch (err) {
      console.error('Failed to update calendar:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteCalendar = async (calendarId) => {
    loading.value = true
    error.value = null

    try {
      await calendarApi.delete(calendarId)
      calendars.value = calendars.value.filter(c => c.id !== calendarId)

      return true
    } catch (err) {
      console.error('Failed to delete calendar:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const setProjectCalendar = async (calendarId) => {
    loading.value = true
    error.value = null

    try {
      const response = await calendarApi.setProjectCalendar(calendarId)
      projectCalendar.value = response.data

      return response.data
    } catch (err) {
      console.error('Failed to set project calendar:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const assignTaskCalendar = async (taskId, calendarId) => {
    loading.value = true
    error.value = null

    try {
      const response = await calendarApi.assignTaskCalendar(taskId, calendarId)
      taskCalendars.value.set(taskId, response.data)

      return response.data
    } catch (err) {
      console.error('Failed to assign task calendar:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const removeTaskCalendar = async (taskId) => {
    loading.value = true
    error.value = null

    try {
      await calendarApi.removeTaskCalendar(taskId)
      taskCalendars.value.delete(taskId)

      return true
    } catch (err) {
      console.error('Failed to remove task calendar:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const addHoliday = async (calendarId, holiday) => {
    loading.value = true
    error.value = null

    try {
      const response = await calendarApi.addHoliday(calendarId, holiday)

      const calendar = calendars.value.find(c => c.id === calendarId)
      if (calendar) {
        if (!calendar.holidays) {
          calendar.holidays = []
        }
        calendar.holidays.push(response.data)
      }

      return response.data
    } catch (err) {
      console.error('Failed to add holiday:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const removeHoliday = async (calendarId, holidayId) => {
    loading.value = true
    error.value = null

    try {
      await calendarApi.removeHoliday(calendarId, holidayId)

      const calendar = calendars.value.find(c => c.id === calendarId)
      if (calendar && calendar.holidays) {
        calendar.holidays = calendar.holidays.filter(h => h.id !== holidayId)
      }

      return true
    } catch (err) {
      console.error('Failed to remove holiday:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const addException = async (calendarId, exception) => {
    loading.value = true
    error.value = null

    try {
      const response = await calendarApi.addException(calendarId, exception)

      const calendar = calendars.value.find(c => c.id === calendarId)
      if (calendar) {
        if (!calendar.exceptions) {
          calendar.exceptions = []
        }
        calendar.exceptions.push(response.data)
      }

      return response.data
    } catch (err) {
      console.error('Failed to add exception:', err)
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  // Calendar calculation utilities
  const isWorkingDay = (date, calendarId = null) => {
    const calendar = calendarId
      ? getCalendarById.value(calendarId)
      : projectCalendar.value || standardPresets.value[0]

    if (!calendar) return true

    const dayOfWeek = new Date(date).getDay()

    // Check if it's a working day
    if (!calendar.workingDays.includes(dayOfWeek)) {
      return false
    }

    // Check if it's a holiday
    if (calendar.holidays) {
      const dateStr = new Date(date).toISOString().split('T')[0]
      if (calendar.holidays.some(h => h.date === dateStr)) {
        return false
      }
    }

    // Check exceptions
    if (calendar.exceptions) {
      const dateStr = new Date(date).toISOString().split('T')[0]
      const exception = calendar.exceptions.find(e => e.date === dateStr)
      if (exception) {
        return exception.isWorkingDay
      }
    }

    return true
  }

  const calculateWorkingDays = (startDate, endDate, calendarId = null) => {
    const start = new Date(startDate)
    const end = new Date(endDate)
    let workingDays = 0

    const current = new Date(start)
    while (current <= end) {
      if (isWorkingDay(current, calendarId)) {
        workingDays++
      }
      current.setDate(current.getDate() + 1)
    }

    return workingDays
  }

  const addWorkingDays = (startDate, days, calendarId = null) => {
    const current = new Date(startDate)
    let addedDays = 0

    while (addedDays < days) {
      current.setDate(current.getDate() + 1)
      if (isWorkingDay(current, calendarId)) {
        addedDays++
      }
    }

    return current
  }

  const subtractWorkingDays = (startDate, days, calendarId = null) => {
    const current = new Date(startDate)
    let subtractedDays = 0

    while (subtractedDays < days) {
      current.setDate(current.getDate() - 1)
      if (isWorkingDay(current, calendarId)) {
        subtractedDays++
      }
    }

    return current
  }

  const getWorkingHours = (date, calendarId = null) => {
    const calendar = calendarId
      ? getCalendarById.value(calendarId)
      : projectCalendar.value || standardPresets.value[0]

    if (!calendar || !isWorkingDay(date, calendarId)) {
      return { start: null, end: null, hours: 0 }
    }

    // Check exceptions
    if (calendar.exceptions) {
      const dateStr = new Date(date).toISOString().split('T')[0]
      const exception = calendar.exceptions.find(e => e.date === dateStr)
      if (exception && exception.workingHours) {
        return exception.workingHours
      }
    }

    return calendar.workingHours || { start: '09:00', end: '17:00', hours: 8 }
  }

  const reset = () => {
    calendars.value = []
    projectCalendar.value = null
    taskCalendars.value = new Map()
    loading.value = false
    error.value = null
  }

  return {
    // State
    calendars,
    projectCalendar,
    taskCalendars,
    standardPresets,
    loading,
    error,

    // Computed
    hasCalendars,
    getCalendarById,
    getProjectCalendar,
    getTaskCalendar,

    // Actions
    fetchCalendars,
    createCalendar,
    updateCalendar,
    deleteCalendar,
    setProjectCalendar,
    assignTaskCalendar,
    removeTaskCalendar,
    addHoliday,
    removeHoliday,
    addException,

    // Utilities
    isWorkingDay,
    calculateWorkingDays,
    addWorkingDays,
    subtractWorkingDays,
    getWorkingHours,
    reset
  }
})
