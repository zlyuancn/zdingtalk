/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/12/8
   Description :
-------------------------------------------------
*/

package user

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zlyuancn/zdingtalk/utils"
)

// 个人免登场景的签名计算
//
// https://ding-doc.dingtalk.com/document/#/org-dev-guide/signature-calculation-for-logon-free-scenarios#topic-1949474
func makeSingleFreeLoginSignature(appSecret string) (timestamp string, signature string) {
	timestamp = strconv.Itoa(int(time.Now().UnixNano() / 1e6))
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(timestamp))
	sum := h.Sum(nil)

	text := base64.StdEncoding.EncodeToString(sum)
	signature = url.QueryEscape(text)
	return
}

// 根据临时授权码获取用户信息返回结果
type DingDing_GetUserInfoByTempCodeResp struct {
	ErrCode  int64  `json:"errcode"` // 错误码, 0表示ok
	ErrMsg   string `json:"errmsg"`  // 错误描述
	UserInfo struct {
		Nick                 string `json:"nick"`    // 用户在钉钉上面的昵称
		UnionId              string `json:"unionid"` // 用户在当前开放应用所属企业的唯一标识
		DingId               string `json:"dingId"`
		Openid               string `json:"openid"`                   // 用户在当前开放应用内的唯一标识
		MainOrgAuthHighLevel bool   `json:"main_org_auth_high_level"` // 用户主企业是否达到高级认证级别
	} `json:"user_info"` // 用户信息
}

