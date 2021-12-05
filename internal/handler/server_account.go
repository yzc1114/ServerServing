package handler

import (
	"ServerServing/api/format"
	SErr "ServerServing/err"
	models "ServerServing/internal/internal_models"
	"ServerServing/internal/service"
	"github.com/gin-gonic/gin"
)

type ServerAccountHandler struct {}

func GetServerAccountHandler() ServerAccountHandler {
	return ServerAccountHandler{}
}

// Create
// @Summary 创建一个服务器的sudo账号。
// @Tags server_account
// @Produce json
// @Router /api/v1/servers/accounts [post]
// @Param serverAccountCreateRequest body internal_models.ServerAccountCreateRequest true "serverAccountCreateRequest"
// @Success 200 {object} internal_models.ServerAccountCreateResponse
func (ServerAccountHandler) Create(c *gin.Context) (*format.JSONRespFormat, *SErr.APIErr) {
	req := &models.ServerAccountCreateRequest{}
	e := c.ShouldBind(req)
	if e != nil {
		return nil, SErr.BadRequestErr
	}

	serversSvc := service.GetServersService()
	sErr := serversSvc.AddAccount(c, req.Host, req.Port, req.AccountName, req.AccountPwd)
	if sErr != nil {
		return nil, sErr
	}
	return format.SimpleOKResp(&models.ServerAccountCreateResponse{}), nil
}

// BackupDirInfo
// @Summary 获取一个账户的backup文件夹的相关信息
// @Tags server_account
// @Produce json
// @Router /api/v1/servers/accounts/backupDir [get]
// @Param serverAccountBackupDirRequest query internal_models.ServerAccountBackupDirRequest true "serverAccountBackupDirRequest"
// @Success 200 {object} internal_models.ServerAccountBackupDirResponse
func (ServerAccountHandler) BackupDirInfo(c *gin.Context) (*format.JSONRespFormat, *SErr.APIErr) {
	req := &models.ServerAccountBackupDirRequest{}
	e := c.ShouldBindQuery(req)
	if e != nil {
		return nil, SErr.BadRequestErr
	}

	serversSvc := service.GetServersService()
	backupInfo, err := serversSvc.BackupDirInfo(c, req.Host, req.Port, req.AccountName)
	if err != nil {
		return nil, err
	}
	return format.SimpleOKResp(&models.ServerAccountBackupDirResponse{
		ServerAccountBackupDirInfo: *backupInfo,
	}), nil
}

// Delete
// @Summary 删除一个服务器的账号。
// @Tags server_account
// @Produce json
// @Router /api/v1/servers/accounts [delete]
// @Param serverAccountDeleteRequest body internal_models.ServerAccountDeleteRequest true "serverAccountDeleteRequest"
// @Success 200 {object} internal_models.ServerAccountDeleteResponse
func (ServerAccountHandler) Delete(c *gin.Context) (*format.JSONRespFormat, *SErr.APIErr) {
	req := &models.ServerAccountDeleteRequest{}
	e := c.ShouldBind(req)
	if e != nil {
		return nil, SErr.BadRequestErr
	}

	sessionsSvc := service.GetSessionsService()
	_, err := sessionsSvc.LoggedInAndIsAdmin(c)
	if err != nil {
		return nil, err
	}

	serversSvc := service.GetServersService()
	targetDir, sErr := serversSvc.DeleteAccount(c, req.Host, req.Port, req.AccountName, req.Backup)
	if sErr != nil {
		return nil, sErr
	}
	return format.SimpleOKResp(&models.ServerAccountCreateResponse{
		BackupDir: targetDir,
	}), nil
}


// Update
// @Summary 更新，恢复一个服务器的账号。
// @Tags server_account
// @Produce json
// @Router /api/v1/servers/accounts [put]
// @Param serverAccountUpdateRequest body internal_models.ServerAccountUpdateRequest true "serverAccountUpdateRequest"
// @Success 200 {object} internal_models.ServerAccountUpdateResponse
func (s ServerAccountHandler) Update(c *gin.Context) (*format.JSONRespFormat, *SErr.APIErr) {
	req := &models.ServerAccountUpdateRequest{}
	e := c.ShouldBind(req)
	if e != nil {
		return nil, SErr.BadRequestErr
	}

	var res *models.ServerAccountUpdateResponse
	var err *SErr.APIErr
	if req.Recover {
		res, err = s.recover(c, req)
	} else {
		res, err = s.update(c, req)
	}
	if err != nil {
		return nil, err
	}
	return format.SimpleOKResp(res), nil
}

func (s ServerAccountHandler) recover(c *gin.Context, req *models.ServerAccountUpdateRequest) (*models.ServerAccountUpdateResponse, *SErr.APIErr) {
	serversSvc := service.GetServersService()
	sErr := serversSvc.RecoverAccount(c, req.Host, req.Port, req.AccountName, req.AccountPwd, req.RecoverBackup)
	if sErr != nil {
		return nil, sErr
	}
	return &models.ServerAccountUpdateResponse{}, nil
}

func (ServerAccountHandler) update(c *gin.Context, req *models.ServerAccountUpdateRequest) (*models.ServerAccountUpdateResponse, *SErr.APIErr) {
	serversSvc := service.GetServersService()
	sErr := serversSvc.UpdateAccount(c, req.Host, req.Port, req.AccountName, req.AccountPwd)
	if sErr != nil {
		return nil, sErr
	}
	return &models.ServerAccountUpdateResponse{}, nil
}