package cmd

import (
	"log"
	"project_pos_app/controller"
	"project_pos_app/model"
	"time"

	"github.com/robfig/cron/v3"
)

func CronJob(ctx *controller.AllController) error {
	// Define a new cron scheduler
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(log.Writer(), "cron: ", log.LstdFlags))))

	// Schedule the task to run once a day at midnight
	_, err := c.AddFunc("*/10 * * * *", func() {
		log.Println("Starting daily report generation...")
		productName := "Nasi Goreng"
		data := model.NotifStock(productName)
		err := ctx.Notif.Service.Notif.CreateNotification(data)
		if err != nil {
			log.Printf("Error generating report: %v\n", err)
		} else {
			log.Printf("Report generation Alert %s completed successfully.", productName)
		}
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
