

- For every resource manager object we create a `ResourceManagerHandler`   
  - we create an informer for the resource kind (ex: Deployment)  
    - for each object in the informer
      - we create a `ObjectHandler` that manage the lifecycle of the object
            - the resource handler is responsible to delete the object on the desired time



### so the flow is as following:

`ResourceManagerHandler` with informer -> for each object we create `ObjectHandler`.