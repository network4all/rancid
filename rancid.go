package rancid

import (
   "fmt"
   l "github.com/network4all/logerror"
   "strings"
   "os/exec"
   "strconv"
)

type ArpEntry struct {
   Address string
   MacAddress string
   Vlan string
}

type RoutingEntry struct {
   Address string
   Subnet string
   Vlan int
   Name string
}

type MacEntry struct {
   Vlan int
   Mac string
   Link string
}

type Nic struct {
   Name string
}


func ExecuteCommand3DES (command string, ipaddress string) string {

    //runned := fmt.Sprintf("Executed on %s command: %s\n", ipaddress, command)

    app := "/usr/local/rancid/bin/clogin"
    cmd := exec.Command(app, "-c",command,ipaddress)
    stdout, err := cmd.Output()
    l.CheckErr(err)
    return string(stdout)
}

func ExecuteNSCommand (command string, ipaddress string) string {
    app := "/usr/libexec/rancid/nslogin"
    cmd := exec.Command(app, "-c",command,ipaddress)
    stdout, err := cmd.Output()
    l.CheckErr(err)
    return string(stdout)
}

// text functions
func Grep (input string, filter string) string {
    var output string
    lines :=strings.Split(string(input),"\n")
    for i := range lines {
       if strings.Contains(lines[i], filter) {
          output = output + lines[i] + "\n"
       }
    }
    return output
}

func DGrep (input string, filter1 string, filter2 string) string {
    var output string
    lines :=strings.Split(string(input),"\n")
    for i := range lines {
       if strings.Contains(lines[i], filter1) || strings.Contains(lines[i], filter2) {
          output = output + lines[i] + "\n"
       }
    }
    return output
}

func NGrep (input string, filter string) string {
    var output string
    lines :=strings.Split(string(input),"\n")
    for i := range lines {
       if !strings.Contains(lines[i], filter) {
          output = output + lines[i] + "\n"
       }
    }
    return output
}

func SpaceFieldsJoin(str string) string {
    return strings.Join(strings.Fields(str), "")
}

func BeforeDelimiter (mydata string, delimiter string) string {
    return strings.Split(mydata, delimiter)[0]
}

func AfterDelimiter (mydata string, delimiter string) string {
    return strings.Join(strings.Split(mydata, delimiter)[1:], delimiter)
}

func Between (mydata string, after string, before string) string {
    return AfterDelimiter(BeforeDelimiter(mydata, before), after)
}

// wireless
type ApEntry struct {
   Name string
   Address string
   MacAddress string
   DeviceType string
   ClientCount int
}

func ExecuteWlcCommand (command string, ipaddress string) string {
   app := "/usr/local/rancid/bin/clogin"
   cmd := exec.Command(app,"-c",command,ipaddress)
   stdout, err := cmd.Output()
   l.CheckErr(err)
   return string(stdout)
}

func GetAPlist (ipaddress string, apTable *[]ApEntry)  {
   // better filter
   aplist := Grep(ExecuteWlcCommand("show ap summary", ipaddress), "NL")

   lines :=strings.Split(string(aplist),"\n")

   for i :=range lines {
      line := lines[i]
      //todo: regex?
      for x:=0; x<15; x++ {
         line = strings.Replace(line, "  ", " ", -1)
      }
      s:= strings.Split(line, " ")
      if len(s)>1 {
         cc, _ := strconv.Atoi(s[len(s)-1])
         *apTable = append(*apTable, ApEntry{Name:s[0], DeviceType:s[2], MacAddress:s[3], Address:s[len(s)-2], ClientCount:cc})
      }
   }
}

func CC (gconfig string, configitem string, mustexist bool) (bool, string) {
   // checkconfig
   if (mustexist) {
     if strings.Contains(gconfig, configitem) {
        return false, ""
     } else {
        return true, fmt.Sprintf("%s missing\n", configitem)
     }
   } else {
     if !strings.Contains(gconfig, configitem) {
        return false, ""
     } else {
        return true, fmt.Sprintf("%s exists\n", configitem)
     }
   }
}


