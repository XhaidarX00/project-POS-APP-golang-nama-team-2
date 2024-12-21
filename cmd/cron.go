package cmd

import (
	"log"
	"project_pos_app/infra"
	"project_pos_app/model"
	"time"

	"github.com/robfig/cron/v3"
)

func CronJob(ctx *infra.IntegrationContext) error {
	// Define a new cron scheduler
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(log.Writer(), "cron: ", log.LstdFlags))))

	// untuk kedepannya pakai websocket atau firebase
	// Schedule the task to check stock and send notifications if stock < LowStock
	_, err := c.AddFunc("0 0 * * *", func() {
		log.Println("Checking product stock levels...")
		products, err := ctx.Ctl.Revenue.Service.Revenue.GetLowStockProducts(ctx.Cfg.LowStock)
		if err != nil {
			log.Printf("Error fetching low stock products: %v\n", err)
			return
		}
		for _, product := range products {
			data := model.Notification{
				Title:     "Low Stock Alert",
				Message:   "Product " + product.Name + " has less than 10 items in stock.",
				CreatedAt: time.Now(),
			}
			err := ctx.Ctl.Notif.Service.Notif.CreateNotification(data)
			if err != nil {
				log.Printf("Error sending notification for product %s: %v\n", product.Name, err)
			} else {
				log.Printf("Notification for product %s sent successfully.\n", product.Name)
			}
		}
	})
	if err != nil {
		return err
	}

	// Schedule the task to generate revenue reports every day at 1 AM
	_, err = c.AddFunc("0 1 * * *", func() {
		// _, err = c.AddFunc("*/1 * * * *", func() {
		log.Println("Starting revenue report generation...")
		// Generate order revenue
		orders, err := ctx.Ctl.Revenue.Service.Revenue.CalculateOrderRevenue()
		if err != nil {
			log.Printf("Error calculating order revenue: %v\n", err)
			return
		}
		for _, order := range orders {
			err := ctx.Ctl.Revenue.Service.Revenue.SaveOrderRevenue(order)
			if err != nil {
				log.Printf("Error saving order revenue for order %d: %v\n", order.ID, err)
			}
		}

		// Generate product revenue
		products, err := ctx.Ctl.Revenue.Service.Revenue.CalculateProductRevenue()
		if err != nil {
			log.Printf("Error calculating product revenue: %v\n", err)
			return
		}
		log.Printf("%v\n", products)

		for _, product := range products {
			product.ProfitMargin = ctx.Cfg.ProfitMargin
			err := ctx.Ctl.Revenue.Service.Revenue.SaveProductRevenue(product)
			if err != nil {
				log.Printf("Error saving product revenue for product %s: %v\n", product.ProductName, err)
			}
		}

		log.Println("Revenue report generation completed successfully.")
	})
	if err != nil {
		return err
	}

	// Start the cron scheduler
	c.Start()

	// Run a blocking loop to keep the cron job running
	go func() {
		select {} // Prevent the function from exiting
	}()

	log.Printf("Cron job initialized and running. Current time: %v\n", time.Now())
	return nil
}
