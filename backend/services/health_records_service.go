package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/clarity/backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HealthRecordsService struct {
	db *gorm.DB
}

func NewHealthRecordsService(db *gorm.DB) *HealthRecordsService {
	return &HealthRecordsService{db: db}
}

// CreateRecord creates a new health record
func (hrs *HealthRecordsService) CreateRecord(userID, recordType, title, description string, metadata map[string]string) (*models.HealthRecord, error) {
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	record := models.HealthRecord{
		ID:          uuid.New().String(),
		UserID:      userID,
		RecordType:  recordType,
		Title:       title,
		Description: description,
		Metadata:    string(metadataJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := hrs.db.Create(&record).Error; err != nil {
		return nil, fmt.Errorf("failed to create record: %w", err)
	}

	return &record, nil
}

// GetRecord retrieves a single record
func (hrs *HealthRecordsService) GetRecord(recordID string) (*models.HealthRecord, error) {
	var record models.HealthRecord
	if err := hrs.db.First(&record, "id = ?", recordID).Error; err != nil {
		return nil, fmt.Errorf("record not found: %w", err)
	}
	return &record, nil
}

// ListRecords retrieves records with pagination
func (hrs *HealthRecordsService) ListRecords(userID string, limit, offset int) ([]models.HealthRecord, int64, error) {
	var records []models.HealthRecord
	var total int64

	if err := hrs.db.Where("user_id = ?", userID).Model(&models.HealthRecord{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count records: %w", err)
	}

	if err := hrs.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list records: %w", err)
	}

	return records, total, nil
}

// UpdateRecord updates an existing record
func (hrs *HealthRecordsService) UpdateRecord(recordID, title, description string, metadata map[string]string) (*models.HealthRecord, error) {
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	record := models.HealthRecord{
		Title:       title,
		Description: description,
		Metadata:    string(metadataJSON),
		UpdatedAt:   time.Now(),
	}

	if err := hrs.db.Model(&models.HealthRecord{}).Where("id = ?", recordID).Updates(record).Error; err != nil {
		return nil, fmt.Errorf("failed to update record: %w", err)
	}

	return hrs.GetRecord(recordID)
}

// DeleteRecord deletes a record
func (hrs *HealthRecordsService) DeleteRecord(recordID string) error {
	if err := hrs.db.Delete(&models.HealthRecord{}, "id = ?", recordID).Error; err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	return nil
}
