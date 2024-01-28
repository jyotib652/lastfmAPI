## lastFMAPI
It is an API which provides Top Song Track from a specific country, Lyrics for the searched Song and Artist's Name, Information & Images.
Depending upon the searched song, artist & lyrics it would also suggest some other song tracks also.

**Requirements**:
* To get the proper response from this api, you have to provide a Country Name. But Country name should be as defined by the ISO 3166-1 country names standard.


**To Run this Golang application** :

1. If you want to run the application in a Docker container then then go to the location where you cloned this git repo and then
   use the following commands to create and run the Docker container:
   
   Suppose you cloned this git repo in your Home directory then run
      ```
   $ cd ~
   $ make build_lastfm
   $ make up      // or, make up_build
      
   ```
   Now the Docker Container is build and running. Next, you have to send request to this Docker Container API.

### Example:
To send request to this API use following examples:

[http://localhost:8080/<"country">](http://localhost:8080/<"country">)

Replcae <"country"> with a country name like:

[http://localhost:8080/United%20States](http://localhost:8080/United%20States)

[http://localhost:8080/Spain](http://localhost:8080/Spain)

[http://localhost:8080/United%20Kingdom](http://localhost:8080/United%20Kingdom)

You can use Web Browser or Postman or Thunder Client or other available such tools to send request to this api.
