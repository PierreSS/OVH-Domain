package main

import (
    "os"
    "fmt"
    "time"
)

//Nettoie les structures en cas de signal d'interreuption
func clear(v *varStruct) {
    v.tabConf = nil
    v.tabRed = nil
    v.tabSite = nil
}

//Catch un signal
func goCatchSignal(c chan os.Signal, vStruct *varStruct) {
    sig := <-c
    clear(vStruct)
    fmt.Printf("\n%sSortie de programme suite à %s\n", time.Now().Format("2006-01-02 15:04:05 : [Program] : "), sig)
    os.Exit(1)
}
