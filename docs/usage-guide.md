# Contribution Steps

> [!TIP]  
> **Do you need to define a new Security Domain?**
> Check [here](new_secdom.md) first to see if you need to define new domain

1. Create a new file under the relevant directory under the taxonomy directory
2. Copy the format of an existing file or check out the [taxonomy](secdomain.go) package as the source of truth.
3. Update the details for the new env or domain
4. Update the images used in the readme
   1. `go run main.go --graph --graphDir=docs/images`

> [!WARNING]
> You will need to install golang and [graphviz](https://graphviz.org/download/).

## Decision Rational

- Rational for decisions around sensitivity or criticality need to be provided in the same file where it is define. The reason for this is that the rational will be present in the PR to provide context when reviewing any changes.
