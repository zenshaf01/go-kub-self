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
    - zarf
        - docker
            - Dockerfile
        - k8s
            - base (this directory acts as a container for configs that will present in all envs (dev, staging, prod)
                - sales (pod)
            - dev
                - sales (pod)

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

# Kubernetes:
Kubernetes is an orchestration tool for containerized applications.
Kubernetes is an asynchronous event driven system.
- Architecture:
    - Control Plane: Sits at the top of clusters and manages it to ensure desired state automatically.
    The control plane runs on multiple nodes across datacenters.
        - Components:
          - API Server (kube-api-server): This is the receiver of all frontend requests either form CLI or UI. If the API server is
          down, the cluster will still be up but it cant accept management requests. REST Api.
          - etcd: distributed key value store. It stores cluster state. 
          All other components of the control plan use the etcd to store information about cluster in etcd
          - Scheduler: The scheduler schedules the creation and management of pods onto worker nodes based on the
          required resources by the pod (and the containers inside) and the available resources on the nodes.
          The scheduler makes decisions based on the availability of the resources.
          - Controller Manager: The Controller manager manages controllers. An example of a controller is a replication
          controller which makes sure the pods have enough replicas running. 
          - Container runtime: This runs the containers in the pods. It downloads the images, start and stops the containers in the pod.
          - kubelet: daemon that runs on teh worker nodes. kubelet is responsible for getting instructions
          from the control plan and makes sure the desired state of the pod is maintained.
          - kube-proxy: this is network proxy to route traffic to the correct pod.
    - Worker Nodes: These execute workloads, running containers inside pods. Pods are inside these.
    - Pods are inside Nodes and contain the containerized applications. Pods are the smallest deployable unit in K8s.

Cluster: A cluster is a enclosed space of resources which you can use to deploy your app to
A cluster might be shared or dedicated.
- So A cluster is the first thing to create for yourself when you want to develop with kubernetes.
- A cluster is always backed by raw computing resource to get work done.
- The above compute power for the cluster comes from something called Nodes.
- Each node is backed by physical compute power of the host machine, and the docker VM. Think of it 
as a 3 layer architecture. The machine (part of it) backs the docker vm and the docker vm provides compute to the applications we run on them.
- In dev cluster should at least have access to 1 Node.
- IN staging there should at least be 3 Nodes.
- Each Node will run our application, so if we have 3 Nodes we have 3 instances of our application running.
- The application is run on these Nodes using something called Pods.
- A one to many Pod can have one to Many services inside it.
- We deploy services in Pods, These Pods define the deployment information and run configurations for the service. 
In the Pod we define how much cpu, ram, disk should this be using, env variables. It also defines the networking for inter service
communication running in the same pod.
- Each Cluster can have many Podes and Each Pod can have its instances being run on many Nodes. Each Pod can have 
many services being run on it.
- We only put things in Pods which need be to restarted on failure. So don't put DS's in Pods in staging.
You can put them in Pods in dev but not staging or prod.
- You can start a kubernetes cluster using kind. KIND stands for KubernetesInDocker.
- `kind` creates a local kubernetes cluster using docker.
- Steps:
  - Create the cluster
  - Create docker images for services (create binaries / images for our application services) (build dockerfiles)
  - Create / Define & Apply the deployments for the Pod
    - You need some tooling for this
    - We will use kustomize to manage deployments

# Kubernetes Quotas:
Kubernetes allows you to tell how much cpu and memory a service (container) gets.
This is done through quotas.
- We want to tell kubernetes that we want to give a service some percentage of the cpu. Like saying S1 takes 25% of CPU, S2 takes 25% and S3 takes 50%.
- But how do you achieve that ?
- What kubernetes or other OS's do is that:
  - For each cpu core, they define 100ms of indefinite time cycles, Now when you say I want a service to take 100% of cpu time. The service will end up having the full 100ms per cycle.
  - If we say we want S1 to take 25% and S2 to take 75%. S1 will take 25ms and S2 will get 75ms of time on each of those 100ms time cycles / slices.
  - We can also specify a value over 100%. If we say S3 take up 150% of cpu, that means (if the machine has 2 cores) the service will get a full 100ms on core1 and 50ms on core 2 each cycle.
  - The amount of time the service gets is not tied to a core. So if a service is allotted 100%, it could get 50ms on one core and 50 ms on the second core. but the guarantee is that 
  it will get a full 100ms.
- With Go we are always matching the number of OS threads with the number of Cores. Co imagine if we had more threads than cores, say we had 2 threads on a single core and we say,
S1 gets only 25ms out of 100ms, that means that since we have 2 threads on 1 core. that 25ms will get halved between the 2 threads, leaving us with less than 12ms per thread, why less
than 12ms ? coz the threads will have to be context switched which will waste time and cpu cycles.


# Background knowledge:
The go scheduler:
- When the go runtime starts, it asks the machine how many cores it has.
- The go run time will create N OS threads for N cores, in order to run goroutines.
- For every OS thread, the runtime will create a P (represents a logical processor (is actually the context and is the scheduler)), The P will be attached to the M (represents OS thread), creates 1 G (the main go routine)
- the main G can create other G's. Once created, they sit in LRQ (Local run Queue) until they get a chance to get executed by the scheduler.
- Each P has its own LRQ.
- At P level only a fixed number of G's can end up on the LRQ. Once all LRQ's of all P's are full, The G's start ending up on the GRQ (Global Run Queue). 
- As we get space in LRQ, G's are transferred from the GRQ to LRQ before they can be put on P tro be executed.
- G is application level thread.
- G has a stack, has those thread states. Runs in user level.
- All application code runs in a G. The G gets attached to the M. The M is attached to the CPU Core which is executing instructions.
- Context switch takes place when an 1 M is taken of the physical core and another M is put onto it so that the work on M2 can be performed. Or at go level,
context switch involves having G1 taken off from M and G2 being put on it (the context switch is 200nano / 2400 instructions).
- Context switches are expensive as they waste cpu cycles as no work is being performed during the switch. (1000 ns or 1 microsecond or 12000 instructions)
- 3Ghz means 3 clock cycles / nanosecond, the clock cycle is the main thing that needs to happen for work to be done. it is the heart beat.
- we can execute 4 instruction per clock cycle.
- so if we have a 3Ghz processor, this means we can execute 3 * 4 = 12 instructions per nanosecond.
- It comes important for us to not waste these cycles and instructions in context switching.
- ms latency is very bad.
- We want our context switch to ideally happen at G level.
- The go scheduler is special in that the context switch happens on the G-M level not on the M-Core level. This means M is always doing work on Core and we waste less instructions 
by only switching at the G-M level. This also results in IO bound work to be converted to CPU bound workload.
- A context switch happens because of 2 reasons:
  - You are doing IO bound work, and the thread need to wait for something else to provide you with data (sys call).
  - The thread has used up its time slice (10 millisecond) and now has to wait for it turn again.
- There are 2 types of workloads:
  - CPU bound
    - A workload which does not require the thread to every wait.
    - We want no context switches for CPU bound tasks as that will waste time during context switch.
    - Try to do a cpu bound work with a single threaded algorithm as that will be faster coz of the absence of context switch.
    - If you want to use a multithreaded algorithm for CPU bound tasks, make sure you never use more threads than the physical cores you have. This will result in no context switches.
  - I/O bound
    - A workload which require a thread to wait for something (OS,disk, net).
    - We want a multithreaded algorithm for I/O bound workload. This allows tus to save any wasted cpu cycles which might be wasted because the thread has to wait for a sys call.
    Instead of waiting we take off the waiting thread and put another one on the M so that it can do its work. 
    - context switch's in the case of I/O bound work is less expensive than allowing a single thread to just wait for the sys call to complete and hog the cpu (wasting cpu cycles).
    - But don't throw a million threads on an I/O bound work. That will be detrimental.
    - use a thread pool (G pool in our case) and try to determine an optimal number of G's that can result in the I/O bound work to be done the fastest.
- Concurrency: Out of order execution
- Parallelism: Physically execute instructions at the same time.

# Key takeaways for above:
- Our go programs are always CPU bound. I/O bound is also converted into CPU bound workload since G's are on M's and there is no context switch between M and Core (which is mroe expensive).
- never use more OS threads than Cores.
- If we are using CPU limits, and we are limiting the service to 25ms or less than the full 100ms:
  - Our Go program should be single threaded (since it is only taking one cpu)

# HTTP 
- Use the `http.Server` to create a new server
- Use the ServeMux to create a handlers and endpoints
- The servemux's job is to:
  - take an http request, see if there is a matching url 
  - See if it has a matching handler for the incoming url path 
  - Create a new goroutine and run that handler in that goroutine
- Never use the DefaultServeMux as it has a huge security vulnerability. Any package or library can insert a handler to it.
- We should liveness and readiness handlers
- Kubernetes needs these handlers to check if the service is alive in the pod
- You should break up your route handler functions into separate packages. The handler functions and the 
routes should live in the same package.
- For all our API's the first parameter should always be the context. We need it for IO or DB or other sys calls

# Packaging
Main Ideas for engineering the project
- package design
  - This creates a firewall around API's (). 
  - We build packages that provide value.
  - There should be no common package for types
  - Each package should have their own types for data which comes and leaves the package
    - You can have 2 types of types:
      - concretes
        - Structs
        - Api accepts the data based on what the data is
      - interfaces
        - interface types
        - Api accepts the data based on how the data behaves (Polymorphic) (runtime polymorphism)
        - runtime polymorphism: when we accept defined interfaces as paramaters
        - compile time polymorphism: Usage of generic type T in go
- horizontal layering
- data relationship
- Accept interfaces, return concrete types.
- The app layer should not have any information regarding the protocol layer (Api). There shouldn't be any code accessing protocol layer code in the app code.



