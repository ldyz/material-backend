package progress

import (
	"log"
	"sync"
	"time"

	"github.com/yourorg/material-backend/backend/internal/api/project"
	"gorm.io/gorm"
)

// ProgressWatcher 监听进度计划变化并自动更新项目进度
type ProgressWatcher struct {
	db            *gorm.DB
	stopChan      chan struct{}
	wg            sync.WaitGroup
	CheckInterval time.Duration
	// 记录上次检查的时间戳，避免重复处理
	LastCheckTime map[uint]time.Time
	mu            sync.RWMutex
}

// NewProgressWatcher 创建新的进度监听器
func NewProgressWatcher(db *gorm.DB) *ProgressWatcher {
	return &ProgressWatcher{
		db:            db,
		stopChan:      make(chan struct{}),
		CheckInterval: 10 * time.Second, // 每10秒检查一次
		LastCheckTime: make(map[uint]time.Time),
	}
}

// Start 启动监听器
func (w *ProgressWatcher) Start() {
	log.Println("进度监听器已启动")

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(w.CheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				w.checkAndUpdateProgress()
			case <-w.stopChan:
				log.Println("进度监听器正在停止...")
				return
			}
		}
	}()
}

// Stop 停止监听器
func (w *ProgressWatcher) Stop() {
	log.Println("正在停止进度监听器...")
	close(w.stopChan)
	w.wg.Wait()
	log.Println("进度监听器已停止")
}

// checkAndUpdateProgress 检查并更新项目进度
func (w *ProgressWatcher) checkAndUpdateProgress() {
	var schedules []ProjectSchedule

	// 查找最近更新的进度计划
	w.mu.Lock()
	if err := w.db.Where("updated_at > ?", time.Now().Add(-w.CheckInterval*2)).
		Find(&schedules).Error; err != nil {
		log.Printf("查询进度计划失败: %v", err)
		w.mu.Unlock()
		return
	}

	for _, schedule := range schedules {
		// 检查是否已经处理过这个时间戳
		if lastTime, exists := w.LastCheckTime[schedule.ProjectID]; exists &&
			lastTime.Equal(schedule.UpdatedAt) || lastTime.After(schedule.UpdatedAt) {
			continue
		}

		// 更新项目进度
		if err := w.updateProjectProgress(schedule.ProjectID, schedule.Data); err != nil {
			log.Printf("更新项目 %d 进度失败: %v", schedule.ProjectID, err)
		} else {
			// 记录处理时间戳
			w.LastCheckTime[schedule.ProjectID] = schedule.UpdatedAt
			log.Printf("已更新项目 %d 的进度", schedule.ProjectID)
		}
	}
	w.mu.Unlock()
}

// updateProjectProgress 更新项目进度
func (w *ProgressWatcher) updateProjectProgress(projectID uint, scheduleData ScheduleData) error {
	// 计算所有活动的平均进度
	totalProgress := 0.0
	count := 0

	for _, activity := range scheduleData.Activities {
		// 跳过虚拟活动
		if activity.IsDummy {
			continue
		}
		totalProgress += activity.Progress
		count++
	}

	var progressPercentage float64
	if count > 0 {
		progressPercentage = totalProgress / float64(count)
	}

	// 更新项目进度
	return w.db.Model(&project.Project{}).
		Where("id = ?", projectID).
		Update("progress_percentage", progressPercentage).Error
}

// UpdateAllProjectsProgress 手动更新所有项目的进度
func (w *ProgressWatcher) UpdateAllProjectsProgress() error {
	log.Println("开始更新所有项目的进度...")

	var schedules []ProjectSchedule
	if err := w.db.Find(&schedules).Error; err != nil {
		return err
	}

	for _, schedule := range schedules {
		if err := w.updateProjectProgress(schedule.ProjectID, schedule.Data); err != nil {
			log.Printf("更新项目 %d 进度失败: %v", schedule.ProjectID, err)
		}
	}

	log.Printf("已更新 %d 个项目的进度", len(schedules))
	return nil
}

// ForceUpdateProjectProgress 强制更新指定项目的进度
func (w *ProgressWatcher) ForceUpdateProjectProgress(projectID uint) error {
	var schedule ProjectSchedule
	if err := w.db.Where("project_id = ?", projectID).First(&schedule).Error; err != nil {
		return err
	}

	return w.updateProjectProgress(projectID, schedule.Data)
}

// GetStatus 获取监听器状态
func (w *ProgressWatcher) GetStatus() map[string]interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return map[string]interface{}{
		"status":         "running",
		"check_interval": w.CheckInterval.String(),
		"last_check_times": w.LastCheckTime,
		"watched_projects_count": len(w.LastCheckTime),
	}
}
