# wiz

wiz is a simply configured wizard providing gui for runing [tasks](https://github.com/go-task/task)

A small fyne application easy to cross-compile and easy to distribuite.

A software wizard or setup assistant is a user interface type that presents a user with a sequence of dialog boxes that lead the user through a series of well-defined steps. Tasks that are complex, infrequently performed, or unfamiliar may be easier to perform using a wizard.

(quoted from **[wikipedia](https://en.wikipedia.org/wiki/Wizard_(software))**)

The wiz configuration is a yaml file
 
 ## wizard defintion

A wizard is defined by a yaml file, by default the wizard definition is calle wiz.yaml.
wiz will look for wiz.yaml in the currente working directory or in the path specified by the -d flag.

every with has tehe folowing attrubute:
+ **title**: The title displaye in the wizard dialog box.
+ **taskfile**: The path (relative ti the wizzard definition) where the taskfile executed shall be.
+ **menu**: The title for the menu diplayd in the  dialog box. The wizard is a collection of panels and each panel has an entry in the menu.
+ **minisize**: The minum size of the dialog box. This is a structured attribute.
    + **w**: The dialog box minimum width (point).
    + **h**: The dialog box minimum heighth (point).
+ **msg**: A string displayed  in the home page the inital panles displayed by the wizard dialog box.
+ **panels**: A list of pannels.


### Panel definition

 Each panel has this attributes:
 

 
  ```yaml
title: wizzard
taskfile: wiztask.yaml
menu: Procedure
minisize:
  w: 600
  h: 600
msg: |
  Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
panels:
  - title: test
    form:
      -
        name: text
        label: text
        type: text
        value: velit rutrum elit
      -
        name: file
        type: file_save
        value: 
      -
        name: dir
        type: dir
        value: .
        options:
          - .shp
      -
        name: select
        type: select
        options: 
          - a
          - b
          - c
          - d
        value: 
      -
        name: avvio
        type: execute
        action: 
          execute: task
        value: 
      
  - title: test1
    form:
      -
        name: text1
        label: text
        type: text
        value: prova
      -
        name: file1
        type: file_open
        value: 
        options:
          - .yaml
          - .yml
          - .json

      -
        name: dir1
        type: dir
        value: .
      -
        name: select1
        type: select
        value: 
        action: 
          execute: selectvalue
      -
        name: VIA
        type: execute
        value: 
        action: 
          execute: task
      
  - title: test2
    form:
      -
        name: text2
        label: text
        type: text
        value: prova
      -
        name: file2
        type: file
        value: 
      -
        name: dir2
        type: dir
        value: .
      -
        name: select
        type: select
        options: 
          - a
          - b
          - c
          - d
        value: 
      -
        name: back
        type: back
        value: 
        action: 
          goto: test
      
  - title: test3
    form:
      -
        name: text
        label: text
        type: text
        value: prova
      -
        name: testo
        label: provas
        type: text
        value: prova
      -
        name: file
        type: file
        value: 
      -
        name: dir
        type: dir
        value: .
      -
        name: select
        type: select
        options: 
          - a
          - b
          - c
          - d
        value: 
      -
        name: back
        type: back
        value: 
        action: 
          goto: test
        

  ``` 