// 根据临时授权码获取用户信息, 注意: 这里需要的appId,appSecret是通过"移动接入应用>登录>创建扫码登录应用授权"获取的
//
// 需要 移动接入应用>登录>创建扫码登录应用授权
// https://ding-doc.dingtalk.com/document/#/org-dev-guide/obtain-the-user-information-based-on-the-sns-temporary-authorization#topic-1995619
func GetUserInfoByTempCode(tmpCode string, appId, appSecret string) (*DingDing_GetUserInfoByTempCodeResp, error) {
	const ApiUrl = "https://oapi.dingtalk.com/sns/getuserinfo_bycode"
	body := fmt.Sprintf(`{"tmp_auth_code": "%s"}`, tmpCode)
	req, err := http.NewRequest("POST", ApiUrl, bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, fmt.Errorf("构建request失败: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	timestamp, signature := makeSingleFreeLoginSignature(appSecret)
	req.URL.RawQuery = fmt.Sprintf("accessKey=%s&timestamp=%s&signature=%s", appId, timestamp, signature)

	var result DingDing_GetUserInfoByTempCodeResp
	err = utils.Request(req, &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return &result, fmt.Errorf("收到错误码: %d: %s", result.ErrCode, result.ErrMsg)
	}

	return &result, nil
}

// 获取access_token返回结果
type DingDing_GetAccessTokenResp struct {
	ErrCode     int64  `json:"errcode"`      // 错误码, 0表示ok
	ErrMsg      string `json:"errmsg"`       // 错误描述
	AccessToken string `json:"access_token"` // AccessToken
	ExpiresIn   int64  `json:"expires_in"`   // 过期时间，单位秒
}

// 获取AccessToken
//
// 需要 创建小程序
// 注意, 不能频繁调用这个接口, api有限流控制.
// 应该根据过期时间缓存起来, 过期后再次调用此接口获取AccessToken
// https://ding-doc.dingtalk.com/document/#/org-dev-guide/obtain-access_token
func GetAccessToken(appKey, appSecret string) (*DingDing_GetAccessTokenResp, error) {
	const ApiUrl = "https://oapi.dingtalk.com/gettoken"
	req, err := http.NewRequest("GET", ApiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("构建request失败: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = fmt.Sprintf("appkey=%s&appsecret=%s", appKey, appSecret)

	var result DingDing_GetAccessTokenResp
	err = utils.Request(req, &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return &result, fmt.Errorf("收到错误码: %d: %s", result.ErrCode, result.ErrMsg)
	}

	return &result, nil
}

// 根据UnionId获取用户id返回结果
type DingDing_GetUserIdByUnionIdResp struct {
	ErrCode   int64  `json:"errcode"`    // 错误码, 0表示ok
	ErrMsg    string `json:"errmsg"`     // 错误描述
	RequestId string `json:"request_id"` // 请求id
	Result    struct {
		ContactType int64  `json:"contact_type"` // 联系类型; 0=企业内部员工; 1=企业外部联系人
		Userid      string `json:"userid"`       // 用户id
	} `json:"result"`
}

// 根据UnionId获取用户id
//
// 需要 应用权限管理>通讯录/通讯录只读权限, 开发管理/服务器出口IP白名单
// https://ding-doc.dingtalk.com/document/#/org-dev-guide/retrieve-user-information-based-on-the-union-id#topic-1960045
func GetUserIdByUnionId(unionId string, accessToken string) (*DingDing_GetUserIdByUnionIdResp, error) {
	const ApiUrl = "https://oapi.dingtalk.com/topapi/user/getbyunionid"
	body := fmt.Sprintf(`{"unionid": "%s"}`, unionId)
	req, err := http.NewRequest("POST", ApiUrl, bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, fmt.Errorf("构建request失败: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = fmt.Sprintf("access_token=%s", accessToken)

	var result DingDing_GetUserIdByUnionIdResp
	err = utils.Request(req, &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return &result, fmt.Errorf("收到错误码: %d: %s", result.ErrCode, result.ErrMsg)
	}

	return &result, nil
}

// 根据用户id获取用户信息返回结果
type DingDing_GetUserInfoByUserIdResp struct {
	ErrCode   int64  `json:"errcode"`    // 错误码, 0表示ok
	ErrMsg    string `json:"errmsg"`     // 错误描述
	RequestId string `json:"request_id"` // 请求id
	Result    struct {
		UnionId string `json:"unionid"` // 员工的userid
		UserId  string `json:"userid"`  // 员工在当前开发者企业账号范围内的唯一标识
		Name    string `json:"name"`    // 员工名称
		Avatar  string `json:"avatar"`  // 头像

		// 以下需要授权 应用权限管理>通讯录>手机号码信息
		StateCode  string `json:"state_code"`  // 国际电话区号
		Mobile     string `json:"mobile"`      // 手机号码
		HideMobile bool   `json:"hide_mobile"` // 是否号码隐藏; 隐藏手机号后，手机号在个人资料页隐藏，但仍可对其发DING、发起钉钉免费商务电话
		Telephone  string `json:"telephone"`   // 分机号
		// 以下需要授权 应用权限管理>通讯录>邮箱等个人信息
		Email string `json:"email"` // 员工邮箱

		JobNumber string `json:"job_number"` // 员工工号
		Title     string `json:"title"`      // 职位
		WorkPlace string `json:"work_place"` // 办公地点
		Remark    string `json:"remark"`     // 备注

		DeptIdList    []int `json:"dept_id_list"` // 所属部门ID列表
		DeptOrderList []struct {
			DeptId int64   `json:"dept_id"` // 部门ID
			Order  float64 `json:"order"`   // 员工在部门中的排序
		} `json:"dept_order_list"`           // 员工在对应的部门中的排序
		Extension string `json:"extension"`  // 扩展属性，最大长度2000个字符
		HiredDate int    `json:"hired_date"` // 入职时间，毫秒时间戳

		Active       bool `json:"active"`      // 是否激活了钉钉
		RealAuthed   bool `json:"real_authed"` // 是否完成了实名认证
		Senior       bool `json:"senior"`      // 是否为企业的高管
		Admin        bool `json:"admin"`       // 是否为企业的管理员
		Boss         bool `json:"boss"`        // 是否为企业的老板
		LeaderInDept []struct {
			Leader bool  `json:"leader"`  // 是否是领导
			DeptId int64 `json:"dept_id"` // 部门id
		} `json:"leader_in_dept"` // 员工在对应的部门中是否领导
		RoleList []struct {
			GroupName string `json:"group_name"` // 角色组名
			Name      string `json:"name"`       // 角色名
			Id        int64  `json:"id"`         // 角色id
		} `json:"role_list"` // 角色列表
	} `json:"result"`
}

// 获取用户信息返回结果
//
// 可选权限 应用权限管理>通讯录>手机号码信息/邮箱等个人信息
// https://ding-doc.dingtalk.com/document/#/org-dev-guide/queries-user-details#topic-1960047
func GetUserInfoByUserId(userId string, accessToken string) (*DingDing_GetUserInfoByUserIdResp, error) {
	const ApiUrl = "https://oapi.dingtalk.com/topapi/v2/user/get"
	body := fmt.Sprintf(`{"userid": "%s", "language": "zh_CN"}`, userId)
	req, err := http.NewRequest("POST", ApiUrl, bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, fmt.Errorf("构建request失败: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = fmt.Sprintf("access_token=%s", accessToken)

	var result DingDing_GetUserInfoByUserIdResp
	err = utils.Request(req, &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return &result, fmt.Errorf("收到错误码: %d: %s", result.ErrCode, result.ErrMsg)
	}

	return &result, nil
}
