package main

//Mes importations
import (
    "time"
    "fmt"
    "io/ioutil"
    "os"
    "os/signal"
    "syscall"
    "gopkg.in/yaml.v2"
)

//Set toutes les structures aux fichiers yaml
func goGetFileContent(ptr *varStruct) {
    //Récupère le fichier URLConfig.yaml
    fileConfig, erreur := ioutil.ReadFile("/etc/hadonis/URLConfig.yaml")
    if (erreur == nil) {
        err := yaml.Unmarshal(fileConfig, &ptr.tabConf)
        if (err != nil) {
            fmt.Printf(time.Now().Format("2006-01-02 15:04:05 : [Program] : Problème avec mise a jour du fichier URLConfig.yaml \n"))
            return
        }
    }
    //Récupère le fichier URLRedir.yaml
    fileRedir, erreur2 := ioutil.ReadFile("/etc/hadonis/URLRedir.yaml")
    if (erreur2 == nil) {
        err2 := yaml.Unmarshal([]byte(fileRedir), &ptr.tabRed)
        if (err2 != nil) {
            fmt.Printf(time.Now().Format("2006-01-02 15:04:05 : [Program] : Problème avec mise a jour du fichier URLRedir.yaml \n"))
            return
        }
    }
    //Récupère le fichier URLSites.yaml
    fileSites, erreur3 := ioutil.ReadFile("/etc/hadonis/URLSites.yaml")
    if (erreur3 == nil) {
        err3 := yaml.Unmarshal(fileSites, &ptr.tabSite)
        if (err3 != nil) {
            fmt.Printf(time.Now().Format("2006-01-02 15:04:05 : [Program] : Problème avec mise a jour du fichier URLSites.yaml \n"))
            return
        }
    }
    //Si une erreur arrive les fichiers anciennement valide sont utilisés
    if (erreur != nil || erreur2 != nil || erreur3 != nil) {
         return
     }
    fmt.Printf(time.Now().Format("2006-01-02 15:04:05 : [Program] : Fichiers de configuration mis à jour\n"))
}

//Créé une redirection 
func helpActiveDns(domaine, sousdomaine, cible, UrlApi, typee string) {
    if goPing(cible) {
        if !isDNSCreated(UrlApi, domaine, cible, sousdomaine) {
            fmt.Printf("%s[%s] Création du DNS : %s en cours\n", time.Now().Format("2006-01-02 15:04:05 : "), typee, cible)
            if create(UrlApi, domaine, sousdomaine, cible, typee) != true {
                 fmt.Printf("%s%s n'a pas pu être activé.\n", time.Now().Format("2006-01-02 15:04:05 : "), cible)
            }
        }
    }
}

//Active les DNS non activé qui ping
func activeDNSnotCreated(ptr *varStruct) {
    dom := ptr.tabRed.Domaine
    var a, b, c, d int = 0, 0, 0, 0

	if isVerbose == true {
        fmt.Printf("[VERBOSE] : %sVérifie si un DNS peut être activé.\n", time.Now().Format("2006-01-02 15:04:05 : [Program] : "))
    }
    //Check chaque redirection
    for ;a < len(dom); a++ {
        for ;b < len(dom[a].Categorie); b++ {
            for ;c < len(dom[a].Categorie[b].Sousdomaine); c++ {
                for ;d < len(dom[a].Categorie[b].Sousdomaine[c].Cible); d++ {
                    if dom[a].Categorie[b].Name == "HA" {
                        helpActiveDns(dom[a].Name, dom[a].Categorie[b].Sousdomaine[c].Name, dom[a].Categorie[b].Sousdomaine[c].Cible[d], ptr.tabConf.UrlApi, dom[a].Categorie[b].Sousdomaine[c].Type)
                    }
                    if dom[a].Categorie[b].Name == "Failover" {
                        if d == 0 {
                            helpActiveDns(dom[a].Name, dom[a].Categorie[b].Sousdomaine[c].Name, dom[a].Categorie[b].Sousdomaine[c].Cible[d], ptr.tabConf.UrlApi, dom[a].Categorie[b].Sousdomaine[c].Type)
                        }
                    }
                }
                d = 0
            }
            c = 0
        }
        b = 0
    }
}

//Main avec un compteur toute les minutes sur un thread
func main() {
    vStruct := varStruct{}
    //Initialisation globale
    goGetFileContent(&vStruct)

    //Récupère l'argument et set sa variable globale 
    isVerbose = false
    arg1 := ""
    if len(os.Args) > 1 {
        arg1 = os.Args[1]
        if arg1 == "-v" || arg1 == "--verbose" {
            isVerbose = true
        }
    }

    //active les DNS qui ping
    activeDNSnotCreated(&vStruct)

    //Création des comptes à rebours pour le check des fichiers, ping des adresses publiques et des dns
    tickFile := time.NewTicker(time.Minute * 1).C
    tickPing := time.NewTicker(time.Second * time.Duration(vStruct.tabConf.PingNsecondes)).C
    GlobalTick := time.NewTicker(time.Second).C

    //Création d'une channel continue pour catch un signal
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
    //Catch un signal
    go goCatchSignal(c, &vStruct)

    //Boucle principal de tick par seconde
    for _ = range GlobalTick {
        select {
            case <-tickFile:
                go goGetFileContent(&vStruct)
            case <-tickPing:
                if goPingPublicServer(&vStruct) {
			tickPing = time.NewTicker(time.Second * time.Duration(vStruct.tabConf.PingNsecondes)).C
		}
            //Si les serveurs publiques ping bien, le ping des dns est lançé
            pingDNS(&vStruct)
            default:
        }
    }
}
