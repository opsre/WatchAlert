package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"
)

type userService struct {
	ctx *ctx.Context
}

type InterUserService interface {
	Search(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	Login(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Register(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	ChangePass(req interface{}) (interface{}, interface{})
}

func newInterUserService(ctx *ctx.Context) InterUserService {
	return &userService{
		ctx: ctx,
	}
}

func (us userService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MemberQuery)
	data, err := us.ctx.DB.User().Search(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (us userService) List(req interface{}) (interface{}, interface{}) {
	data, err := us.ctx.DB.User().List()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (us userService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MemberQuery)
	data, _, err := us.ctx.DB.User().Get(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (us userService) Login(req interface{}) (interface{}, interface{}) {
	r := req.(*models.Member)
	r.Password = tools.GenerateHashPassword(r.Password)

	q := models.MemberQuery{
		UserName: r.UserName,
	}
	data, _, err := us.ctx.DB.User().Get(q)
	if err != nil {
		return nil, err
	}

	switch data.CreateBy {
	case "LDAP":
		if global.Config.Ldap.Enabled {
			err := LdapService.Login(r.UserName, r.Password)
			if err != nil {
				logc.Error(us.ctx.Ctx, fmt.Sprintf("LDAP 用户登陆失败, err: %s", err.Error()))
				return nil, fmt.Errorf("LDAP 用户登陆失败, err: %s", err.Error())
			}
		} else {
			logc.Error(us.ctx.Ctx, "请先开启 LDAP 功能!")
			return nil, fmt.Errorf("请先开启 LDAP 功能!")
		}
	default:
		if data.Password != r.Password {
			return nil, fmt.Errorf("密码错误")
		}
	}

	r.UserId = data.UserId
	tokenData, err := tools.GenerateToken(r.UserId, r.UserName, r.Password)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(global.Config.Jwt.Expire) * time.Second
	us.ctx.Cache.Cache().SetKey("uid-"+data.UserId, tools.JsonMarshal(r), duration)

	return tokenData, nil
}

func (us userService) Register(req interface{}) (interface{}, interface{}) {
	r := req.(*models.Member)

	q := models.MemberQuery{UserName: r.UserName}
	_, ok, _ := us.ctx.DB.User().Get(q)
	if ok {
		return nil, fmt.Errorf("用户已存在")
	}

	// 在初始化admin用户时会固定一个userid，所以这里需要做一下判断；
	if r.UserId == "" {
		r.UserId = tools.RandUid()
	}

	r.Password = tools.GenerateHashPassword(r.Password)
	r.CreateAt = time.Now().Unix()

	if r.CreateBy == "" {
		r.CreateBy = "system"
	}

	err := us.ctx.DB.User().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (us userService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.Member)
	var dbData models.Member

	db := us.ctx.DB.DB().Model(models.Member{})
	db.Where("user_id = ?", r.UserId).First(&dbData)

	if r.Password == "" {
		r.Password = dbData.Password
	} else {
		r.Password = tools.GenerateHashPassword(r.Password)
	}
	err := us.ctx.DB.User().Update(*r)
	if err != nil {
		return nil, err
	}

	us.ctx.DB.User().ChangeCache(r.UserId)

	return nil, nil
}

func (us userService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MemberQuery)
	err := us.ctx.DB.User().Delete(*r)
	if err != nil {
		return nil, err
	}

	us.ctx.DB.User().ChangeCache(r.UserId)

	return nil, nil
}

func (us userService) ChangePass(req interface{}) (interface{}, interface{}) {
	r := req.(*models.Member)

	arr := md5.Sum([]byte(r.Password))
	hashPassword := hex.EncodeToString(arr[:])
	r.Password = hashPassword

	err := us.ctx.DB.User().ChangePass(*r)
	if err != nil {
		return nil, err
	}

	us.ctx.DB.User().ChangeCache(r.UserId)

	return nil, nil
}
