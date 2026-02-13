package repository

import "medical-webhook/internal/domain/line/entity"

type TicketCategoryRepository interface {
	GetTicketCategories() ([]entity.TicketCategory, error)
	FindCategoryByName(name string) (*entity.TicketCategory, error)
	CreateCategory(category *entity.TicketCategory) error
}
