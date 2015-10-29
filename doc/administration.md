Administration
=============
Close server
------------

Close the webserver :

- **Method** : GET
- **URL** : /close
- **Headers** :
    * **Authorization** : Token {{token}}

The response will be :

- **Status** : 200
- **Headers** :
    * **Access-Control-Allow-Headers** : accept, authorization
    * **Access-Control-Allow-Methods** : GET, POST
    * **Access-Control-Allow-Origin** : \*
    * **Content-Type** : application/json
- **Data** :
    ```json
    {
        "token":"55f13bcb744ca96b7ff8fda2492c598d",
        "message":"Closing server"
    }
    ```
