package main

import (
    "http-project/nhlApi"
    "io"
    "log"
    "os"
    "sync"
    "time"
)

func main() {
    now := time.Now()
    rosterFile, err := os.OpenFile("rosters.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening the file rosters.txt: %v", err)
    }

    defer rosterFile.Close()

    wrt := io.MultiWriter(os.Stdout, rosterFile)

    log.SetOutput(wrt)

    teams, err := nhlApi.GetAllTeams()
    if err != nil {
        log.Fatalf("Error getting all the teams: %v", err)
    }

    var wg sync.WaitGroup

    wg.Add(len(teams))

    results := make(chan []nhlApi.Roster, 3)
    for _, team := range teams {
        go func(team nhlApi.Team) {
            roster, err := nhlApi.GetRosters(team.ID)
            if err != nil {
                log.Fatalf("error getting roster: %v", err)
            }
            results <- roster

            wg.Done()
        }(team)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    display(results)

    log.Printf("took %v", time.Now().Sub(now).String())
}

func display(results chan []nhlApi.Roster) {
    for r := range results {
        for _, ros := range r {
            log.Println("----------------")
            log.Printf("ID: %d\n", ros.Person.ID)
            log.Printf("Name: %s\n", ros.Person.FullName)
            log.Printf("Position: %s\n", ros.Position.Abbreviation)
            log.Printf("Jersey: %s\n", ros.JerseyNumber)
            log.Println("----------------")
        }
    }
}
