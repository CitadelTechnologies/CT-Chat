Chatrooms
=============
Get chatroom
------------

Get data about the current chatroom :

- **Method** : GET
- **URL** : /
- **Headers** :
    * **Authorization** : Token {{token}}
    * **X-Origin** : [client-website authorized domain]

The response will be :

- **Status** : 200
- **Headers** :
    * **Access-Control-Allow-Headers** : accept, authorization
    * **Access-Control-Allow-Methods** : GET, POST
    * **Access-Control-Allow-Origin** : [client-website authorized domain]
    * **Content-Type** : application/json
- **Data** :
    ```json
    {
        "token":"55f13bcb744ca96b7ff8fda2492c598d",
        "chatroom":{
            "name":"main",
            "users":[
                {"username":"Kern"}
            ],
            "messages":[
                {
                    "type":"message",
                    "author":"Kern",
                    "content":"test",
                    "chatroom":"main"
                },
                {
                    "type":"message",
                    "author":"Kern",
                    "content":"again",
                    "chatroom":"main"
                }
            ]
        }
    }
    ```

If the domain is missing or not authorized, the following error will be sent :

- **Status** : 403
- **Headers** :
    * **Content-Type** : application/json
- **Data** :
    ```json
    {
        "message": "This domain is not authorized"
    }
    ```