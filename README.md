#Gini

## Description
Any of struct in you config, its Section. Section is global. If you require section two or more times in, he will take his only one time. If you need configs with same structure and different values you can use tag: `ini:"{Name}"`


### Example config
```ini
[Sub2]
Field   =   value
[SubConfig]
Field   =   value
[MySQL]
User    =   root
Pass    =   ****
[Other]
Num      =   123456   
Field1   =   value1
Field2   =   value2
```

### Code
```go
package main 

import (
    "fmt"
    
    "github.com/TheMrViper/gini"
)

type MySQLConfig struct {
    User string
    Pass string
}
type SubConfig struct {
    Field string
}
type OtherConfig struct {
    Num1    int
    Field1  string
    Field2  string
    
    SubConfig SubConfig
}
type MainConfig struct {
    SubConfig1   SubConfig   
    SubConfig2   SubConfig   `ini:"Sub2"`
    MySQLConfig  MySQLConfig `ini:"MySQL"`
    OtherConfig  OtherConfig `ini:"Other"`
}
func main() {
    config := MainConfig{}
    // Reading from file
    if err := gini.ReadConfig("config.ini", &config); err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Println(config)
    
    // Save to file
    if err := gini.WriteConfig("config.ini", &config); err != nil {
        fmt.Println(err)
        return
    }
    return
}
```

