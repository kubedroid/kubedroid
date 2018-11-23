package main

import "syscall"
import "os"
import "fmt"
import "os/exec"

func main() {
    if os.Getenv("GRAPHICS_INITIALIZED") != "1" {
        fmt.Print("Attempting to initialize graphics\n")

        cmd := exec.Command("/bin/bash","/enable-graphics.sh")
        cmdOut, err := cmd.Output()
        if err != nil {
            panic(err)
        }

        fmt.Println(string(cmdOut))

        os.Setenv("GRAPHICS_INITIALIZED", "1")
    } else {
        fmt.Println("Skipping graphics initialization")
    }

    fmt.Print("Starting the virt-launcher wrapper\n")

    args := os.Args
    env := os.Environ()
    binary := "/usr/bin/upstream-virt-launcher"


    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        panic(execErr)
    }
}
