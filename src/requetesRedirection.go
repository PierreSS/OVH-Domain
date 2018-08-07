package main

import (
	"fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "os"
    "bufio"
    "time"
)

//fonction de checking d'erreur de writeFile
func checkError(err error, str string) {
    if err != nil {
        fmt.Printf("Error %s\n", str)
    }
}

//Récupère la réponse de l'apidns et l'écrit dans les logs
func writeFile(body string) {
    //Check body of request

   /* client := &http.Client{}
    resp, _ := client.Do(response)
    body, _ := ioutil.ReadAll(resp.Body)*/
    fmt.Printf("%s%s\n", time.Now().Format("2006-01-02 15:04:05 : "), body)

    //Create & write EOF
    file, err := os.OpenFile("/var/log/hadonis/hadonis.log", os.O_APPEND|os.O_RDWR, 0666)        
    if err != nil {
        file, err = os.Create("/var/log/hadonis/hadonis.log")
        checkError(err, "erreur de création/écriture dans /var/log/hadonis/hadonis.log")
    }

    //Count Lines of log & erase if x lines
    f, _ := os.Open("log")
    fileScanner := bufio.NewScanner(f)
    lineCount := 0
    for fileScanner.Scan() {
        lineCount++
    }
    //Recréer le fichier à plus de 20000 lignes
    if lineCount >= 20000 {
        err := os.Remove("/var/log/hadonis/hadonis.log")
        checkError(err, "erreur d'effacement de /var/log/hadonis/hadonis.log")
        file, err = os.Create("log")
        checkError(err, "erreur de recréation de /var/log/hadonis/hadonis.log")
        defer f.Close()
        fmt.Printf("%d\n", lineCount)
    }

    //Write body in file
    _, err = file.WriteString(body)
    checkError(err, "erreur d'écriture dans /var/log/hadonis/hadonis.log (body)")
    _, err = file.WriteString("\n")
    checkError(err, "erreur d'écriture dans /var/log/hadonis/hadonis.log")
    //Synchronisation des fichiers et fermement de la requête
    file.Sync()
    defer file.Close()
}

//Créer une requête de création de redirection
func create(api, domaine, sousdom, redirection, typee string) (bool) {
    checkGood := false
    response, err := http.NewRequest("POST", api + "dns/" + domaine + "/" + sousdom + "/" + redirection + "/" + typee, nil)
    checkError(err, "request")
    client := &http.Client{}
    resp, _ := client.Do(response)
    body, _ := ioutil.ReadAll(resp.Body)
    
//    fmt.Printf("%s\n", strings.Contains(string(body), "La redirection"));
    if strings.Contains(string(body), "La redirection") == true  {
        checkGood = true
    }
    writeFile(string(body))
    return checkGood
}

//modifie
func put(api, domaine, sousdom, cible, redirection string) (bool) {
    checkGood := false
    response, err := http.NewRequest("PUT", api + "dns/" + domaine + "/" + sousdom + "/" + cible + "/" + redirection, nil)
    checkError(err, "request")
    client := &http.Client{}
    resp, _ := client.Do(response)
    body, _ := ioutil.ReadAll(resp.Body)
    if strings.Contains(string(body), "La redirection") == true  {
        checkGood = true
    }
    writeFile(string(body))
    return checkGood
}

//Créer une requête de suppression de redirection
func delete(api, domaine, redirection, sousdom string) (bool) {
    checkGood := false
    response, err := http.NewRequest("DELETE", api + "dns/" + domaine + "/" + redirection + "/" + sousdom, nil)
    checkError(err, "request")
    client := &http.Client{}
    resp, _ := client.Do(response)
    body, _ := ioutil.ReadAll(resp.Body)
    if strings.Contains(string(body), "La redirection") == true {
        checkGood = true
    }
    writeFile(string(body))
    return checkGood
}
