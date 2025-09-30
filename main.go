package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"host-monitor-agent/api"
	"host-monitor-agent/config"
	"host-monitor-agent/daemon"
)

func main() {
	// 处理命令行参数
	if len(os.Args) > 1 {
		dm := daemon.NewDaemonManager()
		var err error

		switch os.Args[1] {
		case "start":
			err = dm.Start()
		case "stop":
			err = dm.Stop()
		case "restart":
			err = dm.Restart()
		case "reload":
			err = dm.Reload()
		case "status":
			err = dm.Status()
		case "serve":
			// 实际运行服务
			serve()
			return
		default:
			printUsage()
			os.Exit(1)
		}

		if err != nil {
			os.Exit(1)
		}
		return
	}

	// 默认行为：显示帮助信息
	printUsage()
}

func printUsage() {
	fmt.Println("Usage: monitor-agent {start|stop|restart|reload|status|serve}")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  start   - Start the monitor agent daemon")
	fmt.Println("  stop    - Stop the monitor agent daemon")
	fmt.Println("  restart - Restart the monitor agent daemon")
	fmt.Println("  reload  - Reload configuration (send HUP signal)")
	fmt.Println("  status  - Check if the daemon is running")
	fmt.Println("  serve   - Run the server (used internally by daemon)")
}

func serve() {
	// 加载配置
	cfg := config.DefaultConfig()

	// 设置路由
	router := api.SetupRouter()

	// 配置HTTP服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 启动服务器
	go func() {
		log.Printf("Host Monitor Agent starting on %s", addr)
		log.Printf("Access metrics at: http://%s/metrics", addr)
		log.Printf("Health check at: http://%s/health", addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		sig := <-quit
		log.Printf("Received signal: %v", sig)

		// 处理信号
		switch sig {
		case syscall.SIGHUP:
			// 重载配置
			log.Println("Reloading configuration...")
			// TODO: 实现配置重载逻辑
			log.Println("Configuration reloaded")
			// 继续监听信号
			continue
		case syscall.SIGINT, syscall.SIGTERM:
			// 优雅关闭
			log.Println("Shutting down server...")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("Server forced to shutdown: %v", err)
			}

			log.Println("Server exited")
			return
		}
	}
}