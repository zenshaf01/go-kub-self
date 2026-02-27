# Notes

# Initializing project with modules
- When starting a new go project, first you need to initialize for modules.
- Every module needs to have a unique name.
- run: `go mod init {module name}`
- After running the above it will create a go.mod file which you should go and check.

# Creating services:
- For the initial service we are just creating the simplest service.
- The project directory structure we are following is the layered structure.
- Each layer and package should provide purpose
- You can call the file under tha package the same as the directory name.

root:
    - api 
        - services
            - sales
                - main.go
    - foundation
        - logger
            - logger.go

# main function:
- The entry point essentially calls the run function which starts the application.

# Logging:
- You need to make a decision on what sort of logging you want in your app. 
  - Do you need data in log ?
  - Do they just need to have text ?
- The above decision will dictate the type of logger you want.
- Create a new Logger at top of program and pass it around wherever you need the logger.
- Don't create a singleton logger.
- Create a logger and then pass it around the app.
- Don't use the context package to pass the logger around. Many of the layers would
need an empty context for the work they do. You don't want them checking for loggers 
on the context. That is wasted work.
- Logging is CPU, network, disk intensive.
- We generally want control on what and where the logs are coming from. So passing
the logger around manually is not a bad idea. Do it.
- To pass around the logger, you could either pass the logger into a function which
needs logging explicitly as a parameter, or pass it to a construction method and attach 
it to the struct which needs logging. And then you can access the logger from that type
when that type is used as a receiver.

# vendoring:
If you are importing thirparty packages. It is always good to import them into your vendor package.
You can run `go mod tidy && go mod vendor` to import the packages. These will get put into your
vendor folder.
