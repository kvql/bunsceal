# Working with GraphViz

Images with GraphViz are made up of below primary components

- Graphs/subgraphs
- Nodes
- Edges

To fully understand how the graphs are constructed. It's best to view the below image with all the hidden elements shown in orange.  

To generate all images with this shown update the variable `visibility` in the below file:  
[vis-common.go](vis-common.go)

![Example Image](not_auto_updated_example_image.png)

## Other information

We are using the [dot, layout algorithm](https://graphviz.org/docs/outputs/canon/)

Attributes for the graph components define how the final image looks. What attributes can be set and what values are allowed are documented here:  
https://graphviz.org/doc/info/attrs.html
