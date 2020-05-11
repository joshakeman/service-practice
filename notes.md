Everything we do will have 2 to 3 levels of layering...

We want to avoid cognitive overload

Minimize abstarctions, layer to bear minimum, 

"Make things easy to understand, not easy to do"

My recommendation is to bind a project to a repo of code.

A project as we define it is bound to single repo of code... I odn't want to feel a project can only manage one binary if it's an application level project ...

The projects taht I have could be managing multiple binaries... cli, admin, data-fetching tools, etc... an entire team coudl work in that repo...

A proejct gets to define a *policy* for the code that's in it

Each project root will have two folders:
- **cmd**, application layer, responsible for three things. It's a presentation layer (meaning handling requests and responses)
    1. Requests
    2. Business Layer
    3. Handling response
- **internal**-- our business layer
-**platform** - stuff that can be used across projects... ex... logging ... I call this the *platform* layer

Layers are: 1) application level, 2) business level, and 3) foundational

We need a module file, to set up 

to initialize modules: go mod init github.com/joshakeman/service

On disk there will be a mod cache that go pls is using... gopls will cache the module cache

It's critically importnt to have an understanding of packaging in order to make the right decisions... projects often go bad without this

Back in the year 2000 Brian Kernihan (wrote C) was asked: 
    "What do you think is the biggest shortcoming in the C programming language?"

    Kernihan said... I think the biggest shortcoming in C is that it doesn't allow you to create firewalls between the different parts of your program... it's not that you can't do this (simulate object oriented programming) .. the problem is the language and compiler are not providing any help

Go says every single folder in your source tree literally represents a static library of code... a physical compiled unit of code ... the Go cache stores that stuff there.

The linker comes in and takes all those static libraries and binds them into a single program... there really isn't any hierachy in Go, egerything is flattened out in static libraries

There's nothing stopping you from doing this in other languages, but in Go, the language enforces it. What we're doing is not just compiling sourcode, but organizing APIs

Think of every folder in your sourcetree as a physical unit of code with a firewall from every other unit... we need rules for how they interact with eachother

You're not building a monolithic app any more, you're building an app of many parts

We need guidelines that are mechanically sympathetic with packaging

Packages must have a purpose, they must have a self-contained type system.

Go has defined what a unit of code is in the language ...

One discussion I will not get involved with online is what a unit test is ... everyone has their own opinion of what a unit of code is... your opinion will dicatate what you think a unit test is

Go has said you don't get to have your own definition... every folder is a unit of code in Go.

Our unit tests will focus on one unit, one package, at a time.

The more opinionated you are, the simpler things can be. Go is a highly opinionated, convention-over-condfiguration type of language

We need a module file, to set up 

to initialize modules: go mod init github.com/joshakeman/service

On disk there will be a mod cache that go pls is using... gopls will cache the module cache

First question I ask every developer...

What is the business problem you are trying to solve?

Engineering question: **What is the purpose of the logs for this application?**

Bill's answer: The only purpose I want for the logs is to be able to debug the application. I don't want to put data into the logs.

If the purpose of the logs is to debug the app, and you're not putting data into the logs, from my perspective we should have human readable logs... I want a bear minimum trace and enough context to solve the problem. I've seen people put metrics in teh logs as data then they parse those logs out for metrics... I prefer a metrics 

I see people use logging levels... I don't beleive in logging levels... either you need to log it or you don't

If you want to make the logs data that are stored in a database and have tooling around it, then yeah, you need structured logging .. you want to store historical info, metrics, etc ... if logs become data, you need structure logging

How are you gonna do it? JSON, key-value pairs, etc... where do these logs go?

Write logs to standard out ... 

Which layer of your project gets to log and which dont? IMO the application layer can, the business layer can (but should be minimized).. any loggin in platform layer is a NO NO NO

Don't create logger interfaces, it's not going to work for you.

Questions: who can log, who can't, where can we log, where can't we log?

If any function needs to log, the logger should be passed in with the highest level of precision possible... pass it as a funciton paramter first, or maybe as a receiver... we'll never hide it as a context

----

Next big question:

How are we going to handle configuration?

