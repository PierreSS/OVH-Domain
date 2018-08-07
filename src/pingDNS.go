package main

import (
    "fmt"
    "encoding/json"
	"net/http"
    "io/ioutil"
    "time"
)

//Renvoie true si le DNS est créé et false si il ne l'est pas
func isDNSCreated(api, domaine, cible, sousdomaine string) (bool) {
    //récupère la liste des dns pour un domaine
	response, errer := http.Get(api + "dns/" + domaine)
	if errer != nil {
		if isVerbose == true {
            fmt.Printf("[VERBOSE] : %s Erreur en récupérant qu'un dns est créé pour %s.\n", time.Now().Format("2006-01-02 15:04:05 : [Program] : "), cible)
        }
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			if isVerbose == true {
                fmt.Printf("[VERBOSE] : %s Erreur en récupérant qu'un dns est créé pour %s.\n", time.Now().Format("2006-01-02 15:04:05 : [Program] : "), cible)
            }
		}
        //ferme la réponse à la fin de la fonction
        defer response.Body.Close()

        //Remplissage de la structure YAML
		var dns []ObjDNS
		json.Unmarshal(body, &dns)
		for i := 0; i < len(dns); i++ {
			if ((dns[i].Target == cible) && (dns[i].SubDomain == sousdomaine)) {
                if isVerbose == true {
                    fmt.Printf("[VERBOSE] : %sDNS est créé pour %s.\n", time.Now().Format("2006-01-02 15:04:05 : [Program] : "), cible)
                }
				return true
			}
		}
	}
	return false
}

//Active un compte à rebours pour la catégorie HA
func activeCAR(vStruct *varStruct, domaine, sousdomaine, cible, typee string) {
    //Mutex locked pour dire que la redirection est en cours de paramétrage
    inDesactivation.RLock()
    inDesactivation.myMap[sousdomaine + cible] = true
    inDesactivation.RUnlock()
    fmt.Printf("%sCompte à rebours déclenché pour remise en route du DNS %s\n", time.Now().Format("2006-01-02 15:04:05 : "), cible)

    //Création d'un timer en attente de la réativation
    timer := time.NewTimer(time.Second * time.Duration(vStruct.tabConf.TictacHA))
    for {
        //quand le timer fini
        <-timer.C
        if goPing(cible) {
            fmt.Printf("%sTentative de re-création du DNS : %s en cours\n", time.Now().Format("2006-01-02 15:04:05 : [HA] : "), cible)
            //Création d'un timer pour la réactivation
            activationTimer.myMapTimer[sousdomaine + cible]  =  time.NewTimer(time.Second  *  time.Duration(vStruct.tabConf.TictacHA))
            <-activationTimer.myMapTimer[sousdomaine + cible].C
            activationTimer.myMapTimer[sousdomaine + cible] = nil
            if goPing(cible) {
                if create(vStruct.tabConf.UrlApi, domaine, sousdomaine, cible, typee) != false {
                    //Mutex locked pour dire que la redirection est libre d'être tester à nouveau
                    inDesactivation.RLock()
                    inDesactivation.myMap[sousdomaine + cible] = false
                    inDesactivation.RUnlock()
                    return                    
                }
            }
            fmt.Printf("%s%s n'a pas pu être réactivé.\n", time.Now().Format("2006-01-02 15:04:05 : [HA] : "), cible)
        } else { 
            fmt.Printf("%s%s n'est toujours pas joignable.\n", time.Now().Format("2006-01-02 15:04:05 : [HA] : "), cible)
            fmt.Printf("%sReset du timer de réactivation de %s : %s en cours\n", time.Now().Format("2006-01-02 15:04:05 : [HA] : "), sousdomaine, cible)
            timer.Reset(time.Second * time.Duration(vStruct.tabConf.TictacHA))
        }
    }
}

