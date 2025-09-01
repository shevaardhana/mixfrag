package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"mixfrag/database"
	"mixfrag/models"

	"github.com/gin-gonic/gin"
)

type ParfumeResponse struct {
	ID                int                    `json:"id"`
	Name              string                 `json:"name"`
	TotalMl           string                 `json:"total_ml"`
	CategoryParfumeId int                    `json:"category_parfume_id"`
	TotalOilDrop      int                    `json:"total_oil_drop"`
	TotalOil          int                    `json:"total_oil"`
	TotalParfumeBase  int                    `json:"total_parfume_base"`
	CreatedAt         time.Time              `json:"created_at"`
	CreatedBy         string                 `json:"created_by"`
	ModifiedAt        *time.Time             `json:"modified_at"`
	ModifiedBy        *string                `json:"modified_by"`
	Details           []models.ParfumeDetail `json:"details"`
}

type CreateParfumeRequest struct {
	Name              string                 `json:"name"`
	TotalMl           string                 `json:"total_ml"`
	CategoryParfumeId int                    `json:"category_parfume_id"`
	TotalOilDrop      int                    `json:"total_oil_drop"`
	TotalOil          int                    `json:"total_oil"`
	TotalParfumeBase  int                    `json:"total_parfume_base"`
	CreatedBy         string                 `json:"created_by"`
	Details           []models.ParfumeDetail `json:"details"`
}

type CreateParfumeResponse struct {
	ParfumeID int `json:"parfume_id"`
}

// ✅ GET /api/parfumes
func GetParfumes(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, name, total_ml, category_parfume_id, total_oil_drop, total_oil, total_parfume_base, created_at, created_by, modified_at, modified_by
		FROM parfumes`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var parfumes []ParfumeResponse

	for rows.Next() {
		var p ParfumeResponse
		err := rows.Scan(
			&p.ID, &p.Name, &p.TotalMl, &p.CategoryParfumeId,
			&p.TotalOilDrop, &p.TotalOil, &p.TotalParfumeBase,
			&p.CreatedAt, &p.CreatedBy, &p.ModifiedAt, &p.ModifiedBy,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Ambil details
		detailRows, err := database.DB.Query(`
			SELECT id, parfume_id, oil_id, total_drop, rasio, weight_total, created_at, created_by, modified_at, modified_by
			FROM parfume_details WHERE parfume_id = $1`, p.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var details []models.ParfumeDetail
		for detailRows.Next() {
			var d models.ParfumeDetail
			err := detailRows.Scan(
				&d.ID, &d.ParfumeId, &d.OilId, &d.TotalDrop,
				&d.Rasio, &d.WeightTotal, &d.CreatedAt, &d.CreatedBy,
				&d.ModifiedAt, &d.ModifiedBy,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			details = append(details, d)
		}
		detailRows.Close()

		p.Details = details
		parfumes = append(parfumes, p)
	}

	c.JSON(http.StatusOK, parfumes)
}

// ✅ POST /api/parfumes
func CreateParfume(c *gin.Context) {
	var req CreateParfumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	// Insert parfume
	queryParfume := `
		INSERT INTO parfumes (name, total_ml, category_parfume_id, total_oil_drop, total_oil, total_parfume_base, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	var parfumeID int
	err = tx.QueryRow(queryParfume,
		req.Name, req.TotalMl, req.CategoryParfumeId,
		req.TotalOilDrop, req.TotalOil, req.TotalParfumeBase,
		time.Now(), req.CreatedBy,
	).Scan(&parfumeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert parfume failed: " + err.Error()})
		return
	}

	// Insert details
	queryDetail := `
		INSERT INTO parfume_details (parfume_id, oil_id, total_drop, rasio, weight_total, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	for _, d := range req.Details {
		_, err = tx.Exec(queryDetail,
			parfumeID, d.OilId, d.TotalDrop, d.Rasio, d.WeightTotal, time.Now(), req.CreatedBy,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "insert detail failed: " + err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateParfumeResponse{ParfumeID: parfumeID})
}

// ✅ GET /api/parfumes/:id
func GetParfumeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	row := database.DB.QueryRow(`
		SELECT id, name, total_ml, category_parfume_id, total_oil_drop, total_oil, total_parfume_base, created_at, created_by
		FROM parfumes WHERE id = $1`, id)

	var parfume ParfumeResponse
	err = row.Scan(
		&parfume.ID, &parfume.Name, &parfume.TotalMl, &parfume.CategoryParfumeId,
		&parfume.TotalOilDrop, &parfume.TotalOil, &parfume.TotalParfumeBase,
		&parfume.CreatedAt, &parfume.CreatedBy,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "parfume not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ambil detail
	details, _ := database.DB.Query(`
		SELECT id, parfume_id, oil_id, total_drop, rasio, weight_total, created_at, created_by, modified_at, modified_by
		FROM parfume_details WHERE parfume_id = $1`, id)
	defer details.Close()

	for details.Next() {
		var d models.ParfumeDetail
		_ = details.Scan(
			&d.ID, &d.ParfumeId, &d.OilId, &d.TotalDrop,
			&d.Rasio, &d.WeightTotal, &d.CreatedAt, &d.CreatedBy,
			&d.ModifiedAt, &d.ModifiedBy,
		)
		parfume.Details = append(parfume.Details, d)
	}

	c.JSON(http.StatusOK, parfume)
}

// ✅ DELETE /api/parfumes/:id
func DeleteParfume(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	sqlStatement := `DELETE FROM parfumes WHERE id=$1`
	res, err := database.DB.Exec(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete parfum"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "parfume not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "parfume deleted"})
}
