package model

import (
	"blog/library"
	"blog/library/log"
	"blog/orm"
	"errors"
	"strconv"
)

type Admin struct {
	Id string	`from:"id"`
	UserName string	`from:"user_name"`
	Name string	`from:"name"`
	Passwd string	`from:"passwd"`
	Salt string	`from:"salt"`
	IsDel string	`from:"is_del"`
	CreateTime string `from:"create_time"`
	UpdateTIme string	`from:"update_time"`
	DeleteTime string	`from:"delete_time"`
}

func SelectAllAdmin() []Admin {
	Admins := make([]Admin, 0, 0)
	err := orm.DB().Table(AdminTable).Where("is_del", 0).GetModel(&Admins)
	if err != nil {
		log.New().WithField("err", err).Error("GetModel error")
		return nil
	}

	return Admins
}

func FindAdmin(passwd string) (int, error) {
	Admins := SelectAllAdmin()
	if Admins == nil {
		return 0, errors.New("no superman")
	}
	for _, item := range Admins {
		if passwdEncryption(passwd, item.Salt, item.UpdateTIme) == item.Passwd {
			uid, _ := strconv.Atoi(item.Id)
			return uid, nil
		}
	}

	return 0, errors.New("no superman")
}

func passwdEncryption(passwd, salt, update_time string) string {
	return library.Md5Encryption(update_time + library.Md5Encryption(passwd + salt))
}