package service

import (
	"auth/biz/dal/redis"
	"auth/conf"
	"context"
	"fmt"
	"github.com/BeroKiTeer/MyGoMall/common/kitex_gen/auth"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/golang-jwt/jwt/v4"
)

type VerifyTokenByRPCService struct {
	ctx context.Context
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {

	userID, err := GetUserIDFromToken(req.Token)
	if err != nil {
		klog.Error("获取token失败", err)
		return nil, err
	}

	// 查询 userID 在 redis 里面是否存在，如果存在，获取密钥并验证
	var secretKey string
	key := fmt.Sprintf("auth:%d", userID)
	if conf.GetEnv() == "test" {
		if secretKey, err = redis.RedisClient.Get(s.ctx, key).Result(); err != nil {
			return &auth.VerifyResp{Res: false}, err
		}
	} else if conf.GetEnv() == "dev" {
		if secretKey, err = redis.RedisClusterClient.Get(s.ctx, key).Result(); err != nil {
			return &auth.VerifyResp{Res: false}, err
		}
	}

	// 按照查询到的密钥，解析并验证 Token
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			klog.Errorf("不支持的签名方法: %v", token.Header["alg"])
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		}
		// 密钥只支持 []byte 类型，强转。
		return []byte(secretKey), nil
	})

	if err != nil {
		return &auth.VerifyResp{Res: false}, err
	}

	// 断言解析结果，返回是否对得上
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		klog.Info("身份验证成功")
		return &auth.VerifyResp{Res: true}, nil
	} else {
		return &auth.VerifyResp{Res: false}, nil
	}

}
