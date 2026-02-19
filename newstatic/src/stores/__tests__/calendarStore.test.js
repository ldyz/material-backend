import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useCalendarStore } from '../calendarStore.js'

describe('Calendar Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('Calendar Initialization', () => {
    it('should initialize with default calendar', () => {
      const store = useCalendarStore()

      expect(store.calendars).toBeDefined()
      expect(store.currentCalendar).toBeDefined()
    })

    it('should create standard working calendar', () => {
      const store = useCalendarStore()

      expect(store.currentCalendar.workingDays).toEqual([1, 2, 3, 4, 5])
      expect(store.currentCalendar.workingHours.start).toBe('08:00')
      expect(store.currentCalendar.workingHours.end).toBe('17:00')
    })

    it('should initialize with empty holidays', () => {
      const store = useCalendarStore()

      expect(store.currentCalendar.holidays).toEqual([])
    })
  })

  describe('Working Time Calculations', () => {
    it('should calculate working days between dates', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01T08:00:00Z')
      const end = new Date('2024-01-05T17:00:00Z')

      const workingDays = store.calculateWorkingDays(start, end)

      expect(workingDays).toBe(5)
    })

    it('should exclude weekends from working days', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01T08:00:00Z') // Monday
      const end = new Date('2024-01-07T17:00:00Z') // Sunday

      const workingDays = store.calculateWorkingDays(start, end)

      expect(workingDays).toBe(5) // Mon-Fri only
    })

    it('should calculate working hours between dates', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01T08:00:00Z')
      const end = new Date('2024-01-01T17:00:00Z')

      const workingHours = store.calculateWorkingHours(start, end)

      expect(workingHours).toBe(8)
    })

    it('should exclude holidays from working days', () => {
      const store = useCalendarStore()

      store.addHoliday({
        date: '2024-01-03',
        name: 'Test Holiday',
      })

      const start = new Date('2024-01-01T08:00:00Z')
      const end = new Date('2024-01-05T17:00:00Z')

      const workingDays = store.calculateWorkingDays(start, end)

      expect(workingDays).toBe(4) // 5 days minus 1 holiday
    })

    it('should add working days to a date', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01T08:00:00Z')
      const result = store.addWorkingDays(start, 5)

      const expected = new Date('2024-01-08T08:00:00Z') // Monday + 5 work days = next Monday
      expect(result.toDateString()).toBe(expected.toDateString())
    })

    it('should subtract working days from a date', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-08T08:00:00Z')
      const result = store.addWorkingDays(start, -5)

      const expected = new Date('2024-01-01T08:00:00Z')
      expect(result.toDateString()).toBe(expected.toDateString())
    })

    it('should handle partial working days', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01T10:00:00Z')
      const end = new Date('2024-01-01T17:00:00Z')

      const workingHours = store.calculateWorkingHours(start, end)

      expect(workingHours).toBe(7) // 10am to 5pm
    })

    it('should handle non-working hours', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01T19:00:00Z')
      const end = new Date('2024-01-02T08:00:00Z')

      const workingHours = store.calculateWorkingHours(start, end)

      expect(workingHours).toBe(0) // Outside working hours
    })
  })

  describe('Holiday Management', () => {
    it('should add holiday', () => {
      const store = useCalendarStore()

      store.addHoliday({
        date: '2024-12-25',
        name: 'Christmas Day',
      })

      expect(store.currentCalendar.holidays).toContainEqual({
        date: '2024-12-25',
        name: 'Christmas Day',
      })
    })

    it('should add recurring holiday', () => {
      const store = useCalendarStore()

      store.addHoliday({
        date: '12-25',
        name: 'Christmas',
        recurring: true,
      })

      const holiday = store.currentCalendar.holidays[0]
      expect(holiday.recurring).toBe(true)
    })

    it('should remove holiday', () => {
      const store = useCalendarStore()

      store.addHoliday({
        date: '2024-12-25',
        name: 'Christmas Day',
      })

      store.removeHoliday('2024-12-25')

      expect(store.currentCalendar.holidays).not.toContainEqual(
        expect.objectContaining({ date: '2024-12-25' })
      )
    })

    it('should check if date is a holiday', () => {
      const store = useCalendarStore()

      store.addHoliday({
        date: '2024-12-25',
        name: 'Christmas Day',
      })

      const isHoliday = store.isHoliday('2024-12-25')

      expect(isHoliday).toBe(true)
    })

    it('should handle recurring holidays', () => {
      const store = useCalendarStore()

      store.addHoliday({
        date: '12-25',
        name: 'Christmas',
        recurring: true,
      })

      expect(store.isHoliday('2024-12-25')).toBe(true)
      expect(store.isHoliday('2025-12-25')).toBe(true)
    })

    it('should get holidays for date range', () => {
      const store = useCalendarStore()

      store.addHoliday({ date: '2024-01-01', name: 'New Year' })
      store.addHoliday({ date: '2024-01-15', name: 'MLK Day' })
      store.addHoliday({ date: '2024-02-01', name: 'Test Holiday' })

      const holidays = store.getHolidaysInRange(
        new Date('2024-01-01'),
        new Date('2024-01-31')
      )

      expect(holidays).toHaveLength(2)
      expect(holidays[0].name).toBe('New Year')
      expect(holidays[1].name).toBe('MLK Day')
    })
  })

  describe('Working Days Configuration', () => {
    it('should set working days', () => {
      const store = useCalendarStore()

      store.setWorkingDays([1, 2, 3, 4, 5, 6])

      expect(store.currentCalendar.workingDays).toEqual([1, 2, 3, 4, 5, 6])
    })

    it('should validate working days', () => {
      const store = useCalendarStore()

      store.setWorkingDays([1, 2, 3, 4, 5, 6, 7])

      expect(store.currentCalendar.workingDays).toEqual([1, 2, 3, 4, 5, 6, 7])
    })

    it('should check if day is working day', () => {
      const store = useCalendarStore()

      const monday = new Date('2024-01-01') // Monday
      const saturday = new Date('2024-01-06') // Saturday

      expect(store.isWorkingDay(monday)).toBe(true)
      expect(store.isWorkingDay(saturday)).toBe(false)
    })

    it('should handle custom working week', () => {
      const store = useCalendarStore()

      store.setWorkingDays([1, 3, 5]) // Mon, Wed, Fri only

      const monday = new Date('2024-01-01')
      const tuesday = new Date('2024-01-02')
      const wednesday = new Date('2024-01-03')

      expect(store.isWorkingDay(monday)).toBe(true)
      expect(store.isWorkingDay(tuesday)).toBe(false)
      expect(store.isWorkingDay(wednesday)).toBe(true)
    })
  })

  describe('Working Hours Configuration', () => {
    it('should set working hours', () => {
      const store = useCalendarStore()

      store.setWorkingHours({ start: '09:00', end: '18:00' })

      expect(store.currentCalendar.workingHours.start).toBe('09:00')
      expect(store.currentCalendar.workingHours.end).toBe('18:00')
    })

    it('should validate working hours format', () => {
      const store = useCalendarStore()

      expect(() => store.setWorkingHours({ start: '09:00', end: '18:00' })).not.toThrow()
    })

    it('should calculate daily working hours', () => {
      const store = useCalendarStore()

      store.setWorkingHours({ start: '09:00', end: '18:00' })

      const hours = store.getDailyWorkingHours()

      expect(hours).toBe(9)
    })

    it('should handle lunch break', () => {
      const store = useCalendarStore()

      store.setWorkingHours({
        start: '09:00',
        end: '18:00',
        breakStart: '12:00',
        breakEnd: '13:00',
      })

      const hours = store.getDailyWorkingHours()

      expect(hours).toBe(8) // 9 hours minus 1 hour break
    })
  })

  describe('Multiple Calendars', () => {
    it('should create new calendar', () => {
      const store = useCalendarStore()

      const calendar = store.createCalendar({
        name: 'Night Shift',
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: '22:00', end: '06:00' },
      })

      expect(store.calendars).toContainEqual(
        expect.objectContaining({ name: 'Night Shift' })
      )
    })

    it('should switch between calendars', () => {
      const store = useCalendarStore()

      const calendar = store.createCalendar({
        name: 'Custom Calendar',
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: '10:00', end: '19:00' },
      })

      store.setCurrentCalendar(calendar.id)

      expect(store.currentCalendar.id).toBe(calendar.id)
      expect(store.currentCalendar.workingHours.start).toBe('10:00')
    })

    it('should delete calendar', () => {
      const store = useCalendarStore()

      const calendar = store.createCalendar({
        name: 'Temporary Calendar',
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: '08:00', end: '17:00' },
      })

      store.deleteCalendar(calendar.id)

      expect(store.calendars).not.toContainEqual(
        expect.objectContaining({ id: calendar.id })
      )
    })

    it('should not delete default calendar', () => {
      const store = useCalendarStore()

      const defaultId = store.currentCalendar.id

      expect(() => store.deleteCalendar(defaultId)).toThrow()
    })

    it('should get calendar by ID', () => {
      const store = useCalendarStore()

      const calendar = store.createCalendar({
        name: 'Test Calendar',
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: '08:00', end: '17:00' },
      })

      const retrieved = store.getCalendar(calendar.id)

      expect(retrieved).toEqual(calendar)
    })
  })

  describe('Calendar Validation', () => {
    it('should validate calendar configuration', () => {
      const store = useCalendarStore()

      const result = store.validateCalendar({
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: '08:00', end: '17:00' },
      })

      expect(result.valid).toBe(true)
    })

    it('should detect invalid working days', () => {
      const store = useCalendarStore()

      const result = store.validateCalendar({
        workingDays: [8, 9], // Invalid days
        workingHours: { start: '08:00', end: '17:00' },
      })

      expect(result.valid).toBe(false)
      expect(result.errors).toContainEqual(expect.stringContaining('Invalid working day'))
    })

    it('should detect invalid time format', () => {
      const store = useCalendarStore()

      const result = store.validateCalendar({
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: 'invalid', end: '17:00' },
      })

      expect(result.valid).toBe(false)
    })

    it('should detect end time before start time', () => {
      const store = useCalendarStore()

      const result = store.validateCalendar({
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: '17:00', end: '08:00' },
      })

      expect(result.valid).toBe(false)
    })

    it('should detect duplicate holidays', () => {
      const store = useCalendarStore()

      store.addHoliday({ date: '2024-12-25', name: 'Christmas' })
      const result = store.addHoliday({ date: '2024-12-25', name: 'Xmas' })

      expect(result.success).toBe(false)
      expect(result.error).toContain('already exists')
    })
  })

  describe('Date Utilities', () => {
    it('should get next working day', () => {
      const store = useCalendarStore()

      const friday = new Date('2024-01-05')
      const nextWorkingDay = store.getNextWorkingDay(friday)

      expect(nextWorkingDay.getDay()).toBe(1) // Monday
    })

    it('should get previous working day', () => {
      const store = useCalendarStore()

      const monday = new Date('2024-01-08')
      const prevWorkingDay = store.getPreviousWorkingDay(monday)

      expect(prevWorkingDay.getDay()).toBe(5) // Friday
    })

    it('should adjust date to working day if needed', () => {
      const store = useCalendarStore()

      const saturday = new Date('2024-01-06')
      const adjusted = store.adjustToWorkingDay(saturday)

      expect(adjusted.getDay()).toBe(1) // Should move to Monday
    })

    it('should calculate duration in working days', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01T08:00:00Z')
      const end = new Date('2024-01-10T17:00:00Z')

      const duration = store.calculateDuration(start, end)

      expect(duration).toBe(8) // 8 working days
    })

    it('should calculate end date from duration', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01T08:00:00Z')
      const end = store.calculateEndDate(start, 5)

      expect(end.toDateString()).toBe(new Date('2024-01-08').toDateString())
    })
  })

  describe('Calendar Persistence', () => {
    it('should export calendar to JSON', () => {
      const store = useCalendarStore()

      store.addHoliday({ date: '2024-12-25', name: 'Christmas' })

      const exported = store.exportCalendar()

      expect(exported).toBeDefined()
      expect(JSON.parse(exported)).toEqual(
        expect.objectContaining({
          workingDays: expect.any(Array),
          workingHours: expect.any(Object),
          holidays: expect.any(Array),
        })
      )
    })

    it('should import calendar from JSON', () => {
      const store = useCalendarStore()

      const calendarData = {
        name: 'Imported Calendar',
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: '09:00', end: '18:00' },
        holidays: [{ date: '2024-12-25', name: 'Christmas' }],
      }

      store.importCalendar(JSON.stringify(calendarData))

      expect(store.currentCalendar.workingHours.start).toBe('09:00')
      expect(store.currentCalendar.holidays).toHaveLength(1)
    })

    it('should validate imported calendar', () => {
      const store = useCalendarStore()

      const invalidData = {
        workingDays: [1, 2, 3, 4, 5],
        workingHours: { start: 'invalid', end: '17:00' },
      }

      expect(() => store.importCalendar(JSON.stringify(invalidData))).toThrow()
    })
  })

  describe('Edge Cases', () => {
    it('should handle empty working days', () => {
      const store = useCalendarStore()

      store.setWorkingDays([])

      const date = new Date('2024-01-01')
      expect(store.isWorkingDay(date)).toBe(false)
    })

    it('should handle all week as working days', () => {
      const store = useCalendarStore()

      store.setWorkingDays([1, 2, 3, 4, 5, 6, 7])

      const saturday = new Date('2024-01-06')
      expect(store.isWorkingDay(saturday)).toBe(true)
    })

    it('should handle 24-hour working day', () => {
      const store = useCalendarStore()

      store.setWorkingHours({ start: '00:00', end: '23:59' })

      const hours = store.getDailyWorkingHours()
      expect(hours).toBe(24)
    })

    it('should handle zero duration tasks', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01T08:00:00Z')
      const end = new Date('2024-01-01T08:00:00Z')

      const workingHours = store.calculateWorkingHours(start, end)

      expect(workingHours).toBe(0)
    })

    it('should handle large date ranges', () => {
      const store = useCalendarStore()

      const start = new Date('2024-01-01')
      const end = new Date('2024-12-31')

      const workingDays = store.calculateWorkingDays(start, end)

      expect(workingDays).toBeGreaterThan(200)
      expect(workingDays).toBeLessThan(270) // Approximate working days in a year
    })
  })
})
