package main

import (
    "errors"
    "github.com/gin-gonic/gin"
    "github.com/ip2location/ip2location-go/v9"
    "net"
)

const ERRCODE = -1
const SUCCESSCODE = 1
type IpResponse struct {
    Country string
    Region  string
    City    string
}


func main() {
    db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB3.BIN")
    if err != nil {
        panic(err)
        return
    }
    defer  db.Close()
    r := gin.Default()
    r.GET("/ip/:address", func(c *gin.Context) {
        address := c.Param("address")
        if _, err := validateIpv4Address(address); err != nil {
            c.JSON(200, gin.H{
                "error": err.Error(),
                "code": ERRCODE,
            })
            return
        }
        results, err := db.Get_all(address)
        if err != nil {
            c.JSON(200, gin.H{
                "error": err.Error(),
                "code": ERRCODE,
            })
            return
        }
        var resp IpResponse
        resp.Country = results.Country_long
        resp.Region = results.Region
        resp.City = results.City
        c.JSON(200, gin.H{
            "error": nil,
            "code": SUCCESSCODE,
            "data": resp,
        })
        return
    })

    r.Run("0.0.0.0:2222")
}

func validateIpv4Address(ip string) (result bool, err error) {
    address := net.ParseIP(ip)
    if address != nil {
        return true, nil
    } else {
        err = errors.New(ip + "is not a legal ipv4 address")
        return false, err
    }
    return
}