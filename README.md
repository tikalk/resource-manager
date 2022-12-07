

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

   [![Contributors][contributors-shield]][contributors-url]
   [![Forks][forks-shield]][forks-url]
   [![Stargazers][stars-shield]][stars-url]
   [![Issues][issues-shield]][issues-url]
   [![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
  <h3 align="center">Resource Manager Operator</h3>

<p align="center">
    <img src="images/logo.png" alt="Logo" width="200" height="200">
  <p align="center">
    Auto-manage  k8s resources
    <br />
    <a href="https://github.com/tikalk/resource-manager"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/tikalk/resource-manager/issues">Report Bug</a>
    ·
    <a href="https://github.com/tikalk/resource-manager/issues">Request Feature</a>
  </p>
</p>



<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgements">Acknowledgements</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

The Resource-Manager operator has built to automatically manage your kubernetes objects.

Here's why:
* You can easily create scheduled policies that will allow you to manage your kubernetes objects 
* You shouldn't be doing the same tasks over and over like deleting some pods or namespaces


<!-- GETTING STARTED -->
## Getting Started

### Installation

1. add the helm repo
   ```bash
    helm repo add resource-manager https://tikalk.github.io/resource-manager   
    helm repo update
   ```
2. Install the chart
   ```bash
   helm install resource-manager resource-manager/resource-manager-operator \ 
   --create-namespace -n resource-manager
   ```


<!-- USAGE EXAMPLES -->
##  Usage
Let's say you want to give a specific deployments only 8 hours to live.
You can create this kind of policy by applying a *ResourceManager* object.

Delete every deployment in default namespace that has the 'app=nginx' label, after 8 hours from its creation time
```yaml
apiVersion: resource-management.tikalk.com/v1alpha1
kind: ResourceManager
metadata:
  name: resource-manager-example
  namespace: default
spec:
  resourceKind: "Deployment"
  selector:
    matchLabels:
      app: nginx
  action: delete
  expiration:
    after: "8h"
```
###
### Timeframe
You can also delete a resource within a specific hour by using the 'at' key. let's say 12:00

Delete a specific deployment on 12:00, on a daily basis
```yaml
apiVersion: resource-management.tikalk.com/v1alpha1
kind: ResourceManager
metadata:
  name: resource-manager-example
  namespace: default
spec:
  resourceKind: "Deployment"
  selector:
    matchLabels:
      app: nginx
  action: delete
  expiration:
    at: "12:00"
```

### Dry-run

Add the 'dry-run' key for only validate and verify the action

```yaml
apiVersion: resource-management.tikalk.com/v1alpha1
kind: ResourceManager
metadata:
  name: resource-manager-example
  namespace: default
spec:
  disabled: false
  dry-run: false
  selector:
     matchLabels:
        app: nginx
  action: delete
  expiration:
     after: "30m"
```


---

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/tikalk/resource-manager/issues) for a list of proposed features (and known issues).



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the MIT License.



<!-- Creators -->
## Creators

[Gaby Tal](https://github.com/gabytal), [Amit Karni](https://github.com/Amitk3293), [Boris Komraz](https://github.com/bkomraz1), [Daniel Rozner](https://github.com/daniel-ro)






<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/tikalk/resource-manager.svg?style=for-the-badge
[contributors-url]: https://github.com/tikalk/resource-manager/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/tikalk/resource-manager.svg?style=for-the-badge
[forks-url]: https://github.com/tikalk/resource-manager/network/members
[stars-shield]: https://img.shields.io/github/stars/tikalk/resource-manager.svg?style=for-the-badge
[stars-url]: https://github.com/tikalk/resource-manager/stargazers
[issues-shield]: https://img.shields.io/github/issues/tikalk/resource-manager?style=for-the-badge
[issues-url]: https://github.com/tikalk/resource-manager/issues
[license-shield]: https://img.shields.io/github/license/tikalk/resource-manager?style=for-the-badge
[license-url]: https://github.com/tikalk/resource-manager/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://il.linkedin.com/company/tikal-knowledge

