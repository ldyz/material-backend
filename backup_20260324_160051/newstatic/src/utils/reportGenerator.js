/**
 * Report Generator Utility
 * Generates and exports various types of reports from Gantt data
 *
 * @module reportGenerator
 */

import jsPDF from 'jspdf'
import 'jspdf-autotable'
import * as XLSX from 'xlsx'
import { format } from 'date-fns'

/**
 * Column definitions for different report types
 */
const REPORT_COLUMNS = {
  task: [
    { key: 'id', label: 'ID', width: 60 },
    { key: 'name', label: 'Task Name', width: 200 },
    { key: 'status', label: 'Status', width: 100 },
    { key: 'assignee', label: 'Assignee', width: 120 },
    { key: 'priority', label: 'Priority', width: 80 },
    { key: 'progress', label: 'Progress', width: 80 },
    { key: 'startDate', label: 'Start Date', width: 100 },
    { key: 'endDate', label: 'End Date', width: 100 },
    { key: 'duration', label: 'Duration', width: 80 },
    { key: 'budget', label: 'Budget', width: 100 },
    { key: 'actualCost', label: 'Actual Cost', width: 100 }
  ],
  resource: [
    { key: 'name', label: 'Resource Name', width: 150 },
    { key: 'role', label: 'Role', width: 120 },
    { key: 'assignedTasks', label: 'Assigned Tasks', width: 120 },
    { key: 'totalHours', label: 'Total Hours', width: 100 },
    { key: 'capacity', label: 'Capacity', width: 100 },
    { key: 'utilization', label: 'Utilization %', width: 100 }
  ],
  milestone: [
    { key: 'id', label: 'ID', width: 60 },
    { key: 'name', label: 'Milestone', width: 200 },
    { key: 'date', label: 'Date', width: 100 },
    { key: 'status', label: 'Status', width: 100 },
    { key: 'progress', label: 'Progress', width: 80 },
    { key: 'description', label: 'Description', width: 250 }
  ],
  progress: [
    { key: 'name', label: 'Task/Milestone', width: 200 },
    { key: 'startDate', label: 'Start Date', width: 100 },
    { key: 'endDate', label: 'End Date', width: 100 },
    { key: 'progress', label: 'Progress %', width: 100 },
    { key: 'status', label: 'Status', width: 100 },
    { key: 'plannedValue', label: 'Planned Value', width: 120 },
    { key: 'earnedValue', label: 'Earned Value', width: 120 },
    { key: 'actualCost', label: 'Actual Cost', width: 120 }
  ]
}

/**
 * Generate report data structure
 *
 * @param {Object} params - Report parameters
 * @param {Array} params.data - Raw data array
 * @param {Array} params.columns - Column definitions
 * @param {Object} params.config - Report configuration
 * @returns {Object} Formatted report object
 */
export function generateReport({ data, columns, config }) {
  const reportType = config.type || 'task'
  const reportColumns = REPORT_COLUMNS[reportType] || columns

  // Group data if groupBy is specified
  let groupedData = data
  let groups = []

  if (config.groupBy) {
    const groupMap = new Map()

    data.forEach(item => {
      const groupKey = item[config.groupBy] || 'Uncategorized'
      if (!groupMap.has(groupKey)) {
        groupMap.set(groupKey, [])
      }
      groupMap.get(groupKey).push(item)
    })

    groups = Array.from(groupMap.entries()).map(([key, items]) => ({
      key,
      count: items.length,
      items
    }))
  }

  // Calculate summary statistics
  const summary = calculateSummary(data, config)

  return {
    type: reportType,
    title: generateReportTitle(config),
    dateRange: config.dateRange,
    columns: reportColumns,
    data: groupedData,
    groups: groups,
    summary: summary,
    config: config,
    generatedAt: new Date()
  }
}

/**
 * Calculate summary statistics for the report
 *
 * @param {Array} data - Data array
 * @param {Object} config - Report configuration
 * @returns {Object} Summary statistics
 */
