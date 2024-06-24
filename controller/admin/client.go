package admin

import (
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllClients(c *gin.Context) {
	db := database.DB

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Println("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	clients := []model.ClientMiniInfo{}
	err := db.Table("clients").Find(&clients).Error
	if err != nil {
		fmt.Errorf("error: could not find all clients | %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	clientIDs := []uint64{}
	clientIDMap := map[uint64]int{}

	for index, client := range clients {
		clientIDs = append(clientIDs, client.ID)
		clientIDMap[client.ID] = index
	}

	lastDietDates := []model.DietHistory{}

	err = db.Table("diet_histories AS d").
		Select("d.client_id, d.date").
		Joins("JOIN (SELECT client_id, MAX(date) as max_date FROM diet_histories WHERE client_id IN (?) GROUP BY client_id) AS sub ON d.client_id = sub.client_id AND d.date = sub.max_date", clientIDs).
		Find(&lastDietDates).Error
	if err != nil {
		fmt.Errorf("error: could not find the last diet dates client_id: %v | %v", clientIDs, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	for _, res := range lastDietDates {
		index := clientIDMap[res.ClientID]
		clients[index].LastDietDate = res.Date
	}

	c.JSON(http.StatusOK, gin.H{"clients": clients})
	return
}