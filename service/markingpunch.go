package service

import (
	"encoding/json"
	"errors"
	"order_food/global"
	"order_food/model/sea"
	"order_food/model/sea/request"
)

var MkS = MarkingpunchService{}

type MarkingpunchService struct {
	MarkingConf
	Markingpunch
	MarkingTmp
	MarkingBill
}

type MarkingConf struct{}

func (*MarkingConf) Add(req request.MarkingConfAddReq) error {
	conf := sea.MarkingConf{}
	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &conf)

	err := global.GVA_DB.Save(&conf).Error
	return err
}

func (*MarkingConf) Update(req request.MarkingConfUpdateReq) error {
	conf := sea.MarkingConf{}
	_ = global.GVA_DB.Where("id = ?", req.Id).First(&conf).Error
	if conf.ID == 0 {
		return errors.New("配置错误")
	}

	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &conf)

	err := global.GVA_DB.Save(&conf).Error
	return err
}

func (*MarkingConf) Info(id int) (*sea.MarkingConf, error) {
	conf := sea.MarkingConf{}
	err := global.GVA_DB.Where("id = ?", id).First(&conf).Error
	if conf.ID == 0 {
		return nil, errors.New("配置错误")
	}

	return &conf, err
}

func (*MarkingConf) List(req request.MarkingConfListReq) ([]sea.MarkingConf, int64, error) {
	confList := []sea.MarkingConf{}
	confDb := global.GVA_DB.Model(&sea.MarkingConf{})
	if req.DeviceType != nil {
		confDb = confDb.Where("device_type = ?", req.DeviceType)
	}
	var count int64
	err := confDb.Count(&count).Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&confList).Error
	if err != nil {
		return nil, 0, errors.New("获取配置列表错误")
	}

	return confList, count, err
}

type Markingpunch struct{}

func (*Markingpunch) Add(req request.MarkingPunchAddReq) error {
	conf := sea.MarkingPunch{}
	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &conf)

	err := global.GVA_DB.Save(&conf).Error
	return err
}

func (*Markingpunch) Update(req request.MarkingPunchUpdateReq) error {
	info := sea.MarkingPunch{}
	_ = global.GVA_DB.Where("id = ?", req.Id).First(&info).Error
	if info.ID == 0 {
		return errors.New("配置错误")
	}

	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &info)

	err := global.GVA_DB.Save(&info).Error
	return err
}

func (*Markingpunch) Info(id int) (*sea.MarkingPunch, error) {
	info := sea.MarkingPunch{}
	_ = global.GVA_DB.Where("id = ?", id).First(&info).Error
	if info.ID == 0 {
		return nil, errors.New("配置错误")
	}

	return &info, nil
}

func (*Markingpunch) List(req request.MarkingPunchListReq) ([]sea.MarkingPunch, int64, error) {
	list := []sea.MarkingPunch{}
	markingPunchDb := global.GVA_DB.Model(&sea.MarkingPunch{})
	if req.MarkTypeId != 0 {
		markingPunchDb = markingPunchDb.Where("mark_type_id = ?", req.MarkTypeId)
	}
	var count int64
	err := markingPunchDb.Count(&count).Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error

	return list, count, err
}

type MarkingTmp struct{}

func (*MarkingTmp) Add(req request.MarkingTmpAddReq) error {
	conf := sea.MarkingTmp{}
	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &conf)

	err := global.GVA_DB.Save(&conf).Error
	return err
}

func (*MarkingTmp) Update(req request.MarkingTmpUpdateReq) error {
	info := sea.MarkingTmp{}
	_ = global.GVA_DB.Where("id = ?", req.Id).First(&info).Error
	if info.ID == 0 {
		return errors.New("配置错误")
	}

	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &info)

	err := global.GVA_DB.Save(&info).Error
	return err
}
func (*MarkingTmp) Info(id int) (*sea.MarkingTmp, error) {
	info := sea.MarkingTmp{}
	_ = global.GVA_DB.Where("id = ?", id).First(&info).Error
	if info.ID == 0 {
		return nil, errors.New("配置错误")
	}

	return &info, nil
}

func (*MarkingTmp) List(req request.MarkingTmpListReq) ([]sea.MarkingTmp, int64, error) {
	list := []sea.MarkingTmp{}
	markingPunchDb := global.GVA_DB.Model(&sea.Bills{})
	if req.MarkTypeId != 0 {
		markingPunchDb = markingPunchDb.Where("mark_type_id = ?", req.MarkTypeId)
	}
	var count int64
	err := markingPunchDb.Count(&count).Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error

	return list, count, err
}

type MarkingBill struct{}

func (*MarkingBill) Add(req request.MarkingBillAddReq) error {
	conf := sea.Bills{}
	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &conf)

	err := global.GVA_DB.Save(&conf).Error
	return err
}

func (*MarkingBill) Update(req request.MarkingBillUpdateReq) error {
	info := sea.Bills{}
	_ = global.GVA_DB.Where("id = ?", req.Id).First(&info).Error
	if info.ID == 0 {
		return errors.New("配置错误")
	}

	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &info)

	err := global.GVA_DB.Save(&info).Error
	return err
}
func (*MarkingBill) Info(id int) (*sea.Bills, error) {
	info := sea.Bills{}
	_ = global.GVA_DB.Where("id = ?", id).First(&info).Error
	if info.ID == 0 {
		return nil, errors.New("配置错误")
	}

	return &info, nil
}

func (*MarkingBill) List(req request.MarkingBillListReq) ([]sea.Bills, int64, error) {
	list := []sea.Bills{}
	markingBillDb := global.GVA_DB.Model(&sea.Bills{})
	if req.MarkingPunchId != 0 {
		markingBillDb = markingBillDb.Where("marking_punch_id = ?", req.MarkingPunchId)
	}
	if req.BillType != 0 {
		markingBillDb = markingBillDb.Where("bill_type = ?", req.MarkingPunchId)
	}
	var count int64
	err := markingBillDb.Count(&count).Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error

	return list, count, err
}
