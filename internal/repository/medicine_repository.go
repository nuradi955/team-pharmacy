package repository

import (
	"team-pharmacy/internal/config/models"

	"gorm.io/gorm"
)

type MedicineRepository interface {
	Create(medicine *models.Medicine) error
	GetAll() ([]models.Medicine, error)
	GetByID(id uint) (*models.Medicine, error)
	Update(medicine *models.Medicine) error
	Delete(id uint) error
	UpdateAvgRating(medicineId uint, avg float64) error
}
type MedicineRepo struct {
	db *gorm.DB
}


func NewMedicineRepository(db *gorm.DB) MedicineRepository {
	return &MedicineRepo{db: db}
}
func (m *MedicineRepo) Create(medicine *models.Medicine) error {
	return m.db.Create(medicine).Error
}
func (m *MedicineRepo) GetAll() ([]models.Medicine, error) {
	var medicines []models.Medicine
	err := m.db.Find(&medicines).Error
	if err != nil {
		return nil, err
	}
	return medicines, nil
}
func (m *MedicineRepo) GetByID(id uint) (*models.Medicine, error) {
	medicine := models.Medicine{}
	err := m.db.First(&medicine, id).Error
	if err != nil {
		return nil, err
	}
	return &medicine, nil
}
func (m *MedicineRepo) Update(medicine *models.Medicine) error {
	return m.db.Save(medicine).Error
}
func (m *MedicineRepo) Delete(id uint) error {
	return m.db.Delete(&models.Medicine{}, id).Error
}
func (m *MedicineRepo) UpdateAvgRating(medicineID uint, avg float64) error {
	return m.db.Model(&models.Medicine{}).Where("id = ?", medicineID).Update("avg_rating", avg).Error
}
