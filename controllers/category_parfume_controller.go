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

// GET /api/categories
func GetCategorieParfumes(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, name, desc, created_at, created_by, modified_at, modified_by 
		FROM category_parfume`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kategori parfume"})
		return
	}
	defer rows.Close()

	var categories []models.CategoryParfume
	for rows.Next() {
		var cat models.CategoryParfume
		var desc sql.NullString
		var modifiedAt sql.NullTime
		var modifiedBy sql.NullString

		if err := rows.Scan(&cat.ID, &cat.Name, &desc, &cat.CreatedAt, &cat.CreatedBy, &modifiedAt, &modifiedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data kategori", "message": err.Error()})
			return
		}

		if desc.Valid {
			cat.Desc = &desc.String
		}
		if modifiedAt.Valid {
			cat.ModifiedAt = &modifiedAt.Time
		}
		if modifiedBy.Valid {
			cat.ModifiedBy = &modifiedBy.String
		}

		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}

// POST /api/categories
func CreateCategoryParfume(c *gin.Context) {
	var input models.CategoryParfume
	if err := c.ShouldBindJSON(&input); err != nil || input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori parfume harus diisi"})
		return
	}

	input.CreatedAt = time.Now()
	if user, exists := c.Get("username"); exists {
		input.CreatedBy = user.(string)
	} else {
		input.CreatedBy = "unknown"
	}

	sqlStatement := `INSERT INTO category_parfume (name, desc, created_at, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := database.DB.QueryRow(sqlStatement, input.Name, input.Desc, input.CreatedAt, input.CreatedBy).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        id,
		"name":      input.Name,
		"desc":      input.Desc,
		"createdAt": input.CreatedAt,
		"createdBy": input.CreatedBy,
	})
}

// GET /api/categories/:id
func GetCategoryParfume(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var cat models.CategoryParfume
	var desc sql.NullString
	var modifiedAt sql.NullTime
	var modifiedBy sql.NullString

	sqlStatement := `SELECT id, name, desc, created_at, created_by, modified_at, modified_by FROM category_parfume WHERE id=$1`
	err := database.DB.QueryRow(sqlStatement, id).Scan(
		&cat.ID, &cat.Name, &desc, &cat.CreatedAt, &cat.CreatedBy, &modifiedAt, &modifiedBy,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori parfume tidak ditemukan"})
		return
	}

	if desc.Valid {
		cat.Desc = &desc.String
	}
	if modifiedAt.Valid {
		cat.ModifiedAt = &modifiedAt.Time
	}
	if modifiedBy.Valid {
		cat.ModifiedBy = &modifiedBy.String
	}

	c.JSON(http.StatusOK, cat)
}

// DELETE /api/categories/:id
func DeleteCategoryParfume(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	sqlStatement := `DELETE FROM category_parfume WHERE id=$1`
	res, err := database.DB.Exec(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori parfume"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori parfume tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori parfume berhasil dihapus"})
}
