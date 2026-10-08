package task

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, topicID *uuid.UUID, p pagination.Params) ([]Task, int64, error) {
	q := r.db.WithContext(ctx).Model(&Task{})
	if topicID != nil {
		q = q.Where("topic_id = ?", *topicID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []Task
	if err := q.Order("created_at asc").Offset(p.Offset()).Limit(p.Limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	var t Task
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) NextForStudent(ctx context.Context, studentID uuid.UUID, mode string, subjectID, topicID *uuid.UUID, limit int) ([]Task, error) {
	if limit <= 0 {
		limit = 1
	}
	q := r.db.WithContext(ctx).Model(&Task{}).
		Where("id NOT IN (?)",
			r.db.Table("task_attempts").Select("task_id").Where("student_id = ?", studentID),
		)
	if subjectID != nil {
		q = q.Where("subject_id = ?", *subjectID)
	}
	if topicID != nil {
		q = q.Where("topic_id = ?", *topicID)
	}
	if mode == "weak" {
		q = q.Joins("LEFT JOIN progress ON progress.topic_id = tasks.topic_id AND progress.student_id = ?", studentID).
			Order("COALESCE(progress.score, 0) ASC").Order("random()")
	} else {
		q = q.Order("random()")
	}
	var items []Task
	if err := q.Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, apperrors.ErrNotFound
	}
	return items, nil
}

func (r *Repository) TrainerTopics(ctx context.Context) ([]TrainerTopic, error) {
	var topics []TrainerTopic
	err := r.db.WithContext(ctx).Table("topics").
		Select("topics.id, topics.subject_id, subjects.name AS subject_name, topics.name, COUNT(tasks.id) AS count").
		Joins("JOIN subjects ON subjects.id = topics.subject_id").
		Joins("JOIN tasks ON tasks.topic_id = topics.id").
		Group("topics.id, topics.subject_id, subjects.name, topics.name").
		Order("subjects.name, topics.name").Scan(&topics).Error
	return topics, err
}

func (r *Repository) TrainerCatalog(ctx context.Context) ([]TrainerSubject, error) {
	type row struct {
		SubjectID   uuid.UUID
		SubjectName string
		SubjectCode string
		TaskNumber  int
		Count       int
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Table("tasks").
		Select("subjects.id AS subject_id, subjects.name AS subject_name, subjects.code AS subject_code, tasks.task_number, COUNT(tasks.id) AS count").
		Joins("JOIN subjects ON subjects.id = tasks.subject_id").
		Group("subjects.id, subjects.name, subjects.code, tasks.task_number").
		Order("subjects.name, tasks.task_number").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	byID := make(map[uuid.UUID]int)
	items := make([]TrainerSubject, 0)
	for _, row := range rows {
		index, ok := byID[row.SubjectID]
		if !ok {
			index = len(items)
			byID[row.SubjectID] = index
			items = append(items, TrainerSubject{ID: row.SubjectID, Name: row.SubjectName, Code: row.SubjectCode, TaskNumbers: make([]TrainerTaskNumber, 0)})
		}
		items[index].TaskNumbers = append(items[index].TaskNumbers, TrainerTaskNumber{Number: row.TaskNumber, Count: row.Count})
	}
	return items, nil
}

func (r *Repository) TrainingBatch(ctx context.Context, subjectID uuid.UUID, taskNumber, limit int) ([]Task, error) {
	if limit < 1 {
		limit = 1
	}
	if limit > 30 {
		limit = 30
	}
	var items []Task
	err := r.db.WithContext(ctx).
		Where("subject_id = ? AND task_number = ?", subjectID, taskNumber).
		Order("random()").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, apperrors.ErrNotFound
	}
	return items, nil
}

func (r *Repository) MockSession(ctx context.Context, subjectID uuid.UUID) ([]Task, error) {
	var items []Task
	err := r.db.WithContext(ctx).Raw(`
		SELECT * FROM (
			SELECT tasks.*, ROW_NUMBER() OVER (PARTITION BY task_number ORDER BY random()) AS row_number
			FROM tasks WHERE subject_id = ?
		) selected
		WHERE row_number = 1
		ORDER BY task_number
		LIMIT 30`, subjectID).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, apperrors.ErrNotFound
	}
	return items, nil
}
