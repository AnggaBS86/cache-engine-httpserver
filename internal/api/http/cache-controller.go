package http

import (
	"cache_engine_httpserver/internal/api/model"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v3"
)

func GetCache(c fiber.Ctx, ctx *model.CacheAppContext) error {
	key := c.Query("key")
	data, err := ctx.Cache.Get(key)
	if err != nil {
		cacheExists := isCacheExists(err)
		if !cacheExists {
			return c.JSON(fiber.Map{
				"status":  "ERROR",
				"message": "Key not found",
				"cache":   nil,
			})
		}

		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"status":  "ERROR",
			"message": "Failed to decode cache entry",
			"cache":   nil,
		})
	}

	entry := model.CacheEntry{}
	if err := json.Unmarshal(data, &entry); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"status":  "ERROR",
			"message": "Something error with Get cache operation.",
			"cache":   nil,
		})
	}

	if time.Now().After(entry.Expiration) {
		err = ctx.Cache.Delete(key)
		if err != nil {
			log.Println(err.Error())
			return c.JSON(fiber.Map{
				"status":  "ERROR",
				"message": "Something error with cache deletion.",
				"cache":   nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "Key is expired",
			"cache":   nil,
		})
	}

	return c.JSON(fiber.Map{
		"status": "OK",
		"cache": fiber.Map{
			"key":   key,
			"value": entry.Value,
		},
	})
}

func CreateCache(c fiber.Ctx, ctx *model.CacheAppContext) error {
	cacheReq := new(model.CacheCreationRequest)
	if err := c.Bind().Body(cacheReq); err != nil {
		return err
	}

	if valid, err := validateCacheCreate(*cacheReq); valid == false {
		return c.JSON(fiber.Map{
			"status":           "ERROR",
			"message":          "Validation error",
			"cache":            nil,
			"validation_error": err,
		})
	}

	expiration := time.Now().Add(time.Duration(cacheReq.DurationInSeconds) * time.Second)
	entry := model.CacheEntry{
		Value:      cacheReq.Value,
		Expiration: expiration,
	}

	entryData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error when marshaling entry data : %v", err.Error())
		return c.JSON(fiber.Map{
			"status":  "ERROR",
			"message": "Failed to encode cache entry",
			"cache":   nil,
		})
	}

	err = ctx.Cache.Set(cacheReq.Key, entryData)
	if err != nil {
		log.Printf("Error when Set cache value : %v", err.Error())
		return c.JSON(fiber.Map{
			"status":  "ERROR",
			"message": "Something error when set cache",
			"cache":   nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": "Value set successfully",
		"cache": fiber.Map{
			"key":                 cacheReq.Key,
			"value":               cacheReq.Value,
			"duration_in_seconds": cacheReq.DurationInSeconds,
		},
	})
}

func DeleteCache(c fiber.Ctx, ctx *model.CacheAppContext) error {
	key := c.Params("key")
	_, err := ctx.Cache.Get(key)
	if err != nil {
		cacheExists := isCacheExists(err)
		if !cacheExists {
			return c.JSON(fiber.Map{
				"status":  "ERROR",
				"message": "Key not found",
				"cache":   nil,
			})
		}

	}

	err = ctx.Cache.Delete(key)
	if err != nil {
		log.Printf("Error occured when `DeleteCache` : %v", err.Error())
		return c.JSON(fiber.Map{
			"status":  "ERROR",
			"message": "Something error happened when deleting cache with key `" + key + "`",
			"cache":   nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": "Cache deleted successfully",
		"cache": fiber.Map{
			"key": key,
		},
	})
}

func IsCacheExists(c fiber.Ctx, ctx *model.CacheAppContext) error {
	key := c.Params("key")
	_, err := ctx.Cache.Get(key)
	if err != nil {
		cacheExists := isCacheExists(err)
		if !cacheExists {
			return c.JSON(fiber.Map{
				"status":  "OK",
				"message": "Cache does not exists",
				"cache": fiber.Map{
					"exists": false,
				},
			})
		}

	}

	data, err := ctx.Cache.Get(key)
	entry := model.CacheEntry{}
	if err := json.Unmarshal(data, &entry); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"status":  "ERROR",
			"message": "Failed to decode cache entry",
			"cache":   nil,
		})
	}

	if time.Now().After(entry.Expiration) {
		err = ctx.Cache.Delete(key)
		if err != nil {
			log.Println(err.Error())
			return c.JSON(fiber.Map{
				"status":  "ERROR",
				"message": "Something error with cache deletion",
				"cache":   nil,
			})
		}

		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "Cache is expired",
			"cache": fiber.Map{
				"exists": false,
			},
		})
	}

	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": "Cache exists",
		"cache": fiber.Map{
			"exists": true,
		},
	})
}

func isCacheExists(err error) bool {
	return err != bigcache.ErrEntryNotFound
}

func validateCacheCreate(request model.CacheCreationRequest) (bool, map[string]interface{}) {
	validationErr := make(map[string]any)
	if strings.TrimSpace(request.Key) == "" {
		validationErr["key"] = "Cache `key` cannot be empty"
	}

	value, ok := request.Value.(string)
	if !ok {
		log.Println("request.Value is not a string")
		value = ""
	}

	if strings.TrimSpace(value) == "" {
		validationErr["value"] = "Cache `value` cannot be empty"
	}

	if request.DurationInSeconds < 1 {
		validationErr["duration_in_seconds"] = "Value `duration_in_seconds` should be >= 0"
	}

	if len(validationErr) < 1 {
		return true, nil
	}

	return false, validationErr
}