function calculateSummary(data, config) {
  const summary = {
    totalRecords: data.length,
    dateGenerated: format(new Date(), 'yyyy-MM-dd HH:mm:ss')
  }

  if (config.type === 'task' || config.type === 'progress') {
    const completed = data.filter(t => t.status === 'completed').length
    const inProgress = data.filter(t => t.status === 'in_progress').length
    const notStarted = data.filter(t => t.status === 'not_started').length
    const avgProgress = data.reduce((sum, t) => sum + (t.progress || 0), 0) / (data.length || 1)

    Object.assign(summary, {
      completed,
      inProgress,
      notStarted,
      avgProgress: Math.round(avgProgress)
    })

    // Financial metrics
    const totalBudget = data.reduce((sum, t) => sum + (t.budget || 0), 0)
    const totalActualCost = data.reduce((sum, t) => sum + (t.actualCost || 0), 0)

    Object.assign(summary, {
      totalBudget,
      totalActualCost,
      budgetVariance: totalBudget - totalActualCost
    })
  } else if (config.type === 'resource') {
    const totalResources = data.length
    const totalHours = data.reduce((sum, r) => sum + (r.totalHours || 0), 0)
    const totalCapacity = data.reduce((sum, r) => sum + (r.capacity || 0), 0)
    const avgUtilization = data.reduce((sum, r) => sum + (r.utilization || 0), 0) / (data.length || 1)

    Object.assign(summary, {
      totalResources,
      totalHours,
      totalCapacity,
      avgUtilization: Math.round(avgUtilization)
    })
  } else if (config.type === 'milestone') {
    const completed = data.filter(m => m.status === 'completed').length
    const upcoming = data.filter(m => m.status === 'upcoming').length
    const overdue = data.filter(m => m.status === 'overdue').length

    Object.assign(summary, {
      completed,
      upcoming,
      overdue
    })
  }

  return summary
}

/**
 * Generate report title based on configuration
 *
 * @param {Object} config - Report configuration
 * @returns {string} Report title
 */
function generateReportTitle(config) {
  const typeLabels = {
    task: 'Task Report',
    resource: 'Resource Report',
    milestone: 'Milestone Report',
    progress: 'Progress Report'
  }

  let title = typeLabels[config.type] || 'Report'

  if (config.dateRange && config.dateRange.length === 2) {
    const startDate = format(new Date(config.dateRange[0]), 'MMM d, yyyy')
    const endDate = format(new Date(config.dateRange[1]), 'MMM d, yyyy')
    title += ` (${startDate} - ${endDate})`
  }

  return title
}

/**
 * Export report to PDF
 *
 * @param {Object} report - Report object from generateReport()
 * @returns {Promise<void>}
 */
export async function exportToPDF(report) {
  const doc = new jsPDF()
  const pageWidth = doc.internal.pageSize.getWidth()
  const pageHeight = doc.internal.pageSize.getHeight()
  const margin = 20

  let yPosition = margin

  // Add title
  doc.setFontSize(18)
  doc.setFont('helvetica', 'bold')
  doc.text(report.title, margin, yPosition)
  yPosition += 10

  // Add metadata
  doc.setFontSize(10)
  doc.setFont('helvetica', 'normal')
  doc.text(`Generated: ${format(new Date(), 'yyyy-MM-dd HH:mm')}`, margin, yPosition)
  yPosition += 6

  if (report.dateRange && report.dateRange.length === 2) {
    const startDate = format(new Date(report.dateRange[0]), 'MMM d, yyyy')
    const endDate = format(new Date(report.dateRange[1]), 'MMM d, yyyy')
    doc.text(`Period: ${startDate} to ${endDate}`, margin, yPosition)
    yPosition += 6
  }

  doc.text(`Total Records: ${report.summary.totalRecords}`, margin, yPosition)
  yPosition += 15

  // Add summary table
  if (Object.keys(report.summary).length > 2) {
    const summaryData = Object.entries(report.summary)
      .filter(([key]) => key !== 'totalRecords' && key !== 'dateGenerated')
      .map(([key, value]) => [
        formatLabel(key),
        formatValue(value)
      ])

    doc.autoTable({
      startY: yPosition,
      head: [['Metric', 'Value']],
      body: summaryData,
      theme: 'grid',
      headStyles: { fillColor: [64, 158, 255] },
      styles: { fontSize: 10 }
    })

    yPosition = doc.lastAutoTable.finalY + 15
  }

  // Add grouped data or main data table
  if (report.groups && report.groups.length > 0) {
    for (const group of report.groups) {
      // Add group header
      if (yPosition > pageHeight - 50) {
        doc.addPage()
        yPosition = margin
      }

      doc.setFontSize(12)
      doc.setFont('helvetica', 'bold')
      doc.text(`${formatLabel(report.config.groupBy)}: ${group.key} (${group.count} items)`, margin, yPosition)
      yPosition += 10

      // Prepare table data
      const tableData = group.items.map(item =>
        report.columns.map(col => formatValue(item[col.key]))
      )

      const tableHeaders = report.columns.map(col => col.label)

      doc.autoTable({
        startY: yPosition,
        head: [tableHeaders],
        body: tableData,
        theme: 'striped',
        headStyles: { fillColor: [64, 158, 255] },
        styles: { fontSize: 9 },
        columnStyles: report.columns.map(col => ({
          cellWidth: col.width ? col.width / 3 : 'auto'
        }))
      })

      yPosition = doc.lastAutoTable.finalY + 10
    }
  } else {
    // Prepare table data
    const tableData = report.data.map(item =>
      report.columns.map(col => formatValue(item[col.key]))
    )

    const tableHeaders = report.columns.map(col => col.label)

    doc.autoTable({
      startY: yPosition,
      head: [tableHeaders],
      body: tableData,
      theme: 'striped',
      headStyles: { fillColor: [64, 158, 255] },
      styles: { fontSize: 9 },
      columnStyles: report.columns.map(col => ({
        cellWidth: col.width ? col.width / 3 : 'auto'
      }))
    })
  }

  // Save PDF
  const fileName = `${report.title.replace(/\s+/g, '_')}_${format(new Date(), 'yyyyMMdd')}.pdf`
  doc.save(fileName)
}

