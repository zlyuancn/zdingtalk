
# 用户登录/获取用户基本信息

## 根据临时授权码获取用户信息
> 注意: 这里需要的appId,appSecret是通过"移动接入应用>登录>创建扫码登录应用授权"获取的
```go
func GetUserInfoByTempCode(tmpCode string, appId, appSecret string) (*DingDing_GetUserInfoByTempCodeResp, error)
```

## 获取AccessToken
> 注意, 不能频繁调用这个接口, api有限流控制.
```go
func GetAccessToken(appKey, appSecret string) (*DingDing_GetAccessTokenResp, error)
```

## 根据UnionId获取用户id
> 需要 应用权限管理>通讯录/通讯录只读权限, 开发管理/服务器出口IP白名单
```go
func GetUserIdByUnionId(unionId string, accessToken string) (*DingDing_GetUserIdByUnionIdResp, error)
```

## 获取用户信息返回结果
> 可选权限 应用权限管理>通讯录>手机号码信息/邮箱等个人信息
```go
func GetUserInfoByUserId(userId string, accessToken string) (*DingDing_GetUserInfoByUserIdResp, error)
```
