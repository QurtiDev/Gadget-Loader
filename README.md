<h1 align="center">Gadget Loader</h1>
<p align="center">Stealthy Go based shellcode loader!</p>

<p align="center">
<img width="225" height="224" alt="image" src="https://github.com/user-attachments/assets/65521d48-7ae5-416d-9900-43accb5d296b" />
</p>

### Desc.


This program is a Windows shellcode loader written in Go. It downloads a shellcode blob from a provided URL, injects it into the current process' memory and executes it. It benefits from the usage of [BananPhone](https://github.com/C-Sto/BananaPhone) to make indirect system calls instead of calling the Windows API directly, which helps to support the goal of avoiding detection by EDR/XDR solutions! 

>[!WARNING]
>Prerequisites:
>- [Go](https://go.dev/doc/install)
>- [Garble](https://github.com/burrowers/garble)
>- [UPX](https://github.com/upx/upx)
>
> You're expected to have these 3, otherwise you will face problems, please view their respected info and learn how to set them up if you don't have them already installed.

<br>

## Usage:



### Step 0x0
Build your shellcode blob to use, something like Sliver is absolutely amazing to use for C2, Mythic is also great!!
<br>
This step is out of scope, you're expected to have it already.


### Step 0x1
Clone the repo

```
git clone https://github.com/QurtiDev/Gadget-Loader/
```

### Step 0x2
Edit the URL with yours
<img width="1188" height="432" alt="image" src="https://github.com/user-attachments/assets/ba3b99d7-fefc-4ee6-8eba-e99d3672eb6c" />

### Step 0x3

Build it:<br>
if you're on Linux and you need to cross-compile:

### Build command on Linux

Non-garble

```
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-H=windowsgui -s -w" -o updates.exe .
```

Garble
```
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 garble -tiny -literals -seed=random build -trimpath -ldflags="H=windowsgui -s -w" -o updates.exe .
```

<br>
<br>
<br>

### Build command on Windows

<br>
<br>
Non-garble
<img width="1128" height="73" alt="image" src="https://github.com/user-attachments/assets/fe45feba-6480-470e-8ff5-b4abe87a5a4b" />

```
$Env:GOOS="windows"; $Env:GOARCH="amd64"; $Env:CGO_ENABLED="0"; go build -tags windows -trimpath -ldflags="-s -w -H=windowsgui" -o updates.exe . 
```



<br>
<br>
<br>
Garble:

```
$env:GOOS="windows"
$env:GOARCH="amd64"
$env:CGO_ENABLED="0"
garble -tiny -literals -seed=random build -trimpath -ldflags="-H=windowsgui -s -w" -o updates.exe .
```
<br>

> [!TIP]
> Reverse engineers can use any leftover debug symbols that may still be present even after using -s. Running the following strip command can help remove a lot of symbol data!! <br>
> <br>
> - --strip-all - removes As much symbol data as possible <br>
> - -R .note.go.buildinfo -R .gosymtab - Removes specifically Go-related sections very handy for us :)
<br>

Command:
```
strip --strip-all -R .note.go.buildinfo -R .gosymtab updates.exe
```
<br>


![gif](https://tenor.com/tA16roB1dMO.gif)



<br>
### Step 0x4
Optional step:
If you want you can now pack it as Garble built binaries tend to flag by itself but I've noticed that if you pack them with a legitimate and known packer like UPX you can actually bypass static detection really well



```
upx --best --ultra-brute updates.exe -o updatespacked.exe
```


### Step 0x5
Host your payload, you can use something like catbox or filebin for direct hosting or you can also use a python web server such as 

```
python -m http.server 8080
```


<br><br><br>

### Workflow in short:
Downloads shellcode from a remote URL (value is set in the loader.go file on line 111).
<br>Uses the [BananPhone](https://github.com/C-Sto/BananaPhone) library to perform direct system calls without using Windows API directly - Evasion!
<br>Allocates memory with PAGE_EXECUTE_READWRITE permissions.
<br>Writes the downloaded shellcode to the allocated memory.
<br>Executes the shellcode by creating a function pointer to the memory location.


<br><br>
# ⚠️ Usage RULES & Disclaimer ⚠️

This project is intended **ONLY** for authorized use and educational purposes.  
Do **NOT** run this loader on systems you do NOT own, OR without **explicit WRITTEN permission** to test systems.

Any misuse of this loader against unauthorized targets is **strictly prohibited** and may even be illegal.  
I assume **no liability** for any misuse or damages caused by the use of this code. 
