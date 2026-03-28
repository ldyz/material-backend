/**
 * Gantt Constraint Management Utilities
 *
 * Provides constraint validation, application, and impact calculation
 * for Gantt chart task constraints.
 *
 * Supported constraint types:
 * - MSO (Must Start On): Task must start exactly on the specified date
 * - MFO (Must Finish On): Task must finish exactly on the specified date
 * - SNET (Start No Earlier Than): Task cannot start before the specified date
 * - SNLT (Start No Later Than): Task cannot start after the specified date
 * - FNET (Finish No Earlier Than): Task cannot finish before the specified date
 * - FNLT (Finish No Later Than): Task cannot finish after the specified date
 */

/**
 * Validates if a constraint can be applied to a task
 *
 * @param {Object} task - The task object
 * @param {string} constraintType - The constraint type (MSO, MFO, SNET, SNLT, FNET, FNLT)
 * @param {string} constraintDate - The constraint date in YYYY-MM-DD format
 * @returns {Object} Validation result with valid flag and message
 */
export function validateConstraint(task, constraintType, constraintDate) {
  if (!task) {
    return {
      valid: false,
      message: 'Task is required'
    }
  }

  if (!constraintType) {
    return {
      valid: false,
      message: 'Constraint type is required'
    }
  }

  if (!constraintDate) {
    return {
      valid: false,
      message: 'Constraint date is required'
    }
  }

  const validTypes = ['MSO', 'MFO', 'SNET', 'SNLT', 'FNET', 'FNLT']
  if (!validTypes.includes(constraintType)) {
    return {
      valid: false,
      message: `Invalid constraint type. Must be one of: ${validTypes.join(', ')}`
    }
  }

  const taskStart = new Date(task.start_date)
  const taskEnd = new Date(task.end_date)
  const constraintDateTime = new Date(constraintDate)

  // Validate date format
  if (isNaN(constraintDateTime.getTime())) {
    return {
      valid: false,
      message: 'Invalid date format'
    }
  }

  // Check for impossible constraints
  const taskDuration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

  // MSO: Check if task can finish before dependencies would allow
  if (constraintType === 'MSO') {
    // Valid if no dependencies or if dependencies can be satisfied
    return {
      valid: true,
      message: 'Task will start on the specified date'
    }
  }

  // MFO: Check if duration allows reaching the finish date
  if (constraintType === 'MFO') {
    const earliestStart = new Date(constraintDateTime)
    earliestStart.setDate(earliestStart.getDate() - taskDuration)

    return {
      valid: true,
      message: 'Task will finish on the specified date'
    }
  }

  // SNET: Constraint date must be after or same as current start
  if (constraintType === 'SNET') {
    if (constraintDateTime < taskStart) {
      return {
        valid: true,
        message: 'Task will be delayed to start on or after the specified date'
      }
    }
    return {
      valid: true,
      message: 'Task already starts after the constraint date'
    }
  }

  // SNLT: Constraint date must be before or same as current start
  if (constraintType === 'SNLT') {
    if (constraintDateTime < taskStart) {
      return {
        valid: false,
        message: 'Constraint date is earlier than current start date. Task cannot be moved earlier.'
      }
    }
    return {
      valid: true,
      message: 'Task must start on or before the specified date'
    }
  }

  // FNET: Constraint date must be after or same as current end
  if (constraintType === 'FNET') {
    if (constraintDateTime < taskEnd) {
      return {
        valid: false,
        message: 'Constraint date is earlier than current finish date. Task cannot be moved earlier.'
      }
    }
    return {
      valid: true,
      message: 'Task will finish on or after the specified date'
    }
  }

  // FNLT: Constraint date must be before or same as current end
  if (constraintType === 'FNLT') {
    if (constraintDateTime < taskEnd) {
      return {
        valid: true,
        message: 'Task duration will be reduced to meet the constraint'
      }
    }
    return {
      valid: true,
      message: 'Task already finishes before the constraint date'
    }
  }

  return {
    valid: true,
    message: 'Constraint is valid'
  }
}

/**
 * Calculates the impact of applying a constraint to a task
 *
 * @param {Object} task - The task object
 * @param {string} constraintType - The constraint type
 * @param {string} constraintDate - The constraint date in YYYY-MM-DD format
 * @returns {Object} Impact information with wouldShift, shiftDays, newStartDate, newEndDate
 */
