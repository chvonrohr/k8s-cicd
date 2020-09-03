## Training Setup
### Workflow and environment

----

## Web IDE / Theia

* open: http://FIRSTNAME.sk.letsboot.com:3000
* login with:
  * user: letsboot
  * password: Workshoptage.

----

### Gitlab - Training material

* open https://gitlab.com/letsboot/trainings/k8s_wst/FIRSTNAME

Folder structure:
```txt
00-chapter/

01-chapter/
  slides.md (open in gitlab)
  solution/ (copy to step ahead)
    somefile.yaml

02-chapter/
...

todo-app/ (example project)
multistage-demo/ (example project)
project-start/ (starting point of our project)
project-solution/ (NOT the solution)

... other files ...
```

----

## Web IDE / Theia

http://localhost:8080/api = http://FIRSTNAME.sk.letsboot.com:8080/api

Contains:
* all tools
* authentication
* local docker environment
* local cluster
* personal Google Kubernetes Cluster

----

## training material

#### command:

path/to_to_be_in/
```bash
echo "execute in shell"
kubectl exec -it database-TAB -- /bin/sh # -TAB means press tab
```

#### file:

path/filename.yml
```yaml
some:
  thing:
    edit: change
    #... more code in file
    change_me_to: value
    element: value
# ... more code here
```

----

## Tipps & Tricks

* use `tab` key a lot - autocompletion is installed
* less copy & paste more typing & tabbing
  * get used to the commands, you'll use them a lot
* use `ctrl+c` to exit running processes
* use `exit` command if you are in remote shell
* use new terminals if you want to run things in parllel
* click on links in the Theia shell (cmd + click)

----

## Yaml Basics

```yaml
key1: value
object1: { object: { key: value }}
object2:
  object: 
    key: value
list1: [ 'value1', 'value2', object: { key: value }]
list2: 
  - value1
  - value2
  - object: 
    key: value
```