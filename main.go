package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 模拟一扇门
type Door struct {
	hasPrize bool
}

// 全局变量 doors
var doors = []Door{
	{hasPrize: false},
	{hasPrize: true}, // 假设大奖在第二扇门后
	{hasPrize: false},
}

// 模拟一次游戏的结果
type GameResult struct {
	StickWin  bool
	SwitchWin bool
}

// 模拟一次游戏
func simulateGame(resultChan chan<- GameResult, wg *sync.WaitGroup, doors []Door) {
	defer wg.Done()

	// 你的初始选择
	yourChoice := rand.Intn(len(doors))

	// 主持人打开一扇没有大奖的门
	// 这里需要稍微修改逻辑，因为我们需要确保不打开你选择的门和有大奖的门
	// 假设 hostChoice 总是指向没有大奖且不是你选择的门
	var hostChoice int
	for hostChoice == yourChoice || doors[hostChoice].hasPrize {
		hostChoice = rand.Intn(len(doors))
	}

	// 检查坚持选择的结果
	stickWin := doors[yourChoice].hasPrize

	// 改变选择并检查结果
	newChoice := (yourChoice + 1) % len(doors)
	if newChoice == hostChoice {
		newChoice = (newChoice + 1) % len(doors)
	}
	switchWin := doors[newChoice].hasPrize

	// 发送结果
	resultChan <- GameResult{StickWin: stickWin, SwitchWin: switchWin}
}

func montyHallGameConcurrent(numGames int) {
	rand.Seed(time.Now().UnixNano())

	resultChan := make(chan GameResult, numGames)
	var wg sync.WaitGroup

	// 统计结果
	var stickWins, switchWins int

	// 启动goroutines来并行执行游戏
	for i := 0; i < numGames; i++ {
		wg.Add(1)
		go simulateGame(resultChan, &wg, doors) // 传递 doors 到函数
	}

	// 等待所有goroutines完成并收集结果
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.StickWin {
			stickWins++
		}
		if result.SwitchWin {
			switchWins++
		}
	}

	fmt.Printf("玩了 %d 次游戏后（并发模式）：\n", numGames)
	fmt.Printf("坚持最初的选择赢了 %d 次\n", stickWins)
	fmt.Printf("改变选择赢了 %d 次\n", switchWins)
}

func main() {
	montyHallGameConcurrent(541881452) // 运行100,000次游戏来观察结果
}
