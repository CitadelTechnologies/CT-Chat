Chatrooms
=============
Get chatroom
------------

Get data about the current chatroom :

- **Method** : GET
- **URL** : /
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
