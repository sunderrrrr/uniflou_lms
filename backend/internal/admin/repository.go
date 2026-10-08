package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) ListUsers(ctx context.Context, p pagination.Params) ([]User, int64, error) {
	query := r.db.WithContext(ctx).Model(&auth.User{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []User
	err := r.db.WithContext(ctx).Model(&auth.User{}).
		Select("id, email, role, email_verified, created_at").
		Order("created_at DESC").Offset(p.Offset()).Limit(p.Limit).Scan(&users).Error
	return users, total, err
}

func (r *Repository) UpdateUserRole(ctx context.Context, actorID, userID uuid.UUID, role string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user auth.User
		if err := tx.Clauses().Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}
		if user.Role == auth.RoleAdmin && role != auth.RoleAdmin {
			var admins int64
			if err := tx.Model(&auth.User{}).Where("role = ?", auth.RoleAdmin).Count(&admins).Error; err != nil {
				return err
			}
			if admins <= 1 {
				return errLastAdmin
			}
		}
		return tx.Model(&auth.User{}).Where("id = ?", userID).Update("role", role).Error
	})
}

func (r *Repository) DeleteUser(ctx context.Context, actorID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user auth.User
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}
		if user.Role == auth.RoleAdmin {
			var admins int64
			if err := tx.Model(&auth.User{}).Where("role = ?", auth.RoleAdmin).Count(&admins).Error; err != nil {
				return err
			}
			if admins <= 1 {
				return errLastAdmin
			}
		}
		return tx.Delete(&user).Error
	})
}

func (r *Repository) Subjects(ctx context.Context) ([]Subject, error) {
	var subjects []Subject
	if err := r.db.WithContext(ctx).Table("subjects").Select("id, name, code").Order("name").Scan(&subjects).Error; err != nil {
		return nil, err
	}
	var topics []struct {
		ID        uuid.UUID
		SubjectID uuid.UUID
		Name      string
	}
	if err := r.db.WithContext(ctx).Table("topics").Select("id, subject_id, name").Order("name").Scan(&topics).Error; err != nil {
		return nil, err
	}
	bySubject := make(map[uuid.UUID]int, len(subjects))
	for i := range subjects {
		subjects[i].Topics = []Topic{}
		bySubject[subjects[i].ID] = i
	}
	for _, topic := range topics {
		if index, ok := bySubject[topic.SubjectID]; ok {
			subjects[index].Topics = append(subjects[index].Topics, Topic{ID: topic.ID, Name: topic.Name})
		}
	}
	return subjects, nil
}

func (r *Repository) CreateSubject(ctx context.Context, input SubjectInput) (*Subject, error) {
	item := &Subject{ID: uuid.New(), Name: input.Name, Code: input.Code, Topics: []Topic{}}
	if err := r.db.WithContext(ctx).Table("subjects").Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *Repository) CreateTopic(ctx context.Context, input TopicInput) (*Topic, error) {
	item := &Topic{ID: uuid.New(), Name: input.Name}
	if err := r.db.WithContext(ctx).Table("topics").Create(map[string]any{
		"id": item.ID, "subject_id": input.SubjectID, "name": item.Name,
	}).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *Repository) SubjectExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("subjects").Where("id = ?", id).Count(&count).Error
	return count == 1, err
}

func (r *Repository) ListTasks(ctx context.Context, p pagination.Params) ([]Task, int64, error) {
	query := r.db.WithContext(ctx).Model(&Task{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var tasks []Task
	err := query.Order("created_at DESC").Offset(p.Offset()).Limit(p.Limit).Find(&tasks).Error
	return tasks, total, err
}

func (r *Repository) SaveTask(ctx context.Context, item *Task) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) UpdateTask(ctx context.Context, id uuid.UUID, input TaskInput) error {
	result := r.db.WithContext(ctx).Model(&Task{}).Where("id = ?", id).Updates(map[string]any{
		"subject_id": input.SubjectID, "topic_id": input.TopicID, "task_number": input.TaskNumber,
		"type": input.Type, "question": input.Question, "options": input.Options,
		"correct_answer": input.CorrectAnswer, "explanation": input.Explanation, "source": input.Source,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) TopicMatchesSubject(ctx context.Context, topicID, subjectID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("topics").Where("id = ? AND subject_id = ?", topicID, subjectID).Count(&count).Error
	return count == 1, err
}

func (r *Repository) GetTask(ctx context.Context, id uuid.UUID) (*Task, error) {
	var item Task
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
