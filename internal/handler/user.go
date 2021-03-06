package handler

import (
	daModels "ServerServing/da/mysql/da_models"
	SErr "ServerServing/err"
	models "ServerServing/internal/internal_models"
	"ServerServing/internal/service"
	"ServerServing/util"
	"github.com/gin-gonic/gin"
	"log"
)

type UsersHandler struct{}

func GetUserHandler() *UsersHandler {
	return &UsersHandler{}
}

// Create
// @Summary 注册用户
// @Tags user
// @Produce json
// @Router /api/v1/users/ [post]
// @Param userCreateRequest body internal_models.UsersCreateRequest true "userCreateRequest"
// @Success 200 {object} internal_models.UsersCreateResponse
func (h UsersHandler) Create(c *gin.Context) (interface{}, *SErr.APIErr) {
	req := &models.UsersCreateRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		return nil, SErr.BadRequestErr
	}
	s := service.GetUsersService()
	token, sErr := s.Create(c, req.Name, req.Pwd)
	if sErr != nil {
		return nil, sErr
	}
	return &models.UsersCreateResponse{
		Token: token,
	}, nil
}

// Info
// @Summary 获取单个用户信息
// @Tags user
// @Produce json
// @Router /api/v1/users/{id} [get]
// @param id path int true "id"
// @Param x-token header string false "x-token"
// @Success 200 {object} internal_models.UsersInfoResponse
func (h UsersHandler) Info(c *gin.Context) (interface{}, *SErr.APIErr) {
	userIDStr := c.Param("id")
	userID, err := util.ParseInt(userIDStr)
	if err != nil {
		return nil, SErr.BadRequestErr
	}
	sessionSvc := service.GetSessionsService()
	loggedUserID, loggedErr := sessionSvc.GetUserID(c)
	if userID <= 0 && loggedErr != nil {
		return nil, SErr.BadRequestErr.CustomMessageF("待查询的用户ID <= 0，目标用户ID为%d", userID)
	}
	if userID <= 0 {
		userID = loggedUserID
	}
	log.Printf("UsersHandler querying user info for id=[%d]", userID)
	user, sErr := service.GetUsersService().Info(c, userID)
	if sErr != nil {
		return nil, sErr
	}
	return &models.UsersInfoResponse{
		User: h.packUser(user),
	}, nil
}

// Infos
// @Summary 获取多个用户信息，可以添加关键字对姓名搜索。
// @Tags user
// @Produce json
// @Router /api/v1/users/ [get]
// @Param userInfoRequest query internal_models.UsersInfosRequest true "userInfosRequest"
// @Success 200 {object} internal_models.UsersInfosResponse
func (h UsersHandler) Infos(c *gin.Context) (interface{}, *SErr.APIErr) {
	req := &models.UsersInfosRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		return nil, SErr.BadRequestErr
	}
	users, totalCount, sErr := service.GetUsersService().Infos(c, req.From, req.Size, req.SearchKeyword)
	if sErr != nil {
		return nil, sErr
	}
	return &models.UsersInfosResponse{
		Infos:      h.packUsers(users),
		TotalCount: totalCount,
	}, nil
}

// Update
// @Summary 修改用户信息
// @Tags user
// @Produce json
// @Router /api/v1/users/{id} [put]
// @param id path uint true "id"
// @Param updateRequest body internal_models.UsersUpdateRequest true "updateRequest"
// @Param x-token header string true "x-token"
// @Success 200 {object} internal_models.UsersUpdateResponse
func (h UsersHandler) Update(c *gin.Context) (interface{}, *SErr.APIErr) {
	targetIDStr := c.Param("id")
	targetID, err := util.ParseInt(targetIDStr)
	if err != nil {
		return nil, SErr.BadRequestErr
	}
	if targetID <= 0 {
		return nil, SErr.BadRequestErr.CustomMessageF("要修改的用户ID <= 0, 目标的userID为%d", targetID)
	}
	updateReq := &models.UsersUpdateRequest{}
	err = c.Bind(updateReq)
	if err != nil {
		return nil, SErr.BadRequestErr
	}
	sErr := service.GetUsersService().Update(c, targetID, updateReq)
	if sErr != nil {
		return nil, sErr
	}
	return &models.UsersUpdateResponse{}, nil
}

func (h UsersHandler) packUser(user *daModels.User) *models.User {
	return &models.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
		Name:      user.Name,
		Pwd:       user.Pwd,
		Admin:     user.Admin,
	}
}

func (h UsersHandler) packUsers(users []*daModels.User) []*models.User {
	packed := make([]*models.User, 0, len(users))
	for _, user := range users {
		packed = append(packed, h.packUser(user))
	}
	return packed
}
