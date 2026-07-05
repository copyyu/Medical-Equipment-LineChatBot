package usecase

import (
	"context"
	"testing"

	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/entity"
	mock_repository "medical-webhook/internal/mocks/repository"

	"github.com/stretchr/testify/assert"
)

func TestGetTicketList_NormalizesMissingPagination(t *testing.T) {
	repo := mock_repository.NewMockTicketRepository(t)
	// page=0/limit=0 must be normalized to page=1/limit=defaultPageLimit
	repo.EXPECT().
		GetAllTickets(1, defaultPageLimit, "", "", "", "", "").
		Return([]entity.Ticket{}, int64(0), nil).
		Once()

	uc := &TicketUseCase{ticketRepo: repo}
	resp, err := uc.GetTicketList(context.Background(), dto.TicketListRequest{Page: 0, Limit: 0})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, defaultPageLimit, resp.Pagination.Limit)
}

func TestGetTicketList_CapsExcessiveLimit(t *testing.T) {
	repo := mock_repository.NewMockTicketRepository(t)
	repo.EXPECT().
		GetAllTickets(2, maxPageLimit, "", "", "", "", "").
		Return([]entity.Ticket{}, int64(0), nil).
		Once()

	uc := &TicketUseCase{ticketRepo: repo}
	resp, err := uc.GetTicketList(context.Background(), dto.TicketListRequest{Page: 2, Limit: 9999})

	assert.NoError(t, err)
	assert.Equal(t, maxPageLimit, resp.Pagination.Limit)
}

func TestGetTicketByID_NotFoundReturnsNil(t *testing.T) {
	repo := mock_repository.NewMockTicketRepository(t)
	repo.EXPECT().FindTicketByID(uint(42)).Return(nil, nil).Once()

	uc := &TicketUseCase{ticketRepo: repo}
	resp, err := uc.GetTicketByID(context.Background(), 42)

	assert.NoError(t, err)
	assert.Nil(t, resp)
}
