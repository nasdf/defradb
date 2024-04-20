# Global scope

WASM environment has access to anything defined on the global scope.

```go
import "syscall/js"

func main() {
    // https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API/Using_IndexedDB#opening_a_database
    result := js.Global().Get("indexedDB").Call("open", "MyTestDatabase", 3)
}
```

Values can be shared from WASM by setting properties on the global scope.

```go
import "syscall/js"

func main() {
    js.Global().Set("sum", js.FuncOf(func(this js.Value, args []js.Value) any) {
        return args[0].Int() + args[1].Int()
    })
}
```

# Promises

JavaScript promises allow for asynchronous execution of code.

This is important because JavaScript executes on a single thread in the browser.

Promise chaining:

```js
function() {
    fetch()
        .then(res => console.log(res))
        .catch(err => console.error(err))
}
```

Async functions:

```js
async function() {
    try {
        const res = await fetch('http://localhost')
        console.log(res)
    } catch (err) {
        console.error(err)
    }
}
```

# Demo

```js
const db = await defraDB.open()

await db.addSchema('type User { name: String }')

const lensPath = 'http://localhost:8080/rust_wasm32_copy.wasm'

const lensConfig = { lenses: [{ Path: lensPath, Arguments: { src: "name", dst: "fullName" } }] }

await db.addView('User { name }', 'type UserView { fullName: String }', lensConfig)

const users = await db.getCollectionByName('User')

await users.create({ name: "John" })

await users.create({ name: "Fred" })

await db.execRequest('query { UserView { fullName } }')
```
