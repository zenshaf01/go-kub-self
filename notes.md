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

# Kubernetes:
Kubernetes is an orchestration tool for containerized applications.
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
  - Create docker images for services
  - Create & Apply the deployments for the Pod
  - 


