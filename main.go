package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"

	"context"
	"encoding/json"
	"net/http"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		if err := InitCollections(app); err != nil {
			return err
		}

		// create initial admin
		email := os.Getenv("POCKETBASE_INITIAL_ADMIN_EMAIL")
		password := os.Getenv("POCKETBASE_INITIAL_ADMIN_PASSWORD")

		if email != "" && password != "" {
			dao := daos.New(app.Dao().DB())
			admin, err := dao.FindAdminByEmail(email)
			if err == nil && admin != nil {
				log.Println("Initial admin already exists")
			} else {
				admin := &models.Admin{}
				admin.Email = email
				if err := admin.SetPassword(password); err != nil {
					return err
				}

				if err := dao.SaveAdmin(admin); err != nil {
					return err
				}
				log.Println("Initial admin created")
			}
		}

		e.Router.GET("/api/health", func(c echo.Context) error {
			return c.JSON(200, map[string]string{"status": "ok"})
		})

		e.Router.GET("/", func(c echo.Context) error {
			return c.File("public/index.html")
		})

		e.Router.GET("/test", func(c echo.Context) error {
			return c.File("public/test/test.html")
		})

		e.Router.GET("/api/receipts", func(c echo.Context) error {
			dao := app.Dao()
			records, err := dao.FindRecordsByFilter(
				"receipts",
				"1=1",      // filter
				"+created", // sort
				0,          // limit
				0,          // offset
			)
			if err != nil {
				log.Printf("Failed to fetch receipts: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch receipts.")
			}

			result := []map[string]interface{}{}
			for _, record := range records {
				result = append(result, map[string]interface{}{
					"id":             record.Id,
					"title":          record.GetString("title"),
					"filename":       record.GetString("filename"),
					"created":        record.GetCreated(),
					"total":          record.Get("total"),
					"verified_total": record.GetBool("verified_total"),
				})
			}

			return c.JSON(http.StatusOK, result)
		})

		e.Router.GET("/api/receipts/:id/items", func(c echo.Context) error {
			id := c.PathParam("id")
			if id == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Receipt ID is required.")
			}

			dao := app.Dao()
			records, err := dao.FindRecordsByFilter(
				"items",
				"receipt = {:id}",
				"+created",
				0,
				0,
				map[string]interface{}{"id": id},
			)
			if err != nil {
				log.Printf("Failed to fetch items for receipt %s: %v", id, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch items.")
			}

			result := []map[string]interface{}{}
			for _, record := range records {
				result = append(result, map[string]interface{}{
					"name":     record.GetString("name"),
					"price":    record.Get("price"),
					"quantity": record.Get("quantity"),
					"amount":   record.Get("amount"),
				})
			}

			return c.JSON(http.StatusOK, result)
		})

		e.Router.POST("/api/receipts/upload", func(c echo.Context) error {
			file, err := c.FormFile("receipt")
			if err != nil {
				return err
			}

			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			if err := os.MkdirAll("uploads", 0755); err != nil {
				return err
			}

			dst, err := os.Create(filepath.Join("uploads", file.Filename))
			if err != nil {
				return err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

			return c.JSON(200, map[string]interface{}{
				"message":  "Receipt uploaded successfully.",
				"filename": file.Filename,
			})
		})

		e.Router.POST("/api/receipts/analyze", func(c echo.Context) error {
			type AnalyzeRequest struct {
				Filename string `json:"filename"`
			}
			req := new(AnalyzeRequest)
			if err := c.Bind(req); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
			}

			if req.Filename == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Filename is required")
			}

			apiKey := os.Getenv("GEMINI_API_KEY")
			if apiKey == "" {
				log.Println("GEMINI_API_KEY environment variable not set.")
				return echo.NewHTTPError(http.StatusInternalServerError, "Server configuration error.")
			}

			filePath := filepath.Join("uploads", req.Filename)
			imgData, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read file %s: %v", filePath, err)
				return echo.NewHTTPError(http.StatusNotFound, "File not found.")
			}

			ctx := context.Background()
			client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
			if err != nil {
				log.Printf("Failed to create genai client: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to analysis service.")
			}
			defer client.Close()

			model := client.GenerativeModel("gemini-1.5-flash")
			prompt := genai.Text("Analyze this receipt and return a JSON object with four keys: 'title' (the name of the store), 'date' (in YYYY-MM-DD format), 'items' (an array of objects, each with 'name', 'price', and 'quantity'), and 'total' (the total amount). Only return the JSON object.")
			imgPart := genai.ImageData(http.DetectContentType(imgData), imgData)

			resp, err := model.GenerateContent(ctx, imgPart, prompt)
			if err != nil {
				log.Printf("Failed to generate content from model: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to analyze receipt.")
			}

			if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
				return echo.NewHTTPError(http.StatusInternalServerError, "Analysis returned no content.")
			}

			analysisResult := ""
			if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
				analysisResult = string(textPart)
			}

			// Clean the response from Gemini to ensure it's a valid JSON string.
			analysisResult = strings.TrimSpace(analysisResult)
			if strings.HasPrefix(analysisResult, "```json") {
				analysisResult = strings.TrimPrefix(analysisResult, "```json")
				analysisResult = strings.TrimSuffix(analysisResult, "```")
				analysisResult = strings.TrimSpace(analysisResult)
			}

			var resultJSON map[string]interface{}
			if err := json.Unmarshal([]byte(analysisResult), &resultJSON); err != nil {
				log.Printf("Failed to unmarshal analysis result into JSON: %v. Raw result: %s", err, analysisResult)
				// Return the raw text if it's not valid JSON, just in case.
				return c.String(http.StatusOK, analysisResult)
			}

			dao := app.Dao()
			receiptsCollection, err := dao.FindCollectionByNameOrId("receipts")
			if err != nil {
				log.Printf("Failed to find 'receipts' collection: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Database configuration error.")
			}
			itemsCollection, err := dao.FindCollectionByNameOrId("items")
			if err != nil {
				log.Printf("Failed to find 'items' collection: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Database configuration error.")
			}

			receiptRecord := models.NewRecord(receiptsCollection)
			receiptRecord.Set("filename", req.Filename)

			if title, ok := resultJSON["title"].(string); ok {
				receiptRecord.Set("title", title)
			}
			if date, ok := resultJSON["date"].(string); ok {
				receiptRecord.Set("date", date)
			}
			if total, ok := resultJSON["total"].(float64); ok {
				receiptRecord.Set("total", total)
			}

			// Calculate total from items and verify against receipt total
			var calculatedTotal float64
			if items, ok := resultJSON["items"].([]interface{}); ok {
				for _, itemData := range items {
					if itemMap, ok := itemData.(map[string]interface{}); ok {
						price, priceOk := itemMap["price"].(float64)
						quantity, quantityOk := itemMap["quantity"].(float64)
						if !quantityOk {
							quantity = 1
						}
						if priceOk {
							calculatedTotal += price * quantity
						}
					}
				}
			}

			receiptTotal, _ := resultJSON["total"].(float64)
			// Use a small tolerance for float comparison
			if calculatedTotal > 0 && receiptTotal > 0 && (calculatedTotal-receiptTotal < 0.01 && receiptTotal-calculatedTotal < 0.01) {
				receiptRecord.Set("verified_total", true)
			} else {
				receiptRecord.Set("verified_total", false)
			}

			if err := dao.SaveRecord(receiptRecord); err != nil {
				log.Printf("Failed to save receipt record: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save receipt to database.")
			}

			if items, ok := resultJSON["items"].([]interface{}); ok {
				for _, itemData := range items {
					if itemMap, ok := itemData.(map[string]interface{}); ok {
						itemRecord := models.NewRecord(itemsCollection)
						if name, ok := itemMap["name"].(string); ok {
							itemRecord.Set("name", name)
						}
						priceVal := 0.0
						if price, ok := itemMap["price"].(float64); ok {
							itemRecord.Set("price", price)
							priceVal = price
						}
						quantityVal := 1.0
						if quantity, ok := itemMap["quantity"].(float64); ok {
							itemRecord.Set("quantity", quantity)
							quantityVal = quantity
						} else {
							itemRecord.Set("quantity", 1)
						}
						itemRecord.Set("amount", priceVal*quantityVal)
						itemRecord.Set("receipt", receiptRecord.Id)

						if err := dao.SaveRecord(itemRecord); err != nil {
							log.Printf("Failed to save item record: %v", err)
							// Continue saving other items
						}
					}
				}
			}

			return c.JSON(http.StatusOK, resultJSON)
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
