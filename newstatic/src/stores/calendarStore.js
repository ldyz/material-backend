/**
 * Calendar State Management Store
 *
 * Manages project calendars, working days, holidays, and exception dates
 * using Pinia for state management.
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useCalendarStore = defineStore('calendar', () => {
  // State
  const calendars = ref([
    {
      id: 'default',
      name: 'Standard Work Week',
      description: 'Monday-Friday, 8 hours/day',
      workingDays: [1, 2, 3, 4, 5], // Mon-Fri
      workingHours: { start: '08:00', end: '17:00', breakStart: '12:00', breakEnd: '13:00' },
      holidays: []
    }
  ])
  const currentCalendarId = ref('default')
  const loading = ref(false)
  const error = ref(null)

  // Computed
  const hasCalendars = computed(() => calendars.value.length > 0)

  const currentCalendar = computed(() => {
    return calendars.value.find(c => c.id === currentCalendarId.value) || calendars.value[0]
  })

  const getCalendarById = computed(() => {
    return (id) => calendars.value.find(c => c.id === id)
  })

  // Helper function to parse time string
  const parseTime = (timeStr) => {
    const [hours, minutes] = timeStr.split(':').map(Number)
    return hours * 60 + minutes
  }

  // Helper function to format date for holiday checking
  const formatDateKey = (date) => {
    const d = new Date(date)
    const month = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${month}-${day}`
  }

  const formatFullDate = (date) => {
    const d = new Date(date)
    const year = d.getFullYear()
    const month = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  // Holiday Management
  const addHoliday = (holiday) => {
    const calendar = currentCalendar.value
    if (!calendar.holidays) {
      calendar.holidays = []
    }

    // Check for duplicates
    const existing = calendar.holidays.find(h => h.date === holiday.date)
    if (existing) {
      return { success: false, error: 'Holiday already exists' }
    }

    calendar.holidays.push(holiday)
    return { success: true }
  }

  const removeHoliday = (date) => {
    const calendar = currentCalendar.value
    if (!calendar.holidays) return

    calendar.holidays = calendar.holidays.filter(h => h.date !== date)
  }

  const isHoliday = (date) => {
    const calendar = currentCalendar.value
    if (!calendar.holidays || calendar.holidays.length === 0) {
      return false
    }

    const dateStr = typeof date === 'string' ? date : formatFullDate(date)
    const dateKey = typeof date === 'string' ? formatDateKey(date) : formatDateKey(date)

    return calendar.holidays.some(h => {
      if (h.recurring) {
        return h.date === dateKey
      }
      return h.date === dateStr
    })
  }

  const getHolidaysInRange = (startDate, endDate) => {
    const calendar = currentCalendar.value
    if (!calendar.holidays || calendar.holidays.length === 0) {
      return []
    }

    const start = new Date(startDate)
    const end = new Date(endDate)

    return calendar.holidays.filter(holiday => {
      const holidayDate = new Date(holiday.date)
      return holidayDate >= start && holidayDate <= end
    })
  }

  // Working Days Configuration
  const setWorkingDays = (days) => {
    const calendar = currentCalendar.value
    calendar.workingDays = days
  }

  const setWorkingHours = (hours) => {
    const calendar = currentCalendar.value
    calendar.workingHours = hours
  }

  const isWorkingDay = (date) => {
    const calendar = currentCalendar.value
    const d = new Date(date)
    const dayOfWeek = d.getDay()

    // Check if it's a working day
    if (!calendar.workingDays.includes(dayOfWeek)) {
      return false
    }

    // Check if it's a holiday
    return !isHoliday(date)
  }

  const getDailyWorkingHours = () => {
    const calendar = currentCalendar.value
    const { start, end, breakStart, breakEnd } = calendar.workingHours

    const startMinutes = parseTime(start)
    const endMinutes = parseTime(end)
    let dailyHours = (endMinutes - startMinutes) / 60

    // Subtract break time if defined
    if (breakStart && breakEnd) {
      const breakStartMinutes = parseTime(breakStart)
      const breakEndMinutes = parseTime(breakEnd)
      dailyHours -= (breakEndMinutes - breakStartMinutes) / 60
    }

    // Round to handle floating point precision issues
    return Math.round(dailyHours * 100) / 100
  }

  // Working Time Calculations
  const calculateWorkingDays = (startDate, endDate) => {
    const start = new Date(startDate)
    const end = new Date(endDate)
    let workingDays = 0

    const current = new Date(start)
    while (current <= end) {
      if (isWorkingDay(current)) {
        workingDays++
      }
      current.setDate(current.getDate() + 1)
    }

    return workingDays
  }

  const calculateWorkingHours = (startDate, endDate) => {
    const start = new Date(startDate)
    const end = new Date(endDate)
    const calendar = currentCalendar.value
    let totalHours = 0

    const { start: workStart, end: workEnd, breakStart, breakEnd } = calendar.workingHours
    const workStartMinutes = parseTime(workStart)
    const workEndMinutes = parseTime(workEnd)

    // Calculate break minutes if break is defined
    let breakMinutes = 0
    if (breakStart && breakEnd) {
      breakMinutes = parseTime(breakEnd) - parseTime(breakStart)
    }

    // If same day
    if (start.toDateString() === end.toDateString()) {
      if (!isWorkingDay(start)) {
        return 0
      }

      const startMinutes = start.getHours() * 60 + start.getMinutes()
      const endMinutes = end.getHours() * 60 + end.getMinutes()

      const effectiveStart = Math.max(startMinutes, workStartMinutes)
      const effectiveEnd = Math.min(endMinutes, workEndMinutes)

      if (effectiveEnd > effectiveStart) {
        let hours = (effectiveEnd - effectiveStart) / 60

        // Subtract break time if the range includes the break period
        if (breakStart && breakEnd) {
          const breakStartMinutes = parseTime(breakStart)
          const breakEndMinutes = parseTime(breakEnd)

          // Check if the work period overlaps with break
          const periodOverlapsBreak = effectiveStart < breakEndMinutes && effectiveEnd > breakStartMinutes
          if (periodOverlapsBreak) {
            const breakOverlapStart = Math.max(effectiveStart, breakStartMinutes)
            const breakOverlapEnd = Math.min(effectiveEnd, breakEndMinutes)
            const breakOverlapMinutes = Math.max(0, breakOverlapEnd - breakOverlapStart)
            hours -= breakOverlapMinutes / 60
          }
        }

        totalHours = hours
      }
    } else {
      // Multi-day calculation
      const current = new Date(start)
      while (current <= end) {
        if (isWorkingDay(current)) {
          if (current.toDateString() === start.toDateString()) {
            const startMinutes = start.getHours() * 60 + start.getMinutes()
            const effectiveStart = Math.max(startMinutes, workStartMinutes)
            const dayHours = (workEndMinutes - effectiveStart) / 60

            // Subtract break if first day includes it
            if (breakStart && breakEnd && effectiveStart < parseTime(breakEnd)) {
              const breakStartMinutes = parseTime(breakStart)
              const breakEndMinutes = parseTime(breakEnd)
              if (effectiveStart < breakEndMinutes) {
                const breakOverlapStart = Math.max(effectiveStart, breakStartMinutes)
                const breakOverlapEnd = breakEndMinutes
                const breakOverlapMinutes = Math.max(0, breakOverlapEnd - breakOverlapStart)
                totalHours -= breakOverlapMinutes / 60
              }
            }

            totalHours += dayHours
          } else if (current.toDateString() === end.toDateString()) {
            const endMinutes = end.getHours() * 60 + end.getMinutes()
            const effectiveEnd = Math.min(endMinutes, workEndMinutes)
            const dayHours = (effectiveEnd - workStartMinutes) / 60

            // Subtract break if last day includes it
            if (breakStart && breakEnd && effectiveEnd > parseTime(breakStart)) {
              const breakStartMinutes = parseTime(breakStart)
              const breakEndMinutes = parseTime(breakEnd)
              if (effectiveEnd > breakStartMinutes) {
                const breakOverlapStart = breakStartMinutes
                const breakOverlapEnd = Math.min(effectiveEnd, breakEndMinutes)
                const breakOverlapMinutes = Math.max(0, breakOverlapEnd - breakOverlapStart)
                totalHours -= breakOverlapMinutes / 60
              }
            }

            totalHours += dayHours
          } else {
            // Full day - subtract break if defined
            totalHours += (workEndMinutes - workStartMinutes - breakMinutes) / 60
          }
        }
        current.setDate(current.getDate() + 1)
      }
    }

    return Math.max(0, Math.round(totalHours * 100) / 100)
  }

  const addWorkingDays = (startDate, days) => {
    const current = new Date(startDate)
    let addedDays = 0
    const direction = days >= 0 ? 1 : -1
    const absDays = Math.abs(days)

    current.setDate(current.getDate() + direction)

    while (addedDays < absDays) {
      if (isWorkingDay(current)) {
        addedDays++
      }
      if (addedDays < absDays) {
        current.setDate(current.getDate() + direction)
      }
    }

    return current
  }

  // Date Utilities
  const getNextWorkingDay = (date) => {
    const current = new Date(date)
    current.setDate(current.getDate() + 1)

    while (!isWorkingDay(current)) {
      current.setDate(current.getDate() + 1)
    }

    return current
  }

  const getPreviousWorkingDay = (date) => {
    const current = new Date(date)
    current.setDate(current.getDate() - 1)

    while (!isWorkingDay(current)) {
      current.setDate(current.getDate() - 1)
    }

    return current
  }

  const adjustToWorkingDay = (date) => {
    if (isWorkingDay(date)) {
      return new Date(date)
    }

    // Move forward to next working day
    return getNextWorkingDay(date)
  }

  const calculateDuration = (startDate, endDate) => {
    return calculateWorkingDays(startDate, endDate)
  }

  const calculateEndDate = (startDate, workingDays) => {
    return addWorkingDays(startDate, workingDays)
  }

  // Multiple Calendars Management
  const createCalendar = (calendarData) => {
    const id = `cal-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    const newCalendar = {
      id,
      name: calendarData.name || 'New Calendar',
      description: calendarData.description || '',
      workingDays: calendarData.workingDays || [1, 2, 3, 4, 5],
      workingHours: calendarData.workingHours || { start: '08:00', end: '17:00' },
      holidays: calendarData.holidays || []
    }

    calendars.value.push(newCalendar)
    return newCalendar
  }

  const setCurrentCalendar = (calendarId) => {
    const calendar = calendars.value.find(c => c.id === calendarId)
    if (calendar) {
      currentCalendarId.value = calendarId
    }
  }

  const deleteCalendar = (calendarId) => {
    if (calendarId === 'default') {
      throw new Error('Cannot delete default calendar')
    }

    const index = calendars.value.findIndex(c => c.id === calendarId)
    if (index !== -1) {
      calendars.value.splice(index, 1)

      // If we deleted the current calendar, switch to default
      if (currentCalendarId.value === calendarId) {
        currentCalendarId.value = 'default'
      }
    }
  }

  const getCalendar = (calendarId) => {
    return calendars.value.find(c => c.id === calendarId)
  }

  // Calendar Validation
  const validateCalendar = (config) => {
    const errors = []

    // Validate working days
    if (!config.workingDays || !Array.isArray(config.workingDays)) {
      errors.push('Working days must be an array')
    } else {
      const invalidDays = config.workingDays.filter(d => d < 0 || d > 7)
      if (invalidDays.length > 0) {
        errors.push(`Invalid working day values: ${invalidDays.join(', ')}`)
      }
    }

    // Validate working hours
    if (!config.workingHours) {
      errors.push('Working hours configuration is required')
    } else {
      const { start, end } = config.workingHours

      // Validate time format
      const timeRegex = /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/
      if (!timeRegex.test(start)) {
        errors.push('Invalid start time format (use HH:MM)')
      }
      if (!timeRegex.test(end)) {
        errors.push('Invalid end time format (use HH:MM)')
      }

      // Validate time range
      if (start && end && timeRegex.test(start) && timeRegex.test(end)) {
        const startMinutes = parseTime(start)
        const endMinutes = parseTime(end)

        if (endMinutes <= startMinutes) {
          errors.push('End time must be after start time')
        }
      }
    }

    return {
      valid: errors.length === 0,
      errors
    }
  }

  // Export/Import
  const exportCalendar = () => {
    const calendar = currentCalendar.value
    return JSON.stringify({
      name: calendar.name,
      description: calendar.description,
      workingDays: calendar.workingDays,
      workingHours: calendar.workingHours,
      holidays: calendar.holidays
    }, null, 2)
  }

  const importCalendar = (jsonData) => {
    let data
    try {
      data = JSON.parse(jsonData)
    } catch (err) {
      throw new Error('Invalid JSON format')
    }

    const validation = validateCalendar(data)
    if (!validation.valid) {
      throw new Error(`Invalid calendar configuration: ${validation.errors.join(', ')}`)
    }

    const calendar = currentCalendar.value
    calendar.name = data.name || calendar.name
    calendar.description = data.description || calendar.description
    calendar.workingDays = data.workingDays || calendar.workingDays
    calendar.workingHours = data.workingHours || calendar.workingHours
    calendar.holidays = data.holidays || []
  }

  // Reset
  const reset = () => {
    calendars.value = [
      {
        id: 'default',
        name: 'Standard Work Week',
        description: 'Monday-Friday, 8 hours/day',
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: '08:00', end: '17:00', breakStart: '12:00', breakEnd: '13:00' },
        holidays: []
      }
    ]
    currentCalendarId.value = 'default'
    loading.value = false
    error.value = null
  }

  return {
    // State
    calendars,
    currentCalendarId,
    loading,
    error,

    // Computed
    hasCalendars,
    currentCalendar,
    getCalendarById,

    // Holiday Management
    addHoliday,
    removeHoliday,
    isHoliday,
    getHolidaysInRange,

    // Working Days Configuration
    setWorkingDays,
    setWorkingHours,
    isWorkingDay,
    getDailyWorkingHours,

    // Working Time Calculations
    calculateWorkingDays,
    calculateWorkingHours,
    addWorkingDays,

    // Date Utilities
    getNextWorkingDay,
    getPreviousWorkingDay,
    adjustToWorkingDay,
    calculateDuration,
    calculateEndDate,

    // Multiple Calendars
    createCalendar,
    setCurrentCalendar,
    deleteCalendar,
    getCalendar,

    // Validation
    validateCalendar,

    // Export/Import
    exportCalendar,
    importCalendar,

    // Reset
    reset
  }
})
