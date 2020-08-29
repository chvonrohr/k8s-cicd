# Check Setup
# How the training works

----

## v1: Web IDE
## v2: local setup

----

## Web IDE / Theia

* open url from email (different for everyone)
* login with:
  * user: participant
  * password: BrowserCoding

----

## Web IDE / Theia

> Important: In your Browser instead of http://localhost:8080/ use http://url-to-your-theia:8080/. In the terminal still use localhost.

----

## overall hints

* use `tab` key a lot - autocompletion is installed
* write the commands, don't copy paste
* get used to the commands, you'll use them a lot
* use `ctrl+c` to exit running processes
* use `exit` command if you are in remote shell
* use new terminals if you want to run things in parllel
* click on links in the Theia shell

----

## training material

path/you/should/be/in
```bash
echo "execute in shell"
kubectl exec -it database-TAB -- /bin/sh # -TAB means press tab
```

path/filename.yml
```yaml
some:
  thing:
    edit: change

# ... (there is more code on the next page)
```

----

## Windows users

* if possible use subsystem for linux shell (bash)
* or use powershell - commands may have to be sligthly adapte (ie. \ )



----

## Yaml Basics

```yaml
element1: value
object1: { object: { variable: value }}
object2:
  object2: 
    variable: value
array1: [ 'value1', 'value2', object: { variable: value }]
array2: 
  - value1
  - value2
  - object: 
    varialbe: value
```