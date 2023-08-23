package logic

import (
	"chapter4.1.bluebell/dao/mysql"
	"chapter4.1.bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查到数据库 查找到所有的community 并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
