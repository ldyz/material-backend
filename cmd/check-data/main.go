package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yourorg/material-backend/backend/internal/db"
)

type StockInfo struct {
	ID         uint
	MaterialID *uint
	Quantity   float64
	Unit       string
}

type MaterialInfo struct {
	ID   uint
	Name string
}

func main() {
	// 从环境变量或配置文件获取数据库连接字符串
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// 默认PostgreSQL连接
		dsn = "host=localhost port=5432 user=postgres dbname=material sslmode=disable"
	}

	database, err := db.New(dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	fmt.Println("========== 库存数据完整性检查 ==========\n")

	// 1. 统计总库存数
	var totalStocks int64
	database.Table("stocks").Count(&totalStocks)
	fmt.Printf("总库存记录数: %d\n", totalStocks)

	// 2. 统计有material_name的库存数
	var withMaterial int64
	database.Table("stocks s").
		Joins("LEFT JOIN materials m ON s.material_id = m.id").
		Where("m.id IS NOT NULL AND m.name IS NOT NULL AND m.name != ''").
		Count(&withMaterial)
	fmt.Printf("有材料名称的库存记录数: %d\n", withMaterial)

	// 3. 统计孤立的库存记录（material_id指向不存在的记录）
	var orphanedStocks int64
	database.Table("stocks s").
		Joins("LEFT JOIN materials m ON s.material_id = m.id").
		Where("s.material_id IS NOT NULL AND m.id IS NULL").
		Count(&orphanedStocks)
	fmt.Printf("孤立的库存记录数: %d\n", orphanedStocks)

	// 4. 查看孤立的库存详情
	if orphanedStocks > 0 {
		fmt.Println("\n========== 孤立的库存记录详情 ==========")
		var orphanedList []StockInfo
		database.Table("stocks s").
			Select("s.id, s.material_id, s.quantity, s.unit").
			Joins("LEFT JOIN materials m ON s.material_id = m.id").
			Where("s.material_id IS NOT NULL AND m.id IS NULL").
			Scan(&orphanedList)

		for _, stock := range orphanedList {
			materialID := "NULL"
			if stock.MaterialID != nil {
				materialID = fmt.Sprintf("%d", *stock.MaterialID)
			}
			fmt.Printf("  Stock ID: %d, Material ID: %s, Quantity: %.2f %s\n",
				stock.ID, materialID, stock.Quantity, stock.Unit)
		}
	}

	// 5. 查看没有material_id的库存
	var nullMaterialID int64
	database.Table("stocks").Where("material_id IS NULL").Count(&nullMaterialID)
	fmt.Printf("\nmaterial_id为NULL的库存记录数: %d\n", nullMaterialID)

	// 6. 统计materials总数
	var totalMaterials int64
	database.Table("materials").Count(&totalMaterials)
	fmt.Printf("\n材料记录总数: %d\n", totalMaterials)

	// 7. 检查material_id的最大值和最小值
	var maxMaterialID, minMaterialID int
	database.Table("stocks").
		Select("COALESCE(MAX(material_id), 0)").
		Scan(&maxMaterialID)
	database.Table("stocks").
		Select("COALESCE(MIN(material_id), 0)").
		Where("material_id IS NOT NULL").
		Scan(&minMaterialID)
	fmt.Printf("库存中material_id范围: %d - %d\n", minMaterialID, maxMaterialID)

	// 8. 检查materials表中ID的最大值
	var maxIDInMaterials int
	database.Table("materials").Select("COALESCE(MAX(id), 0)").Scan(&maxIDInMaterials)
	fmt.Printf("materials表中最大ID: %d\n", maxIDInMaterials)

	fmt.Println("\n========== 建议 ==========")
	if orphanedStocks > 0 {
		fmt.Printf("发现 %d 条孤立的库存记录，建议清理\n", orphanedStocks)
		fmt.Println("\n清理SQL:")
		fmt.Println("-- 方案1: 删除孤立记录")
		fmt.Println("DELETE FROM stocks WHERE material_id IS NOT NULL AND material_id NOT IN (SELECT id FROM materials);")
		fmt.Println("\n-- 方案2: 将material_id设为NULL")
		fmt.Println("UPDATE stocks SET material_id = NULL WHERE material_id IS NOT NULL AND material_id NOT IN (SELECT id FROM materials);")
	} else {
		fmt.Println("✅ 数据完整性良好，没有孤立的库存记录")
	}
}
