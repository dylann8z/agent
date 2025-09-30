package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

const (
	PidFile = "monitor-agent.pid"
	LogFile = "monitor-agent.log"
)

// DaemonManager 守护进程管理器
type DaemonManager struct {
	pidFile string
	logFile string
}

// NewDaemonManager 创建守护进程管理器
func NewDaemonManager() *DaemonManager {
	return &DaemonManager{
		pidFile: PidFile,
		logFile: LogFile,
	}
}

// GetPID 获取进程PID
func (dm *DaemonManager) GetPID() (int, error) {
	data, err := os.ReadFile(dm.pidFile)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return pid, nil
}

// IsRunning 检查进程是否运行
func (dm *DaemonManager) IsRunning() bool {
	pid, err := dm.GetPID()
	if err != nil {
		return false
	}

	// 检查进程是否存在
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// 发送信号0检查进程是否存活
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// Start 启动守护进程
func (dm *DaemonManager) Start() error {
	if dm.IsRunning() {
		pid, _ := dm.GetPID()
		return fmt.Errorf("monitor-agent is already running (PID: %d)", pid)
	}

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// 打开日志文件
	logFile, err := os.OpenFile(dm.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	// 启动守护进程
	cmd := exec.Command(execPath, "serve")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // 创建新会话
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %v", err)
	}

	// 写入PID文件
	pid := cmd.Process.Pid
	if err := os.WriteFile(dm.pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to write pid file: %v", err)
	}

	// 等待一下确认进程启动
	time.Sleep(time.Second)

	if !dm.IsRunning() {
		os.Remove(dm.pidFile)
		return fmt.Errorf("daemon failed to start")
	}

	fmt.Printf("monitor-agent started successfully (PID: %d)\n", pid)
	return nil
}

// Stop 停止守护进程
func (dm *DaemonManager) Stop() error {
	if !dm.IsRunning() {
		fmt.Println("monitor-agent is not running")
		os.Remove(dm.pidFile)
		return nil
	}

	pid, err := dm.GetPID()
	if err != nil {
		return fmt.Errorf("failed to get pid: %v", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	fmt.Printf("Stopping monitor-agent (PID: %d)...\n", pid)

	// 发送TERM信号
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM: %v", err)
	}

	// 等待进程退出（最多10秒）
	for i := 0; i < 10; i++ {
		if !dm.IsRunning() {
			fmt.Println("monitor-agent stopped")
			os.Remove(dm.pidFile)
			return nil
		}
		time.Sleep(time.Second)
	}

	// 强制杀死
	fmt.Println("Process did not stop gracefully, forcing...")
	if err := process.Signal(syscall.SIGKILL); err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}

	time.Sleep(time.Second)
	os.Remove(dm.pidFile)
	fmt.Println("monitor-agent stopped")
	return nil
}

// Restart 重启守护进程
func (dm *DaemonManager) Restart() error {
	fmt.Println("Restarting monitor-agent...")

	if err := dm.Stop(); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	return dm.Start()
}

// Reload 重载配置
func (dm *DaemonManager) Reload() error {
	if !dm.IsRunning() {
		return fmt.Errorf("monitor-agent is not running")
	}

	pid, err := dm.GetPID()
	if err != nil {
		return fmt.Errorf("failed to get pid: %v", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	fmt.Printf("Reloading monitor-agent (PID: %d)...\n", pid)

	// 发送HUP信号
	if err := process.Signal(syscall.SIGHUP); err != nil {
		return fmt.Errorf("failed to send SIGHUP: %v", err)
	}

	fmt.Println("monitor-agent reloaded successfully")
	return nil
}

// Status 查看状态
func (dm *DaemonManager) Status() error {
	if dm.IsRunning() {
		pid, _ := dm.GetPID()
		fmt.Printf("monitor-agent is running (PID: %d)\n", pid)
		return nil
	}

	fmt.Println("monitor-agent is not running")
	return fmt.Errorf("not running")
}