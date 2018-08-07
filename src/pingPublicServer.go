package main

import (
    "fmt"
	"time"
    "os/exec"
)

//Ping le domaine passé en parametre et renvoie true si il est pingable
func goPing(redirection string) (bool) {
	err := exec.Command("ping", "-c 3", "-q", "-w 3", redirection).Run()
	if err != nil {
        if isVerbose == true {
            fmt.Printf("[VERBOSE] : %s%s n'est pas pingable.\n", time.Now().Format("2006-01-02 15:04:05 : [Program] : "), redirection)
        }
		return false
	}
    if isVerbose == true {
        fmt.Printf("[VERBOSE] : %s%s est pingable.\n", time.Now().Format("2006-01-02 15:04:05 : [Program] : "), redirection)
    }
	return true
}

//Ping les serveurs séléctionnés dans URLSites
func goPingPublicServer(ptr *varStruct) (bool) {
    timer := time.NewTimer(time.Second * time.Duration(ptr.tabConf.PingNsecondes))
    for {
        compt := 0
        //Compte le nombre de domaine publiques à ping
        for i := 0; i < len(ptr.tabSite.Adresse); i++ {
            if goPing(ptr.tabSite.Adresse[i]) {
                if compt > (len(ptr.tabSite.Adresse) / 2) {
                    fmt.Printf("%s%d serveurs publics joignables\n", time.Now().Format("2006-01-02 15:04:05 : [Program] : "), compt)
                    timer.Stop()
                    return true
                }
                compt++
            }
        }
        //Si aucun serveur publiques ping une boucle de ping est créée
        <-timer.C
        timer.Reset(time.Second * time.Duration(ptr.tabConf.PingNsecondes))
        fmt.Printf("%s%d serveurs publics injoignables\n", time.Now().Format("2006-01-02 15:04:05 : [Program] : "), len(ptr.tabSite.Adresse) - compt)
        goGetFileContent(ptr)
    }
}