//Active un compte à rebours pour la catégorie failOver
func activeCARfo(vStruct *varStruct, domaine, sousdomaine, cible, redirection string) {
    //Mutex locked pour dire que la redirection est en cours de paramétrage
    inDesactivation.RLock()
    inDesactivation.myMap[sousdomaine + cible] = true
    inDesactivation.RUnlock()
    fmt.Printf("%sCompte à rebours déclenché pour remise en route du DNS %s\n", time.Now().Format("2006-01-02 15:04:05 : [Failover] : "), cible)

    //Création d'un timer en attente de la réativation
    timer := time.NewTimer(time.Second * time.Duration(vStruct.tabConf.TictacHA))
    for {
        //quand le timer fini
        <-timer.C
        if goPing(cible) {
            fmt.Printf("%sReswitching %s à son maitre : %s en cours\n", time.Now().Format("2006-01-02 15:04:05 : [Failover] : "), redirection, cible)
            //Création d'un timer pour la réactivation
            activationTimer.myMapTimer[sousdomaine + cible]  =  time.NewTimer(time.Second  *  time.Duration(vStruct.tabConf.TictacHA))
            //timer.Reset(time.Second * time.Duration(vStruct.tabConf.TictacHA))
            <-activationTimer.myMapTimer[sousdomaine + cible].C
            activationTimer.myMapTimer[sousdomaine + cible] = nil
            if goPing(cible) {
                if put(vStruct.tabConf.UrlApi, domaine, sousdomaine, redirection, cible) != false {
                    //Mutex locked pour dire que la redirection est libre d'être tester à nouveau
                    inDesactivation.RLock()
                    inDesactivation.myMap[sousdomaine + cible] = false
                    inDesactivation.RUnlock()
                    return
                }
            }
            //Mutex locked pour dire que la redirection est libre d'être tester à nouveau
            fmt.Printf("%s%s n'a pas pu être réactivé.\n", time.Now().Format("2006-01-02 15:04:05 : [Failover] : "), cible)
        } else {
            fmt.Printf("%s%s n'est toujours pas pingable.\n", time.Now().Format("2006-01-02 15:04:05 : [Failover] : "), cible)
            fmt.Printf("%sReset du timer de réactivation de %s : %s en cours\n", time.Now().Format("2006-01-02 15:04:05 : [Failover] : "), redirection, cible)
            timer.Reset(time.Second * time.Duration(vStruct.tabConf.TictacHA))
        }
    }
}


//Check de la catégorie HA pour un domaine
func HA(vStruct *varStruct, cible, domaine, sousdomaine, typee string) {
    //Check si la redirection n'est pas en cours de paramétrage
    inDesactivation.RLock()
    n := inDesactivation.myMap[sousdomaine + cible]
    inDesactivation.RUnlock()
    if n == false {
        if !goPing(cible) {
            if isDNSCreated(vStruct.tabConf.UrlApi, domaine, cible, sousdomaine) {
                //Si la redirection n'est plus pingable et que le DNS est créé
                compt := 0
                    for i := 0; i < 3; i++ {
                        if !goPing(cible) {
                            compt++
                            if compt == 3 {
                            //Si il n'arrive pas à ping 3 fois                                
                                fmt.Printf("%s%s.%s Tentative de suppression %s en cours.\n", time.Now().Format("2006-01-02 15:04:05 : [HA] : "), sousdomaine, domaine, cible)
                                if delete(vStruct.tabConf.UrlApi, domaine, cible, sousdomaine) != false {
                                    go activeCAR(vStruct, domaine, sousdomaine, cible, typee)                                
                                }
                                fmt.Printf("%s%s.%s Suppression %s effectuée.\n", time.Now().Format("2006-01-02 15:04:05 : [HA] : "), sousdomaine, domaine, cible)
                            }
                        } else {
                            fmt.Printf("%s%s.%s ne supprime pas %s car au moins un ping sur trois a réussi.\n", time.Now().Format("2006-01-02 15:04:05 : [HA] : "), sousdomaine, domaine, cible)              
                        }
                    }
            }
        }
    }
    if n == true {
        if !goPing(cible) {
            //Si la redirection est en cours de redémarrage et que la cible ne ping pas le timer de redémarrage est reset
            if activationTimer.myMapTimer[sousdomaine + cible] != nil {
                activationTimer.myMapTimer[sousdomaine + cible].Reset(time.Second * time.Duration(vStruct.tabConf.TictacHA))
                fmt.Printf("%sLe compte a rebours a été réinitialisé à cause d'une mauvaise connexion à %s .\n", time.Now().Format("2006-01-02 15:04:05 : [HA] : "), cible)
            }
        }
    }
}

