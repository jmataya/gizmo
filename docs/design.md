
# Gizmo Design

Gizmo is a library that allows for the creation of versioned, view-aware
objects. Those components can either be used as raw objects or converted to a
typed entity that can be used in a relational style. The aim is to make working
with flexible, branchable, and versioned objects easy for both ourselves and our
customers.

## Background and Motivation

This library is a direct descendent of [The Enlightened, Post-Modern Product
Model](https://github.com/FoxComm/highlander/blob/master/documents/design/product/product_model.pdf)
and [Product Model Versioning](
  https://github.com/FoxComm/highlander/blob/master/documents/design/product/product_model_with_versioning.pdf)
concepts that were implemented in Phoenix across entities like products, SKUs,
promotions, and discounts.

This model builds on that legacy by attempting to fix some of the more glaring
issues in the previous implementation, such as:

* Inability to reason about full entities;
* Poor encapsulation of implementation details;
* Lossy versioning between entities;
* Confusing naming conventions.

### Inability to Reason About Full Entities

There was no clear way to read and understand which properties and entities made
up, for example, a `Product`. Much of the pattern of working with entities
involved manually manipulating the properties, through form and shadow, and
managing relationships with object links.

Then only way to get a full understanding of the entity was to read through the
logic for retrieving and manipulating the entity. No bueno.

### Poor Encapsulation Of Implementation Details

An inadequate encapsulation of implementation details, such as `ObjectForm`,
`ObjectShadow`, and `Illumination` _(for details on those concepts, see the
documents linked above)_, which resulted in convoluted, duplicated code that
as difficult to debug.

### Lossy Versioning Between Objects

A lossy concept of how versioned objects relate to each other, meaning that
we had complicated algorithms to keep versioned objects in sync, or we lost
the data necessary to reconstruct full objects.

### Confusing Naming Conventions

We have a number of poorly fleshed out names and concepts. For example, what is
the difference between a `headID`, `formID`, and `commitID`? When is the correct
time to use them?

## Principles

Gizmo is designed around the idea that a flexible, versioned entity whose
structure can be entirely controlled by the implementing code is beneficial,
whether that's a service that implements product catalog, or an generic object
service presented to customers.

It aims to build on prior implementations by adhering to a simple set of design
philosophies:

* An entity is more than just a collection of properties;
* The content and relations of an entity are separate;
* The implementation details of branching and versioning should be opaque to the
  user;
* It should be easy to write understandable and safe abstractions atop the
  generic entity.

Let's go through each principle one-by-one.

### An Entity Is More Than Just a Collection of Properties

Let's start by considering an example that's near and dear to our hearts:
products. A product is a collection of properties, such as its name,
description, and slug. It also contains a reference to the item in inventory
that it represents, which we call a SKU. For the purposes of this example, we'll
ignore details like custom properties and variants.

Since this library is in Go, here's an example of how a product might be
structured:

```Go
type Product struct {
  gizmo.EntityObject

  Title       string
  Description string
  Slug        string

  SKU SKU
}

type SKU struct {
  gizmo.EntityObject

  Code        string
  InStock     bool
  RetailPrice Money
  SalePrice   Money
}
```

Notice that glancing at the above example gives a very clear picture of how the
product is structured. Beyond knowing its basic properties, we know that there
is an object called SKU that contains inventory and price information.

### Content and Relations of an Entity are Separate

The versioned, post-modern product model that is described above is perfect for
handling the storage and manipulation of simple properties on an entity. In
order to have a more complete picture of the entity, both the properties and
relations need to be versioned.

Consider the code in the previous example, in that example we would have four
artifacts versioned:

1. SKU Content
  - Contains `Code`, `InStock`, `RetailPrice`, and `SalePrice`
  - Uses `ObjectForm`, `ObjectShadow`, and `ObjectCommit`
1. SKU Entity
  - Identifies which version of the SKU content to use
  - Stores no relation, because SKU has no relation
1. Product Content
  - Contains `Title`, `Description` and `Slug`
  - Uses `ObjectForm`, `ObjectShadow`, and `ObjectCommit`
1. Product Entity  
  - Identifies which version of the Product content to use
  - Identifies which version of the SKU entity to use

### Implementation of Branching and Versioning Should be Opaque

Pretty simple: users should only deal with the generic illuminated entity, which
might have a JSON representation looking like:

```JSON
{
  "id": 1,
  "commit_id": 13,
  "view_id": 1,
  "attributes": {
    "title": {
      "type": "string",
      "value": "My Product"
    },
    "description": {
      "type": "text",
      "value": "I'm a product!"
    },
    "slug": {
      "type": "string",
      "value": "my-product"
    }
  },
  "relations": {
    "skus": [23]
  }
}
```

Users could also interact with structs like `Product` and `SKU` above.

### Understandable and Safe Abstractions

We want it to be super easy to use entities, so user should be able to interact
with the types of object illustrated by `Product` and `SKU`.

## Outline

* Background
* Objective
  * Reasoning
  * Goals
  * Non-Goals
* Architecture
* Example Usage
* Roadmap
