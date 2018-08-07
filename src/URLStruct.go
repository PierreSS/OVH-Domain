package main

import (
	"sync"
	"time"
)

//Correspond à un fichier
type varStruct struct {
	tabConf 	*tabConfig
	tabRed  	*tabRedir
	tabSite		*tabSites
}

//Génére un objet grâce à URLConfig.yaml
type tabConfig struct {
	UrlApi		string `yaml:"api"`
	PingNsecondes    int `yaml:"pingNsecondes"`
	TictacHA    	int `yaml:"tictacHA"`
	TictacFailOver    int `yaml:"tictacFailOver"`
}

//Génére un objet grâce à URLSites.yaml
type tabSites struct {
    Adresse     []string `yaml:"adresse"`
}

//Génére un objet grâce à URLRedir.yaml
type tabRedir struct {
	Domaine		[]tabRedirDomaine
}

type tabRedirDomaine struct {
	Name	string
	Categorie []tabRedirCategorie
}

type tabRedirCategorie struct {
	Name	string
	Sousdomaine []tabRedirSousDomaine
}

type tabRedirSousDomaine struct {
	Name	string
	Type	string
	Cible []string
}

//Structure de récupération de dns
type ObjDNS struct {
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
	Zone      string `json:"zone"`
	FieldType string `json:"fieldType"`
	ID        int    `json:"id"`
	SubDomain string `json:"subDomain"`
}

//Le DNS est-il entrain de se supprimer
var inDesactivation = struct {
    sync.RWMutex
    myMap map[string]bool
}{myMap: make(map[string]bool)}

//Le Timer de réactivation est-il activer 
var  activationTimer = struct  {
	myMapTimer  map[string]*time.Timer
}{myMapTimer:  make(map[string]*time.Timer)}

//Check si verbose
var isVerbose bool