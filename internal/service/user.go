package service

import (
	"ServerServing/da/mysql/da_models"
	SErr "ServerServing/err"
	"ServerServing/internal/dal"
	"ServerServing/internal/internal_models"
	"github.com/gin-gonic/gin"
)

type UsersService struct{}

func GetUsersService() *UsersService {
	return &UsersService{}
}

func (s *UsersService) Create(c *gin.Context, name string, pwd string) (*da_models.User, *SErr.APIErr) {
	user := &da_models.User{
		Name: name,
		Pwd:  pwd,
	}
	sErr := dal.GetUserDal().Create(user)
	if sErr != nil {
		return nil, sErr
	}
	sErr = GetSessionsService().Create(c, name, pwd)
	if sErr != nil {
		return nil, sErr
	}
	return user, nil
}

func (s *UsersService) Info(c *gin.Context, userID int) (*da_models.User, *SErr.APIErr) {
	user, sErr := dal.GetUserDal().GetByID(userID)
	if sErr != nil {
		return nil, sErr
	}
	return user, sErr
}

func (s *UsersService) CheckPwd(c *gin.Context, name string, pwd string) (*da_models.User, *SErr.APIErr) {
	user, sErr := dal.GetUserDal().GetByName(name)
	if sErr != nil {
		return nil, sErr
	}
	if user.Pwd != pwd {
		return nil, SErr.WrongPwdErr
	}
	return user, sErr
}

func (s *UsersService) Infos(c *gin.Context, from, size int, searchKeyword *string) ([]*da_models.User, int, *SErr.APIErr) {
	var users []*da_models.User
	var totalCount int
	var sErr *SErr.APIErr
	if searchKeyword == nil {
		users, totalCount, sErr = dal.GetUserDal().List(from, size)
	} else {
		users, totalCount, sErr = dal.GetUserDal().SearchByName(*searchKeyword, from, size)
	}
	return users, totalCount, sErr
}

func (s *UsersService) Update(c *gin.Context, targetID int, updateReq *internal_models.UsersUpdateRequest) *SErr.APIErr {
	userID, sErr := GetSessionsService().GetUserID(c)
	if sErr != nil {
		return SErr.NeedLoginErr.CustomMessageF("更新用户信息前，必须要登录。")
	}
	// 只有管理员可以更新其他人的数据，只有本人可以修改本人的数据
	userInfo, sErr := s.Info(c, userID)
	if sErr != nil {
		return sErr
	}
	if !userInfo.Admin && userID != targetID {
		return SErr.AdminOnlyActionErr
	}
	return dal.GetUserDal().Update(targetID, updateReq)
}

func (s *UsersService) IsAdmin(c *gin.Context, targetID int) (bool, *SErr.APIErr) {
	userInfo, err := s.Info(c, targetID)
	if err != nil {
		return false, err
	}
	return userInfo.Admin, nil
}