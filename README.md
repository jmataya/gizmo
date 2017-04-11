# Gizmo

Gizmo is a library that allows for the creation of versioned, view-aware
objects written in a relational style. It is a simple, safe, and productive
framework for creating content that will be present data different in different
situations and is automatically versioned on every change.

## Getting Started

1. Download and install it:

    ```sh
    $ go get github.com/FoxComm/gizmo
    ```

1. Run the database migrations:

    ```sh
    $ gizmo setup
    ```

1. Import gizmo to your code:

    ```go
    import "github.com/FoxComm/gizmo"
    ```

1. Add `gizmo.EntityObject` to the struct you want to save:

    ```go
    type Product struct {
      gizmo.EntityObject

      Title       string
      Description string
    }
    ```

1. Create `Manager` and connect to your database:

    ```go
    dbHost := "127.0.0.1"
    dbName := "gizmo"
    dbUser := "gizmo"
    dbPassword := "I<3Gizmo"

    mgr, err := gizmo.NewManager(dbHost, dbName, dbUser, dbPassword)
    ```

1. Save a struct implementing `EntityObject` using `Manager`:

    ```go
    p := Product{
      Title:       "Some product",
      Description: "A product for demo purposes",
    }

    viewID := 1
    saved, err := mgr.Create(p, viewID)
    ```

1. Retrieve the struct based on it's ID and view ID:

    ```go
    id := 1
    viewID := 1

    product := Product{}
    err := mgr.FindByID(id, viewID, &product)
    ```
