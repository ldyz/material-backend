package db

import (
	"fmt"
	"gorm.io/gorm"
)

// SyncPostgreSQLSequences synchronizes PostgreSQL sequences with actual table data
// This fixes issues where sequences are out of sync with table data
func SyncPostgreSQLSequences(db *gorm.DB) error {
	// List of tables to sync
	tables := []string{
		"inbound_orders",
		"inbound_items",
		"requisitions",
		"requisition_items",
		"materials",
		"material_master",
		"material_plans",
		"material_plan_items",
		"projects",
		"stocks",
		"stock_logs",
		"notifications",
		"users",
		"workflow_instances",
		"workflow_tasks",
	}

	for _, table := range tables {
		var seqName string
		var maxID uint

		// Get sequence name for this table
		err := db.Raw(`
			SELECT pg_get_serial_sequence($1, 'id')
		`, table).Scan(&seqName).Error

		if err != nil || seqName == "" {
			// Table doesn't have a serial sequence, skip
			continue
		}

		// Get max ID from table
		err = db.Raw(fmt.Sprintf("SELECT COALESCE(MAX(id), 0) FROM %s", table)).Scan(&maxID).Error
		if err != nil {
			continue
		}

		// Reset sequence to max ID + 1
		sql := fmt.Sprintf("SELECT setval('%s', %d, false)", seqName, maxID+1)
		if err := db.Exec(sql).Error; err != nil {
			fmt.Printf("Warning: Failed to sync sequence for table %s: %v\n", table, err)
		}
	}

	return nil
}
