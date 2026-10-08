package repositories

import (
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/models"
)

type RecognitionTaskRepository struct{ BaseRepository }

func NewRecognitionTaskRepository(tx *gorm.DB) RecognitionTaskRepository {
	return RecognitionTaskRepository{BaseRepository{DB: GetDB(), TX: tx}}
}

func (r RecognitionTaskRepository) Get(id uint) (models.RecognitionTask, error) {
	var task models.RecognitionTask
	err := r.GetDB().First(&task, id).Error
	return task, err
}

func (r RecognitionTaskRepository) GetByRequest(owner uint, request string) (models.RecognitionTask, error) {
	var task models.RecognitionTask
	err := r.GetDB().Where("owner_user_id = ? AND client_request_id = ?", owner, request).First(&task).Error
	return task, err
}

func (r RecognitionTaskRepository) Update(query *gorm.DB, values map[string]interface{}) *gorm.DB {
	values["version"] = gorm.Expr("version + 1")
	return query.Model(&models.RecognitionTask{}).Updates(values)
}
