# GSheet-Proxy
## Introduction  
A RestAPI proxy to Google Sheet (Golang). This is a WebHook Proxy. This golang app will act as a proxy.   
 - It will pre-load a remote database (using REST API) and cache locally.   
 - It then quickly return lookup info       
   
  
## Usage  
### Install Golang
   Download source: https://golang.org/doc/install 
   Config:
   ```
   export GOROOT=$HOME/Data/Softwares/go  # Path to go source directory
   export GOPATH=$HOME/go                 # Path to working directory
   export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

   ```
   
### Install dependencies

    ```
    go mod vendor
    ```

### Combine

    ./release.sh

  
Combined artifacts will be copied into **release** folder.  
  
### Deploy  
 - Copy **release** folder to the deployed machine.  
 - Run **mkcert_auto_install.bat** to deploy RootCA to your machine  
 - Modify **config.yaml** to match with your Goolge Sheet REST API  
 - Run appropriate binary file to execute EventHub.  
  
### Config EventHub  
**config.yaml**

    server: ":8443"  
      
    storage:  
      Adapter: mem  
      Folder: ./data  
      
    dbsources:  
     - Name: nestle  
        FetchingUrl: https://script.google.com/macros/s/AKfycbzMagS4EswslEawGoZg-HOilKozja6tFWbDDgt9e-hdVipYfQ/exec?path=/checkin&order=QRCode&offset=%offset%&limit=%size%  
        FetchingFormat: json  
        UpdateUrl: https://script.google.com/macros/s/AKfycbzMagS4EswslEawGoZg-HOilKozja6tFWbDDgt9e-hdVipYfQ/exec?path=/checkin/%key%&method=POST  
        UpdateMethod: GET  
        IdField: QRCode  

  
**server**: is the listening address of the server  
**dbsources**: is a collection of proxied-database  
 - name: is the unique name of the db  
 - FetchingUrl: Google Sheet REST (fetch) url  (with %offset% & %size% options)
 - UpdateUrl: Google Sheet REST (push) url (with %key% options - the IdField value)
 - IdField: is the column-name of the unique value field
 
**storage**: currently support **mem** and **bow**  
 - for most machines, **mem** storage is the best choice. But data will not be persisted to disk.  
 - for high memory machines, **bow** storage is the best choice and data will be persisted to disk.  
  
### Create Self-Signed certificate for EventHub  

    # set CAROOT to current folder
    export CAROOT=$(pwd)
    
    # install CARoot to local machine
    mkcert install -CAROOT  
    
    # create new certificate for with custom domain & ip using current CARoot
    mkcert <replaced-local-domain-name> <replaced-your-machine-IP> localhost 127.0.0.1  

Rename the newly generated files to **eventHub.pem** & **eventHub-key.pem** and copy to **releases** folder

### WebHook API for EventHub  (JSON Response)

 - GET **/api/db/:dbName** show **:dbName** content
 - GET **/api/db/:dbName/import** send a trigger to Start import the database **:dbName**
 - GET **/api/item/:dbName/:itemKey** show an item detail with **:itemKey** on **:dbName**
 - GET **/api/qr-check/:dbName/:itemKey**  OR **/qr-check/:dbName**
// url: /qr-check/:dbName?key=%qrData%&activityName=checkin&gateway=cong2&key2=val2  
// params:  
//     key: unique id  
//     activityName: checkin
//     extra_fields: key1=val1,key2=val2

### Admin Url for EventHub  (HTML Response)

 - GET **/admin/db/:dbName** show **:dbName** content  (beautiful table)
 - GET **/admin/item/:dbName/:itemKey** show an item detail with **:itemKey** on **:dbName**
 - GET **/admin/qr-check/:dbName/:itemKey**  OR **admin/qr-check/:dbName** , scan, check and show item detail
// url: /qr-check/:dbName?Key=%qrData%&activityName=asdfadsf&key1=val1&key2=val2  
// params:  
//     Key: unique id  
//     activityName: xyz  
//     extra_fields: key1=val1,key2=val2

  
### Https self-signed Certificate using Mkcert  
 - https://12bit.vn/articles/tao-https-cho-localhost-su-dung-mkcert/
 - https://github.com/FiloSottile/mkcert

 