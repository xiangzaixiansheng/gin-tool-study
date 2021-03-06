package service

import (
	"gin-tool-study/model"
	"gin-tool-study/pkg/e"
	"gin-tool-study/serializer"
	"strconv"

	logging "github.com/sirupsen/logrus"
)

type AddressService struct {
	Name    string `form:"name" json:"name"`
	Phone   string `form:"phone" json:"phone"`
	Address string `form:"address" json:"address"`
}

func (service AddressService) Create(id uint) serializer.Response {
	//插入数据
	code := e.SUCCESS
	address := model.Address{
		UserID:  id,
		Name:    service.Name,
		Phone:   service.Phone,
		Address: service.Address,
	}
	err := model.DB.Create(&address).Error
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//返回数据库中这个用户最新的地址信息
	var addresses []model.Address
	if err = model.DB.Model(model.Address{}).Where("user_id = ?", id).Order("created_at DESC").Find(&addresses).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildAddresses(addresses),
		Msg:    e.GetMsg(code),
	}

}

func (service AddressService) List(id uint) serializer.Response {
	code := e.SUCCESS
	var addresses []model.Address
	if err := model.DB.Model(model.Address{}).Where("user_id", id).Order("create_time DESC").Find(&addresses).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildAddresses(addresses),
		Msg:    e.GetMsg(code),
	}

}

func (service AddressService) Delete(id string) serializer.Response {
	//可以先找出来对应id的数据
	code := e.SUCCESS
	var address model.Address
	if err := model.DB.Where("id = ?", id).Find(&address).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	err := model.DB.Delete(&address).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}
func (service AddressService) Update(uid uint, aid string) serializer.Response {
	code := e.SUCCESS
	address := model.Address{
		Address: service.Address,
		Phone:   service.Phone,
		Name:    service.Name,
		UserID:  uid,
	}
	Id, _ := strconv.Atoi(aid)
	address.ID = uint(Id)

	//字段全部更新
	err := model.DB.Save(&address).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	var addresses []model.Address
	err = model.DB.Model(&model.Address{}).Where("user_id = ?", uid).Order("create_time DESC").Find(&addresses).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildAddresses(addresses),
	}

}
