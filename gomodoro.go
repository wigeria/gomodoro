package main

import (
    "flag"
    "fmt"
    "time"
    "runtime"
    "log"
    "os/exec"
)

type Settings struct {
    taskName string
    workTime string
    restTime string
    repeat int
}

func openBrowser(url string) {
    // Taken straight off from https://gist.github.com/hyg/9c4afcd91fe24316cbf0 
    var err error

    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        err = fmt.Errorf("unsupported platform")
    }
    if err != nil {
        log.Fatal(err)
    }
}

func parseSettings () *Settings {
    // Parses all of the CLI flags through the `flags` library
    settings := new(Settings)

    flag.StringVar(&settings.taskName, "taskName",
        "Unnamed", "Name for the Task")
    flag.StringVar(&settings.workTime, "work",
        "25m", "Work Duration (as ParseDuration args)")
    flag.StringVar(&settings.restTime, "rest",
        "5m", "Rest Duration (as ParseDuration args)")
    flag.IntVar(&settings.repeat, "reps",
        1, "Number of times to repeat the pomodoro")

    flag.Parse()

    return settings
}

func tickUntilTime(ticker *time.Ticker, until time.Time) {
    // Ticks until the given time while printing out the ticker intervals
    for {
        tick := <-ticker.C
        timeLeft := until.Sub(tick).Round(time.Second)
        // Timer that prints over the same line; the spaces are to accomodate
        // different string lengths
        fmt.Printf("\r%s          ", timeLeft.String())
        if timeLeft.Seconds() == 0 {
            break
        }
    }
}

func main () {
    tickInterval := time.Second
    settings := parseSettings()

    workDur, _ := time.ParseDuration(settings.workTime)
    restDur, _ := time.ParseDuration(settings.restTime)
    // Initiating the ticker instance which includes it's interval channel
    ticker := time.NewTicker(tickInterval)
    defer ticker.Stop()

    fmt.Printf("Now working on %s...", settings.taskName)
    for reps := 0; reps < settings.repeat; reps++ {
        // Estimating the finish times for both the work and rest for this rep
        workFinishTime := time.Now().Add(workDur)
        restFinishTime := workFinishTime.Add(restDur)

        fmt.Println("\nWork now!")
        tickUntilTime(ticker, workFinishTime)

        openBrowser("https://www.reddit.com/r/aww/")
        fmt.Println("\nRest Now!")
        tickUntilTime(ticker, restFinishTime)
    }

    fmt.Println("\nSuccessfully finished all reps!")
}
