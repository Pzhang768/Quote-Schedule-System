package db

import (
	"github.com/melfish/br-api/internal/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var managerCount int64
	db.Model(&models.Manager{}).Count(&managerCount)
	if managerCount == 0 {
		managers := []models.Manager{
			{Name: "Sarah Chen", Email: "sarah.chen@br.com"},
			{Name: "James Okafor", Email: "james.okafor@br.com"},
			{Name: "Maria Santos", Email: "maria.santos@br.com"},
		}
		if err := db.Create(&managers).Error; err != nil {
			return err
		}
	}

	var technicianCount int64
	db.Model(&models.Technician{}).Count(&technicianCount)
	if technicianCount == 0 {
		technicians := []models.Technician{
			{Name: "Tom Brennan", Email: "tom.brennan@br.com"},
			{Name: "Lisa Nguyen", Email: "lisa.nguyen@br.com"},
			{Name: "Carlos Vega", Email: "carlos.vega@br.com"},
			{Name: "Priya Sharma", Email: "priya.sharma@br.com"},
			{Name: "Jack Wilson", Email: "jack.wilson@br.com"},
		}
		if err := db.Create(&technicians).Error; err != nil {
			return err
		}
	}

	var quoteCount int64
	db.Model(&models.Quote{}).Count(&quoteCount)
	if quoteCount == 0 {
		quotes := []models.Quote{
			{CustomerName: "Bob Fletcher", Address: "12 Maple St, Sydney", Description: "Fix leaking roof"},
			{CustomerName: "Ellen Moore", Address: "45 River Rd, Melbourne", Description: "Replace hot water system"},
			{CustomerName: "Kevin Hartley", Address: "8 Ocean Ave, Brisbane", Description: "Install solar panels"},
		}
		if err := db.Create(&quotes).Error; err != nil {
			return err
		}
	}

	return nil
}