//Check de la catégorie failOver pour un domaine
func failOver(vStruct *varStruct, domaine, sousdomaine, cible, redirection string) {
    //Check si la redirection n'est pas en cours de paramétrage
    inDesactivation.RLock()
    n := inDesactivation.myMap[sousdomaine + cible]
    inDesactivation.RUnlock()
    if n == false {
        if !goPing(cible) {
                if isDNSCreated(vStruct.tabConf.UrlApi, domaine, cible, sousdomaine) {
                    //Si la redirection n'est plus pingable et que le DNS est créé
                    compt := 0
                        for i := 0; i < 3; i++ {
                            if !goPing(cible) {
                                compt++
                                if compt == 3 {
                                    //Si il n'arrive pas à ping 3 fois 
                                    compt = 0
                                    if goPing(redirection) {
                                        fmt.Printf("%s%s.%s bascule de %s sur l'esclave %s en cours.\n", time.Now().Format("2006-01-02 15:04:05 : [Failover] : "), sousdomaine, domaine, cible, redirection)
                                        if put(vStruct.tabConf.UrlApi, domaine, sousdomaine, cible, redirection) != false {
                                            go activeCARfo(vStruct, domaine, sousdomaine, cible, redirection)                                        
                                        }
                                    } else {
                                        fmt.Printf("%s%s.%s, impossible de basculer %s car %s n'est pas joignable.\n", time.Now().Format("2006-01-02 15:04:05 : [Failover] : "), sousdomaine, domaine, cible, redirection)
                                    }
                                }
                            } else {
                                fmt.Printf("%s%s.%s ne supprime pas %s car au moins un ping sur trois a réussi ou %s ne répond pas.\n", time.Now().Format("2006-01-02 15:04:05 : [HA] : "), sousdomaine, domaine, cible, redirection)                    
                            }
                        }
                }
        }
    }
    if n == true {
        if !goPing(cible) {
            //Si la redirection est en cours de redémarrage et que la cible ne ping pas le timer de redémarrage est reset
            if activationTimer.myMapTimer[sousdomaine + cible] != nil {
                activationTimer.myMapTimer[sousdomaine + cible].Reset(time.Second * time.Duration(vStruct.tabConf.TictacHA))
                fmt.Printf("%ssLe compte a rebours a été réinitialisé à cause d'une mauvaise connexion à %s.\n", time.Now().Format("2006-01-02 15:04:05 : [Failover] : "), cible)
            }
        }
    }
}

//Ping tout les DNS de chaque domaine 
func pingDNS(vStruct *varStruct) () {
    dom := vStruct.tabRed.Domaine
    var a, b, c, d int = 0, 0, 0, 0

    //Check chaque redirection
    for ;a < len(dom); a++ {
        for ;b < len(dom[a].Categorie); b++ {
            for ;c < len(dom[a].Categorie[b].Sousdomaine); c++ {
                for ;d < len(dom[a].Categorie[b].Sousdomaine[c].Cible); d++ {
                    //La réactivation n'est pas la même si la catégorie est HA ou FailOver
                    if dom[a].Categorie[b].Name == "HA" {
                        HA(vStruct, dom[a].Categorie[b].Sousdomaine[c].Cible[d], dom[a].Name, dom[a].Categorie[b].Sousdomaine[c].Name, dom[a].Categorie[b].Sousdomaine[c].Type)
                    }
                    if dom[a].Categorie[b].Name == "Failover" {
                        if d == 0 {
                            failOver(vStruct, dom[a].Name, dom[a].Categorie[b].Sousdomaine[c].Name, dom[a].Categorie[b].Sousdomaine[c].Cible[d], dom[a].Categorie[b].Sousdomaine[c].Cible[d + 1])
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
