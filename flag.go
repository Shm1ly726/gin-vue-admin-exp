package main

import (
    "flag"
    "fmt"
)

const banner = `
┌─┐┬┌┐┌   ┬  ┬┬ ┬┌─┐  ┌─┐┌┬┐┌┬┐┬┌┐┌   ┌─┐─┐ ┬┌─┐
│ ┬││││───└┐┌┘│ │├┤───├─┤ │││││││││───├┤ ┌┴┬┘├─┘
└─┘┴┘└┘    └┘ └─┘└─┘  ┴ ┴─┴┘┴ ┴┴┘└┘   └─┘┴ └─┴  
                     gin-vue-admin-exp version: v0.1
                     
`

func Banner() {
    fmt.Print(banner)
}

func userHelp() {
    fmt.Fprintf(flag.CommandLine.Output(), `使用方法: gin-vue-admin-exp.exe`)
    flag.PrintDefaults()
}

func Flag(Info *HostInfo) error {
    Banner()
    flag.StringVar(&Info.Url, "u", "", "URL address of the host you want to scan,for example：https://www.baidu.com")
    flag.StringVar(&Info.Token, "x", "", "Token for authentication")
    flag.Parse()
    if Info.Url == "" {
        fmt.Println("URL is none")
        flag.Usage = userHelp
        flag.Usage()
        return fmt.Errorf("URL is none")
    }
    if Info.Token == "" {
        fmt.Println("Token is none, using frontEndBypass")
        frontEndBypass(Info)
    } else {
        fmt.Println("Token is provided, using exp")
        exp(Info)
    }
    return nil
}