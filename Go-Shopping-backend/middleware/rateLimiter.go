package middleware

import (
	"Go-Shopping-backend/initializers"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Client struct {
	Limit int `json:"limit"`
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip, _, _ := net.SplitHostPort(ctx.Request.RemoteAddr)
		key := "client:" + ip

		exists, err := initializers.RedisClient.Exists(key).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while checking data of client in redis",
			})
			return
		}

		if exists == 1 {
			clientJson, err := initializers.RedisClient.Get(key).Result()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error while getting data of client in redis",
				})
				return
			}

			var client Client

			err = json.Unmarshal([]byte(clientJson), &client)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error while unmarshaling data of client ",
				})
				return
			}

			client.Limit++
			log.Println(client.Limit)
			remainingTime := initializers.RedisClient.TTL(key).Val()

			if client.Limit > 5 {
				ctx.JSON(http.StatusTooManyRequests, gin.H{
					"message": "Too many requests please retry after some time",
				})
				ctx.Abort()
				return
			}

			ClientJson, err := json.Marshal(client)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error while marshal new client data",
				})
			}

			initializers.RedisClient.Set(key, ClientJson, remainingTime)

			ctx.Next()
			return
		} else {

			newClient := Client{Limit: 1}
			newClientJson, err := json.Marshal(newClient)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error while marshal new client data",
				})
			}

			err = initializers.RedisClient.Set(key, newClientJson, 30*time.Second).Err()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error storing new client to redis",
				})
			}

			ctx.Next()

			return

		}
	}
}