/**
 * Export report to Excel
 *
 * @param {Object} report - Report object from generateReport()
 * @returns {Promise<void>}
 */
export async function exportToExcel(report) {
  const workbook = XLSX.utils.book_new()

  // Summary sheet
  if (Object.keys(report.summary).length > 0) {
    const summaryData = Object.entries(report.summary).map(([key, value]) => ({
      Metric: formatLabel(key),
      Value: formatValue(value)
    }))

    const summarySheet = XLSX.utils.json_to_sheet(summaryData)
    XLSX.utils.book_append_sheet(workbook, summarySheet, 'Summary')
  }

  // Data sheet
  let dataSheet

  if (report.groups && report.groups.length > 0) {
    // Flatten grouped data
    const flatData = []
    report.groups.forEach(group => {
      group.items.forEach(item => {
        flatData.push({
          ...item,
          _group: group.key
        })
      })
    })

    const rowData = flatData.map(item => {
      const row = {}
      report.columns.forEach(col => {
        row[col.label] = formatValue(item[col.key])
      })
      return row
    })

    dataSheet = XLSX.utils.json_to_sheet(rowData)
  } else {
    const rowData = report.data.map(item => {
      const row = {}
      report.columns.forEach(col => {
        row[col.label] = formatValue(item[col.key])
      })
      return row
    })

    dataSheet = XLSX.utils.json_to_sheet(rowData)
  }

  // Set column widths
  dataSheet['!cols'] = report.columns.map(col => ({
    wch: Math.ceil((col.width || 100) / 7)
  }))

  XLSX.utils.book_append_sheet(workbook, dataSheet, 'Data')

  // Save Excel file
  const fileName = `${report.title.replace(/\s+/g, '_')}_${format(new Date(), 'yyyyMMdd')}.xlsx`
  XLSX.writeFile(workbook, fileName)
}

/**
 * Generate print layout HTML
 *
 * @param {Object} report - Report object from generateReport()
 * @returns {string} HTML string for printing
 */
