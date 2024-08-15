# Dovelop AI - A way for developers to get things done!
A CLI application for developers to improve productivity by having an interactive AI-powered to do list conveniently at ones finger tips (with some typing).
Dovelop AI also helps users understand how to take smaller steps to complete bigger things by using AI to break down big tasks into smaller ones!

## Demo
![dovelop-cli](https://github.com/user-attachments/assets/53efffa6-2695-4c8e-a23f-222463651200)

## Purpose

## How to setup

Make sure that Go is installed: https://go.dev/doc/install

### MacOS
1. Make sure you're in the dovelop-CLI directory</br>
2. Create the executable:</br>
```
go build -o dovelop
```
3. Move executable to bin:</br>
```
sudo mv dovelop /usr/local/bin/
```
4. Check the PATH:</br>
```
export PATH=$PATH:/usr/local/bin
```
5. Apply new changes:</br>
```
source ~/.zshrc
```
### Windows
1. Make sure you're in the dovelop-CLI directory
2. Create the executable:</br>
```
go build -o dovelop.exe
```
3. Make custom directory for the executable</br>
```
mkdir C:\Users\{your-username}\bin
```
4. Move executable to directory</br>
```
move dovelop.exe C:\Users\{your-username}\bin
```
5. Add directory to PATH manually
```
- Search for 'edit the system environment variables' or 'PATH'
- Click 'Environment Variables'
- Click 'Path' in User variables, then click Edit
- Create new environment variable, and insert the directory you made in step 4
```
#### How to use:
Simply call dovelop in any directory to run the executable Go program: </br>
```
$ dovelop
```

## Use Cases
- Creating a new task for your to do list
- Deleting a task
- Editing a task/name of task
- View your to do list
- Prompt ChatGPT a suggestion on how to break down a task into three new and more concrete tasks
## Limitations (to be changed!)
Current limited in being able to instantly add the ChatGPT response (3 new tasks) and to replace the prompted task with those new three, to be added very soon! Will also add bash scripting for ease of installation
