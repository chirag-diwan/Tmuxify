# Tmuxify

## The motivation
I have always wanted a utility that could be intelligent regarding starting my tmux session i.e. it can have diffrent tmux session setup based on diffrent projects , opening diffrent stuff in diffrent windows , running specific commands etc.

Tmuxify full fills all those needs , it uses a .tmuxify.toml file per project that tell it exactly how to setup the session when ran.

## .tmuxify.toml
The session setup file is just a single toml file with a session filed and a windows array
```toml
    [session]
    name = "tmuxify" #specific the session name (REQUIRED)
    main = 0 #the window that will be selected when new session will be created (0 indexed)
    
    #The window , every field in window is OPTIONAL


    [[window]]
    name = "nvim" 
    cmds = ["nvim ."]
    
    [[window]]
    name = "cmd" 
    cmds = ["ls"]
```

## Application configuration
Tmuxify needs the home dir path to be added to the `PATH` variable , if not then Tmuxify will fall back to its internal defaults (or it may error , depending on how much information it was able to get)

Tmuxify needs a top level configuration at `/home/user/.tmuxify-conf.toml` .
All the options are OPTIONAL you may or may not specific them though specifying them is prefered.(If not specified then tmuxify will use the internal defaults)

```toml

roots = ["Projects" , ".config" , "Streaming"] # The search paths , dont need to prefix them with /home/user/
ignore = [ ".git", "node_modules", ".cache", ".bun", ".cargo", ".wrangler"] # The paths to ignore IN ANY DIRECTORY
max_depth = 4 # Maximum depth to search

```

# Installation
```bash

git clone https://github.com/chirag-diwan/Tmuxify.git
cd Tmuxify
go build -o out/tmuxify
mv out/tmuxify ~/.local/bin/ # move the binary to search path for shell

```

## Reporting Issues
If you encounter a bug or have a feature request, please open an issue describing:

What you expected to happen.
What actually happened.
Steps to reproduce the issue.
Your operating system and Go version (if relevant).
