package main

import (
	"encoding/json"
	//"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"log"
	"os"
)

func main() {
	content, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	var config map[string]interface{}
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error parsing config file: ", err)
	}

	client := resty.New()

	app := fiber.New()
	app.Use(cache.New())

	app.Get("/file/:identifier", func(c *fiber.Ctx) error {
		if identifierInFiles(config, c.Params("identifier")) {
			fileURL, ok := getURLFromFiles(config, c.Params("identifier"))
			if !ok {
				return c.SendStatus(500)
			}

			resp, err := client.R().
				Get(fileURL)
	
			if err != nil {
				return c.SendStatus(500)
			}

			if resp.Status() != "200 OK" {
				return c.SendStatus(500)
			}

			c.Set(fiber.HeaderContentType, resp.Header().Get(fiber.HeaderContentType))
			return c.Send(resp.Body())
		} else {
			// raise 404
			return c.SendStatus(404)
		}
	})

	app.Listen(":3000")
}

func identifierInFiles(config map[string]interface{}, identifier string) bool {
	// Check if the "files" key exists in the config map
	if files, ok := config["files"]; ok {
		// Check if "files" is a map[string]interface{} type
		if filesMap, ok := files.(map[string]interface{}); ok {
			// Check if the "identifier" key exists in the "files" map
			_, identifierExists := filesMap[identifier]
			return identifierExists
		}
	}
	return false
}

func getURLFromFiles(config map[string]interface{}, identifier string) (string, bool) {
	files, ok := config["files"]
	if !ok {
		return "", false
	}
	filesMap, ok := files.(map[string]interface{})
	if !ok {
		return "", false
	}
	value, ok := filesMap[identifier]
	if !ok {
		return "", false
	}
	identifierMap:= value.(map[string]interface{})
	url := identifierMap["url"]
	urlString, ok := url.(string)
	if ok {
		return urlString, true
	}
	return "", false
}
