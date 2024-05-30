package admin

import (
	"errors"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func SaveProfileByDietitian(c *gin.Context) {
	db := database.DB

	req := model.Client{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := model.Client{}
	err := db.Where("email = ? and deleted_at IS NULL", req.Email).Find(&client).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		// cannot find your record, please contact us
		c.JSON(http.StatusNotFound, gin.H{"error": errors.New("record_not_found")})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//err := db.Save(client)
}
