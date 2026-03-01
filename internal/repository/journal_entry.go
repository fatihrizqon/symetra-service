package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var journalSortColumns = map[string]string{
	"journal_number": "journal_entries.journal_number",
	"date":           "journal_entries.date",
	"description":    "journal_entries.description",
	"status":         "journal_entries.status",
	"type":           "journal_entries.type",
	"total_debit":    "journal_entries.total_debit",
	"total_credit":   "journal_entries.total_credit",
	"created_at":     "journal_entries.created_at",
	"updated_at":     "journal_entries.updated_at",
}

type IJournalEntryRepository interface {
	Create(entry entity.JournalEntry, lines []entity.JournalLine) (entity.JournalEntry, error)
	FindAll(qp *util.QueryParams) ([]entity.JournalEntry, int, error)
	FindById(id uuid.UUID) (entity.JournalEntry, error)
	Update(entry entity.JournalEntry, lines []entity.JournalLine) (entity.JournalEntry, error)
	Delete(id uuid.UUID) error
	Post(id uuid.UUID) error
	Void(id uuid.UUID) error
	GenerateJournalNumber(journalType entity.JournalType) (string, error)
	HasTransactions(coaId uuid.UUID) (bool, error)
}

type JournalEntryRepository struct {
	Db *gorm.DB
}

func NewJournalEntryRepository(db *gorm.DB) IJournalEntryRepository {
	return &JournalEntryRepository{Db: db}
}

// GenerateJournalNumber creates a sequential journal number with a type-specific prefix.
// general → JE-202506-0001
// revenue → RV-202506-0001
// expense → EX-202506-0001
func (r *JournalEntryRepository) GenerateJournalNumber(journalType entity.JournalType) (string, error) {
	prefix, ok := entity.JournalNumberPrefix[journalType]
	if !ok {
		prefix = "JE"
	}

	now := time.Now()
	monthPrefix := fmt.Sprintf("%s-%d%02d", prefix, now.Year(), now.Month())

	var count int64
	if err := r.Db.Model(&entity.JournalEntry{}).
		Where("journal_number LIKE ?", monthPrefix+"%").
		Count(&count).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%04d", monthPrefix, count+1), nil
}

// HasTransactions checks if a COA account has been used in any journal line.
func (r *JournalEntryRepository) HasTransactions(coaId uuid.UUID) (bool, error) {
	var count int64
	if err := r.Db.Model(&entity.JournalLine{}).
		Where("coa_id = ?", coaId).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create persists a journal entry + its lines atomically.
func (r *JournalEntryRepository) Create(entry entity.JournalEntry, lines []entity.JournalLine) (entity.JournalEntry, error) {
	tx := r.Db.Begin()

	if err := tx.Create(&entry).Error; err != nil {
		tx.Rollback()
		return entry, err
	}

	for i := range lines {
		lines[i].JournalEntryId = entry.Id
	}

	if err := tx.Create(&lines).Error; err != nil {
		tx.Rollback()
		return entry, err
	}

	tx.Commit()

	if err := r.Db.Preload("Lines.COA").First(&entry, "id = ?", entry.Id).Error; err != nil {
		return entry, err
	}

	return entry, nil
}

// FindAll retrieves paginated journal entries with optional search, status, and type filters.
func (r *JournalEntryRepository) FindAll(qp *util.QueryParams) ([]entity.JournalEntry, int, error) {
	var entities []entity.JournalEntry
	var totalCount int64

	query := r.Db.Model(&entity.JournalEntry{})
	query = util.ApplySearch(query, qp)
	query = entity.JournalEntry{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, journalSortColumns, "journal_entries.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(totalCount), nil
}

// FindById retrieves a single journal entry with all lines preloaded.
func (r *JournalEntryRepository) FindById(id uuid.UUID) (entity.JournalEntry, error) {
	var entry entity.JournalEntry
	if err := r.Db.
		Preload("Lines.COA").
		First(&entry, "id = ?", id).Error; err != nil {
		return entry, err
	}
	return entry, nil
}

// Update replaces journal entry header + deletes old lines + inserts new lines atomically.
func (r *JournalEntryRepository) Update(entry entity.JournalEntry, lines []entity.JournalLine) (entity.JournalEntry, error) {
	tx := r.Db.Begin()

	if err := tx.Model(&entry).Updates(map[string]interface{}{
		"date":         entry.Date,
		"description":  entry.Description,
		"total_debit":  entry.TotalDebit,
		"total_credit": entry.TotalCredit,
		"updated_at":   time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return entry, err
	}

	if err := tx.Where("journal_entry_id = ?", entry.Id).Delete(&entity.JournalLine{}).Error; err != nil {
		tx.Rollback()
		return entry, err
	}

	for i := range lines {
		lines[i].Id = uuid.New()
		lines[i].JournalEntryId = entry.Id
	}

	if err := tx.Create(&lines).Error; err != nil {
		tx.Rollback()
		return entry, err
	}

	tx.Commit()

	if err := r.Db.Preload("Lines.COA").First(&entry, "id = ?", entry.Id).Error; err != nil {
		return entry, err
	}

	return entry, nil
}

// Delete removes a draft journal entry and its lines.
func (r *JournalEntryRepository) Delete(id uuid.UUID) error {
	tx := r.Db.Begin()

	if err := tx.Where("journal_entry_id = ?", id).Delete(&entity.JournalLine{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&entity.JournalEntry{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Post transitions a draft journal entry to posted status.
func (r *JournalEntryRepository) Post(id uuid.UUID) error {
	result := r.Db.Model(&entity.JournalEntry{}).
		Where("id = ? AND status = ?", id, entity.JournalStatusDraft).
		Updates(map[string]interface{}{
			"status":     entity.JournalStatusPosted,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("journal entry not found or is not in draft status")
	}
	return nil
}

// Void transitions a posted journal entry to void status.
func (r *JournalEntryRepository) Void(id uuid.UUID) error {
	result := r.Db.Model(&entity.JournalEntry{}).
		Where("id = ? AND status = ?", id, entity.JournalStatusPosted).
		Updates(map[string]interface{}{
			"status":     entity.JournalStatusVoid,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("journal entry not found or is not in posted status")
	}
	return nil
}
