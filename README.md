# mdns-todo-api
Simple to-do list REST API. Uses boltDB for storage.

### Install dependencies:

    go mod download

### Run app:
    go run main.go 

URL: http://localhost:8090

### Endpoints:

| Name  | Method | Path | Paylaod |
| --- | --- | --- | --- |
| Health  | GET | /health | - |
| Get list | GET | /items | - |
| Add item | POST | /items | {content: "content here" } |
| Delete item | DELETE | /items/{id} | - |


#### Other -half finished- parts of the project:
- UI: https://github.com/srgyrn/mdns-todo-ui
- Automation: https://github.com/srgyrn/mdns-todo-ui-automation
