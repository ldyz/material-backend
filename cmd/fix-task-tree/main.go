package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:200"`
	ParentID *uint  `gorm:"index"`
	ProjectID uint `gorm:"index"`
}

func main() {
	dsn := "host=127.0.0.1 port=5432 user=materials password=julei1984 dbname=materials sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 获取所有任务
	var tasks []Task
	if err := db.Order("id ASC").Find(&tasks).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 当前任务层级关系 ===")
	for _, task := range tasks {
		parentID := "null"
		if task.ParentID != nil {
			parentID = fmt.Sprintf("%d", *task.ParentID)
		}
		fmt.Printf("任务 %d (%s): parent_id = %s\n", task.ID, task.Name, parentID)
	}

	// 检测循环依赖
	fmt.Println("\n=== 检测循环依赖 ===")
	visited := make(map[uint]bool)
	inStack := make(map[uint]bool)

	var detectCycle func(uint) bool
	detectCycle = func(taskID uint) bool {
		visited[taskID] = true
		inStack[taskID] = true

		// 找到该任务
		var task Task
		if err := db.First(&task, taskID).Error; err != nil {
			return false
		}

		if task.ParentID != nil {
			if inStack[*task.ParentID] {
				fmt.Printf("发现循环: 任务 %d -> 任务 %d\n", taskID, *task.ParentID)
				return true
			}
			if !visited[*task.ParentID] {
				if detectCycle(*task.ParentID) {
					return true
				}
			}
		}

		inStack[taskID] = false
		return false
	}

	hasCycle := false
	for _, task := range tasks {
		if !visited[task.ID] {
			if detectCycle(task.ID) {
				hasCycle = true
			}
		}
	}

	if !hasCycle {
		fmt.Println("未检测到循环依赖")
		return
	}

	// 修复策略：
	// 1. 找出被引用最多的任务（最可能是根任务）
	// 2. 找出循环中的任务
	// 3. 断开循环：将循环中ID最小的任务的 parent_id 设为 null

	fmt.Println("\n=== 修复循环依赖 ===")

	// 统计每个任务被作为父任务的次数
	parentCount := make(map[uint]int)
	for _, task := range tasks {
		if task.ParentID != nil {
			parentCount[*task.ParentID]++
		}
	}

	fmt.Println("任务被引用次数:")
	for id, count := range parentCount {
		fmt.Printf("  任务 %d: 被引用 %d 次\n", id, count)
	}

	// 策略：找出所有涉及循环的任务，将 ID 最小的设为根任务
	// 根据前面的分析，循环是 4 -> 5 -> 4
	// 我们将任务 4 设为根任务（parent_id = null）
	// 然后确保任务 5 的 parent_id = 4

	fmt.Println("\n修复方案:")
	fmt.Println("1. 将任务 4 设为根任务 (parent_id = NULL)")
	fmt.Println("2. 将任务 5 的 parent_id 改为 4")
	fmt.Println("3. 其他任务的 parent_id 保持不变")

	// 执行修复
	updates := 0

	// 将任务4设为根任务
	result := db.Model(&Task{}).Where("id = ?", 4).Update("parent_id", nil)
	if result.Error != nil {
		log.Printf("更新任务4失败: %v", result.Error)
	} else {
		fmt.Printf("✓ 任务 4 已设为根任务\n")
		updates++
	}

	// 确保任务5的parent_id是4
	var task5 Task
	if err := db.First(&task5, 5).Error; err == nil {
		if task5.ParentID == nil || *task5.ParentID != 4 {
			result := db.Model(&task5).Update("parent_id", 4)
			if result.Error != nil {
				log.Printf("更新任务5失败: %v", result.Error)
			} else {
				fmt.Printf("✓ 任务 5 的 parent_id 已改为 4\n")
				updates++
			}
		}
	}

	// 显示修复后的结果
	fmt.Println("\n=== 修复后的任务层级关系 ===")
	var tasksAfter []Task
	if err := db.Order("id ASC").Find(&tasksAfter).Error; err != nil {
		log.Fatal(err)
	}

	for _, task := range tasksAfter {
		parentID := "null"
		if task.ParentID != nil {
			parentID = fmt.Sprintf("%d", *task.ParentID)
		}
		fmt.Printf("任务 %d (%s): parent_id = %s\n", task.ID, task.Name, parentID)
	}

	fmt.Printf("\n总计修复了 %d 个任务\n", updates)
	fmt.Println("完成！")
}
