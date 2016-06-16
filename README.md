#Gini
Simple and evil library for configuration files


## Description
Any of struct in you config, its Section. Section is global. If you require section two or more times in, he will take his only one time. If you need configs with same structure and different values you can use tag: `ini-name:"{Name}"`


## Tags
`ini` = tag for ignoring fields in structure

`ini-name` = tag for set another name to fields

`ini-default` = tag for set default value if value missing from config file




## Example usage
You can look how i use this in my projects here: 
[Gini-Example](https://github.com/TheMrViper/gini-example)