export function generatePrintLayout(report) {
  const styles = `
    <style>
      * { margin: 0; padding: 0; box-sizing: border-box; }
      body {
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        font-size: 12px;
        color: #333;
        line-height: 1.4;
        padding: 20px;
      }
      .report-header {
        text-align: center;
        margin-bottom: 30px;
        border-bottom: 2px solid #409EFF;
        padding-bottom: 20px;
      }
      .report-title {
        font-size: 24px;
        font-weight: bold;
        color: #303133;
        margin-bottom: 10px;
      }
      .report-meta {
        font-size: 11px;
        color: #909399;
      }
      .report-summary {
        margin-bottom: 30px;
        padding: 15px;
        background: #f5f7fa;
        border-radius: 4px;
      }
      .summary-title {
        font-size: 14px;
        font-weight: 600;
        margin-bottom: 10px;
        color: #303133;
      }
      .summary-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
        gap: 10px;
      }
      .summary-item {
        display: flex;
        flex-direction: column;
      }
      .summary-label {
        font-size: 11px;
        color: #909399;
        margin-bottom: 4px;
      }
      .summary-value {
        font-size: 16px;
        font-weight: 600;
        color: #303133;
      }
      table {
        width: 100%;
        border-collapse: collapse;
        margin-bottom: 20px;
      }
      th, td {
        border: 1px solid #dcdfe6;
        padding: 8px 12px;
        text-align: left;
      }
      th {
        background: #409EFF;
        color: #fff;
        font-weight: 600;
      }
      tr:nth-child(even) {
        background: #f5f7fa;
      }
      .group-header {
        background: #e8f4ff;
        padding: 10px;
        font-weight: 600;
        margin: 20px 0 10px 0;
        border-left: 4px solid #409EFF;
      }
      @media print {
        body { padding: 0; }
        .page-break { page-break-after: always; }
      }
    </style>
  `

  let html = `
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="UTF-8">
      <title>${report.title}</title>
      ${styles}
    </head>
    <body>
  `

  // Header
  html += `
    <div class="report-header">
      <div class="report-title">${report.title}</div>
      <div class="report-meta">
        Generated: ${format(new Date(), 'MMMM d, yyyy HH:mm')}
        ${report.dateRange && report.dateRange.length === 2 ?
          ` | Period: ${format(new Date(report.dateRange[0]), 'MMM d, yyyy')} - ${format(new Date(report.dateRange[1]), 'MMM d, yyyy')}` :
          ''
        }
      </div>
    </div>
  `

  // Summary
  if (Object.keys(report.summary).length > 0) {
    html += `
      <div class="report-summary">
        <div class="summary-title">Summary</div>
        <div class="summary-grid">
    `

    Object.entries(report.summary)
      .filter(([key]) => key !== 'totalRecords' && key !== 'dateGenerated')
      .forEach(([key, value]) => {
        html += `
          <div class="summary-item">
            <span class="summary-label">${formatLabel(key)}</span>
            <span class="summary-value">${formatValue(value)}</span>
          </div>
        `
      })

    html += `
        </div>
      </div>
    `
  }

  // Data tables
  if (report.groups && report.groups.length > 0) {
    report.groups.forEach(group => {
      html += `
        <div class="group-header">
          ${formatLabel(report.config.groupBy)}: ${group.key} (${group.count} items)
        </div>
      `
      html += generateTableHTML(group.items, report.columns)
    })
  } else {
    html += generateTableHTML(report.data, report.columns)
  }

  html += `
    </body>
    </html>
  `

  return html
}

/**
 * Generate HTML table from data
 *
 * @param {Array} data - Data array
 * @param {Array} columns - Column definitions
 * @returns {string} HTML table string
 */
function generateTableHTML(data, columns) {
  let html = '<table>'

  // Header
  html += '<thead><tr>'
  columns.forEach(col => {
    html += `<th style="width: ${col.width || 100}px">${col.label}</th>`
  })
  html += '</tr></thead>'

  // Body
  html += '<tbody>'
  data.forEach(item => {
    html += '<tr>'
    columns.forEach(col => {
      html += `<td>${formatValue(item[col.key])}</td>`
    })
    html += '</tr>'
  })
  html += '</tbody>'

  html += '</table>'
  return html
}

/**
 * Format label for display
 *
 * @param {string} key - Key in camelCase or snake_case
 * @returns {string} Formatted label
 */
function formatLabel(key) {
  return key
    .replace(/([A-Z])/g, ' $1')
    .replace(/_/g, ' ')
    .replace(/^./, str => str.toUpperCase())
    .trim()
}

/**
 * Format value for display
 *
 * @param {any} value - Value to format
 * @returns {string} Formatted value
 */
function formatValue(value) {
  if (value === null || value === undefined) return '-'
  if (typeof value === 'number') {
    return value.toLocaleString('en-US', { maximumFractionDigits: 2 })
  }
  if (typeof value === 'boolean') return value ? 'Yes' : 'No'
  return String(value)
}

/**
 * Save report configuration as template
 *
 * @param {Object} config - Report configuration
 * @param {string} templateName - Name for the template
 */
export function saveReportTemplate(config, templateName) {
  const templates = JSON.parse(localStorage.getItem('reportTemplates') || '{}')
  templates[templateName] = {
    config,
    savedAt: new Date().toISOString()
  }
  localStorage.setItem('reportTemplates', JSON.stringify(templates))
}

/**
 * Load report template
 *
 * @param {string} templateName - Name of the template
 * @returns {Object|null} Template configuration or null
 */
export function loadReportTemplate(templateName) {
  const templates = JSON.parse(localStorage.getItem('reportTemplates') || '{}')
  return templates[templateName]?.config || null
}

/**
 * List all saved report templates
 *
 * @returns {Array} Array of template names
 */
export function listReportTemplates() {
  const templates = JSON.parse(localStorage.getItem('reportTemplates') || '{}')
  return Object.keys(templates).map(name => ({
    name,
    savedAt: templates[name].savedAt
  }))
}

export default {
  generateReport,
  exportToPDF,
  exportToExcel,
  generatePrintLayout,
  saveReportTemplate,
  loadReportTemplate,
  listReportTemplates
}
