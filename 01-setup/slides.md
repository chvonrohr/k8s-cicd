## Training Setup
### Workflow and environment

----

## Workflow

1. We show and explain
2. You get all the slides and infos
3. You walk through it during Exercise slots

> we focus on explaining while hands-on

----

## Web IDE / Theia

* open: http://FIRSTNAME.sk.letsboot.com:3443
* login with:
  * user: FIRSTNAME
  * password: `Workshoptage.`

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
project-start/ (starting point of our project)
project-vision/ (NOT the solution)

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
* gitlab project for CI/CD (forked)
* gitlab authentication

----

## training material

#### command:

project-start/web/
```bash
echo "execute in shell"
kubectl exec -it database-TAB -- /bin/sh # -TAB means press tab
```
<small>no path means you can run it in any folder</small>

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

Note:
* careful with copy paste

----

## Tips & Tricks

* careful with "common" shortcuts: cmd+w
* use `tab` key a lot - autocompletion is installed
* less copy & paste => more typing & tabbing
* use `ctrl+c` to exit running processes
* use `exit` command if you are in remote shell
* use multiple terminals if you want to run things in parallel
* click on links in the Theia shell (cmd + click)
* ignore theia extension errors (bottom right)

Note:
* use `alt+arrows left right` to navigate the commandline cursor

----

## Yaml Basics

```yaml
key1: value
object1:
  object: 
    key: value
object2: { object: { key: value }}
list1: 
  - value1
  - value2
  - object: 
    key: value
list2: [ value1, value2, object: { key: value }]
--- # new document
key1: value
```

> Careful with indentation! (2 spaces)

Note:
* easier to write version of json
* https://docs.ansible.com/ansible/latest/reference_appendices/YAMLSyntax.html#yaml-basics