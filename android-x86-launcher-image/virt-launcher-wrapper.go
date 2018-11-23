package main

import "syscall"
import "os"
import "fmt"
import "os/exec"

func main() {
    if os.Getenv("GRAPHICS_INITIALIZED") != "1" {
        fmt.Println("Granting qemu permissions to /dev/dri/renderD128")
        cmd := exec.Command("/usr/bin/setfacl","-m", "u:qemu:rw", "/dev/dri/renderD128")
        cmdOut, err := cmd.Output()

        fmt.Println(string(cmdOut))

        if err != nil {
            panic(err)
        }

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
