package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const REPLAY_GAME int = 1
const EXIT_GAME int = 2

const EASY_DIFFICULT int = 1
const MEDIUM_DIFFICULT int = 2
const HARD_DIFFICULT int = 3

const HOT int = 5
const WARM int = 15

func printDifficultyMenu() {
	fmt.Println("Выберите сложность игры:")
	fmt.Printf("%v. Лёгкая\n", EASY_DIFFICULT)
	fmt.Printf("%v. Средняя\n", MEDIUM_DIFFICULT)
	fmt.Printf("%v. Сложная\n", HARD_DIFFICULT)
}

func chooseDifficuilt() (int, int) {
	var input string

	for {
		fmt.Print("Ваш выбор: ")
		fmt.Scan(&input)
		choose, err := strconv.Atoi(input)

		if err != nil {
			fmt.Println("Нужно ввести целое число! (1, 2, ...)")
			continue
		}

		if choose < 0 || choose > HARD_DIFFICULT {
			fmt.Printf("Можно выбрать %v, %v или %v пункт\n", EASY_DIFFICULT, MEDIUM_DIFFICULT, HARD_DIFFICULT)
			continue
		}

		switch choose {
		case EASY_DIFFICULT:
			return 50, 15

		case MEDIUM_DIFFICULT:
			return 100, 10

		case HARD_DIFFICULT:
			return 200, 5

		}

	}
}

func startGame(maxNum, maxAttemps int) {
	randomInt := rand.Intn(maxNum) + 1
	lastInputs := []int{}

	var number int
	var isFound bool

	fmt.Printf("Игра 'Угадай число' - от 1 до %v началась!\n", maxNum)
	fmt.Printf("Угадай число за %v попыток!\n", maxAttemps)

	for i := 0; i < maxAttemps; i++ {
		fmt.Printf("Попытка %v\n", i+1)

		printLastNums(lastInputs)

		number = inputNum(maxNum)
		lastInputs = append(lastInputs, number)

		if number > randomInt {
			fmt.Println("Секретное число меньше👇")

		} else if number < randomInt {
			fmt.Println("Секретное число больше👆")

		} else if number == randomInt {
			fmt.Println("\x1b[32mВы угадали!\x1b[0m") // green
			fmt.Println("Игра закончена!")
			isFound = true
			saveToFile(isFound, i)
			break
		}

		printHelp(number, randomInt)
	}

	if !isFound {
		fmt.Println("\x1b[35mВы проиграли!\x1b[0m") // red
		fmt.Println("Секретное число было: ", randomInt)
		saveToFile(isFound, maxAttemps)
	}
}

type Record struct {
	Date    string `json:"date"`
	Result  string `json:"result"`
	Attemps int    `json:"attemps"`
}

func saveToFile(isFound bool, attemps int) {
	filename := "result.json"

	winOrLose := "lose"
	if isFound {
		winOrLose = "win"
	}

	record := Record{
		Date:    time.Now().String(),
		Result:  winOrLose,
		Attemps: attemps,
	}

	var records []Record

	if _, err := os.Stat(filename); err == nil {
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		if len(data) > 0 {
			err = json.Unmarshal(data, &records)
		}
	}

	records = append(records, record)

	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("result.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}
}

func printLastNums(lastInputs []int) {
	if len(lastInputs) == 0 {
		return
	}

	fmt.Print("Последнии введённые числа: ")
	for _, num := range lastInputs {
		if num != 0 {
			fmt.Print(num, " ")
		}
	}
	fmt.Println()
}

func printHelp(num, randNum int) {
	abs := num - randNum
	abs = max(abs, -abs)

	if abs <= HOT {
		fmt.Println("🔥 Горячо")
		return
	}

	if abs <= WARM {
		fmt.Println("🙂 Тепло")
		return
	}

	fmt.Println("❄️  Холодно")
}

func inputNum(maxNum int) int {
	var input string

	for {
		fmt.Print("Введите число: ")
		fmt.Scan(&input)
		number, err := strconv.Atoi(input)

		if err != nil {
			fmt.Println("Нужно ввести целое число! (1, 2, ...)")
			continue
		}

		if number < 1 || number > maxNum {
			fmt.Printf("Число должно быть в диапозоне от 1 до %v!\n", maxNum)
			continue
		}
		return number
	}
}

func printExitMenu() {
	fmt.Printf("%v. Сыграть повторно\n", REPLAY_GAME)
	fmt.Printf("%v. Выйти\n", EXIT_GAME)
}

func isGameExit() bool {
	var input string
	for {
		fmt.Print("Ваш выбор: ")
		fmt.Scan(&input)

		choose, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Нужно ввести целое число! (1, 2, ...)")
			continue
		}

		if choose < 0 || choose > EXIT_GAME {
			fmt.Printf("Можно выбрать %v или %v пункт\n", REPLAY_GAME, EXIT_GAME)
			continue
		}

		switch choose {
		case REPLAY_GAME:
			return false

		case EXIT_GAME:
			return true

		default:
			return false
		}
	}
}

func gameLoop(gameIsRunning bool) {
	for gameIsRunning {
		printDifficultyMenu()
		maxNum, maxAttemps := chooseDifficuilt()
		startGame(maxNum, maxAttemps)
		printExitMenu()
		if isGameExit() {
			gameIsRunning = false
		}
	}

}

func main() {
	gameIsRunning := true
	gameLoop(gameIsRunning)
}
