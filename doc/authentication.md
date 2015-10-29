Authentication
=============
Introduction
------------

To authenticate to the chat, you need to send the following request :

- **Method** : POST
- **URL** : /
- **Headers** :
    * **Content-Type** : application/x-www-form-urlencoded; charset=UTf-8
- **Data Params** : username=[string]

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
    	"token": "sq6d23qg1htr1sfd1f25ez",
        "message": "You are connected"
    }
    ```
In case of error, the response will be :

- **Status** : 401
- **Headers** :
    * **Access-Control-Allow-Headers** : accept, authorization
    * **Access-Control-Allow-Methods** : GET, POST
    * **Access-Control-Allow-Origin** : \*
    * **Content-Type** : application/json
- **Data** :
	```json
    {
        "message": "You are not connected"
    }
    ```