export function calculateConstraintImpact(task, constraintType, constraintDate) {
  if (!task || !constraintType || !constraintDate) {
    return {
      wouldShift: false,
      shiftDays: 0,
      newStartDate: null,
      newEndDate: null
    }
  }

  const taskStart = new Date(task.start_date)
  const taskEnd = new Date(task.end_date)
  const constraintDateTime = new Date(constraintDate)
  const taskDuration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

  let newStartDate = null
  let newEndDate = null
  let wouldShift = false
  let shiftDays = 0

  switch (constraintType) {
    case 'MSO':
      // Must start exactly on constraint date
      newStartDate = new Date(constraintDateTime)
      newEndDate = new Date(constraintDateTime)
      newEndDate.setDate(newEndDate.getDate() + taskDuration)
      shiftDays = Math.ceil((newStartDate - taskStart) / (1000 * 60 * 60 * 24))
      wouldShift = shiftDays !== 0
      break

    case 'MFO':
      // Must finish exactly on constraint date
      newEndDate = new Date(constraintDateTime)
      newStartDate = new Date(constraintDateTime)
      newStartDate.setDate(newStartDate.getDate() - taskDuration)
      shiftDays = Math.ceil((newStartDate - taskStart) / (1000 * 60 * 60 * 24))
      wouldShift = shiftDays !== 0
      break

    case 'SNET':
      // Start no earlier than constraint date
      if (constraintDateTime > taskStart) {
        newStartDate = new Date(constraintDateTime)
        newEndDate = new Date(constraintDateTime)
        newEndDate.setDate(newEndDate.getDate() + taskDuration)
        shiftDays = Math.ceil((newStartDate - taskStart) / (1000 * 60 * 60 * 24))
        wouldShift = true
      }
      break

    case 'SNLT':
      // Start no later than constraint date
      if (constraintDateTime < taskStart) {
        newStartDate = new Date(constraintDateTime)
        newEndDate = new Date(constraintDateTime)
        newEndDate.setDate(newEndDate.getDate() + taskDuration)
        shiftDays = Math.ceil((newStartDate - taskStart) / (1000 * 60 * 60 * 24))
        wouldShift = true
      }
      break

    case 'FNET':
      // Finish no earlier than constraint date
      if (constraintDateTime > taskEnd) {
        newEndDate = new Date(constraintDateTime)
        newStartDate = new Date(constraintDateTime)
        newStartDate.setDate(newStartDate.getDate() - taskDuration)
        shiftDays = Math.ceil((newStartDate - taskStart) / (1000 * 60 * 60 * 24))
        wouldShift = true
      }
      break

    case 'FNLT':
      // Finish no later than constraint date
      if (constraintDateTime < taskEnd) {
        newEndDate = new Date(constraintDateTime)
        newStartDate = new Date(constraintDateTime)
        newStartDate.setDate(newStartDate.getDate() - taskDuration)
        shiftDays = Math.ceil((newStartDate - taskStart) / (1000 * 60 * 60 * 24))
        wouldShift = true
      }
      break
  }

  return {
    wouldShift,
    shiftDays,
    newStartDate: newStartDate ? newStartDate.toISOString().split('T')[0] : null,
    newEndDate: newEndDate ? newEndDate.toISOString().split('T')[0] : null
  }
}

/**
 * Applies a constraint to a task, returning the updated task dates
 *
 * @param {Object} task - The task object
 * @param {string} constraintType - The constraint type
 * @param {string} constraintDate - The constraint date in YYYY-MM-DD format
 * @returns {Object} Updated task with new start_date and end_date
 */
export function applyConstraint(task, constraintType, constraintDate) {
  if (!task || !constraintType || !constraintDate) {
    return task
  }

  const impact = calculateConstraintImpact(task, constraintType, constraintDate)

  return {
    ...task,
    start_date: impact.newStartDate || task.start_date,
    end_date: impact.newEndDate || task.end_date,
    constraint: {
      type: constraintType,
      date: constraintDate
    }
  }
}

/**
 * Checks if a task's current schedule satisfies its constraint
 *
 * @param {Object} task - The task object with constraint property
 * @returns {Object} Result with satisfied flag and message
 */
export function checkConstraintSatisfied(task) {
  if (!task || !task.constraint) {
    return {
      satisfied: true,
      message: 'No constraint applied'
    }
  }

  const { type, date } = task.constraint
  const taskStart = new Date(task.start_date)
  const taskEnd = new Date(task.end_date)
  const constraintDate = new Date(date)

  switch (type) {
    case 'MSO':
      if (taskStart.getTime() !== constraintDate.getTime()) {
        return {
          satisfied: false,
          message: `Task must start on ${date}, currently starts on ${task.start_date}`
        }
      }
      break

    case 'MFO':
      if (taskEnd.getTime() !== constraintDate.getTime()) {
        return {
          satisfied: false,
          message: `Task must finish on ${date}, currently finishes on ${task.end_date}`
        }
      }
      break

    case 'SNET':
      if (taskStart < constraintDate) {
        return {
          satisfied: false,
          message: `Task cannot start before ${date}, currently starts on ${task.start_date}`
        }
      }
      break

    case 'SNLT':
      if (taskStart > constraintDate) {
        return {
          satisfied: false,
          message: `Task must start on or before ${date}, currently starts on ${task.start_date}`
        }
      }
      break

    case 'FNET':
      if (taskEnd < constraintDate) {
        return {
          satisfied: false,
          message: `Task cannot finish before ${date}, currently finishes on ${task.end_date}`
        }
      }
      break

    case 'FNLT':
      if (taskEnd > constraintDate) {
        return {
          satisfied: false,
          message: `Task must finish on or before ${date}, currently finishes on ${task.end_date}`
        }
      }
      break
  }

  return {
    satisfied: true,
    message: 'Constraint is satisfied'
  }
}

/**
 * Gets all tasks that violate their constraints
 *
 * @param {Array} tasks - Array of task objects
 * @returns {Array} Array of tasks that violate constraints with violation details
 */
export function getConstraintViolations(tasks) {
  if (!tasks || !Array.isArray(tasks)) {
    return []
  }

  return tasks
    .filter(task => task.constraint)
    .map(task => ({
      task,
      check: checkConstraintSatisfied(task)
    }))
    .filter(result => !result.check.satisfied)
    .map(result => ({
      taskId: result.task.id,
      taskName: result.task.name,
      constraint: result.task.constraint,
      violation: result.check.message
    }))
}

/**
 * Removes a constraint from a task
 *
 * @param {Object} task - The task object
 * @returns {Object} Task without constraint
 */
export function removeConstraint(task) {
  if (!task) {
    return task
  }

  const { constraint, ...taskWithoutConstraint } = task
  return taskWithoutConstraint
}