Do we get to restart the service if configuration changes? The answer should be yes. We want to be able to say the service can stop and be restarted if configuration changes.
This is an area of our server we want to simplify as much as possible.
Which packages access config?

Our rule will be we only have one sourcecode file that can touch configuration: main.go

You should only have to go into main to see all the configuration options

-----

Modules have privacy considerations...
You sometimes have to customize dependencies...

You dev enevironment talks to a Go Proxy, or the Module Mirror (means the same thing) ... this is sitting in Google's Cloud ... when you run go mod tidy, the go tool goes out to the Go Proxy server and asks what it knows about the packages we need ... if it does, it will send back a zip file with a module of code for a given version ... the Go tooling then decompressed the zip file in the module cache ...

Now when you run a go build and the compiler goes looking for the packages, it knows the look in the module cache for a given version.

What if Go Proxy didn't have an assembled module for any of these packages? Then it's going to go out to the version control system (in this case github) and ask for it and bring it back for that version, and it gets stored there (the Go Proxy) for ever. Every version of that code ever released will be stored in the Go Proxy server.

How does version selection work? There's basically two algorithms out there ... the one algorithm is a SAT solver ... prior to go modules, deps (the tool we all used) was a SAT solver ... it would try to identify the latest version of every module you need to build your project ... the idea is the latest, greatest version of all dependencies, and dependencies dependencies, should give you the most stable codebase.

Go doesn't use a SAT solver, it uses a different algorithm called MVS... Minimal Version Selection

The theory here is the latest but not necessarily ghreatest version of any dependency that will give you the best build ... 

Imagine our app imports package A, and a dependency for package A has a dependency package D of version 1.0... package A has its own go.mod file saying it needs 1.0 of D

When we go out looking for packages, we find the latest greatest version of D is 1.12... a SAT solver would pull the latest greatest version of D (1.12) ... Go's MVS tool instead pulls what specifically is asked for, which is 1.0 ...

Imagine we have dependency B as well which also wants D as a dependency, but B asks for version 1.6 ...

What does our algorithm do? It will change to use 1.6 for both modules.

From the compiler's point of view to build this app it needs A, B, and D 1.6 ...

But, what if you want to use 1.12 version of D? You can do that. You can use Go Get to make sure you're using that version. But there's a lot of considerations in updating dependencies... do you want to upgrade only direct dependencies? All dependencies?

Bill's Opinion: I upgrade everything. This allows you to run on the latest greatest of every package. But who's to say that's as stable as respecting your dependency tree?

Go mod tidy will always find the greatest (not necessarily greatest) then I can use Go Get to upgrade everything later on

--->
Privacy considerations:
You may not want Google to know what packages/versions you're using... how do you do that?

One thing you can do is tell the Go tooling to never go to the public server... just go direct to the Version Control System (these pulls will take longer)

The other option is you can run your own private proxy ... there's two private proxies out there today ... there's the community one called Athens, and then there's one from Jfrog part of a bigger package called ArtFactory ...

You could privately in your own network run your own Athens Proxy Server, which will in turn go to the VCS ... it may even be faster than Go Proxy and you're maintaining higher levels of privacy

What if you're running your own private VCS ? You've got a top secret code there you don't want anyone to know about ... you forget how modules work and in your dev environment you forget to tell Athens to directly go there and instead it goes to Go Proxy, which says it can't get to that private code... Go Proxy then has a record for at least 30 days that you have code that exists there (even if it couldn't pull it)

So you want to make sure you configure to either go direct or private

There are a few relevant variables of you type GO ENV...
    GOPROXY: "https://proxy.golang.org,direct"
    *notes:* Above value says if the proxy server returns a 404 or 401 (msising or gone), then if will try to go direct... if you want to always go direct just get rid of that URL, or instead put your athens proxy server... you must make sure your dev environment is configured right
    GOPRIVATE: "http://bitbucket.org/anthem"
    *notes:* maybe I'm ok with going to Go Proxy for public libraries, but for my company packages, I want it to go direct. So I will list the base URL for those packages (ie http://bitbucket.org/anthem) 
    GONOPROXY: 
    GONOSUMDB:
