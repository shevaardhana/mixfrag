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

// GET /api/essential-oils
func GetEssentialOils(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, name, category_smell_id, category_note_id, rank_note,
		       created_at, created_by, modified_at, modified_by
		FROM essential_oils`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data essential oil"})
		return
	}
	defer rows.Close()

	var oils []models.EssentialOil
	for rows.Next() {
		var eo models.EssentialOil
		var modifiedAt sql.NullTime
		var modifiedBy sql.NullString

		if err := rows.Scan(&eo.ID, &eo.Name, &eo.CategorySmellId, &eo.CategoryNoteId, &eo.RankNote,
			&eo.CreatedAt, &eo.CreatedBy, &modifiedAt, &modifiedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data essential oil", "message": err.Error()})
			return
		}

		if modifiedAt.Valid {
			eo.ModifiedAt = &modifiedAt.Time
		}
		if modifiedBy.Valid {
			eo.ModifiedBy = &modifiedBy.String
		}

		oils = append(oils, eo)
	}

	c.JSON(http.StatusOK, oils)
}

// POST /api/essential-oils
func CreateEssentialOil(c *gin.Context) {
	var input models.EssentialOil
	if err := c.ShouldBindJSON(&input); err != nil || input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama essential oil harus diisi"})
		return
	}

	input.CreatedAt = time.Now()
	if user, exists := c.Get("username"); exists {
		input.CreatedBy = user.(string)
	} else {
		input.CreatedBy = "unknown"
	}

	sqlStatement := `
		INSERT INTO essential_oils (name, category_smell_id, category_note_id, rank_note, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int
	err := database.DB.QueryRow(sqlStatement,
		input.Name, input.CategorySmellId, input.CategoryNoteId, input.RankNote, input.CreatedAt, input.CreatedBy,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":              id,
		"name":            input.Name,
		"categorySmellId": input.CategorySmellId,
		"categoryNoteId":  input.CategoryNoteId,
		"rankNote":        input.RankNote,
		"createdAt":       input.CreatedAt,
		"createdBy":       input.CreatedBy,
	})
}

// GET /api/essential-oils/:id
func GetEssentialOil(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var eo models.EssentialOil
	var modifiedAt sql.NullTime
	var modifiedBy sql.NullString

	sqlStatement := `
		SELECT id, name, category_smell_id, category_note_id, rank_note,
		       created_at, created_by, modified_at, modified_by
		FROM essential_oils WHERE id=$1`
	err := database.DB.QueryRow(sqlStatement, id).Scan(
		&eo.ID, &eo.Name, &eo.CategorySmellId, &eo.CategoryNoteId, &eo.RankNote,
		&eo.CreatedAt, &eo.CreatedBy, &modifiedAt, &modifiedBy,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Essential oil tidak ditemukan"})
		return
	}

	if modifiedAt.Valid {
		eo.ModifiedAt = &modifiedAt.Time
	}
	if modifiedBy.Valid {
		eo.ModifiedBy = &modifiedBy.String
	}

	c.JSON(http.StatusOK, eo)
}

// DELETE /api/essential-oils/:id
func DeleteEssentialOil(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	sqlStatement := `DELETE FROM essential_oils WHERE id=$1`
	res, err := database.DB.Exec(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus essential oil"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Essential oil tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Essential oil berhasil dihapus"})
}
