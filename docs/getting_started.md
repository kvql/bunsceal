# Getting Started Guide for Anyone new to the topic

> Just need to make a change to a Security Environment or Domain?  
> See here:  
> [How to contribute or make changes](contributing.md)

**Key Concept:**

- Security Environments shouldn't be able to connect with each other (except Shared Service), they are completely isolated
- Security Domains just represent subsets of an SecEnv and should be able to connect with each other. However, those allowed connections should be explicit and denied by default.
  - Security Domain is just the name given to internal segments of our infra which have different compliance/security requirements.

TODO OSCAL definition of requirements


## Criticality

Before going into Security Domains, it would be useful to read up on what we mean by criticality.  
[Criticality Overview](criticality.md)

## What is a Security Domain?

It's easier to start by thinking of this from a network perspective

TLDR:  

- Think of security domain as a VPC (or a group of VPCs) within an SecEnv.
- It is logical name for a segment within our Security Environments
- Why we need a name for these segments is due to them having different security and compliance requirements and we need to have some logic concept to map those requirements too.

### When should I define a new security domain?

[Defining a new Security Domain](new_secdom.md)
