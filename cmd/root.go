/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-rate-limiter",
	Short: "Rate limit testing server",
	Long:  `Run a rate limit testing server based on the algorithm of your choice`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var tokenBucketCmd = &cobra.Command{
	Use:   "tokenBucket",
	Short: "Token bucket rate limit algorithm",
	Long:  `Run a new server with the token bucket rate limit algorithm`,
	Run: func(cmd *cobra.Command, args []string) {
		NewServer(map[string]string{
			"algorithm":  "token-bucket",
			"capacity":   cmd.Flag("capacity").Value.String(),
			"refillRate": cmd.Flag("refillRate").Value.String(),
		}).Run(cmd.Flag("addr").Value.String())
	},
}

var fixedWindowCmd = &cobra.Command{
	Use:   "fixedWindow",
	Short: "Fixed window rate limit algorithm",
	Long:  `Run a new server with the fixed window rate limit algorithm`,
	Run: func(cmd *cobra.Command, args []string) {
		NewServer(map[string]string{
			"algorithm": "fixed-window",
			"capacity":  cmd.Flag("capacity").Value.String(),
			"duration":  cmd.Flag("duration").Value.String(),
		}).Run(cmd.Flag("addr").Value.String())
	},
}

var slidingWindowLogCmd = &cobra.Command{
	Use:   "slidingWindowLog",
	Short: "Sliding window log rate limit algorithm",
	Long:  `Run a new server with the sliding window log rate limit algorithm`,
	Run: func(cmd *cobra.Command, args []string) {
		NewServer(map[string]string{
			"algorithm": "sliding-window-log",
			"capacity":  cmd.Flag("capacity").Value.String(),
			"duration":  cmd.Flag("duration").Value.String(),
		}).Run(cmd.Flag("addr").Value.String())
	},
}

var slidingWindowCounterCmd = &cobra.Command{
	Use:   "slidingWindowCounter",
	Short: "Sliding window counter rate limit algorithm",
	Long:  `Run a new server with the sliding window counter rate limit algorithm`,
	Run: func(cmd *cobra.Command, args []string) {
		NewServer(map[string]string{
			"algorithm": "sliding-window-counter",
			"capacity":  cmd.Flag("capacity").Value.String(),
			"duration":  cmd.Flag("duration").Value.String(),
			"weight":    cmd.Flag("weight").Value.String(),
		}).Run(cmd.Flag("addr").Value.String())
	},
}

var redisSlidingWindowCounterCmd = &cobra.Command{
	Use:   "redisSlidingWindowCounter",
	Short: "Sliding window counter rate limit algorithm using Redis",
	Long:  `Run a new server with the sliding window counter rate limit algorithm using Redis to store the data`,
	Run: func(cmd *cobra.Command, args []string) {
		NewServer(map[string]string{
			"algorithm": "redis-sliding-window-counter",
			"capacity":  cmd.Flag("capacity").Value.String(),
			"duration":  cmd.Flag("duration").Value.String(),
			"weight":    cmd.Flag("weight").Value.String(),
			"redisURL":  cmd.Flag("redisURL").Value.String(),
		}).Run(cmd.Flag("addr").Value.String())
	},
}

func init() {
	// Server config
	rootCmd.PersistentFlags().String("addr", ":8080", "The address to listen on")

	// Token bucket rate limit algorithm
	rootCmd.AddCommand(tokenBucketCmd)
	tokenBucketCmd.Flags().Int32("capacity", 10, "The maximum number of requests allowed in the time window")
	tokenBucketCmd.Flags().Int32("refillRate", 1, "The number of requests to add per second")

	// Fixed window rate limit algorithm
	rootCmd.AddCommand(fixedWindowCmd)
	fixedWindowCmd.Flags().Int32("capacity", 60, "The maximum number of requests allowed in the time window")
	fixedWindowCmd.Flags().Int32("duration", 60, "The duration of the window in seconds")

	// Sliding window log rate limit algorithm
	rootCmd.AddCommand(slidingWindowLogCmd)
	slidingWindowLogCmd.Flags().Int32("capacity", 60, "The maximum number of requests allowed in the time window")
	slidingWindowLogCmd.Flags().Int32("duration", 60, "The duration of the window in seconds")

	// Sliding window counter rate limit algorithm
	rootCmd.AddCommand(slidingWindowCounterCmd)
	slidingWindowCounterCmd.Flags().Int32("capacity", 60, "The maximum number of requests allowed in the time window")
	slidingWindowCounterCmd.Flags().Int32("duration", 60, "The duration of the window in seconds")
	slidingWindowCounterCmd.Flags().Float64("weight", 0.4, "The weight of the current window in the average calculation")

	// Sliding window counter rate limit algorithm using Redis
	rootCmd.AddCommand(redisSlidingWindowCounterCmd)
	redisSlidingWindowCounterCmd.Flags().Int32("capacity", 60, "The maximum number of requests allowed in the time window")
	redisSlidingWindowCounterCmd.Flags().Int32("duration", 60, "The duration of the window in seconds")
	redisSlidingWindowCounterCmd.Flags().Float64("weight", 0.4, "The weight of the current window in the average calculation")
	redisSlidingWindowCounterCmd.Flags().String("redisURL", "redis://localhost:6379/0", "The URL of the Redis server")
}
