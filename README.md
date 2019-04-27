/* Made by Pierre */

/* INSTALLING GO PKG */
apt-get install golang

Pour installer les dépendances: make deps
Pour compiler: make

/* HOW TO HAVE LOG */
sudo chmod -R a+rw /var/log/hadonis/


/* Notice d'utilisation */
- Les trois fichiers de configuration doivent se trouver dans /etc/hadonis/ et doivent se nommer URLConfig.yaml, URLRedir.yaml, URLSites.yaml
- TYPENAME peut être soit A, soit AAAA soit CNAME.
- CATEGORIENAME peut être soit Failover soit HA avec HA ayant que des adresses de type A ou AAAA et CNAME ayant des adresses de type CNAME.
- Le fichier de redirection doit être syntaxé comme ceci pour un seul domaine, une seule catégorie, un seul sous-domaine et une seule redirection :

domaine:
   -
     name: DOMAINNAME
     categorie:
        -
          name: CATEGORIENAME
          sousdomaine:
            -
              name: SOUSDOMAINENAME
              type: TYPENAME
              cible:
                - "REDIRECTION"


ATTENTION : Le yaml ne supporte pas les tabulations, veuillez donc utilisez uniquement des espaces lors de sa configuration.



/* Notice d'utilisation avancé */

- Ajouter un domaine -

domaine:
   -
     name: DOMAINNAME
     categorie:
        -
          name: CATEGORIENAME
          sousdomaine:
            -
              name: SOUSDOMAINENAME
              type: TYPENAME
              cible:
                - "REDIRECTION"
   -
     name: DOMAINNAME 2
     categorie:
        -
          name: CATEGORIENAME
          sousdomaine:
            -
              name: SOUSDOMAINENAME
              type: TYPENAME
              cible:
                - "REDIRECTION"


- Ajouter une categorie -

domaine:
   -
     name: DOMAINNAME
     categorie:
        -
          name: CATEGORIENAME
          sousdomaine:
            -
              name: SOUSDOMAINENAME
              type: TYPENAME
              cible:
                - "REDIRECTION"
        -
          name: CATEGORIENAME 2
          sousdomaine:
            -
              name: SOUSDOMAINENAME
              type: TYPENAME
              cible:
                - "REDIRECTION"


- Ajouter un sousdomaine -

domaine:
   -
     name: DOMAINNAME
     categorie:
        -
          name: CATEGORIENAME
          sousdomaine:
            -
              name: SOUSDOMAINENAME
              type: TYPENAME
              cible:
                - "REDIRECTION"
            -
              name: SOUSDOMAINENAME 2
              type: TYPENAME
              cible:
                - "REDIRECTION"


- Ajouter une redirection -

domaine:
   -
     name: DOMAINNAME
     categorie:
        -
          name: CATEGORIENAME
          sousdomaine:
            -
              name: SOUSDOMAINENAME
              type: TYPENAME
              cible:
                - "REDIRECTION"
                - "REDIRECTION 2"

