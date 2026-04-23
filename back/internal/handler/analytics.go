package handler

import (...)

// ProjectStatsResponse - структура ответа со статистикой
type ProjectStatsResponse struct {
    TotalAmount     float64 `json:"total_amount"`
    DonorsCount     int64   `json:"donors_count"`
    LikesCount      int64   `json:"likes_count"`
    CommentsCount   int64   `json:"comments_count"`
    AvgDonation     float64 `json:"avg_donation"`
    ProgressPercent float64 `json:"progress_percent"`
    DaysRemaining   int     `json:"days_remaining"`
}

// GetProjectStats - handler для получения статистики проекта
func (h *Handler) GetProjectStats(c *gin.Context) {
    // реализация
}
