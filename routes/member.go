package routes

import (
	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/gin-gonic/gin"
)

//GetProfile to view list article
func GetProfile(c *gin.Context) {
	userID := uint(c.MustGet("jwt_user_id").(float64))
	var user models.Member

	transactions := config.DB.Where("id = ?", userID).Preload("Transactions", "member_id = ?", userID).Find(&user)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   transactions,
	})
}

//MemberIndex to view member list
func MemberIndex(c *gin.Context) {
	var user models.Member
	memberRole := false

	members := config.DB.Where("role = ?", memberRole).Find(&user)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   members,
	})
}

//MemberProfile to view member detail
func MemberProfile(c *gin.Context) {
	var user models.Member
	userID := c.Param("id")
	memberRole := false

	member := config.DB.Select([]string{"id", "username", "full_name", "phone", "email", "address", "avatar"}).Where("role = ?", memberRole).Preload("Transactions", "member_id = ?", userID).First(&user, userID)

	if member.RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Member not found",
		})
		c.Abort()
		return
	}

	// fmt.Println(member)
	// var transaction models.Transaction
	// trans := config.DB.Find(&transaction, "member_id = ?", member.Value)

	c.JSON(200, gin.H{
		"status": 200,
		"member": member,
		"id":     "learn",
	})
}
